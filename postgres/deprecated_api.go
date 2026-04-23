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
	"github.com/stephenafamo/bob/dialect/psql/im"
)

type deprecatedAPIRepo struct {
	db bob.DB
}

func NewDeprecatedAPIRepo(db *sql.DB) matrix.DeprecatedAPIRepo {
	return &deprecatedAPIRepo{db: bob.NewDB(db)}
}

func (r *deprecatedAPIRepo) StoreAll(ctx context.Context, deprecatedAPIs ...*matrix.DeprecatedAPI) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	for _, api := range deprecatedAPIs {
		cols := dbinfo.DeprecatedApis.Columns
		if _, err := bobmodel.DeprecatedApis.Insert(
			r.d2b(api),
			im.OnConflict(
				cols.ClusterID.Name,
				cols.CurrentGroupVersion.Name,
				cols.Kind.Name,
			).DoUpdate(
				im.SetExcluded(
					cols.Name.Name,
					cols.ClusterK8SVersion.Name,
					cols.Deprecated.Name,
					cols.DeprecatedIn.Name,
					cols.RemovedIn.Name,
					cols.ReplacementVersion.Name,
				),
			),
		).One(ctx, tx); err != nil {
			return err
		}
	}

	if err := refreshClusterResourceCount(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *deprecatedAPIRepo) FindByClusterID(ctx context.Context, clusterID matrix.ClusterID) (*matrix.DeprecatedAPI, error) {
	where := bobmodel.SelectWhere.DeprecatedApis
	b, err := bobmodel.DeprecatedApis.Query(
		where.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b), nil
}

func (r *deprecatedAPIRepo) DeleteByClusterID(ctx context.Context, clusterID matrix.ClusterID) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	where := bobmodel.DeleteWhere.DeprecatedApis
	if _, err := bobmodel.DeprecatedApis.Delete(
		where.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
	).All(ctx, r.db); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	if err := refreshClusterResourceCount(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *deprecatedAPIRepo) d2b(d *matrix.DeprecatedAPI) *bobmodel.DeprecatedAPISetter {
	return &bobmodel.DeprecatedAPISetter{
		ClusterID:           omit.From(uuid.FromStringOrNil(d.ClusterID)),
		OrganizationID:      omit.From(uuid.FromStringOrNil(d.OrganizationID)),
		CurrentGroupVersion: omit.From(d.CurrentGroupVersion),
		Kind:                omit.From(d.Kind),
		Name:                omit.From(d.Name),
		ClusterK8SVersion:   omit.From(d.ClusterK8sVersion),
		Deprecated:          omit.From(d.Deprecated),
		DeprecatedIn:        omit.From(d.DeprecatedIn),
		RemovedIn:           omit.From(d.RemovedIn),
		ReplacementVersion:  omit.From(d.ReplacementVersion),
	}
}

func (r *deprecatedAPIRepo) b2d(b *bobmodel.DeprecatedAPI) *matrix.DeprecatedAPI {
	return &matrix.DeprecatedAPI{
		ClusterID:           b.ClusterID.String(),
		OrganizationID:      b.OrganizationID.String(),
		CurrentGroupVersion: b.CurrentGroupVersion,
		Kind:                b.Kind,
		Name:                b.Name,
		ClusterK8sVersion:   b.ClusterK8SVersion,
		Deprecated:          b.Deprecated,
		DeprecatedIn:        b.DeprecatedIn,
		RemovedIn:           b.RemovedIn,
		ReplacementVersion:  b.ReplacementVersion,
	}
}
