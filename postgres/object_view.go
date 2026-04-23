package postgres

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

type objectViewRepo struct {
	db bob.DB
}

func NewObjectViewRepo(db *sql.DB) matrix.ObjectViewRepo {
	return &objectViewRepo{db: bob.NewDB(db)}
}

func (r *objectViewRepo) Find(ctx context.Context, clusterID matrix.ClusterID, id string) (*matrix.ObjectView, error) {
	where := bobmodel.SelectWhere.Objects
	b, err := bobmodel.Objects.Query(
		where.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		where.ID.EQ(uuid.FromStringOrNil(id)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b, nil)
}

func (r *objectViewRepo) List(
	ctx context.Context,
	fields []string,
	filter matrix.ObjectViewFilter,
	sort matrix.ObjectViewSort,
	pagination matrix.CursorPagination,
) (*matrix.ObjectViewList, error) {
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		bobmodel.SelectThenLoad.Object.Cluster(),
	}

	if len(fields) > 0 && !slices.Contains(fields, "*") {
		queryMods = append(queryMods, r.withFieldSelect(fields)...)
	}

	if pagination.PageSize == 0 {
		count, err := bobmodel.Objects.Query(r.withFilter(filter)...).Count(ctx, r.db)
		if err != nil {
			return nil, err
		}
		pagination.PageSize = int(count)
	}

	if cursor := pagination.NextCursor; cursor != "" {
		nextCursorQuery, err := r.withNextCursor(cursor, sort)
		if err != nil {
			return nil, err
		}
		queryMods = append(queryMods, nextCursorQuery...)
	}

	if cursor := pagination.PrevCursor; cursor != "" {
		prevCursorQuery, err := r.withPrevCursor(cursor, pagination.PageSize, filter, sort)
		if err != nil {
			return nil, err
		}
		queryMods = append(queryMods, prevCursorQuery...)
	}

	if pagination.PrevCursor == "" {
		queryMods = append(queryMods, r.withFilter(filter)...)
		queryMods = append(queryMods, sm.Limit(pagination.PageSize))
	}
	queryMods = append(queryMods, r.withSort(sort)...)

	bs, err := bobmodel.Objects.Query(queryMods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	nextPagination, err := r.nextPagination(
		ctx,
		filter,
		sort,
		pagination,
		bs,
	)
	if err != nil {
		return nil, err
	}

	var ds = make([]matrix.ObjectView, 0, len(bs))
	for _, b := range bs {
		d, err := r.b2d(b, fields)
		if err != nil {
			return nil, err
		}
		ds = append(ds, *d)
	}

	return &matrix.ObjectViewList{
		Pagination: nextPagination,
		Items:      ds,
	}, nil
}

// GetManifest
func (r *objectViewRepo) GetManifest(ctx context.Context, clusterID matrix.ClusterID, id string) ([]byte, error) {
	where := bobmodel.SelectWhere.Objects
	b, err := bobmodel.Objects.Query(
		bobmodel.SelectThenLoad.Object.Cluster(),
		where.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		where.ID.EQ(uuid.FromStringOrNil(id)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	manifest, err := yaml.JSONToYAML(b.Raw.GetOrZero().Val)
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

func (r *objectViewRepo) b2d(b *bobmodel.Object, fields []string) (*matrix.ObjectView, error) {
	objectView := &matrix.ObjectView{}

	var organizationID string
	if b.R.Cluster != nil {
		organizationID = b.R.Cluster.OrganizationID.String()
	}

	var raw unstructured.Unstructured
	if !b.Raw.IsNull() {
		if err := json.Unmarshal(b.Raw.GetOrZero().Val, &raw); err != nil {
			return nil, err
		}
	}

	var spec json.RawMessage
	rawSpec, found := raw.Object["spec"]
	if found {
		j, err := json.Marshal(rawSpec)
		if err != nil {
			return nil, err
		}
		spec = j
	}

	var status json.RawMessage
	rawStatus, found := raw.Object["status"]
	if found {
		j, err := json.Marshal(rawStatus)
		if err != nil {
			return nil, err
		}
		status = j
	}

	if len(fields) == 0 || slices.Contains(fields, "*") {
		objectView = &matrix.ObjectView{
			ClusterID:      b.ClusterID.String(),
			ID:             b.ID.String(),
			OrganizationID: organizationID,
			Namespace:      b.Namespace,
			Name:           b.Name,
			Kind:           b.Kind,
			HeathStatus:    b.HealthStatus,
			APIVersion:     raw.GetAPIVersion(),
			Labels:         raw.GetLabels(),
			Annotations:    raw.GetAnnotations(),
			Spec:           spec,
			Status:         status,
			CreatedAt:      b.CreatedAt,
		}
	} else {
		for _, f := range fields {
			// ID and CreatedAt are mandatory fields.
			objectView.ID = b.ID.String()
			objectView.CreatedAt = b.CreatedAt

			switch f {
			case matrix.ObjectFieldClusterID:
				objectView.ClusterID = b.ClusterID.String()
			case matrix.ObjectFieldOrganizationID:
				objectView.OrganizationID = organizationID
			case matrix.ObjectFieldNamespace:
				objectView.Namespace = b.Namespace
			case matrix.ObjectFieldName:
				objectView.Name = b.Name
			case matrix.ObjectFieldKind:
				objectView.Kind = b.Kind
			case matrix.ObjectFieldHealthStatus:
				objectView.HeathStatus = b.HealthStatus
			case matrix.ObjectFieldAPIVersion:
				objectView.APIVersion = raw.GetAPIVersion()
			case matrix.ObjectFieldLabels:
				objectView.Labels = raw.GetLabels()
			case matrix.ObjectFieldAnnotaions:
				objectView.Annotations = raw.GetAnnotations()
			case matrix.ObjectFieldSpec:
				objectView.Spec = spec
			case matrix.ObjectFieldStatus:
				objectView.Status = status
			}
		}
	}
	return objectView, nil
}

func (r *objectViewRepo) withFilter(filter matrix.ObjectViewFilter) []bob.Mod[*dialect.SelectQuery] {
	queryMods := []bob.Mod[*dialect.SelectQuery]{}
	where := bobmodel.SelectWhere.Objects

	if f := filter.ObjectIDs; f != nil && len(f) > 0 {
		ids := make([]uuid.UUID, 0, len(f))
		for _, id := range f {
			ids = append(ids, uuid.FromStringOrNil(id))
		}
		queryMods = append(queryMods, where.ID.In(ids...))
	}

	if f := filter.OrganizationID; f != nil {
		cols := bobmodel.Objects.Columns
		clusterCols := bobmodel.Clusters.Columns
		queryMods = append(
			queryMods,
			sm.Where(
				cols.ClusterID.In(
					psql.Select(
						sm.Columns(
							clusterCols.ID,
						),
						sm.Where(clusterCols.OrganizationID.EQ(psql.Arg(*f))),
						sm.From(dbinfo.Clusters.Name),
					),
				),
			),
		)
	}

	if f := filter.ClusterID; f != nil {
		queryMods = append(queryMods, where.ClusterID.EQ(uuid.FromStringOrNil(*f)))
	}

	if f := filter.Namespaces; f != nil && len(f) > 0 {
		queryMods = append(queryMods, where.Namespace.In(f...))
	}

	if f := filter.Namespaced; f != nil {
		if *f {
			queryMods = append(queryMods, where.Namespace.NE(""))
		} else {
			queryMods = append(queryMods, where.Namespace.EQ(""))
		}
	}

	if f := filter.Kind; f != nil {
		queryMods = append(queryMods, where.Kind.ILike(*f))
	}

	return queryMods
}

func (r *objectViewRepo) withSort(sort matrix.ObjectViewSort) []bob.Mod[*dialect.SelectQuery] {
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	tableName := dbinfo.Objects.Name
	cols := dbinfo.Objects.Columns

	sortCreatedAt := ensureSort(sort.SortCreatedAt)
	sortID := ensureSort(sort.SortID)

	queryMods = append(
		queryMods,
		sm.OrderBy(fmt.Sprintf("%s.%s %s", tableName, cols.CreatedAt.Name, sortCreatedAt)),
		sm.OrderBy(fmt.Sprintf("%s.%s %s", tableName, cols.ID.Name, sortID)),
	)
	return queryMods
}

func (r *objectViewRepo) withFieldSelect(fields []string) []bob.Mod[*dialect.SelectQuery] {
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	tableName := dbinfo.Objects.Name
	cols := dbinfo.Objects.Columns

	isRawSelected := false

	for _, f := range fields {
		switch f {
		case matrix.ObjectFieldNamespace:
			queryMods = append(queryMods, sm.Columns(cols.Namespace.Name))
		case matrix.ObjectFieldName:
			queryMods = append(queryMods, sm.Columns(cols.Name.Name))
		case matrix.ObjectFieldKind:
			queryMods = append(queryMods, sm.Columns(cols.Kind.Name))
		case matrix.ObjectFieldHealthStatus:
			queryMods = append(queryMods, sm.Columns(cols.HealthStatus.Name))
		case matrix.ObjectFieldAPIVersion,
			matrix.ObjectFieldLabels,
			matrix.ObjectFieldAnnotaions,
			matrix.ObjectFieldStatus,
			matrix.ObjectFieldSpec:
			isRawSelected = true
		}
	}

	if isRawSelected {
		queryMods = append(queryMods, sm.Columns(cols.Raw.Name))
	}

	// ClusterID must be always selected, as it is needed for relationship loading.
	queryMods = append(queryMods, sm.Columns(cols.ClusterID.Name))

	// Always select id and created_at as they are part of cursor.
	queryMods = append(queryMods, sm.Columns(
		psql.Quote(tableName, cols.ID.Name),
	))
	queryMods = append(queryMods, sm.Columns(
		psql.Quote(tableName, cols.CreatedAt.Name),
	))

	return queryMods
}

func (r *objectViewRepo) withNextCursor(cursor string, sort matrix.ObjectViewSort) ([]bob.Mod[*dialect.SelectQuery], error) {
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	cols := bobmodel.Objects.Columns
	c, err := r.decodeCursor(cursor)
	if err != nil {
		return nil, err
	}

	sortCreatedAt := ensureSort(sort.SortCreatedAt)
	sortID := ensureSort(sort.SortID)

	createdAtExpression := cols.CreatedAt.GT(psql.Arg(c.CreatedAt))
	if sortCreatedAt == "DESC" {
		createdAtExpression = cols.CreatedAt.LT(psql.Arg(c.CreatedAt))
	}

	idExpression := psql.And(
		cols.ID.GT(psql.Arg(c.ID)),
		cols.CreatedAt.EQ(psql.Arg(c.CreatedAt)),
	)
	if sortID == "DESC" {
		idExpression = psql.And(
			cols.ID.LT(psql.Arg(c.ID)),
			cols.CreatedAt.EQ(psql.Arg(c.CreatedAt)),
		)
	}

	queryMods = append(
		queryMods,
		sm.Where(
			psql.Or(
				createdAtExpression,
				idExpression,
			),
		),
	)

	return queryMods, nil
}

func (r *objectViewRepo) withPrevCursor(cursor string, limit int, filter matrix.ObjectViewFilter, sort matrix.ObjectViewSort) ([]bob.Mod[*dialect.SelectQuery], error) {
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	cols := bobmodel.Objects.Columns
	c, err := r.decodeCursor(cursor)
	if err != nil {
		return nil, err
	}

	sortCreatedAt := ensureSort(sort.SortCreatedAt)
	sortID := ensureSort(sort.SortID)

	reverseSortCreatedAt := "DESC"
	reverseSortID := "DESC"

	createdAtExpression := cols.CreatedAt.LT(psql.Arg(c.CreatedAt))
	if sortCreatedAt == "DESC" {
		createdAtExpression = cols.CreatedAt.GT(psql.Arg(c.CreatedAt))
		reverseSortCreatedAt = "ASC"
	}

	idExpression := psql.And(
		cols.ID.LT(psql.Arg(c.ID)),
		cols.CreatedAt.EQ(psql.Arg(c.CreatedAt)),
	)
	if sortID == "DESC" {
		idExpression = psql.And(
			cols.ID.GT(psql.Arg(c.ID)),
			cols.CreatedAt.EQ(psql.Arg(c.CreatedAt)),
		)
		reverseSortID = "ASC"
	}

	queryForPrevCursor := r.withFilter(filter)
	queryForPrevCursor = append(
		queryForPrevCursor,
		sm.From(bobmodel.Objects.Name()),
		sm.Where(
			psql.Or(
				createdAtExpression,
				idExpression,
			),
		),
		sm.OrderBy(fmt.Sprintf("%s %s", cols.CreatedAt, reverseSortCreatedAt)),
		sm.OrderBy(fmt.Sprintf("%s %s", cols.ID, reverseSortID)),
		sm.Limit(limit),
	)

	queryMods = append(
		queryMods,
		sm.From(
			psql.Select(
				queryForPrevCursor...,
			),
		).As("object"),
	)
	return queryMods, nil
}

func (r *objectViewRepo) nextPagination(
	ctx context.Context,
	filter matrix.ObjectViewFilter,
	sort matrix.ObjectViewSort,
	current matrix.CursorPagination,
	bObject []*bobmodel.Object,
) (matrix.CursorPagination, error) {
	objectSize := len(bObject)

	nextCursor, err := func() (string, error) {
		if objectSize == 0 {
			return "", nil
		}

		queryMods := []bob.Mod[*dialect.SelectQuery]{}
		queryMods = append(queryMods, r.withFilter(filter)...)
		queryMods = append(queryMods, sm.Limit(current.PageSize+1))

		cols := bobmodel.Objects.Columns

		sortCreatedAt := ensureSort(sort.SortCreatedAt)
		sortID := ensureSort(sort.SortID)

		createdAt := bObject[objectSize-1].CreatedAt
		id := bObject[objectSize-1].ID

		createdAtExpression := cols.CreatedAt.GT(psql.Arg(createdAt))
		if sortCreatedAt == "DESC" {
			createdAtExpression = cols.CreatedAt.LT(psql.Arg(createdAt))
		}

		idExpression := psql.And(
			cols.ID.GT(psql.Arg(id)),
			cols.CreatedAt.EQ(psql.Arg(createdAt)),
		)
		if sortID == "DESC" {
			idExpression = psql.And(
				cols.ID.LT(psql.Arg(id)),
				cols.CreatedAt.EQ(psql.Arg(createdAt)),
			)
		}

		queryMods = append(
			queryMods,
			sm.Where(
				psql.Or(
					createdAtExpression,
					idExpression,
				),
			),
		)

		count, err := bobmodel.Objects.Query(queryMods...).Count(ctx, r.db)
		if err != nil {
			return "", err
		}

		nextCursor := ""
		if count > 0 {
			c, err := r.encodeCursor(&objectCursor{
				ID:        id.String(),
				CreatedAt: createdAt,
			})
			if err != nil {
				return "", err
			}
			nextCursor = c
		}
		return nextCursor, nil
	}()
	if err != nil {
		return matrix.CursorPagination{}, err
	}

	prevCursor, err := func() (string, error) {
		if objectSize == 0 {
			return "", nil
		}

		cols := bobmodel.Objects.Columns

		sortCreatedAt := ensureSort(sort.SortCreatedAt)
		sortID := ensureSort(sort.SortID)

		queryMods := []bob.Mod[*dialect.SelectQuery]{}
		queryMods = append(queryMods, r.withFilter(filter)...)
		queryMods = append(queryMods, sm.Limit(current.PageSize))
		queryMods = append(
			queryMods,
			sm.OrderBy(fmt.Sprintf(
				"%s %s",
				cols.CreatedAt,
				sortCreatedAt,
			)),
			sm.OrderBy(fmt.Sprintf(
				"%s %s",
				cols.ID,
				sortID,
			)),
		)

		createdAt := bObject[0].CreatedAt
		id := bObject[0].ID

		createdAtExpression := cols.CreatedAt.LT(psql.Arg(createdAt))
		if sortCreatedAt == "DESC" {
			createdAtExpression = cols.CreatedAt.GT(psql.Arg(createdAt))
		}

		idExpression := psql.And(
			cols.ID.LT(psql.Arg(id)),
			cols.CreatedAt.EQ(psql.Arg(createdAt)),
		)
		if sortID == "DESC" {
			idExpression = psql.And(
				cols.ID.GT(psql.Arg(id)),
				cols.CreatedAt.EQ(psql.Arg(createdAt)),
			)
		}

		queryMods = append(queryMods,
			sm.Where(
				psql.Or(
					createdAtExpression,
					idExpression,
				),
			),
			sm.Limit(current.PageSize+1),
		)

		count, err := bobmodel.Objects.Query(queryMods...).Count(ctx, r.db)
		if err != nil {
			return "", err
		}

		prevCursor := ""
		if count >= int64(objectSize) {
			createdAt := bObject[0].CreatedAt
			id := bObject[0].ID
			c, err := r.encodeCursor(&objectCursor{
				ID:        id.String(),
				CreatedAt: createdAt,
			})
			if err != nil {
				return "", err
			}
			prevCursor = c
		}
		return prevCursor, nil
	}()
	if err != nil {
		return matrix.CursorPagination{}, err
	}

	return matrix.CursorPagination{
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		PageSize:   objectSize,
	}, nil
}

type objectCursor struct {
	CreatedAt time.Time `json:"createdAt"`
	ID        string    `json:"id"`
}

func (r *objectViewRepo) encodeCursor(cursor *objectCursor) (string, error) {
	b, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (r *objectViewRepo) decodeCursor(cursorStr string) (*objectCursor, error) {
	b, err := base64.StdEncoding.DecodeString(cursorStr)
	if err != nil {
		return nil, err
	}

	var cursor objectCursor
	if err := json.Unmarshal(b, &cursor); err != nil {
		return nil, err
	}

	return &cursor, nil
}
