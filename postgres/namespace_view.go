package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	dbinfo "github.com/kubegrade/matrix/postgres/bob_dbinfo"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	v1 "k8s.io/api/core/v1"
)

type namespaceViewRepo struct {
	db bob.DB
}

func NewNamespaceViewRepo(db *sql.DB) matrix.NamespaceViewRepo {
	return &namespaceViewRepo{
		db: bob.NewDB(db),
	}
}

func (r *namespaceViewRepo) Find(ctx context.Context, clusterID matrix.ClusterID, id string) (*matrix.NamespaceView, error) {
	where := bobmodel.SelectWhere.Objects
	b, err := bobmodel.Objects.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
		where.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		bobmodel.SelectThenLoad.Object.Cluster(),
		where.Kind.EQ("Namespace"),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b)
}

func (r *namespaceViewRepo) List(ctx context.Context, filter matrix.NamespaceViewFilter) (*matrix.NamespaceViewList, error) {
	where := bobmodel.SelectWhere.Objects
	join := bobmodel.SelectJoins.Objects
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		bobmodel.SelectThenLoad.Object.Cluster(),
		where.Kind.EQ("Namespace"),
	}

	if f := filter.OrganizationID; f != nil {
		clusterCols := dbinfo.Clusters.Columns
		queryMods = append(
			queryMods,
			join.InnerJoin.Cluster,
			sm.Where(
				psql.And(
					psql.Quote(clusterCols.OrganizationID.Name).EQ(psql.Arg(*f)),
				),
			),
		)
	}

	if f := filter.ClusterID; f != nil {
		queryMods = append(queryMods, where.ClusterID.EQ(uuid.FromStringOrNil(*f)))
	}

	bs, err := bobmodel.Objects.Query(queryMods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	var ds = make([]matrix.NamespaceView, 0, len(bs))
	for _, b := range bs {
		d, err := r.b2d(b)
		if err != nil {
			return nil, err
		}
		ds = append(ds, *d)
	}
	return &matrix.NamespaceViewList{
		Namespaces: ds,
	}, nil
}

func (r *namespaceViewRepo) b2d(b *bobmodel.Object) (*matrix.NamespaceView, error) {
	var organizationID string
	if b.R.Cluster != nil {
		organizationID = b.R.Cluster.OrganizationID.String()
	}

	var namespace v1.Namespace
	if err := json.Unmarshal(b.Raw.GetOrZero().Val, &namespace); err != nil {
		return nil, err
	}

	age := toAge(namespace.CreationTimestamp.Time)
	return &matrix.NamespaceView{
		ID:             b.ID.String(),
		ClusterID:      b.ClusterID.String(),
		OrganizationID: organizationID,
		Name:           b.Name,
		Status:         string(namespace.Status.Phase),
		Age:            age,
		CreatedAt:      b.CreatedAt,
	}, nil
}
