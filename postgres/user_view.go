package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

type userViewRepo struct {
	db bob.DB
}

func NewUserViewRepo(db *sql.DB) matrix.UserViewRepo {
	return &userViewRepo{
		db: bob.NewDB(db),
	}
}

func (r *userViewRepo) Find(ctx context.Context, id matrix.UserAccountID) (*matrix.UserView, error) {
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

func (r *userViewRepo) FindByEmail(ctx context.Context, email string) (*matrix.UserView, error) {
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

func (r *userViewRepo) List(ctx context.Context, filter matrix.UserViewFilter) (*matrix.UserViewList, error) {
	where := bobmodel.SelectWhere.UserAccounts
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(queryMods, where.OrganizationID.EQ(uuid.FromStringOrNil(*f)))
	}

	if f := filter.Verified; f != nil {
		queryMods = append(queryMods, where.Verified.EQ(*f))
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

	ds := make([]matrix.UserView, 0, len(bs))
	for _, b := range bs {
		d, err := r.b2d(b)
		if err != nil {
			return nil, err
		}
		ds = append(ds, *d)
	}
	return &matrix.UserViewList{
		Users: ds,
	}, nil
}

func (r *userViewRepo) b2d(b *bobmodel.UserAccount) (*matrix.UserView, error) {
	role, ok := matrix.RoleMap[b.RoleName]
	if !ok {
		return nil, errors.Join(ekg.ErrNotFound, errors.New("role not found"))
	}

	return &matrix.UserView{
		ID:             b.ID.String(),
		OrganizationID: b.OrganizationID.String(),
		Role:           role.Name,
		Email:          b.Email,
		Verified:       b.Verified,
		FirstName:      b.FirstName,
		LastName:       b.LastName,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}, nil
}
