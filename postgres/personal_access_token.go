package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	dberrors "github.com/kubegrade/matrix/postgres/bob_dberrors"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
)

type personalAccessTokenRepo struct {
	db bob.DB
}

func NewPersonalAccessTokenRepo(db *sql.DB) matrix.PersonalAccessTokenRepo {
	return &personalAccessTokenRepo{
		db: bob.NewDB(db),
	}
}

func (r *personalAccessTokenRepo) Create(ctx context.Context, personalAccessToken *matrix.PersonalAccessToken) error {
	if _, err := bobmodel.PersonalAccessTokens.Insert(r.d2b(personalAccessToken)).One(ctx, r.db); err != nil {
		if errors.Is(dberrors.PersonalAccessTokenErrors.ErrUniqueUniqPersonalAccessTokenOrgName, err) {
			return ekg.ErrDuplicatePersonalAccessTokenName
		}
		return err
	}
	return nil
}

func (r *personalAccessTokenRepo) Update(ctx context.Context, update matrix.PersonalAccessTokenUpdate) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	id := uuid.FromStringOrNil(update.ID)
	b, err := bobmodel.PersonalAccessTokens.Query(
		bobmodel.SelectWhere.PersonalAccessTokens.ID.EQ(id),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	setter := &bobmodel.PersonalAccessTokenSetter{
		ID: omit.From(id),
	}

	if name := update.Payload.Name; name != nil {
		setter.Name = omit.From(*name)
	}

	if value := update.Payload.EncryptedValue; value != nil {
		setter.EncryptedValue = omit.From(*value)
	}

	if expiredAt := update.Payload.ExpiredAt; expiredAt != nil {
		setter.ExpiredAt = omitnull.From(*expiredAt)
	}

	if err := b.Update(ctx, tx, setter); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *personalAccessTokenRepo) Find(ctx context.Context, id matrix.PersonalAccessTokenID) (*matrix.PersonalAccessToken, error) {
	where := bobmodel.SelectWhere.PersonalAccessTokens
	b, err := bobmodel.PersonalAccessTokens.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b), nil
}

func (r *personalAccessTokenRepo) Delete(ctx context.Context, id matrix.PersonalAccessTokenID) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	where := bobmodel.SelectWhere.PersonalAccessTokens
	b, err := bobmodel.PersonalAccessTokens.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
		bobmodel.SelectThenLoad.PersonalAccessToken.RepositoryPersonalAccessTokens(),
	).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	if err := b.R.RepositoryPersonalAccessTokens.DeleteAll(ctx, tx); err != nil {
		return err
	}
	if err := b.Delete(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *personalAccessTokenRepo) d2b(d *matrix.PersonalAccessToken) *bobmodel.PersonalAccessTokenSetter {
	owner := "unknown"
	if d.Owner != "" {
		owner = d.Owner
	}
	return &bobmodel.PersonalAccessTokenSetter{
		ID:             omit.From(uuid.FromStringOrNil(d.ID)),
		OrganizationID: omit.From(uuid.FromStringOrNil(d.OrganizationID)),
		Provider:       omit.From(d.Provider),
		Owner:          omit.From(owner),
		Name:           omit.From(d.Name),
		EncryptedValue: omit.From(d.EncryptedValue),
		ExpiredAt:      omitnull.FromPtr(d.ExpiredAt),
		CreatedBy:      omitnull.From(uuid.FromStringOrNil(d.CreatedBy)),
		UpdatedBy:      omitnull.From(uuid.FromStringOrNil(d.UpdatedBy)),
		CreatedAt:      omit.From(d.CreatedAt),
		UpdatedAt:      omit.From(d.UpdatedAt),
	}
}

func (r *personalAccessTokenRepo) b2d(b *bobmodel.PersonalAccessToken) *matrix.PersonalAccessToken {
	return &matrix.PersonalAccessToken{
		ID:             b.ID.String(),
		OrganizationID: b.OrganizationID.String(),
		Owner:          b.Owner,
		Provider:       b.Provider,
		Name:           b.Name,
		EncryptedValue: b.EncryptedValue,
		ExpiredAt:      b.ExpiredAt.Ptr(),
		CreatedBy:      b.CreatedBy.GetOrZero().String(),
		UpdatedBy:      b.UpdatedBy.GetOrZero().String(),
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}
