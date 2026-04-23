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
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

type userAccountRepo struct {
	db bob.DB
}

func NewUserAccountRepo(db *sql.DB) matrix.UserAccountRepo {
	return &userAccountRepo{
		db: bob.NewDB(db),
	}
}

func (r *userAccountRepo) Create(ctx context.Context, user *matrix.UserAccount) error {
	if _, err := bobmodel.UserAccounts.Insert(
		r.d2b(user),
	).One(ctx, r.db); err != nil {
		return err
	}
	return nil
}

func (r *userAccountRepo) Update(ctx context.Context, update matrix.UserAccountUpdate) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	userID := uuid.FromStringOrNil(update.ID)
	user, err := bobmodel.FindUserAccount(ctx, tx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	setter := &bobmodel.UserAccountSetter{}
	if organizationID := update.Payload.OrganizationID; organizationID != nil {
		setter.OrganizationID = omit.From(uuid.FromStringOrNil(*organizationID))
	}

	if firstName := update.Payload.FirstName; firstName != nil {
		setter.FirstName = omit.FromPtr(firstName)
	}

	if lastName := update.Payload.LastName; lastName != nil {
		setter.LastName = omit.FromPtr(lastName)
	}

	if verified := update.Payload.Verified; verified != nil {
		setter.Verified = omit.FromPtr(verified)
	}

	if role := update.Payload.Role; role != nil {
		setter.RoleName = omit.From(role.Name)
	}

	if err := user.Update(ctx, tx, setter); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *userAccountRepo) Find(ctx context.Context, id matrix.UserAccountID) (*matrix.UserAccount, error) {
	where := bobmodel.SelectWhere.UserAccounts
	b, err := bobmodel.UserAccounts.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b)
}

func (r *userAccountRepo) FindByEmail(ctx context.Context, email string) (*matrix.UserAccount, error) {
	where := bobmodel.SelectWhere.UserAccounts
	b, err := bobmodel.UserAccounts.Query(
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

func (r *userAccountRepo) List(ctx context.Context, filter matrix.UserAccountFilter) ([]*matrix.UserAccount, error) {
	where := bobmodel.SelectWhere.UserAccounts
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(queryMods, where.OrganizationID.EQ(uuid.FromStringOrNil(*f)))
	}
	if f := filter.RoleName; len(f) > 0 {
		queryMods = append(queryMods, where.RoleName.In(f...))
	}

	bs, err := bobmodel.UserAccounts.Query(
		queryMods...,
	).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	ds := make([]*matrix.UserAccount, 0, len(bs))
	for _, b := range bs {
		d, err := r.b2d(b)
		if err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func (r *userAccountRepo) Delete(ctx context.Context, id matrix.UserAccountID) error {
	where := bobmodel.SelectWhere.UserAccounts
	b, err := bobmodel.UserAccounts.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
		bobmodel.Preload.UserAccount.UserForgotPasswordToken(),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	whereEmailVerificationToken := bobmodel.DeleteWhere.EmailVerificationTokens
	if _, err := bobmodel.EmailVerificationTokens.Delete(
		whereEmailVerificationToken.Email.EQ(b.Email),
	).One(ctx, r.db); err != nil {
		return err
	}

	if b.R.UserForgotPasswordToken != nil {
		if err := b.R.UserForgotPasswordToken.Delete(ctx, r.db); err != nil {
			return err
		}
	}
	return b.Delete(ctx, r.db)
}

func (r *userAccountRepo) DeleteByEmail(ctx context.Context, email string) error {
	where := bobmodel.SelectWhere.UserAccounts
	b, err := bobmodel.UserAccounts.Query(
		where.Email.EQ(email),
		bobmodel.Preload.UserAccount.UserForgotPasswordToken(),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	whereEmailVerificationToken := bobmodel.DeleteWhere.EmailVerificationTokens
	if _, err := bobmodel.EmailVerificationTokens.Delete(
		whereEmailVerificationToken.Email.EQ(b.Email),
	).One(ctx, r.db); err != nil {
		return err
	}

	if b.R.UserForgotPasswordToken != nil {
		if err := b.R.UserForgotPasswordToken.Delete(ctx, r.db); err != nil {
			return err
		}
	}
	return b.Delete(ctx, r.db)
}

// ExistsByEmail
func (r *userAccountRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	where := bobmodel.SelectWhere.UserAccounts
	return bobmodel.UserAccounts.Query(where.Email.EQ(email)).Exists(ctx, r.db)
}

func (r *userAccountRepo) d2b(d *matrix.UserAccount) *bobmodel.UserAccountSetter {
	return &bobmodel.UserAccountSetter{
		ID:             omit.From(uuid.FromStringOrNil(d.ID)),
		OrganizationID: omit.From(uuid.FromStringOrNil(d.OrganizationID)),
		RoleName:       omit.From(d.Role.Name),
		Email:          omit.From(d.Email),
		Verified:       omit.From(d.Verified),
		FirstName:      omit.From(d.FirstName),
		LastName:       omit.From(d.LastName),
		CreatedAt:      omit.From(d.CreatedAt),
		UpdatedAt:      omit.From(d.UpdatedAt),
	}
}

func (r *userAccountRepo) b2d(b *bobmodel.UserAccount) (*matrix.UserAccount, error) {
	role, ok := matrix.RoleMap[b.RoleName]
	if !ok {
		return nil, errors.Join(ekg.ErrNotFound, errors.New("role not found"))
	}

	return &matrix.UserAccount{
		ID:             b.ID.String(),
		OrganizationID: b.OrganizationID.String(),
		Role:           role,
		Email:          b.Email,
		Verified:       b.Verified,
		FirstName:      b.FirstName,
		LastName:       b.LastName,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}, nil
}
