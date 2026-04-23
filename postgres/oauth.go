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

type oauthRepo struct {
	db bob.DB
}

func NewOAuthRepo(db *sql.DB) matrix.OAuthRepo {
	return &oauthRepo{
		db: bob.NewDB(db),
	}
}

// Store
func (r *oauthRepo) Store(ctx context.Context, login *matrix.OAuth) error {
	if _, err := bobmodel.OAuths.Insert(
		r.d2b(login),
	).One(ctx, r.db); err != nil {
		return err
	}
	return nil
}

// Find
func (r *oauthRepo) FindByProviderAndEmail(ctx context.Context, provider string, email string) (*matrix.OAuth, error) {
	where := bobmodel.SelectWhere.OAuths
	b, err := bobmodel.OAuths.Query(
		where.Provider.EQ(provider),
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

func (r *oauthRepo) d2b(d *matrix.OAuth) *bobmodel.OAuthSetter {
	return &bobmodel.OAuthSetter{
		Email:     omit.From(d.Email),
		Provider:  omit.From(d.Provider),
		CreatedAt: omit.From(d.CreatedAt),
		UpdatedAt: omit.From(d.UpdatedAt),
	}
}

func (r *oauthRepo) b2d(b *bobmodel.OAuth) (*matrix.OAuth, error) {
	return &matrix.OAuth{
		Email:     b.Email,
		Provider:  b.Provider,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}, nil
}
