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

type InvitedUserViewRepo struct {
	db bob.DB
}

func NewInvitedUserViewRepo(db *sql.DB) matrix.InvitedUserViewRepo {
	return &InvitedUserViewRepo{
		db: bob.NewDB(db),
	}
}

func (r *InvitedUserViewRepo) Find(ctx context.Context, email string) (*matrix.InvitedUserView, error) {
	where := bobmodel.SelectWhere.InvitedUsers
	b, err := bobmodel.InvitedUsers.Query(
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

func (r *InvitedUserViewRepo) List(ctx context.Context, filter matrix.InvitedUserViewFilter) (*matrix.InvitedUserViewList, error) {
	where := bobmodel.SelectWhere.InvitedUsers
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(queryMods, where.OrganizationID.EQ(uuid.FromStringOrNil(*f)))
	}

	bs, err := bobmodel.InvitedUsers.Query(queryMods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	var ds = make([]matrix.InvitedUserView, 0, len(bs))
	for _, b := range bs {
		d, err := r.b2d(b)
		if err != nil {
			return nil, err
		}
		ds = append(ds, *d)
	}
	return &matrix.InvitedUserViewList{
		InvitedUsers: ds,
	}, nil
}

func (r *InvitedUserViewRepo) b2d(b *bobmodel.InvitedUser) (*matrix.InvitedUserView, error) {
	role, ok := matrix.RoleMap[b.RoleName]
	if !ok {
		return nil, errors.Join(ekg.ErrNotFound, errors.New("role not found"))
	}

	return &matrix.InvitedUserView{
		Email:          b.Email,
		OrganizationID: b.OrganizationID.String(),
		Role:           role.Name,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}, nil
}
