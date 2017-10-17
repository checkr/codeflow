package resolvers

import (
	"context"
	"fmt"

	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	graphql "github.com/neelance/graphql-go"
)

func (r *Resolver) Release(ctx context.Context, args *struct{ ID graphql.ID }) *ReleaseResolver {
	release := models.Release{}
	return &ReleaseResolver{db: r.db, Release: release}
}

type ReleaseResolver struct {
	db      *gorm.DB
	Release models.Release
}

type ReleaseInput struct {
	ID        *string
	FeatureId *string
}

func (r *Resolver) CreateRelease(args *struct{ Release *ReleaseInput }) (*ReleaseResolver, error) {
	fmt.Println("CreateRelease")
	// var release models.Release

	spew.Dump(*args.Release)
	return nil, nil
}

func (r *ReleaseResolver) ID() graphql.ID {
	return graphql.ID(r.Release.Model.ID.String())
}

func (r *ReleaseResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	var project models.Project

	r.db.Model(r.Release).Related(&project)

	return &ProjectResolver{db: r.db, Project: project}, nil
}

func (r *ReleaseResolver) User(ctx context.Context) (*UserResolver, error) {
	var user models.User

	r.db.Model(r.User).Related(&user)

	return &UserResolver{db: r.db, User: user}, nil
}

func (r *ReleaseResolver) HeadFeature() (*FeatureResolver, error) {
	var feature models.Feature

	r.db.Where("id = ?", r.Release.HeadFeatureId).First(&feature)

	return &FeatureResolver{db: r.db, Feature: feature}, nil
}

func (r *ReleaseResolver) TailFeature() (*FeatureResolver, error) {
	var feature models.Feature

	r.db.Where("id = ?", r.Release.TailFeatureId).First(&feature)

	return &FeatureResolver{db: r.db, Feature: feature}, nil
}

func (r *ReleaseResolver) State() string {
	return string(r.Release.State)
}

func (r *ReleaseResolver) StateMessage() string {
	return r.Release.StateMessage
}

func (r *ReleaseResolver) Created() graphql.Time {
	return graphql.Time{Time: r.Release.Created}
}
