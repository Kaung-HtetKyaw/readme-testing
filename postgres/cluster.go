package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/opt/omit"
	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	dbinfo "github.com/kubegrade/matrix/postgres/bob_dbinfo"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/im"
)

type clusterRepo struct {
	db bob.DB
}

func NewClusterRepo(db *sql.DB) matrix.ClusterRepo {
	return &clusterRepo{
		db: bob.NewDB(db),
	}
}

func (r *clusterRepo) Create(ctx context.Context, cluster *matrix.Cluster) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	if _, err := bobmodel.Clusters.Insert(
		r.d2b(cluster),
	).One(ctx, tx); err != nil {
		return err
	}

	if info := cluster.Info; info != nil {
		if err := r.upsertInfo(ctx, tx, cluster.ID, info); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *clusterRepo) Store(ctx context.Context, cluster *matrix.Cluster) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	cols := dbinfo.Clusters.Columns
	if _, err := bobmodel.Clusters.Insert(
		r.d2b(cluster),
		im.OnConflict(cols.ID.Name).DoUpdate(
			im.SetExcluded(
				cols.UpdatedAt.Name,
				cols.LastObserved.Name,
				cols.ClusterGroupID.Name,
			),
		),
	).One(ctx, tx); err != nil {
		return err
	}

	if info := cluster.Info; info != nil {
		if err := r.upsertInfo(ctx, tx, cluster.ID, info); err != nil {
			return err
		}
	} else {
		bInfo, err := bobmodel.ClusterInfos.Query(
			bobmodel.SelectWhere.ClusterInfos.ClusterID.EQ(uuid.FromStringOrNil(cluster.ID)),
		).One(ctx, tx)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}

		if bInfo != nil {
			if err := bInfo.Delete(ctx, tx); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *clusterRepo) UpdateAll(ctx context.Context, updates ...matrix.ClusterUpdate) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	for _, u := range updates {
		clusterID := uuid.FromStringOrNil(u.ID)
		cluster, err := bobmodel.FindCluster(ctx, tx, clusterID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.Join(ekg.ErrNotFound, err)
			}
			return err
		}

		setter := &bobmodel.ClusterSetter{
			ID: omit.From(clusterID),
		}
		if clusterGroupID := u.Payload.ClusterGroupID; clusterGroupID != nil {
			setter.ClusterGroupID = omit.From(uuid.FromStringOrNil(*clusterGroupID))
		}

		if info := u.Payload.Info; info != nil {
			r.upsertInfo(ctx, tx, u.ID, info)
		}

		if err := cluster.Update(ctx, tx, setter); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *clusterRepo) Find(ctx context.Context, id matrix.ClusterID) (*matrix.Cluster, error) {
	where := bobmodel.SelectWhere.Clusters
	b, err := bobmodel.Clusters.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
		bobmodel.SelectThenLoad.Cluster.ClusterInfo(),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b), nil
}

func (r *clusterRepo) List(ctx context.Context, filter matrix.ClusterFilter) ([]*matrix.Cluster, error) {
	where := bobmodel.SelectWhere.Clusters
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		bobmodel.SelectThenLoad.Cluster.ClusterInfo(),
	}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(queryMods, where.OrganizationID.EQ(uuid.FromStringOrNil(*f)))
	}

	if f := filter.ClusterGroupID; f != nil {
		queryMods = append(queryMods, where.ClusterGroupID.EQ(uuid.FromStringOrNil(*f)))
	}

	bs, err := bobmodel.Clusters.Query(
		queryMods...,
	).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	ds := make([]*matrix.Cluster, 0, len(bs))
	for _, b := range bs {
		d := r.b2d(b)
		if err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func (r *clusterRepo) Delete(ctx context.Context, id matrix.ClusterID) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	where := bobmodel.SelectWhere.Clusters
	load := bobmodel.SelectThenLoad.Cluster
	b, err := bobmodel.Clusters.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
		load.ClusterInfo(),
		load.ClusterTags(),
		load.DeprecatedApis(),
		load.Objects(),
	).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	if b.R.ClusterInfo != nil {
		if err := b.R.ClusterInfo.Delete(ctx, tx); err != nil {
			return err
		}
	}

	if b.R.ClusterTags != nil {
		if err := b.R.ClusterTags.DeleteAll(ctx, tx); err != nil {
			return err
		}
	}

	if b.R.DeprecatedApis != nil {
		if err := b.R.DeprecatedApis.DeleteAll(ctx, tx); err != nil {
			return err
		}
	}

	if b.R.Objects != nil {
		if err := b.R.Objects.DeleteAll(ctx, tx); err != nil {
			return err
		}
	}

	if err := b.Delete(ctx, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *clusterRepo) Exist(ctx context.Context, id matrix.ClusterID) (bool, error) {
	return bobmodel.ClusterExists(ctx, r.db, uuid.FromStringOrNil(id))
}

// TODO: Create at and Updated at are unclear
func (r *clusterRepo) d2b(d *matrix.Cluster) *bobmodel.ClusterSetter {
	return &bobmodel.ClusterSetter{
		ID:             omit.From(uuid.FromStringOrNil(d.ID)),
		OrganizationID: omit.From(uuid.FromStringOrNil(d.OrganizationID)),
		ClusterGroupID: omit.From(uuid.FromStringOrNil(d.ClusterGroupID)),
	}
}

func (r *clusterRepo) b2d(b *bobmodel.Cluster) *matrix.Cluster {
	cluster := &matrix.Cluster{
		ID:             b.ID.String(),
		OrganizationID: b.OrganizationID.String(),
		ClusterGroupID: b.ClusterGroupID.String(),
	}

	if b.R.ClusterInfo != nil {
		cluster.Info = &matrix.ClusterInfo{
			Name:     b.R.ClusterInfo.Name,
			Version:  b.R.ClusterInfo.Version,
			Platform: b.R.ClusterInfo.Platform,
			Provider: b.R.ClusterInfo.Provider,
		}
	}
	return cluster
}

func (r *clusterRepo) upsertInfo(ctx context.Context, tx bob.Transaction, clusterID matrix.ClusterID, info *matrix.ClusterInfo) error {
	cols := dbinfo.ClusterInfos.Columns
	setter := &bobmodel.ClusterInfoSetter{
		ClusterID: omit.From(uuid.FromStringOrNil(clusterID)),
		Name:      omit.From(info.Name),
		Version:   omit.From(info.Version),
		Platform:  omit.From(info.Platform),
		Provider:  omit.From(info.Provider),
	}
	if _, err := bobmodel.ClusterInfos.Insert(
		setter,
		im.OnConflict(cols.ClusterID.Name).DoUpdate(
			im.SetExcluded(
				cols.Name.Name,
				cols.Version.Name,
				cols.Platform.Name,
				cols.Provider.Name,
			),
		),
	).One(ctx, tx); err != nil {
		return err
	}
	return nil
}
