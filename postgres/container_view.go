package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

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

type containerViewRepo struct {
	db bob.DB
}

func NewContainerViewRepo(db *sql.DB) matrix.ContainerViewRepo {
	return &containerViewRepo{
		db: bob.NewDB(db),
	}
}

func (r *containerViewRepo) Find(ctx context.Context, clusterID matrix.ClusterID, namespace string, podID string, name string) (*matrix.ContainerView, error) {
	where := bobmodel.SelectWhere.Objects
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		where.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		where.Namespace.EQ(namespace),
		where.ID.EQ(uuid.FromStringOrNil(podID)),
		where.Raw.IsNotNull(),
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

	var pod v1.Pod
	if err := json.Unmarshal(b.Raw.GetOrZero().Val, &pod); err != nil {
		return nil, err
	}

	var organizationID string
	if b.R.Cluster != nil {
		organizationID = b.R.Cluster.OrganizationID.String()
	}

	for _, c := range pod.Spec.Containers {
		if c.Name == name {
			var state v1.ContainerState
			var status v1.ContainerStatus
			for _, s := range pod.Status.ContainerStatuses {
				if s.Name == c.Name {
					state = s.State
					status = s
				}
			}

			return r.b2d(
				b.ClusterID.String(),
				organizationID,
				b.Namespace,
				b.ID.String(),
				b.Name,
				c,
				state,
				status,
			)
		}
	}
	return nil, ekg.ErrNotFound
}

func (r *containerViewRepo) List(ctx context.Context, filter matrix.ContainerViewFilter) (*matrix.ContainerViewList, error) {
	objectTable := dbinfo.Objects.Name
	cols := dbinfo.Objects.Columns

	where := bobmodel.SelectWhere.Objects
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		where.Raw.IsNotNull(),
		bobmodel.SelectThenLoad.Object.Cluster(),
		where.Kind.EQ("Pod"),
	}

	if f := filter.OrganizationID; f != nil {
		clusterTable := dbinfo.Clusters.Name
		clusterCols := dbinfo.Clusters.Columns

		queryMods = append(
			queryMods,
			sm.InnerJoin(clusterTable).On(
				psql.Quote(objectTable, cols.ClusterID.Name).EQ(psql.Quote(clusterTable, "id")),
			),
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

	if len(filter.PodIDs) > 0 {
		ids := make([]uuid.UUID, 0, len(filter.PodIDs))
		for _, id := range filter.PodIDs {
			ids = append(ids, uuid.FromStringOrNil(id))
		}
		queryMods = append(queryMods, where.ID.In(ids...))
	}

	bs, err := bobmodel.Objects.Query(queryMods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	var ds = make([]matrix.ContainerView, 0, len(bs))
	for _, b := range bs {
		var organizationID string
		if b.R.Cluster != nil {
			organizationID = b.R.Cluster.OrganizationID.String()
		}

		var pod v1.Pod
		if err := json.Unmarshal(b.Raw.GetOrZero().Val, &pod); err != nil {
			return nil, err
		}

		for _, c := range pod.Spec.Containers {
			var state v1.ContainerState
			var status v1.ContainerStatus
			for _, s := range pod.Status.ContainerStatuses {
				if s.Name == c.Name {
					state = s.State
					status = s
				}
			}

			d, err := r.b2d(
				b.ClusterID.String(),
				organizationID,
				b.Namespace,
				b.ID.String(),
				b.Name,
				c,
				state,
				status,
			)
			if err != nil {
				return nil, err
			}
			ds = append(ds, *d)
		}
	}
	return &matrix.ContainerViewList{
		Containers: ds,
	}, nil
}

func (r *containerViewRepo) b2d(
	clusterID string,
	organizationID string,
	namespace string,
	podID string,
	podName string,
	container v1.Container,
	containerState v1.ContainerState,
	containerStatus v1.ContainerStatus,
) (*matrix.ContainerView, error) {
	image := matrix.FormatImage(container.Image)

	var createdAt time.Time
	if containerState.Running != nil {
		createdAt = containerState.Running.StartedAt.Time
	}

	return &matrix.ContainerView{
		OrganizationID:  organizationID,
		ClusterID:       clusterID,
		Namespace:       namespace,
		PodID:           podID,
		PodName:         podName,
		Name:            container.Name,
		ImageName:       image.Name,
		ImageVersion:    image.Version,
		ImageRepository: image.Repository,
		Ready:           containerStatus.Ready,
		CreatedAt:       createdAt,
	}, nil
}
