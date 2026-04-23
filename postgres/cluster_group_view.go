package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

type clusterGroupViewRepo struct {
	db bob.DB
}

func NewClusterGroupViewRepo(db *sql.DB) matrix.ClusterGroupViewRepo {
	return &clusterGroupViewRepo{
		db: bob.NewDB(db),
	}
}

func (r *clusterGroupViewRepo) Find(ctx context.Context, id matrix.ClusterGroupID) (*matrix.ClusterGroupView, error) {
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

func (r *clusterGroupViewRepo) List(ctx context.Context, filter matrix.ClusterGroupViewFilter) (*matrix.ClusterGroupViewList, error) {
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

	ds := make([]matrix.ClusterGroupView, 0, len(bs))
	for _, b := range bs {
		d := r.b2d(b)
		if err != nil {
			return nil, err
		}
		ds = append(ds, *d)
	}
	return &matrix.ClusterGroupViewList{ClusterGroups: ds}, nil
}

func (r *clusterGroupViewRepo) b2d(b *bobmodel.ClusterGroup) *matrix.ClusterGroupView {
	return &matrix.ClusterGroupView{
		ID:             b.ID.String(),
		OrganizationID: b.OrganizationID.String(),
		Name:           b.Name,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}
