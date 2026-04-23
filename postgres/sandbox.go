package postgres

import (
	"context"
	"database/sql"

	"github.com/aarondl/opt/omit"
	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/matrix"
	dbinfo "github.com/kubegrade/matrix/postgres/bob_dbinfo"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

type sandboxRepo struct {
	db bob.DB
}

func NewSandboxRepo(db *sql.DB) matrix.SandboxRepo {
	return &sandboxRepo{
		db: bob.NewDB(db),
	}
}

// SyncClusterToOrganization syncs all sandbox cluster data to
// the specific organization.
func (r *sandboxRepo) SyncClusterToOrganization(ctx context.Context, organizationID matrix.OrganizationID) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	bClusterGroup, err := bobmodel.ClusterGroups.Query(
		bobmodel.SelectWhere.ClusterGroups.OrganizationID.EQ(uuid.FromStringOrNil(organizationID)),
		bobmodel.SelectWhere.ClusterGroups.Name.EQ(matrix.ClusterGroupDefaultName),
	).One(ctx, tx)
	if err != nil {
		return err
	}

	sandboxes, err := bobmodel.Sandboxes.Query().All(ctx, tx)
	if err != nil {
		return err
	}

	for _, s := range sandboxes {
		clusterID, err := uuid.NewV4()
		if err != nil {
			return err
		}
		if err := r.insertCluster(
			ctx,
			tx,
			clusterID,
			organizationID,
			bClusterGroup.ID,
		); err != nil {
			return err
		}

		if err := r.syncClusterInfo(ctx, tx, s.ID, clusterID); err != nil {
			return err
		}

		if err := r.syncObject(ctx, tx, s.ID, clusterID); err != nil {
			return err
		}

		if err := r.syncDeprecatedAPIGroup(ctx, tx, s.ID, clusterID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *sandboxRepo) insertCluster(
	ctx context.Context,
	tx bob.Transaction,
	clusterID uuid.UUID,
	organizationID string,
	clusterGroupID uuid.UUID,
) error {
	clusterSetter := &bobmodel.ClusterSetter{
		ID:             omit.From(clusterID),
		OrganizationID: omit.From(uuid.FromStringOrNil(organizationID)),
		ClusterGroupID: omit.From(clusterGroupID),
	}
	if _, err := bobmodel.Clusters.Insert(clusterSetter).One(ctx, tx); err != nil {
		return err
	}
	return nil
}

func (r *sandboxRepo) syncClusterInfo(
	ctx context.Context,
	tx bob.Transaction,
	sandboxID uuid.UUID,
	clusterID uuid.UUID,
) error {
	cols := dbinfo.ClusterInfos.Columns
	sandboxCols := dbinfo.SandboxClusterInfos.Columns
	if _, err := psql.Insert(
		im.Into(
			"cluster_info",
			cols.ClusterID.Name,
			cols.Name.Name,
			cols.Version.Name,
			cols.Platform.Name,
		),
		im.Query(
			psql.Select(
				sm.Columns(
					psql.Arg(clusterID),
					sandboxCols.Name.Name,
					sandboxCols.Version.Name,
					sandboxCols.Platform.Name,
				),
				sm.From("sandbox_cluster_info"),
				sm.Where(
					psql.Quote(sandboxCols.SandboxID.Name).EQ(psql.Arg(sandboxID)),
				),
			),
		),
	).Exec(ctx, tx); err != nil {
		return err
	}
	return nil
}

func (r *sandboxRepo) syncObject(
	ctx context.Context,
	tx bob.Transaction,
	sandboxID uuid.UUID,
	clusterID uuid.UUID,
) error {
	cols := dbinfo.Objects.Columns
	sandboxCols := dbinfo.SandboxObjects.Columns
	if _, err := psql.Insert(
		im.Into(
			"object",
			cols.ClusterID.Name,
			cols.ID.Name,
			cols.Name.Name,
			cols.Namespace.Name,
			cols.ResourceVersion.Name,
			cols.Kind.Name,
			cols.Raw.Name,
		),
		im.Query(
			psql.Select(
				sm.Columns(
					psql.Arg(clusterID),
					sandboxCols.ID.Name,
					sandboxCols.Name.Name,
					sandboxCols.Namespace.Name,
					sandboxCols.ResourceVersion.Name,
					sandboxCols.Kind.Name,
					sandboxCols.Raw.Name,
				),
				sm.From("sandbox_object"),
				sm.Where(
					psql.Quote(sandboxCols.SandboxID.Name).EQ(psql.Arg(sandboxID)),
				),
			),
		),
	).Exec(ctx, tx); err != nil {
		return err
	}
	return nil
}

func (r *sandboxRepo) syncDeprecatedAPIGroup(
	ctx context.Context,
	tx bob.Transaction,
	sandboxID uuid.UUID,
	clusterID uuid.UUID,
) error {
	cols := dbinfo.DeprecatedApis.Columns
	sandboxCols := dbinfo.SandboxDeprecatedAPIGroups.Columns
	if _, err := psql.Insert(
		im.Into(
			"deprecated_api_group",
			cols.ClusterID.Name,
		),
		im.Query(
			psql.Select(
				sm.Columns(
					psql.Arg(clusterID),
					sandboxCols.DeprecatedApis.Name,
				),
				sm.From("sandbox_deprecated_api_group"),
				sm.Where(
					psql.Quote(sandboxCols.SandboxID.Name).EQ(psql.Arg(sandboxID)),
				),
			),
		),
	).Exec(ctx, tx); err != nil {
		return err
	}
	return nil
}
