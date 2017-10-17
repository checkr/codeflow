package resolvers

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	log "github.com/codeamp/logger"

	"github.com/codeamp/circuit/plugins"
	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/codeamp/transistor"
	"github.com/extemporalgenome/slug"
	"github.com/jinzhu/gorm"
	graphql "github.com/neelance/graphql-go"
	"golang.org/x/crypto/ssh"
)

type ProjectInput struct {
	ID          *string
	GitProtocol string
	GitUrl      string
	Bookmarked  *bool
}

func (r *Resolver) Project(ctx context.Context, args *struct {
	ID   *graphql.ID
	Slug *string
	Name *string
}) (*ProjectResolver, error) {
	var project models.Project
	var query *gorm.DB

	if args.ID != nil {
		query = r.db.Where("id = ?", *args.ID)
	} else if args.Slug != nil {
		query = r.db.Where("slug = ?", *args.Slug)
	} else if args.Name != nil {
		query = r.db.Where("name = ?", *args.Name)
	} else {
		return nil, fmt.Errorf("Missing argument id or slug")
	}

	if err := query.First(&project).Error; err != nil {
		return nil, err
	}

	return &ProjectResolver{db: r.db, Project: project}, nil
}

func (r *Resolver) UpdateProject(args *struct{ Project *ProjectInput }) (*ProjectResolver, error) {

	var project models.Project

	if args.Project.ID == nil {
		return nil, fmt.Errorf("Missing argument id")
	}

	if r.db.Where("id = ?", args.Project.ID).First(&project).RecordNotFound() {
		log.InfoWithFields("Project not found", log.Fields{
			"id": args.Project.ID,
		})
		return nil, fmt.Errorf("Project not found.")
	}

	protocol := "HTTPS"
	switch args.Project.GitProtocol {
	case "private", "PRIVATE", "ssh", "SSH":
		protocol = "SSH"
	case "public", "PUBLIC", "https", "HTTPS":
		protocol = "HTTPS"
	}

	res := plugins.GetRegexParams("(?P<host>(git@|https?:\\/\\/)([\\w\\.@]+)(\\/|:))(?P<owner>[\\w,\\-,\\_]+)\\/(?P<repo>[\\w,\\-,\\_]+)(.git){0,1}((\\/){0,1})", args.Project.GitUrl)
	repository := fmt.Sprintf("%s/%s", res["owner"], res["repo"])

	project.GitUrl = args.Project.GitUrl

	// Check if project already exists with same name
	if r.db.Unscoped().Where("id != ? and repository = ?", args.Project.ID, repository).First(&models.Project{}).RecordNotFound() == false {
		return nil, fmt.Errorf("Project with repository name already exists.")
	}

	project.GitUrl = args.Project.GitUrl
	project.GitProtocol = protocol
	project.Repository = repository
	project.Name = repository
	project.Slug = slug.Slug(repository)
	r.db.Save(project)

	// Cascade delete all features and releases related to old git url
	r.db.Where("projectId = ?", project.ID).Delete(models.Feature{})
	r.db.Where("projectId = ?", project.ID).Delete(models.Release{})
	return &ProjectResolver{db: r.db, Project: project}, nil
}

