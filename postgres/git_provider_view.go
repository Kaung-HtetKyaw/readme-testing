package postgres

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

type gitProviderViewRepo struct {
	db bob.DB
}

func NewGitProviderViewRepo(db *sql.DB) matrix.GitProviderViewRepo {
	return &gitProviderViewRepo{db: bob.NewDB(db)}
}

func (r *gitProviderViewRepo) List(ctx context.Context, filter matrix.GitProviderViewFilter) (*matrix.GitProviderViewList, error) {
	patMap, err := r.personalAccessTokens(ctx, filter)
	if err != nil {
		return nil, err
	}

	repoMap, err := r.repositories(ctx, filter)
	if err != nil {
		return nil, err
	}

	items := make([]matrix.GitProviderView, 0, len(matrix.GitProviderNames))
	for _, n := range matrix.GitProviderNames {
		orgIDs := make([]matrix.OrganizationID, 0, len(patMap[n])+len(repoMap[n]))
		for id := range patMap[n] {
			orgIDs = append(orgIDs, id)
		}
		for id := range repoMap[n] {
			if !slices.Contains(orgIDs, id) {
				orgIDs = append(orgIDs, id)
			}
		}

		for _, orgID := range orgIDs {
			provider := &matrix.GitProviderView{
				Name:                 n,
				OrganizationID:       orgID,
				PersonalAccessTokens: patMap[n][orgID],
				Repositories:         repoMap[n][orgID],
			}
			items = append(items, *provider)
		}
	}

	return &matrix.GitProviderViewList{
		Items: items,
	}, nil
}

type personalAccessTokenMap = map[matrix.GitProviderName]map[matrix.OrganizationID][]matrix.PersonalAccessTokenView

func (r *gitProviderViewRepo) personalAccessTokens(
	ctx context.Context,
	filter matrix.GitProviderViewFilter,
) (personalAccessTokenMap, error) {
	where := bobmodel.SelectWhere.PersonalAccessTokens
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		where.Provider.In(matrix.GitProviderNames...),
		bobmodel.SelectThenLoad.Repository.RepositoryPersonalAccessTokens(),
	}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(
			queryMods,
			where.OrganizationID.EQ(uuid.FromStringOrNil(*f)),
		)
	}

	bs, err := bobmodel.PersonalAccessTokens.Query(
		queryMods...,
	).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	providerMap := make(personalAccessTokenMap)
	for _, b := range bs {
		orgMap, ok := providerMap[b.Provider]
		if !ok {
			orgMap = make(map[matrix.OrganizationID][]matrix.PersonalAccessTokenView)
		}

		orgMap[b.OrganizationID.String()] = append(
			orgMap[b.OrganizationID.String()],
			matrix.PersonalAccessTokenView{
				ID:        b.ID.String(),
				Name:      b.Name,
				Owner:     b.Owner,
				Value:     redactValue(b.EncryptedValue),
				CreatedBy: b.CreatedBy.GetOrZero().String(),
				UpdatedBy: b.UpdatedBy.GetOrZero().String(),
				ExpiredAt: b.ExpiredAt.Ptr(),
				CreatedAt: b.CreatedAt,
				UpdatedAt: b.UpdatedAt,
			},
		)
		providerMap[b.Provider] = orgMap
	}
	return providerMap, nil
}

type repositoryMap = map[matrix.GitProviderName]map[matrix.OrganizationID][]matrix.RepositoryView

func (r *gitProviderViewRepo) repositories(ctx context.Context, filter matrix.GitProviderViewFilter) (repositoryMap, error) {
	where := bobmodel.SelectWhere.Repositories
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		where.Provider.In(matrix.GitProviderNames...),
		bobmodel.SelectThenLoad.Repository.RepositoryPersonalAccessTokens(),
	}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(
			queryMods,
			where.OrganizationID.EQ(uuid.FromStringOrNil(*f)),
		)
	}

	bs, err := bobmodel.Repositories.Query(
		queryMods...,
	).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	providerMap := make(repositoryMap)
	for _, b := range bs {
		orgMap, ok := providerMap[b.Provider]
		if !ok {
			orgMap = map[matrix.OrganizationID][]matrix.RepositoryView{}
		}

		allowedPATs := make([]matrix.PersonalAccessTokenID, 0, len(b.R.RepositoryPersonalAccessTokens))
		for _, pat := range b.R.RepositoryPersonalAccessTokens {
			allowedPATs = append(allowedPATs, pat.PersonalAccessTokenID.String())
		}
		orgMap[b.OrganizationID.String()] = append(
			orgMap[b.OrganizationID.String()],
			matrix.RepositoryView{
				ID: b.ID.String(),
				// OrganizationID do not set
				// Provider do not set
				Namespace:              b.Namespace,
				Name:                   b.Name,
				Description:            b.Description.GetOrZero(),
				PersonalAccessTokenIDs: allowedPATs,
				CreatedBy:              b.CreatedBy.GetOrZero().String(),
				UpdatedBy:              b.UpdatedBy.GetOrZero().String(),
				CreatedAt:              b.CreatedAt,
				UpdatedAt:              b.UpdatedAt,
			},
		)
		providerMap[b.Provider] = orgMap
	}
	return providerMap, nil
}
