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

type deprecatedAPIViewRepo struct {
	db bob.DB
}

func NewDeprecatedAPIViewRepo(db *sql.DB) matrix.DeprecatedAPIViewRepo {
	return &deprecatedAPIViewRepo{db: bob.NewDB(db)}
}

func (r *deprecatedAPIViewRepo) List(ctx context.Context, filter matrix.FilterDeprecatedAPIView) (*matrix.DeprecatedAPIViewList, error) {
	where := bobmodel.SelectWhere.DeprecatedApis
	queryMods := []bob.Mod[*dialect.SelectQuery]{}

	if f := filter.ClusterID; f != nil {
		queryMods = append(queryMods, where.ClusterID.EQ(uuid.FromStringOrNil(*f)))
	}

	if f := filter.OrganizationID; f != nil {
		queryMods = append(queryMods, where.OrganizationID.EQ(uuid.FromStringOrNil(*f)))
	}

	bs, err := bobmodel.DeprecatedApis.Query(queryMods...).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}

	var ds = make([]matrix.DeprecatedAPIView, 0, len(bs))
	for _, b := range bs {
		where := bobmodel.SelectWhere.Objects
		queryMods := []bob.Mod[*dialect.SelectQuery]{
			where.ClusterID.EQ(b.ClusterID),
			where.Kind.EQ(b.Kind),
		}

		bObjects, err := bobmodel.Objects.Query(queryMods...).All(ctx, r.db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors.Join(ekg.ErrNotFound, err)
			}
			return nil, err
		}
		if bObjects == nil || len(bObjects) == 0 {
			continue
		}

		objects := make([]matrix.AffectedObjectView, 0, len(bObjects))
		for _, bo := range bObjects {
			objects = append(
				objects,
				matrix.AffectedObjectView{
					ID:           bo.ID.String(),
					Name:         bo.Name,
					Kind:         bo.Kind,
					Namespace:    bo.Namespace,
					HealthStatus: bo.HealthStatus,
				},
			)
		}

		d, err := r.b2d(b, objects)
		if err != nil {
			return nil, err
		}
		ds = append(ds, *d)
	}
	return &matrix.DeprecatedAPIViewList{
		Items: ds,
	}, nil
}

func (r *deprecatedAPIViewRepo) b2d(
	b *bobmodel.DeprecatedAPI,
	affectedObjects []matrix.AffectedObjectView,
) (*matrix.DeprecatedAPIView, error) {
	view := &matrix.DeprecatedAPIView{
		ClusterID:           b.ClusterID.String(),
		OrganizationID:      b.OrganizationID.String(),
		CurrentGroupVersion: b.CurrentGroupVersion,
		Kind:                b.Kind,
		Name:                b.Name,
		ClusterK8sVersion:   b.ClusterK8SVersion,
		Deprecated:          b.Deprecated,
		DeprecatedIn:        b.DeprecatedIn,
		RemovedIn:           b.RemovedIn,
		ReplacementVersion:  b.ReplacementVersion,
		AffectedObjects:     affectedObjects,
	}
	return view, nil
}
