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

type podViewRepo struct {
	db bob.DB
}

func NewPodViewRepo(db *sql.DB) matrix.PodViewRepo {
	return &podViewRepo{
		db: bob.NewDB(db),
	}
}

func (r *podViewRepo) Find(ctx context.Context, clusterID matrix.ClusterID, namespace string, id string) (*matrix.PodView, error) {
	where := bobmodel.SelectWhere.Objects
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		where.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		where.Namespace.EQ(namespace),
		where.ID.EQ(uuid.FromStringOrNil(id)),
		bobmodel.SelectThenLoad.Object.Cluster(),
		where.Kind.EQ("Pod"),
	}

	b, err := bobmodel.Objects.Query(queryMods...).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b)
}

func (r *podViewRepo) List(ctx context.Context, filter matrix.PodViewFilter) (*matrix.PodViewList, error) {
	where := bobmodel.SelectWhere.Objects
	join := bobmodel.SelectJoins.Objects
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		bobmodel.SelectThenLoad.Object.Cluster(),
		where.Kind.EQ("Pod"),
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

	if len(filter.Namespaces) > 0 {
		queryMods = append(queryMods, where.Namespace.In(filter.Namespaces...))
	}

	bs, err := bobmodel.Objects.Query(queryMods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	var ds = make([]matrix.PodView, 0, len(bs))
	for _, b := range bs {
		d, err := r.b2d(b)
		if err != nil {
			return nil, err
		}
		ds = append(ds, *d)
	}
	return &matrix.PodViewList{
		Pods: ds,
	}, nil
}

func (r *podViewRepo) b2d(b *bobmodel.Object) (*matrix.PodView, error) {
	var organizationID string
	if b.R.Cluster != nil {
		organizationID = b.R.Cluster.OrganizationID.String()
	}

	var pod v1.Pod
	if err := json.Unmarshal(b.Raw.GetOrZero().Val, &pod); err != nil {
		return nil, err
	}

	restartCount := 0
	containerCount := 0
	for _, s := range pod.Status.ContainerStatuses {
		restartCount += int(s.RestartCount)
		containerCount += 1
	}

	age := toAge(pod.CreationTimestamp.Time)

	return &matrix.PodView{
		OrganizationID: organizationID,
		ClusterID:      b.ClusterID.String(),
		Namespace:      b.Namespace,
		ID:             b.ID.String(),
		Name:           b.Name,
		Status:         string(pod.Status.Phase),
		RestartCount:   restartCount,
		ContainerCount: containerCount,
		Age:            age,
		CreatedAt:      b.CreatedAt,
	}, nil
}
