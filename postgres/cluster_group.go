package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/opt/omit"
	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	dberrors "github.com/kubegrade/matrix/postgres/bob_dberrors"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

type clusterGroupRepo struct {
	db bob.DB
}

func NewClusterGroupRepo(db *sql.DB) matrix.ClusterGroupRepo {
	return &clusterGroupRepo{
		db: bob.NewDB(db),
	}
}

func (r *clusterGroupRepo) Create(ctx context.Context, clusterGroup *matrix.ClusterGroup) error {
	if _, err := bobmodel.ClusterGroups.Insert(
		r.d2b(clusterGroup),
	).One(ctx, r.db); err != nil {
		if errors.Is(dberrors.ClusterGroupErrors.ErrUniqueUniqClusterGroupOrgName, err) {
			return ekg.ErrDuplicateClusterGroupName
		}
		return err
	}
	return nil
}

func (r *clusterGroupRepo) Find(ctx context.Context, id matrix.ClusterGroupID) (*matrix.ClusterGroup, error) {
	where := bobmodel.SelectWhere.ClusterGroups
	b, err := bobmodel.ClusterGroups.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b), nil
}

func (r *clusterGroupRepo) List(ctx context.Context, filter matrix.ClusterGroupFilter) ([]*matrix.ClusterGroup, error) {
	where := bobmodel.SelectWhere.ClusterGroups
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(queryMods, where.OrganizationID.EQ(uuid.FromStringOrNil(*f)))
	}

	if f := filter.Name; f != nil {
		queryMods = append(queryMods, where.Name.EQ(*f))
	}
	bs, err := bobmodel.ClusterGroups.Query(
		queryMods...,
	).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	ds := make([]*matrix.ClusterGroup, 0, len(bs))
	for _, b := range bs {
		d := r.b2d(b)
		if err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func (r *clusterGroupRepo) Update(ctx context.Context, id matrix.ClusterGroupID, update *matrix.ClusterGroupUpdate) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	b, err := bobmodel.FindClusterGroup(ctx, tx, uuid.FromStringOrNil(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	where := bobmodel.SelectWhere.ClusterGroups
	setter := &bobmodel.ClusterGroupSetter{}

	if name := update.Name; name != nil {
		// Check if any other cluster group has the same name in the same organization
		if _, err := bobmodel.ClusterGroups.Query(
			where.OrganizationID.EQ(b.OrganizationID),
			where.Name.EQ(*name),
		).One(ctx, tx); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			// Continue operation
		} else {
			return ekg.ErrDuplicateClusterGroupName
		}

		setter.Name = omit.FromPtr(name)
	}
	if err := b.Update(ctx, tx, setter); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *clusterGroupRepo) Delete(ctx context.Context, id matrix.ClusterGroupID) error {
	where := bobmodel.SelectWhere.ClusterGroups
	b, err := bobmodel.ClusterGroups.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}
	return b.Delete(ctx, r.db)
}

func (r *clusterGroupRepo) d2b(d *matrix.ClusterGroup) *bobmodel.ClusterGroupSetter {
	return &bobmodel.ClusterGroupSetter{
		ID:             omit.From(uuid.FromStringOrNil(d.ID)),
		OrganizationID: omit.From(uuid.FromStringOrNil(d.OrganizationID)),
		Name:           omit.From(d.Name),
		CreatedAt:      omit.From(d.CreatedAt),
		UpdatedAt:      omit.From(d.UpdatedAt),
	}
}

func (r *clusterGroupRepo) b2d(b *bobmodel.ClusterGroup) *matrix.ClusterGroup {
	return &matrix.ClusterGroup{
		ID:             b.ID.String(),
		OrganizationID: b.OrganizationID.String(),
		Name:           b.Name,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}
