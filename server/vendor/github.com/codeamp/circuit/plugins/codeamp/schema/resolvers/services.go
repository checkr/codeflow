package resolvers

import (
	"context"

	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/codeamp/circuit/plugins/codeamp/utils"
)

func (r *Resolver) Services(ctx context.Context) ([]*ServiceResolver, error) {
	if _, err := utils.CheckAuth(ctx, []string{}); err != nil {
		return nil, err
	}

	var rows []models.Service
	var results []*ServiceResolver

	r.db.Order("created desc").Find(&rows)
	for _, service := range rows {
		results = append(results, &ServiceResolver{db: r.db, Service: service})
	}

	return results, nil
}
