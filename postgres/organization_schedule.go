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
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

type organizationScheduleRepo struct {
	db bob.DB
}

func NewOrganizationScheduleRepo(db *sql.DB) matrix.OrganizationScheduleRepo {
	return &organizationScheduleRepo{
		db: bob.NewDB(db),
	}
}

// CreateAll
func (r *organizationScheduleRepo) CreateAll(ctx context.Context, schedules ...*matrix.OrganizationSchedule) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	for _, s := range schedules {
		if _, err := bobmodel.OrganizationSchedules.Insert(
			r.d2b(s),
		).All(ctx, tx); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// Store
func (r *organizationScheduleRepo) Store(ctx context.Context, schedule *matrix.OrganizationSchedule) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	cols := dbinfo.OrganizationSchedules.Columns
	if _, err := bobmodel.Clusters.Insert(
		r.d2b(schedule),
		im.OnConflict(cols.ID.Name).DoUpdate(
			im.SetExcluded(
				cols.Processing.Name,
			),
		),
	).One(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// List
// NOTE I don't like the idea of letting List function update the processing flag,
// better come up with the better solution.
func (r *organizationScheduleRepo) List(ctx context.Context, filter matrix.OrganizationScheduleFilter) ([]*matrix.OrganizationSchedule, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	where := bobmodel.SelectWhere.OrganizationSchedules
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		where.Processing.EQ(false),
		sm.ForUpdate().SkipLocked(),
	}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(queryMods, where.OrganizationID.EQ(uuid.FromStringOrNil(*f)))
	}

	if f := filter.Names; len(f) > 0 {
		queryMods = append(queryMods, where.Name.In(f...))
	}

	if f := filter.ScheduledForBefore; f != nil {
		queryMods = append(queryMods, where.ScheduledFor.LTE(*f))
	}

	bs, err := bobmodel.OrganizationSchedules.Query(
		queryMods...,
	).All(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	for _, b := range bs {
		if b.Update(
			ctx,
			tx,
			&bobmodel.OrganizationScheduleSetter{
				ID:         omit.From(b.ID),
				Processing: omit.From(true),
			},
		); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	ds := make([]*matrix.OrganizationSchedule, 0, len(bs))
	for _, b := range bs {
		d := r.b2d(b)
		ds = append(ds, d)
	}
	return ds, nil
}

// Delete
func (r *organizationScheduleRepo) Delete(ctx context.Context, name string, organizationID matrix.OrganizationID) error {
	where := bobmodel.DeleteWhere.OrganizationSchedules
	if _, err := bobmodel.OrganizationSchedules.Delete(
		where.Name.EQ(name),
		where.OrganizationID.EQ(uuid.FromStringOrNil(organizationID)),
	).One(ctx, r.db); err != nil {
		return err
	}
	return nil
}

func (r *organizationScheduleRepo) d2b(d *matrix.OrganizationSchedule) *bobmodel.OrganizationScheduleSetter {
	return &bobmodel.OrganizationScheduleSetter{
		ID:             omit.From(uuid.FromStringOrNil(d.ID)),
		OrganizationID: omit.From(uuid.FromStringOrNil(d.OrganizationID)),
		Name:           omit.From(d.Name),
		ScheduledFor:   omit.From(d.ScheduledFor),
		CreatedAt:      omit.From(d.CreatedAt),
		// Keep payload null for now.
	}
}

func (r *organizationScheduleRepo) b2d(b *bobmodel.OrganizationSchedule) *matrix.OrganizationSchedule {
	return &matrix.OrganizationSchedule{
		ID:             b.ID.String(),
		OrganizationID: b.OrganizationID.String(),
		Name:           b.Name,
		Payload:        b.Payload,
		ScheduledFor:   b.ScheduledFor,
		CreatedAt:      b.CreatedAt,
	}
}
