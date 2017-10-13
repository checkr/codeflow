package resolvers

import (
	"context"

	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/codeamp/circuit/plugins/codeamp/utils"
)

func (r *Resolver) ServiceSpecs(ctx context.Context) ([]*ServiceSpecResolver, error) {
	if _, err := utils.CheckAuth(ctx, []string{}); err != nil {
		return nil, err
	}

	var rows []models.ServiceSpec
	var results []*ServiceSpecResolver

	r.db.Order("created desc").Find(&rows)
	for _, serviceSpec := range rows {
		results = append(results, &ServiceSpecResolver{db: r.db, ServiceSpec: serviceSpec})
	}

	return results, nil
}
