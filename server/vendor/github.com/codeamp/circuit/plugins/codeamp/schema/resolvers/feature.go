package resolvers

import (
	"context"

	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/jinzhu/gorm"
	graphql "github.com/neelance/graphql-go"
)

func (r *Resolver) Feature(ctx context.Context, args *struct{ ID graphql.ID }) *FeatureResolver {
	feature := models.Feature{}
	return &FeatureResolver{db: r.db, Feature: feature}
}

type FeatureResolver struct {
	db      *gorm.DB
	Feature models.Feature
}

func (r *FeatureResolver) ID() graphql.ID {
	return graphql.ID(r.Feature.Model.ID.String())
}

func (r *FeatureResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	var project models.Project

	r.db.Model(r.Feature).Related(&project)

	return &ProjectResolver{db: r.db, Project: project}, nil
}

func (r *FeatureResolver) Message() string {
	return r.Feature.Message
}

func (r *FeatureResolver) User() string {
	return r.Feature.User
}

func (r *FeatureResolver) Hash() string {
	return r.Feature.Hash
}

func (r *FeatureResolver) ParentHash() string {
	return r.Feature.ParentHash
}

func (r *FeatureResolver) Ref() string {
	return r.Feature.Ref
}

func (r *FeatureResolver) Created() graphql.Time {
	return graphql.Time{Time: r.Feature.Created}
}
