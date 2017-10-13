package resolvers

import (
	"context"

	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/codeamp/circuit/plugins/codeamp/utils"
)

func (r *Resolver) Users(ctx context.Context) ([]*UserResolver, error) {
	if _, err := utils.CheckAuth(ctx, []string{"admin"}); err != nil {
		return nil, err
	}

	var rows []models.User
	var results []*UserResolver

	r.db.Find(&rows)

	for _, user := range rows {
		results = append(results, &UserResolver{db: r.db, User: user})
	}

	return results, nil
}
