package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/codeamp/circuit/plugins/codeamp/utils"

	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	graphql "github.com/neelance/graphql-go"
	uuid "github.com/satori/go.uuid"
)

type EnvironmentVariableInput struct {
	ID        *string
	Key       string
	Value     string
	Type      *string
	ProjectId *string
}

type EnvironmentVariableResolver struct {
	db                  *gorm.DB
	EnvironmentVariable models.EnvironmentVariable
}

func (r *Resolver) EnvironmentVariable(ctx context.Context, args *struct{ ID graphql.ID }) (*EnvironmentVariableResolver, error) {
	envVar := models.EnvironmentVariable{}
	if err := r.db.Where("id = ?", args.ID).First(&envVar).Error; err != nil {
		return nil, err
	}

	return &EnvironmentVariableResolver{db: r.db, EnvironmentVariable: envVar}, nil
}

func (r *Resolver) CreateEnvironmentVariable(ctx context.Context, args *struct{ EnvironmentVariable *EnvironmentVariableInput }) (*EnvironmentVariableResolver, error) {
	projectId, err := uuid.FromString(*args.EnvironmentVariable.ProjectId)
	if err != nil {
		return &EnvironmentVariableResolver{}, err
	}

	userIdString, err := utils.CheckAuth(ctx, []string{})
	if err != nil {
		return &EnvironmentVariableResolver{}, err
	}

	userId, err := uuid.FromString(userIdString)
	if err != nil {
		return &EnvironmentVariableResolver{}, err
	}

	var existingEnvVar models.EnvironmentVariable
	if r.db.Where("key = ? and project_id = ? and deleted_at is null", args.EnvironmentVariable.Key, args.EnvironmentVariable.ProjectId).Find(&existingEnvVar).RecordNotFound() {
		spew.Dump(args.EnvironmentVariable)
		envVar := models.EnvironmentVariable{
			Key:       args.EnvironmentVariable.Key,
			Value:     args.EnvironmentVariable.Value,
			ProjectId: projectId,
			Version:   int32(0),
			Type:      *args.EnvironmentVariable.Type,
			UserId:    userId,
			Created:   time.Now(),
		}

		r.db.Create(&envVar)

		r.actions.EnvironmentVariableCreated(&envVar)

		return &EnvironmentVariableResolver{db: r.db, EnvironmentVariable: envVar}, nil
	} else {
		return nil, fmt.Errorf("CreateEnvironmentVariable: key already exists")
	}
}

func (r *Resolver) UpdateEnvironmentVariable(ctx context.Context, args *struct{ EnvironmentVariable *EnvironmentVariableInput }) (*EnvironmentVariableResolver, error) {

	var existingEnvVar models.EnvironmentVariable
	if r.db.Where("id = ?", args.EnvironmentVariable.ID).Find(&existingEnvVar).RecordNotFound() {
		return nil, fmt.Errorf("UpdateEnvironmentVariable: key doesn't exist.")
	} else {
		envVar := models.EnvironmentVariable{
			Key:       args.EnvironmentVariable.Key,
			Value:     args.EnvironmentVariable.Value,
			ProjectId: existingEnvVar.ProjectId,
			Version:   existingEnvVar.Version + int32(1),
			Type:      existingEnvVar.Type,
			UserId:    existingEnvVar.UserId,
			Created:   time.Now(),
		}
		r.db.Delete(&existingEnvVar)
		r.db.Create(&envVar)
		r.actions.EnvironmentVariableUpdated(&envVar)

		return &EnvironmentVariableResolver{db: r.db, EnvironmentVariable: envVar}, nil
	}
}

func (r *Resolver) DeleteEnvironmentVariable(ctx context.Context, args *struct{ EnvironmentVariable *EnvironmentVariableInput }) (*EnvironmentVariableResolver, error) {

	var existingEnvVar models.EnvironmentVariable
	if r.db.Where("id = ?", args.EnvironmentVariable.ID).Find(&existingEnvVar).RecordNotFound() {
		return nil, fmt.Errorf("UpdateEnvironmentVariable: key doesn't exist.")
	} else {
		var rows []models.EnvironmentVariable

		r.db.Where("project_id = ? and key = ?", existingEnvVar.ProjectId, existingEnvVar.Key).Find(&rows)
		for _, envVar := range rows {
			r.db.Delete(&envVar)
		}
		r.actions.EnvironmentVariableDeleted(&existingEnvVar)

		return &EnvironmentVariableResolver{db: r.db, EnvironmentVariable: existingEnvVar}, nil
	}
}

func (r *EnvironmentVariableResolver) ID() graphql.ID {
	return graphql.ID(r.EnvironmentVariable.Model.ID.String())
}

func (r *EnvironmentVariableResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	var project models.Project
	r.db.Model(r.EnvironmentVariable).Related(&project)
	return &ProjectResolver{db: r.db, Project: project}, nil
}

func (r *EnvironmentVariableResolver) Key() string {
	return r.EnvironmentVariable.Key
}

func (r *EnvironmentVariableResolver) Value() string {
	return r.EnvironmentVariable.Value
}

func (r *EnvironmentVariableResolver) Version() int32 {
	return r.EnvironmentVariable.Version
}

func (r *EnvironmentVariableResolver) Type() string {
	return r.EnvironmentVariable.Type
}

func (r *EnvironmentVariableResolver) Created() graphql.Time {
	return graphql.Time{Time: r.EnvironmentVariable.Created}
}

func (r *EnvironmentVariableResolver) User() (*UserResolver, error) {
	var user models.User
	r.db.Model(r.EnvironmentVariable).Related(&user)
	return &UserResolver{db: r.db, User: user}, nil
}

func (r *EnvironmentVariableResolver) Versions(ctx context.Context) ([]*EnvironmentVariableResolver, error) {
	spew.Dump("VERsIONs!")
	if _, err := utils.CheckAuth(ctx, []string{}); err != nil {
		return nil, err
	}
	spew.Dump("MADE IT!")
	var rows []models.EnvironmentVariable
	var results []*EnvironmentVariableResolver

	r.db.Unscoped().Where("project_id = ? and key = ?", r.EnvironmentVariable.ProjectId, r.EnvironmentVariable.Key).Order("version desc").Find(&rows)

	spew.Dump(rows)
	for _, envVar := range rows {
		results = append(results, &EnvironmentVariableResolver{db: r.db, EnvironmentVariable: envVar})
	}

	return results, nil
}
