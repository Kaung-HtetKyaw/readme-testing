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

type organizationViewRepo struct {
	db bob.DB
}

func NewOrganizationViewRepo(db *sql.DB) matrix.OrganizationViewRepo {
	return &organizationViewRepo{
		db: bob.NewDB(db),
	}
}

func (r *organizationViewRepo) Find(ctx context.Context, id matrix.OrganizationID) (*matrix.OrganizationView, error) {
	b, err := bobmodel.FindOrganization(ctx, r.db, uuid.FromStringOrNil(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b), nil
}

func (r *organizationViewRepo) List(ctx context.Context, filter matrix.OrganizationViewFilter) (*matrix.OrganizationViewList, error) {
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	bs, err := bobmodel.Organizations.Query(queryMods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	ds := make([]matrix.OrganizationView, 0, len(bs))
	for _, b := range bs {
		d := r.b2d(b)
		ds = append(ds, *d)
	}
	return &matrix.OrganizationViewList{
		Organizations: ds,
	}, nil
}

func (r *organizationViewRepo) b2d(b *bobmodel.Organization) *matrix.OrganizationView {
	return &matrix.OrganizationView{
		ID:        b.ID.String(),
		Name:      b.Name,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
