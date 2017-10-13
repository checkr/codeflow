package resolvers

import (
	"context"

	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/codeamp/circuit/plugins/codeamp/utils"
)

func (r *Resolver) Features(ctx context.Context) ([]*FeatureResolver, error) {
	if _, err := utils.CheckAuth(ctx, []string{}); err != nil {
		return nil, err
	}

	var rows []models.Feature
	var results []*FeatureResolver

	r.db.Order("created desc").Find(&rows)
	for _, feature := range rows {
		results = append(results, &FeatureResolver{db: r.db, Feature: feature})
	}

	return results, nil
}
