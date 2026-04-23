package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/matrix"
	dbinfo "github.com/kubegrade/matrix/postgres/bob_dbinfo"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/im"
)

type upgradeAssessmentRepo struct {
	db bob.DB
}

func NewUpgradeAssessmentRepo(db *sql.DB) matrix.UpgradeAssessmentRepo {
	return &upgradeAssessmentRepo{
		db: bob.NewDB(db),
	}
}

func (r *upgradeAssessmentRepo) Store(ctx context.Context, upgradeAssessment *matrix.UpgradeAssessment) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	if upgradeAssessment.Matched {
		cols := dbinfo.UpgradeableComponents.Columns
		if _, err := bobmodel.UpgradeableComponents.Insert(
			r.toUpgradeableComponent(upgradeAssessment),
			im.OnConflict(
				cols.ClusterID.Name,
				cols.ObjectID.Name,
				cols.Name.Name,
			).DoUpdate(
				im.SetExcluded(
					cols.Version.Name,
					cols.NextCompatible.Name,
					cols.MinCompatibleVersion.Name,
					cols.MaxCompatibleVersion.Name,
					cols.UpdatedAt.Name,
				),
			),
		).One(ctx, tx); err != nil {
			return err
		}
	} else {
		cols := dbinfo.UnmatchedComponents.Columns
		if _, err := bobmodel.UnmatchedComponents.Insert(
			r.toUnmatchedComponent(upgradeAssessment),
			im.OnConflict(
				cols.ClusterID.Name,
				cols.ObjectID.Name,
				cols.Name.Name,
			).DoUpdate(
				im.SetExcluded(
					cols.Version.Name,
					cols.UpdatedAt.Name,
				),
			),
		).One(ctx, tx); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *upgradeAssessmentRepo) Find(ctx context.Context, clusterID matrix.ClusterID, objectID matrix.ObjectID, name string) (*matrix.UpgradeAssessment, error) {
	whereUpgradeableComponent := bobmodel.SelectWhere.UpgradeableComponents
	bUpgradeableComponent, err := bobmodel.UpgradeableComponents.Query(
		whereUpgradeableComponent.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		whereUpgradeableComponent.ObjectID.EQ(uuid.FromStringOrNil(objectID)),
		whereUpgradeableComponent.Name.EQ(name),
	).One(ctx, r.db)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if bUpgradeableComponent != nil {
		return r.b2d(bUpgradeableComponent, nil)
	} else {
		whereUnmatchedComponent := bobmodel.SelectWhere.UnmatchedComponents
		bUnmatchedComponent, err := bobmodel.UnmatchedComponents.Query(
			whereUnmatchedComponent.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
			whereUnmatchedComponent.ObjectID.EQ(uuid.FromStringOrNil(objectID)),
			whereUnmatchedComponent.Name.EQ(name),
		).One(ctx, r.db)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
		return r.b2d(nil, bUnmatchedComponent)
	}
}

func (r *upgradeAssessmentRepo) Delete(ctx context.Context, clusterID matrix.ClusterID, objectID matrix.ObjectID, name string) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	whereUpgradeableComponent := bobmodel.DeleteWhere.UpgradeableComponents
	if _, err := bobmodel.UpgradeableComponents.Delete(
		whereUpgradeableComponent.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		whereUpgradeableComponent.ObjectID.EQ(uuid.FromStringOrNil(objectID)),
		whereUpgradeableComponent.Name.EQ(name),
	).One(ctx, r.db); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	whereUnmatchedComponent := bobmodel.DeleteWhere.UnmatchedComponents
	if _, err := bobmodel.UnmatchedComponents.Delete(
		whereUnmatchedComponent.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		whereUnmatchedComponent.ObjectID.EQ(uuid.FromStringOrNil(objectID)),
		whereUnmatchedComponent.Name.EQ(name),
	).One(ctx, r.db); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	return tx.Commit(ctx)
}