func (r *Resolver) CreateProject(args *struct{ Project *ProjectInput }) (*ProjectResolver, error) {
	protocol := "HTTPS"
	switch args.Project.GitProtocol {
	case "private", "PRIVATE", "ssh", "SSH":
		protocol = "SSH"
	case "public", "PUBLIC", "https", "HTTPS":
		protocol = "HTTPS"
	}

	project := models.Project{
		GitProtocol: protocol,
		GitUrl:      args.Project.GitUrl,
		Secret:      transistor.RandomString(30),
	}

	res := plugins.GetRegexParams("(?P<host>(git@|https?:\\/\\/)([\\w\\.@]+)(\\/|:))(?P<owner>[\\w,\\-,\\_]+)\\/(?P<repo>[\\w,\\-,\\_]+)(.git){0,1}((\\/){0,1})", args.Project.GitUrl)
	repository := fmt.Sprintf("%s/%s", res["owner"], res["repo"])

	// Check if project already exists with same name
	existingProject := models.Project{}

	if r.db.Unscoped().Where("repository = ?", repository).First(&existingProject).RecordNotFound() {
		log.InfoWithFields("Project not found", log.Fields{
			"repository": repository,
		})
	} else {
		//return nil, fmt.Errorf("This repository already exists. Try again with a different git url.")
	}

	project.Name = repository
	project.Repository = repository
	project.Slug = slug.Slug(repository)

	deletedProject := models.Project{}
	if err := r.db.Unscoped().Where("repository = ?", repository).First(&deletedProject).Error; err != nil {
		project.Model.ID = deletedProject.Model.ID
	}

	// priv *rsa.PrivateKey;
	priv, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		return nil, err
	}

	err = priv.Validate()
	if err != nil {
		return nil, err
	}

	// Get der format. priv_der []byte
	priv_der := x509.MarshalPKCS1PrivateKey(priv)

	// pem.Block
	priv_blk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   priv_der,
	}

	// Public Key generation
	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, err
	}

	project.RsaPrivateKey = string(pem.EncodeToMemory(&priv_blk))
	project.RsaPublicKey = string(ssh.MarshalAuthorizedKey(pub))

	r.db.Create(&project)

	r.actions.ProjectCreated(&project)

	return &ProjectResolver{db: r.db, Project: project}, nil
}

type ProjectResolver struct {
	db      *gorm.DB
	Project models.Project
}

func (r *ProjectResolver) ID() graphql.ID {
	return graphql.ID(r.Project.Model.ID.String())
}

func (r *ProjectResolver) Name() string {
	return r.Project.Name
}

func (r *ProjectResolver) Slug() string {
	return r.Project.Slug
}

func (r *ProjectResolver) Repository() string {
	return r.Project.Repository
}

func (r *ProjectResolver) Secret() string {
	return r.Project.Secret
}

func (r *ProjectResolver) GitUrl() string {
	return r.Project.GitUrl
}

func (r *ProjectResolver) GitProtocol() string {
	return r.Project.GitProtocol
}

func (r *ProjectResolver) RsaPrivateKey() string {
	return r.Project.RsaPrivateKey
}

func (r *ProjectResolver) RsaPublicKey() string {
	return r.Project.RsaPublicKey
}

func (r *ProjectResolver) Features(ctx context.Context) ([]*FeatureResolver, error) {
	var rows []models.Feature
	var results []*FeatureResolver

	r.db.Where("project_id = ?", r.Project.ID).Order("created desc").Find(&rows)

	for _, feature := range rows {
		results = append(results, &FeatureResolver{db: r.db, Feature: feature})
	}

	return results, nil
}

func (r *ProjectResolver) Services(ctx context.Context) ([]*ServiceResolver, error) {
	var rows []models.Service
	var results []*ServiceResolver

	r.db.Where("project_id = ?", r.Project.ID).Find(&rows)

	for _, service := range rows {
		results = append(results, &ServiceResolver{db: r.db, Service: service})
	}

	return results, nil
}

func (r *ProjectResolver) Releases(ctx context.Context) ([]*ReleaseResolver, error) {
	var rows []models.Release
	var results []*ReleaseResolver

	r.db.Model(r.Project).Related(&rows)

	for _, release := range rows {
		results = append(results, &ReleaseResolver{db: r.db, Release: release})
	}

	return results, nil
}

func (r *ProjectResolver) EnvironmentVariables(ctx context.Context) ([]*EnvironmentVariableResolver, error) {
	var rows []models.EnvironmentVariable
	var results []*EnvironmentVariableResolver

	r.db.Select("key, version, id, value, created, type, user_id, project_id, deleted_at").Where("project_id = ?", r.Project.ID).Order("key, version, created desc").Find(&rows)

	for _, envVar := range rows {
		results = append(results, &EnvironmentVariableResolver{db: r.db, EnvironmentVariable: envVar})
	}

	return results, nil
}
