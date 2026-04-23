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

type invitedUserRepo struct {
	db bob.DB
}

func NewInvitedUserRepo(db *sql.DB) matrix.InvitedUserRepo {
	return &invitedUserRepo{
		db: bob.NewDB(db),
	}
}

func (r *invitedUserRepo) Store(ctx context.Context, invite *matrix.InvitedUser) error {
	cols := dbinfo.InvitedUsers.Columns
	if _, err := bobmodel.InvitedUsers.Insert(
		r.d2b(invite),
		im.OnConflict(cols.Email.Name).DoUpdate(
			im.SetExcluded(
				cols.RoleName.Name,
				cols.Expiration.Name,
				cols.UpdatedAt.Name,
			),
		),
	).One(ctx, r.db); err != nil {
		return err
	}
	return nil
}

func (r *invitedUserRepo) Find(ctx context.Context, email string) (*matrix.InvitedUser, error) {
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

func (r *invitedUserRepo) Delete(ctx context.Context, email string) error {
	where := bobmodel.SelectWhere.InvitedUsers
	b, err := bobmodel.InvitedUsers.Query(
		where.Email.EQ(email),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}
	return b.Delete(ctx, r.db)
}

func (r *invitedUserRepo) d2b(d *matrix.InvitedUser) *bobmodel.InvitedUserSetter {
	return &bobmodel.InvitedUserSetter{
		Email:          omit.From(d.Email),
		OrganizationID: omit.From(uuid.FromStringOrNil(d.OrganizationID)),
		RoleName:       omit.From(d.Role.Name),
		Expiration:     omit.From(d.Expiration),
		CreatedAt:      omit.From(d.CreatedAt),
		UpdatedAt:      omit.From(d.UpdatedAt),
	}
}

func (r *invitedUserRepo) b2d(b *bobmodel.InvitedUser) (*matrix.InvitedUser, error) {
	role, ok := matrix.RoleMap[b.RoleName]
	if !ok {
		return nil, errors.Join(ekg.ErrNotFound, errors.New("role not found"))
	}

	return &matrix.InvitedUser{
		Email:          b.Email,
		OrganizationID: b.OrganizationID.String(),
		Role:           role,
		Expiration:     b.Expiration,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}, nil
}
