package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/opt/omit"
	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

type clusterTagRepo struct {
	db bob.DB
}

func NewClusterTagRepo(db *sql.DB) matrix.ClusterTagRepo {
	return &clusterTagRepo{
		db: bob.NewDB(db),
	}
}

func (r *clusterTagRepo) Create(ctx context.Context, clusterTags ...*matrix.ClusterTag) error {
	bs := make([]bob.Mod[*dialect.InsertQuery], 0, len(clusterTags))
	for _, t := range clusterTags {
		bs = append(bs, r.d2b(t))
	}
	if _, err := bobmodel.ClusterTags.Insert(bs...).One(ctx, r.db); err != nil {
		return err
	}
	return nil
}

func (r *clusterTagRepo) Delete(ctx context.Context, clusterTags ...*matrix.ClusterTag) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	for _, t := range clusterTags {
		bTag, err := bobmodel.FindClusterTag(ctx, tx, uuid.FromStringOrNil(t.ClusterID), t.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.Join(ekg.ErrNotFound, err)
			}
			return err
		}
		if err := bTag.Delete(ctx, tx); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *clusterTagRepo) d2b(d *matrix.ClusterTag) *bobmodel.ClusterTagSetter {
	return &bobmodel.ClusterTagSetter{
		ClusterID: omit.From(uuid.FromStringOrNil(d.ClusterID)),
		Name:      omit.From(d.Name),
	}
}
