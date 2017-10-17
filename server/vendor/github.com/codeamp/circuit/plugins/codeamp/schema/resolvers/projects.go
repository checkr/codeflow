package resolvers

import (
	"context"

	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/codeamp/circuit/plugins/codeamp/utils"
)

func (r *Resolver) Projects(ctx context.Context) ([]*ProjectResolver, error) {
	if _, err := utils.CheckAuth(ctx, []string{}); err != nil {
		return nil, err
	}

	var rows []models.Project
	var results []*ProjectResolver

	r.db.Find(&rows)
	for _, project := range rows {
		results = append(results, &ProjectResolver{db: r.db, Project: project})
	}

	return results, nil
}
