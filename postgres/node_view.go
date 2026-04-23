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

type nodeViewRepo struct {
	db bob.DB
}

func NewNodeViewRepo(db *sql.DB) matrix.NodeViewRepo {
	return &nodeViewRepo{
		db: bob.NewDB(db),
	}
}

func (r *nodeViewRepo) Find(ctx context.Context, clusterID matrix.ClusterID, id string) (*matrix.NodeView, error) {
	where := bobmodel.SelectWhere.Objects
	b, err := bobmodel.Objects.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
		where.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		bobmodel.SelectThenLoad.Object.Cluster(),
		where.Kind.EQ("Node"),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b)
}

func (r *nodeViewRepo) List(ctx context.Context, filter matrix.NodeViewFilter) (*matrix.NodeViewList, error) {
	where := bobmodel.SelectWhere.Objects
	join := bobmodel.SelectJoins.Objects
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		bobmodel.SelectThenLoad.Object.Cluster(),
		where.Kind.EQ("Node"),
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

	var ds = make([]matrix.NodeView, 0, len(bs))
	for _, b := range bs {
		d, err := r.b2d(b)
		if err != nil {
			return nil, err
		}
		ds = append(ds, *d)
	}
	return &matrix.NodeViewList{
		Nodes: ds,
	}, nil
}

func (r *nodeViewRepo) b2d(b *bobmodel.Object) (*matrix.NodeView, error) {
	var organizationID string
	if b.R.Cluster != nil {
		organizationID = b.R.Cluster.OrganizationID.String()
	}

	var node v1.Node
	if err := json.Unmarshal(b.Raw.GetOrZero().Val, &node); err != nil {
		return nil, err
	}

	condition := matrix.NodeViewCondition{
		Status: "unknown",
		Type:   "unknown",
		Reason: "unknown",
	}

	var lastTransitionTime *time.Time = nil
	for _, c := range node.Status.Conditions {
		if lastTransitionTime == nil || c.LastTransitionTime.Time.After(*lastTransitionTime) {
			lastTransitionTime = &c.LastTransitionTime.Time
			condition.Status = string(c.Status)
			condition.Type = string(c.Type)
			condition.Reason = string(c.Reason)
		}
	}

	age := toAge(node.CreationTimestamp.Time)
	return &matrix.NodeView{
		ID:               b.ID.String(),
		ClusterID:        b.ClusterID.String(),
		OrganizationID:   organizationID,
		Name:             b.Name,
		KubeletVersion:   node.Status.NodeInfo.KubeletVersion,
		CurrentCondition: condition,
		Age:              age,
		CreatedAt:        b.CreatedAt,
	}, nil
}
