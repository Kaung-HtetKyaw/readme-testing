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
)

type organizationRepo struct {
	db bob.DB
}

func NewOrganizationRepo(db *sql.DB) matrix.OrganizationRepo {
	return &organizationRepo{
		db: bob.NewDB(db),
	}
}

func (r *organizationRepo) Create(ctx context.Context, org *matrix.Organization) error {
	if _, err := bobmodel.Organizations.Insert(
		r.d2b(org),
	).One(ctx, r.db); err != nil {
		return err
	}
	return nil
}

func (r *organizationRepo) Find(ctx context.Context, id matrix.OrganizationID) (*matrix.Organization, error) {
	b, err := bobmodel.FindOrganization(ctx, r.db, uuid.FromStringOrNil(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b), nil
}

func (r *organizationRepo) Update(ctx context.Context, update matrix.OrganizationUpdate) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	orgID := uuid.FromStringOrNil(update.ID)
	org, err := bobmodel.FindOrganization(ctx, tx, orgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	setter := &bobmodel.OrganizationSetter{}
	if name := update.Payload.Name; name != nil {
		setter.Name = omit.FromPtr(name)
	}

	if err := org.Update(ctx, tx, setter); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *organizationRepo) Delete(ctx context.Context, id matrix.OrganizationID) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	orgID := uuid.FromStringOrNil(id)
	org, err := bobmodel.FindOrganization(ctx, tx, orgID)
	if err != nil {
		return err
	}

	if err := org.Delete(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *organizationRepo) d2b(d *matrix.Organization) *bobmodel.OrganizationSetter {
	return &bobmodel.OrganizationSetter{
		ID:        omit.From(uuid.FromStringOrNil(d.ID)),
		Name:      omit.From(d.Name),
		CreatedAt: omit.From(d.CreatedAt),
		UpdatedAt: omit.From(d.UpdatedAt),
	}
}

func (r *organizationRepo) b2d(b *bobmodel.Organization) *matrix.Organization {
	return &matrix.Organization{
		ID:        b.ID.String(),
		Name:      b.Name,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
