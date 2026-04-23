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

type verifyTokenRepo struct {
	db bob.DB
}

func NewEmailVerificationTokenRepo(db *sql.DB) matrix.EmailVerificationTokenRepo {
	return &verifyTokenRepo{
		db: bob.NewDB(db),
	}
}

func (r *verifyTokenRepo) Store(ctx context.Context, token *matrix.EmailVerificationToken) error {
	cols := dbinfo.EmailVerificationTokens.Columns
	if _, err := bobmodel.EmailVerificationTokens.Insert(
		r.d2b(token),
		im.OnConflict(cols.Email.Name).DoUpdate(
			im.SetExcluded(
				cols.Value.Name,
				cols.Expiration.Name,
				cols.UpdatedAt.Name,
			),
		),
	).One(ctx, r.db); err != nil {
		return err
	}
	return nil
}

func (r *verifyTokenRepo) Find(ctx context.Context, value string) (*matrix.EmailVerificationToken, error) {
	where := bobmodel.SelectWhere.EmailVerificationTokens
	b, err := bobmodel.EmailVerificationTokens.Query(
		where.Value.EQ(uuid.FromStringOrNil(value)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b), nil
}

func (r *verifyTokenRepo) Delete(ctx context.Context, value string) error {
	where := bobmodel.SelectWhere.EmailVerificationTokens
	b, err := bobmodel.EmailVerificationTokens.Query(
		where.Value.EQ(uuid.FromStringOrNil(value)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}
	return b.Delete(ctx, r.db)
}

func (r *verifyTokenRepo) d2b(d *matrix.EmailVerificationToken) *bobmodel.EmailVerificationTokenSetter {
	return &bobmodel.EmailVerificationTokenSetter{
		Value:      omit.From(uuid.FromStringOrNil(d.Value)),
		Email:      omit.From(d.Email),
		Expiration: omit.From(d.Expiration),
	}
}

func (r *verifyTokenRepo) b2d(b *bobmodel.EmailVerificationToken) *matrix.EmailVerificationToken {
	return &matrix.EmailVerificationToken{
		Value:      b.Value.String(),
		Email:      b.Email,
		Expiration: b.Expiration,
	}
}
