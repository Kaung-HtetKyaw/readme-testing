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
)

type clusterTokenRepo struct {
	db bob.DB
}

func NewClusterTokenRepo(db *sql.DB) matrix.ClusterTokenRepo {
	return &clusterTokenRepo{
		db: bob.NewDB(db),
	}
}

func (r *clusterTokenRepo) Create(ctx context.Context, token *matrix.ClusterToken) error {
	if _, err := bobmodel.ClusterTokens.Insert(
		r.d2b(token),
	).One(ctx, r.db); err != nil {
		if errors.Is(dberrors.ClusterTokenErrors.ErrUniqueClusterTokenPkey, err) {
			return ekg.ErrDuplicateIdentity
		}
		if errors.Is(dberrors.ClusterTokenErrors.ErrUniqueUniqClusterTokenOrgName, err) {
			return ekg.ErrDuplicateClusterTokenName
		}
		return err
	}
	return nil
}

func (r *clusterTokenRepo) Update(ctx context.Context, value string, update matrix.ClusterTokenUpdate) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	b, err := bobmodel.FindClusterToken(ctx, tx, value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	where := bobmodel.SelectWhere.ClusterTokens
	setter := &bobmodel.ClusterTokenSetter{}

	if name := update.Name; name != nil {
		if _, err := bobmodel.ClusterTokens.Query(
			where.OrganizationID.EQ(b.OrganizationID),
			where.Name.EQ(*name),
		).One(ctx, tx); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			// Continue operation
		} else {
			return ekg.ErrDuplicateClusterTokenName
		}

		setter.Name = omit.FromPtr(name)
	}

	if err := b.Update(ctx, tx, setter); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *clusterTokenRepo) Find(ctx context.Context, value string) (*matrix.ClusterToken, error) {
	where := bobmodel.SelectWhere.ClusterTokens
	b, err := bobmodel.ClusterTokens.Query(
		where.Value.EQ(value),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b), nil
}

func (r *clusterTokenRepo) Delete(ctx context.Context, value string) error {
	where := bobmodel.SelectWhere.ClusterTokens
	b, err := bobmodel.ClusterTokens.Query(
		where.Value.EQ(value),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}
	return b.Delete(ctx, r.db)
}

func (r *clusterTokenRepo) d2b(d *matrix.ClusterToken) *bobmodel.ClusterTokenSetter {
	return &bobmodel.ClusterTokenSetter{
		Value:          omit.From(d.Value),
		OrganizationID: omit.From(uuid.FromStringOrNil(d.OrganizationID)),
		Name:           omit.From(d.Name),
		CreatedAt:      omit.From(d.CreatedAt),
		UpdatedAt:      omit.From(d.UpdatedAt),
	}
}

func (r *clusterTokenRepo) b2d(b *bobmodel.ClusterToken) *matrix.ClusterToken {
	return &matrix.ClusterToken{
		Value:          b.Value,
		OrganizationID: b.OrganizationID.String(),
		Name:           b.Name,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}
