package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/gofrs/uuid/v5"
	"github.com/kubegrade/ekg"
	"github.com/kubegrade/matrix"
	"github.com/kubegrade/matrix/postgres/bobmodel"
	"github.com/stephenafamo/bob"
)

type repositoryRepo struct {
	db bob.DB
}

func NewRepositoryRepo(db *sql.DB) matrix.RepositoryRepo {
	return &repositoryRepo{
		db: bob.NewDB(db),
	}
}

func (r *repositoryRepo) Create(ctx context.Context, repository *matrix.Repository) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	if _, err := bobmodel.Repositories.Insert(r.d2b(repository)).One(ctx, tx); err != nil {
		return err
	}

	for _, id := range repository.PersonalAccessTokenIDs {
		if _, err := bobmodel.RepositoryPersonalAccessTokens.Insert(
			&bobmodel.RepositoryPersonalAccessTokenSetter{
				RepositoryID:          omit.From(uuid.FromStringOrNil(repository.ID)),
				PersonalAccessTokenID: omit.From(uuid.FromStringOrNil(id)),
				CreatedAt:             omit.From(time.Now()),
			},
		).One(ctx, tx); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *repositoryRepo) Update(ctx context.Context, update matrix.RepositoryUpdate) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	id := uuid.FromStringOrNil(update.ID)
	b, err := bobmodel.Repositories.Query(
		bobmodel.SelectWhere.Repositories.ID.EQ(id),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	setter := &bobmodel.RepositorySetter{
		ID: omit.From(id),
	}

	if name := update.Payload.Name; name != nil {
		setter.Name = omit.From(*name)
	}

	if description := update.Payload.Description; description != nil {
		setter.Description = omitnull.From(*description)
	}

	if namespace := update.Payload.Namespace; namespace != nil {
		setter.Name = omit.From(*namespace)
	}

	if err := b.Update(ctx, tx, setter); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *repositoryRepo) Find(ctx context.Context, id string) (*matrix.Repository, error) {
	where := bobmodel.SelectWhere.Repositories
	b, err := bobmodel.Repositories.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(ekg.ErrNotFound, err)
		}
		return nil, err
	}
	return r.b2d(b), nil
}

func (r *repositoryRepo) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	where := bobmodel.SelectWhere.Repositories
	b, err := bobmodel.Repositories.Query(
		where.ID.EQ(uuid.FromStringOrNil(id)),
		bobmodel.SelectThenLoad.Repository.RepositoryPersonalAccessTokens(),
	).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ekg.ErrNotFound, err)
		}
		return err
	}

	if err := b.R.RepositoryPersonalAccessTokens.DeleteAll(ctx, tx); err != nil {
		return err
	}
	if err := b.Delete(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AddPersonalAccessTokens
func (r *repositoryRepo) AddPersonalAccessTokens(ctx context.Context, repositoryID matrix.RepositoryID, patIDs ...matrix.PersonalAccessTokenID) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	for _, id := range patIDs {
		if _, err := bobmodel.RepositoryPersonalAccessTokens.Insert(
			&bobmodel.RepositoryPersonalAccessTokenSetter{
				RepositoryID:          omit.From(uuid.FromStringOrNil(repositoryID)),
				PersonalAccessTokenID: omit.From(uuid.FromStringOrNil(id)),
				CreatedAt:             omit.From(time.Now()),
			},
		).One(ctx, tx); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// RemovePersonalAccessTokens
func (r *repositoryRepo) RemovePersonalAccessTokens(ctx context.Context, repositoryID matrix.RepositoryID, patIDs ...matrix.PersonalAccessTokenID) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)

	patUUIDs := make([]uuid.UUID, 0, len(patIDs))
	for _, id := range patIDs {
		patUUIDs = append(patUUIDs, uuid.FromStringOrNil(id))
	}

	where := bobmodel.DeleteWhere.RepositoryPersonalAccessTokens
	if _, err := bobmodel.RepositoryPersonalAccessTokens.Delete(
		where.RepositoryID.EQ(uuid.FromStringOrNil(repositoryID)),
		where.PersonalAccessTokenID.In(patUUIDs...),
	).All(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *repositoryRepo) d2b(d *matrix.Repository) *bobmodel.RepositorySetter {
	return &bobmodel.RepositorySetter{
		ID:                  omit.From(uuid.FromStringOrNil(d.ID)),
		OrganizationID:      omit.From(uuid.FromStringOrNil(d.OrganizationID)),
		Name:                omit.From(d.Name),
		Provider:            omit.From(d.Provider),
		Namespace: omit.From(d.Namespace),
		Description:         omitnull.From(d.Description),
		CreatedBy:           omitnull.From(uuid.FromStringOrNil(d.CreatedBy)),
		UpdatedBy:           omitnull.From(uuid.FromStringOrNil(d.UpdatedBy)),
		CreatedAt:           omit.From(d.CreatedAt),
		UpdatedAt:           omit.From(d.UpdatedAt),
	}
}

func (r *repositoryRepo) b2d(b *bobmodel.Repository) *matrix.Repository {
	return &matrix.Repository{
		ID:                  b.ID.String(),
		OrganizationID:      b.OrganizationID.String(),
		Name:                b.Name,
		Provider:            b.Provider,
		Namespace: b.Namespace,
		Description:         b.Description.GetOrZero(),
		CreatedBy:           b.CreatedBy.GetOrZero().String(),
		UpdatedBy:           b.UpdatedBy.GetOrZero().String(),
		CreatedAt:           b.CreatedAt,
		UpdatedAt:           b.UpdatedAt,
	}
}