// PurgeStale
func (r *upgradeAssessmentRepo) PurgeStale(ctx context.Context, clusterID matrix.ClusterID, currentRunVersion string) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	whereUpgradeableComponent := bobmodel.DeleteWhere.UpgradeableComponents
	if _, err := bobmodel.UpgradeableComponents.Delete(
		whereUpgradeableComponent.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		whereUpgradeableComponent.RunVersion.NE(uuid.FromStringOrNil(currentRunVersion)),
	).One(ctx, r.db); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	whereUnmatchedComponent := bobmodel.DeleteWhere.UnmatchedComponents
	if _, err := bobmodel.UnmatchedComponents.Delete(
		whereUnmatchedComponent.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		whereUnmatchedComponent.RunVersion.NE(uuid.FromStringOrNil(currentRunVersion)),
	).One(ctx, r.db); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *upgradeAssessmentRepo) toUpgradeableComponent(d *matrix.UpgradeAssessment) *bobmodel.UpgradeableComponentSetter {
	var minSupportedVersion string
	var maxSupportedVersion string

	if d.MinSupportedVersion != nil {
		minSupportedVersion = *d.MinSupportedVersion
	}
	if d.MaxSupportedVersion != nil {
		maxSupportedVersion = *d.MaxSupportedVersion
	}

	now := time.Now()
	return &bobmodel.UpgradeableComponentSetter{
		ClusterID:            omit.From(uuid.FromStringOrNil(d.ClusterID)),
		ObjectID:             omit.From(uuid.FromStringOrNil(d.ObjectID)),
		Name:                 omit.From(d.Name),
		Version:              omit.From(d.Version),
		NextCompatible:       omit.From(d.NextCompatible),
		MinCompatibleVersion: omit.From(minSupportedVersion),
		MaxCompatibleVersion: omit.From(maxSupportedVersion),
		CreatedAt:            omit.From(now),
		UpdatedAt:            omit.From(now),
	}
}

func (r *upgradeAssessmentRepo) toUnmatchedComponent(d *matrix.UpgradeAssessment) *bobmodel.UnmatchedComponentSetter {
	now := time.Now()
	return &bobmodel.UnmatchedComponentSetter{
		ClusterID: omit.From(uuid.FromStringOrNil(d.ClusterID)),
		ObjectID:  omit.From(uuid.FromStringOrNil(d.ObjectID)),
		Name:      omit.From(d.Name),
		Version:   omit.From(d.Version),
		CreatedAt: omit.From(now),
		UpdatedAt: omit.From(now),
	}
}

func (r *upgradeAssessmentRepo) b2d(
	bUpgradeableComponent *bobmodel.UpgradeableComponent,
	bUnmatchedComponent *bobmodel.UnmatchedComponent,
) (*matrix.UpgradeAssessment, error) {
	if bUpgradeableComponent != nil {
		return &matrix.UpgradeAssessment{
			ClusterID:           bUpgradeableComponent.ClusterID.String(),
			ObjectID:            bUpgradeableComponent.ObjectID.String(),
			Name:                bUpgradeableComponent.Name,
			Version:             bUpgradeableComponent.Version,
			Matched:             true,
			NextCompatible:      bUpgradeableComponent.NextCompatible,
			MinSupportedVersion: &bUpgradeableComponent.MinCompatibleVersion,
			MaxSupportedVersion: &bUpgradeableComponent.MaxCompatibleVersion,
		}, nil
	} else if bUnmatchedComponent != nil {
		return &matrix.UpgradeAssessment{
			ClusterID: bUnmatchedComponent.ClusterID.String(),
			ObjectID:  bUnmatchedComponent.ObjectID.String(),
			Name:      bUnmatchedComponent.Name,
			Version:   bUnmatchedComponent.Version,
			Matched:   false,
		}, nil
	}
	return nil, errors.New("invalid inputs for upgrade assessment domain model building")
}
