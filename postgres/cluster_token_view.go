package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
)

type clusterTokenViewRepo struct {
	db bob.DB
}

func NewClusterTokenViewRepo(db *sql.DB) matrix.ClusterTokenViewRepo {
	return &clusterTokenViewRepo{
		db: bob.NewDB(db),
	}
}

func (r *clusterTokenViewRepo) Find(ctx context.Context, value string) (*matrix.ClusterTokenView, error) {
	where := bobmodel.SelectWhere.ClusterTokens
	b, err := bobmodel.ClusterTokens.Query(
		where.Value.EQ(value),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b), nil
}

func (r *clusterTokenViewRepo) List(ctx context.Context, filter matrix.ClusterTokenViewFilter) (*matrix.ClusterTokenViewList, error) {
	where := bobmodel.SelectWhere.ClusterTokens
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(queryMods, where.OrganizationID.EQ(uuid.FromStringOrNil(*f)))
	}

	if f := filter.Name; f != nil {
		likePattern := fmt.Sprintf("%%%s%%", *f)
		queryMods = append(queryMods, where.Name.Like(likePattern))
	}

	bs, err := bobmodel.ClusterTokens.Query(queryMods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	var ds = make([]matrix.ClusterTokenView, 0, len(bs))
	for _, b := range bs {
		d := r.b2d(b)
		ds = append(ds, *d)
	}
	return &matrix.ClusterTokenViewList{
		ClusterTokens: ds,
	}, nil
}

func (r *clusterTokenViewRepo) b2d(b *bobmodel.ClusterToken) *matrix.ClusterTokenView {
	return &matrix.ClusterTokenView{
		Value:          b.Value,
		Name:           b.Name,
		OrganizationID: b.OrganizationID.String(),
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}
