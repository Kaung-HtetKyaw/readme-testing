package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/opt/omit"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
)

type passwordAuthRepo struct {
	db bob.DB
}

func NewPasswordAuthRepo(db *sql.DB) matrix.PasswordAuthRepo {
	return &passwordAuthRepo{
		db: bob.NewDB(db),
	}
}

// Store
func (r *passwordAuthRepo) Store(ctx context.Context, auth *matrix.PasswordAuth) error {
	if _, err := bobmodel.PasswordAuths.Insert(
		r.d2b(auth),
	).One(ctx, r.db); err != nil {
		return err
	}
	return nil
}

// Find
func (r *passwordAuthRepo) Find(ctx context.Context, email string) (*matrix.PasswordAuth, error) {
	where := bobmodel.SelectWhere.PasswordAuths
	b, err := bobmodel.PasswordAuths.Query(
		where.Email.EQ(email),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b)
}

// Exists
func (r *passwordAuthRepo) Exists(ctx context.Context, email string) (bool, error) {
	where := bobmodel.SelectWhere.PasswordAuths
	return bobmodel.PasswordAuths.Query(where.Email.EQ(email)).Exists(ctx, r.db)
}

func (r *passwordAuthRepo) d2b(d *matrix.PasswordAuth) *bobmodel.PasswordAuthSetter {
	return &bobmodel.PasswordAuthSetter{
		Email:        omit.From(d.Email),
		PasswordHash: omit.From(d.PasswordHash),
		PasswordSalt: omit.From(d.PasswordSalt),
		CreatedAt:    omit.From(d.CreatedAt),
		UpdatedAt:    omit.From(d.UpdatedAt),
	}
}

func (r *passwordAuthRepo) b2d(b *bobmodel.PasswordAuth) (*matrix.PasswordAuth, error) {
	return &matrix.PasswordAuth{
		Email:        b.Email,
		PasswordHash: b.PasswordHash,
		PasswordSalt: b.PasswordSalt,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}, nil
}
