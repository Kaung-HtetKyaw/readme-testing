package postgres

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	dbinfo "github.com/kubegrade/matrix/postgres/bob_dbinfo"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/types"
)

type objectRepo struct {
	db bob.DB
}

func NewObjectRepo(db *sql.DB) matrix.ObjectRepo {
	return &objectRepo{db: bob.NewDB(db)}
}

func (r *objectRepo) Create(ctx context.Context, object *matrix.Object) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	b := r.d2b(object)

	healthStatus, err := objectHealthStatus(object.Kind, object.Raw)
	if err != nil {
		return err
	}
	b.HealthStatus = omit.From(healthStatus)

	if _, err := bobmodel.Objects.Insert(b).One(ctx, tx); err != nil {
		return err
	}
	if err := refreshClusterResourceCount(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *objectRepo) Store(ctx context.Context, object *matrix.Object) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	b := r.d2b(object)

	healthStatus, err := objectHealthStatus(object.Kind, object.Raw)
	if err != nil {
		return err
	}
	b.HealthStatus = omit.From(healthStatus)

	cols := dbinfo.Objects.Columns
	if _, err := bobmodel.Objects.Insert(
		b,
		im.OnConflict(
			cols.ClusterID.Name,
			cols.ID.Name,
		).DoUpdate(
			im.SetExcluded(
				cols.RunVersion.Name,
			),
		),
	).One(ctx, tx); err != nil {
		return err
	}
	if err := refreshClusterResourceCount(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *objectRepo) Find(ctx context.Context, clusterID matrix.ClusterID, id matrix.ObjectID) (*matrix.Object, error) {
	where := bobmodel.SelectWhere.Objects
	b, err := bobmodel.Objects.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
		where.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b), nil
}

