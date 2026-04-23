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

type personalAccessTokenViewRepo struct {
	db bob.DB
}

func NewPersonalAccessTokenViewRepo(db *sql.DB) matrix.PersonalAccessTokenViewRepo {
	return &personalAccessTokenViewRepo{db: bob.NewDB(db)}
}

func (r *personalAccessTokenViewRepo) Find(ctx context.Context, redacted matrix.PersonalAccessTokenRedacted, id matrix.PersonalAccessTokenID) (*matrix.PersonalAccessTokenView, error) {
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
	return r.b2d(redacted, b)
}

func (r *personalAccessTokenViewRepo) List(ctx context.Context, redacted matrix.PersonalAccessTokenRedacted, filter matrix.PersonalAccessTokenViewFilter) (*matrix.PersonalAccessTokenViewList, error) {
	where := bobmodel.SelectWhere.PersonalAccessTokens
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(queryMods, where.OrganizationID.EQ(uuid.FromStringOrNil(*f)))
	}
	if f := filter.Provider; f != nil {
		queryMods = append(queryMods, where.Provider.EQ(*f))
	}

	bs, err := bobmodel.PersonalAccessTokens.Query(queryMods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	var ds = make([]matrix.PersonalAccessTokenView, 0, len(bs))
	for _, b := range bs {
		d, err := r.b2d(redacted, b)
		if err != nil {
			return nil, err
		}
		ds = append(ds, *d)
	}
	return &matrix.PersonalAccessTokenViewList{
		Items: ds,
	}, nil
}

// TODO encrypt value
func (r *personalAccessTokenViewRepo) b2d(redacted matrix.PersonalAccessTokenRedacted, b *bobmodel.PersonalAccessToken) (*matrix.PersonalAccessTokenView, error) {
	if redacted.Value {
		b.EncryptedValue = redactValue(b.EncryptedValue)
	}

	return &matrix.PersonalAccessTokenView{
		ID:             b.ID.String(),
		OrganizationID: b.OrganizationID.String(),
		Owner:          b.Owner,
		Provider:       b.Provider,
		Name:           b.Name,
		Value:          b.EncryptedValue,
		CreatedBy:      b.CreatedBy.GetOrZero().String(),
		UpdatedBy:      b.UpdatedBy.GetOrZero().String(),
		ExpiredAt:      b.ExpiredAt.Ptr(),
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}, nil
}
