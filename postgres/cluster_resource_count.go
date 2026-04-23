package postgres

import (
	"context"

	"github.com/stephenafamo/bob"
)

func refreshClusterResourceCount(ctx context.Context, db bob.Executor) error {
	sqlStr := "REFRESH MATERIALIZED VIEW cluster_resource_count"
	_, err := db.ExecContext(ctx, sqlStr)
	if err != nil {
		return err
	}
	return nil
}