func (r *objectRepo) List(
	ctx context.Context,
	filter matrix.ObjectFilter,
	sort matrix.ObjectSort,
	pagination matrix.CursorPagination,
) ([]*matrix.Object, matrix.CursorPagination, error) {
	queryMods := []bob.Mod[*dialect.SelectQuery]{
		bobmodel.SelectThenLoad.Object.Cluster(),
	}

	if pagination.PageSize == 0 {
		count, err := bobmodel.Objects.Query(r.withFilter(filter)...).Count(ctx, r.db)
		if err != nil {
			return nil, matrix.CursorPagination{}, err
		}
		pagination.PageSize = int(count)
	}

	if cursor := pagination.NextCursor; cursor != "" {
		nextCursorQuery, err := r.withNextCursor(cursor, sort)
		if err != nil {
			return nil, matrix.CursorPagination{}, err
		}
		queryMods = append(queryMods, nextCursorQuery...)
	}

	if cursor := pagination.PrevCursor; cursor != "" {
		prevCursorQuery, err := r.withPrevCursor(cursor, pagination.PageSize, filter, sort)
		if err != nil {
			return nil, matrix.CursorPagination{}, err
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
			return nil, matrix.CursorPagination{}, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, matrix.CursorPagination{}, err
	}

	nextPagination, err := r.nextPagination(
		ctx,
		filter,
		sort,
		pagination,
		bs,
	)
	if err != nil {
		return nil, matrix.CursorPagination{}, err
	}

	var ds = make([]*matrix.Object, 0, len(bs))
	for _, b := range bs {
		d := r.b2d(b)
		ds = append(ds, d)
	}

	return ds, nextPagination, nil
}

func (r *objectRepo) Update(ctx context.Context, update matrix.ObjectUpdate) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	id := uuid.FromStringOrNil(update.ID)
	clusterID := uuid.FromStringOrNil(update.ClusterID)

	where := bobmodel.SelectWhere.Objects
	b, err := bobmodel.Objects.Query(
		where.ID.EQ(id),
		where.ClusterID.EQ(clusterID),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	setter := &bobmodel.ObjectSetter{
		ID:        omit.From(id),
		ClusterID: omit.From(clusterID),
	}

	if name := update.Payload.Name; name != nil {
		setter.Name = omit.From(*name)
	}

	if namespace := update.Payload.Namespace; namespace != nil {
		setter.Namespace = omit.From(*namespace)
	}

	if resourceVersion := update.Payload.ResourceVersion; resourceVersion != nil {
		setter.ResourceVersion = omit.From(*resourceVersion)
	}

	if raw := update.Payload.Raw; raw != nil {
		rawJSON := omitnull.Val[types.JSON[json.RawMessage]]{}
		rawJSON.Set(types.JSON[json.RawMessage]{raw})
		setter.Raw = rawJSON

		healthStatus, err := objectHealthStatus(b.Kind, raw)
		if err != nil {
			return err
		}
		setter.HealthStatus = omit.From(healthStatus)
	}

	if err := b.Update(ctx, tx, setter); err != nil {
		return err
	}

	if err := refreshClusterResourceCount(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *objectRepo) Delete(ctx context.Context, clusterID matrix.ClusterID, id matrix.ObjectID) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	where := bobmodel.SelectWhere.Objects
	b, err := bobmodel.Objects.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
		where.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	if err := b.Delete(ctx, tx); err != nil {
		return err
	}

	if err := refreshClusterResourceCount(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// PurgeStaleObject
func (r *objectRepo) PurgeStaleObject(ctx context.Context, clusterID matrix.ClusterID, currentRunVersion string) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	deleteWhere := bobmodel.DeleteWhere.Objects
	if _, err := bobmodel.Objects.Delete(
		deleteWhere.ClusterID.EQ(uuid.FromStringOrNil(clusterID)),
		deleteWhere.RunVersion.NE(uuid.FromStringOrNil(currentRunVersion)),
	).All(ctx, r.db); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *objectRepo) d2b(d *matrix.Object) *bobmodel.ObjectSetter {
	rawJSON := omitnull.Val[types.JSON[json.RawMessage]]{}
	rawJSON.Set(types.JSON[json.RawMessage]{d.Raw})

	return &bobmodel.ObjectSetter{
		ID:              omit.From(uuid.FromStringOrNil(d.ID)),
		ClusterID:       omit.From(uuid.FromStringOrNil(d.ClusterID)),
		Namespace:       omit.From(d.Namespace),
		Name:            omit.From(d.Name),
		ResourceVersion: omit.From(d.ResourceVersion),
		RunVersion:      omit.From(uuid.FromStringOrNil(d.RunVersion)),
		Raw:             rawJSON,
		Kind:            omit.From(d.Kind),
		CreatedAt:       omit.From(d.CreatedAt),
	}
}

func (r *objectRepo) b2d(b *bobmodel.Object) *matrix.Object {
	return &matrix.Object{
		ID:              b.ID.String(),
		ClusterID:       b.ClusterID.String(),
		Namespace:       b.Namespace,
		Name:            b.Name,
		ResourceVersion: b.ResourceVersion,
		RunVersion:      b.RunVersion.String(),
		Raw:             b.Raw.GetOrZero().Val,
		Kind:            b.Kind,
		CreatedAt:       b.CreatedAt,
	}
}

func (r *objectRepo) withFilter(filter matrix.ObjectFilter) []bob.Mod[*dialect.SelectQuery] {
	queryMods := []bob.Mod[*dialect.SelectQuery]{}
	where := bobmodel.SelectWhere.Objects

	if f := filter.Kind; f != nil {
		queryMods = append(queryMods, where.Kind.EQ(*f))
	}

	return queryMods
}

func (r *objectRepo) withSort(sort matrix.ObjectSort) []bob.Mod[*dialect.SelectQuery] {
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

func (r *objectRepo) withNextCursor(cursor string, sort matrix.ObjectSort) ([]bob.Mod[*dialect.SelectQuery], error) {
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

func (r *objectRepo) withPrevCursor(cursor string, limit int, filter matrix.ObjectFilter, sort matrix.ObjectSort) ([]bob.Mod[*dialect.SelectQuery], error) {
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

func (r *objectRepo) nextPagination(
	ctx context.Context,
	filter matrix.ObjectFilter,
	sort matrix.ObjectSort,
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

func (r *objectRepo) encodeCursor(cursor *objectCursor) (string, error) {
	b, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (r *objectRepo) decodeCursor(cursorStr string) (*objectCursor, error) {
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
