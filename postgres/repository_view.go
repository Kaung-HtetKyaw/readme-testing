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

type repositoryViewRepo struct {
	db bob.DB
}

func NewRepositoryViewRepo(db *sql.DB) matrix.RepositoryViewRepo {
	return &repositoryViewRepo{db: bob.NewDB(db)}
}

func (r *repositoryViewRepo) List(ctx context.Context, filter matrix.RepositoryViewFilter) (*matrix.RepositoryViewList, error) {
	where := bobmodel.SelectWhere.Repositories
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		bobmodel.SelectThenLoad.Repository.RepositoryPersonalAccessTokens(),
	}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(queryMods, where.OrganizationID.EQ(uuid.FromStringOrNil(*f)))
	}

	bs, err := bobmodel.Repositories.Query(queryMods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	var ds = make([]matrix.RepositoryView, 0, len(bs))
	for _, b := range bs {
		d, err := r.b2d(b)
		if err != nil {
			return nil, err
		}
		ds = append(ds, *d)
	}
	return &matrix.RepositoryViewList{
		Items: ds,
	}, nil
}

// TODO encrypt value
func (r *repositoryViewRepo) b2d(b *bobmodel.Repository) (*matrix.RepositoryView, error) {
	patIDs := make([]matrix.PersonalAccessTokenID, 0, len(b.R.RepositoryPersonalAccessTokens))
	for _, b := range b.R.RepositoryPersonalAccessTokens {
		patIDs = append(patIDs, b.PersonalAccessTokenID.String())
	}

	return &matrix.RepositoryView{
		ID:                     b.ID.String(),
		OrganizationID:         b.OrganizationID.String(),
		Provider:               b.Provider,
		Namespace:    b.Namespace,
		Name:                   b.Name,
		Description:            b.Description.GetOrZero(),
		PersonalAccessTokenIDs: patIDs,
		CreatedAt:              b.CreatedAt,
		UpdatedAt:              b.UpdatedAt,
	}, nil
}
