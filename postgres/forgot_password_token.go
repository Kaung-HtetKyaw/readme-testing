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

type forgotPasswordRepo struct {
	db bob.DB
}

func NewForgotPasswordTokenRepo(db *sql.DB) matrix.ForgotPasswordTokenRepo {
	return &forgotPasswordRepo{
		db: bob.NewDB(db),
	}
}

func (r *forgotPasswordRepo) Store(ctx context.Context, token *matrix.ForgotPasswordToken) error {
	cols := dbinfo.ForgotPasswordTokens.Columns
	if _, err := bobmodel.ForgotPasswordTokens.Insert(
		r.d2b(token),
		im.OnConflict(cols.UserID.Name).DoUpdate(
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

func (r *forgotPasswordRepo) Find(ctx context.Context, value string) (*matrix.ForgotPasswordToken, error) {
	where := bobmodel.SelectWhere.ForgotPasswordTokens
	b, err := bobmodel.ForgotPasswordTokens.Query(
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

func (r *forgotPasswordRepo) Delete(ctx context.Context, value string) error {
	where := bobmodel.SelectWhere.ForgotPasswordTokens
	b, err := bobmodel.ForgotPasswordTokens.Query(
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

func (r *forgotPasswordRepo) d2b(d *matrix.ForgotPasswordToken) *bobmodel.ForgotPasswordTokenSetter {
	return &bobmodel.ForgotPasswordTokenSetter{
		Value:      omit.From(uuid.FromStringOrNil(d.Value)),
		UserID:     omit.From(uuid.FromStringOrNil(d.UserAccountID)),
		Expiration: omit.From(d.Expiration),
	}
}

func (r *forgotPasswordRepo) b2d(b *bobmodel.ForgotPasswordToken) *matrix.ForgotPasswordToken {
	return &matrix.ForgotPasswordToken{
		Value:         b.Value.String(),
		UserAccountID: b.UserID.String(),
		Expiration:    b.Expiration,
	}
}
