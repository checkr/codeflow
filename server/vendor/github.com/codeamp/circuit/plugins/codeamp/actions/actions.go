package actions

import (
	"fmt"

	"github.com/codeamp/circuit/plugins"
	"github.com/codeamp/circuit/plugins/codeamp/models"
	log "github.com/codeamp/logger"
	"github.com/codeamp/transistor"
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
)

type Actions struct {
	events chan transistor.Event
	db     *gorm.DB
}

func NewActions(events chan transistor.Event, db *gorm.DB) *Actions {
	return &Actions{
		events: events,
		db:     db,
	}
}

func (x *Actions) HeartBeat(tick string) {
	var projects []models.Project

	x.db.Find(&projects)
	for _, project := range projects {
		if tick == "minute" {
			x.GitSync(&project)
		}
	}
}

func (x *Actions) GitSync(project *models.Project) {
	var feature models.Feature
	var release models.Release
	hash := ""

	// Get latest release and deployed feature hash
	if x.db.Where("project_id = ?", project.ID).Order("created_at DESC").First(&release).RecordNotFound() {
		// get latest feature if there is no releases
		x.db.Where("project_id = ?", project.ID).Order("created DESC").First(&feature)
		hash = feature.Hash
	} else {
		hash = release.HeadFeature.Hash
	}

	gitSync := plugins.GitSync{
		Action: plugins.Update,
		State:  plugins.Waiting,
		Project: plugins.Project{
			Slug:       project.Slug,
			Repository: project.Repository,
		},
		Git: plugins.Git{
			Url:           project.GitUrl,
			Protocol:      project.GitProtocol,
			Branch:        "master",
			RsaPrivateKey: project.RsaPrivateKey,
			RsaPublicKey:  project.RsaPublicKey,
		},
		From: hash,
	}

	x.events <- transistor.NewEvent(gitSync, nil)
}

func (x *Actions) GitCommit(commit plugins.GitCommit) {
	project := models.Project{}
	feature := models.Feature{}

	if x.db.Where("repository = ?", commit.Repository).First(&project).RecordNotFound() {
		log.InfoWithFields("project not found", log.Fields{
			"repository": commit.Repository,
		})
		return
	}

	if x.db.Where("project_id = ? AND hash = ?", project.ID, commit.Hash).First(&feature).RecordNotFound() {
		feature = models.Feature{
			ProjectId:  project.ID,
			Message:    commit.Message,
			User:       commit.User,
			Hash:       commit.Hash,
			ParentHash: commit.ParentHash,
			Ref:        commit.Ref,
			Created:    commit.Created,
		}
		x.db.Save(&feature)

		wsMsg := plugins.WebsocketMsg{
			Event:   fmt.Sprintf("projects/%s/features", project.Slug),
			Payload: feature,
		}
		x.events <- transistor.NewEvent(wsMsg, nil)
	} else {
		log.InfoWithFields("feature already exists", log.Fields{
			"repository": commit.Repository,
			"hash":       commit.Hash,
		})
	}
}

func (x *Actions) ProjectCreated(project *models.Project) {
	wsMsg := plugins.WebsocketMsg{
		Event:   "projects",
		Payload: project,
	}
	x.events <- transistor.NewEvent(wsMsg, nil)
}

func (x *Actions) ServiceCreated(service *models.Service) {
	project := models.Project{}
	if x.db.Where("id = ?", service.ProjectId).First(&project).RecordNotFound() {
		log.InfoWithFields("project not found", log.Fields{
			"service": service,
		})
	}

	wsMsg := plugins.WebsocketMsg{
		Event:   fmt.Sprintf("projects/%s/services/new", project.Slug),
		Payload: service,
	}
	x.events <- transistor.NewEvent(wsMsg, nil)
}

func (x *Actions) ServiceUpdated(service *models.Service) {
	project := models.Project{}
	if x.db.Where("id = ?", service.ProjectId).First(&project).RecordNotFound() {
		log.InfoWithFields("project not found", log.Fields{
			"service": service,
		})
	}

	wsMsg := plugins.WebsocketMsg{
		Event:   fmt.Sprintf("projects/%s/services/updated", project.Slug),
		Payload: service,
	}
	x.events <- transistor.NewEvent(wsMsg, nil)
}

func (x *Actions) ServiceDeleted(service *models.Service) {
	project := models.Project{}
	if x.db.Where("id = ?", service.ProjectId).First(&project).RecordNotFound() {
		log.InfoWithFields("project not found", log.Fields{
			"service": service,
		})
	}

	wsMsg := plugins.WebsocketMsg{
		Event:   fmt.Sprintf("projects/%s/services/deleted", project.Slug),
		Payload: service,
	}
	x.events <- transistor.NewEvent(wsMsg, nil)
}

func (x *Actions) ServiceSpecCreated(service *models.ServiceSpec) {
	wsMsg := plugins.WebsocketMsg{
		Event:   fmt.Sprintf("serviceSpecs/new"),
		Payload: service,
	}
	x.events <- transistor.NewEvent(wsMsg, nil)
}

func (x *Actions) ServiceSpecDeleted(service *models.ServiceSpec) {
	wsMsg := plugins.WebsocketMsg{
		Event:   fmt.Sprintf("serviceSpecs/deleted"),
		Payload: service,
	}
	x.events <- transistor.NewEvent(wsMsg, nil)
}

func (x *Actions) ServiceSpecUpdated(service *models.ServiceSpec) {
	wsMsg := plugins.WebsocketMsg{
		Event:   fmt.Sprintf("serviceSpecs/updated"),
		Payload: service,
	}
	x.events <- transistor.NewEvent(wsMsg, nil)
}

func (x *Actions) EnvironmentVariableCreated(envVar *models.EnvironmentVariable) {
	project := models.Project{}
	if x.db.Where("id = ?", envVar.ProjectId).First(&project).RecordNotFound() {
		log.InfoWithFields("project not found", log.Fields{
			"service": envVar,
		})
	}

	spew.Dump(fmt.Sprintf("projects/%s/environmentVariables/created", project.Slug))

	wsMsg := plugins.WebsocketMsg{
		Event:   fmt.Sprintf("projects/%s/environmentVariables/created", project.Slug),
		Payload: envVar,
	}
	x.events <- transistor.NewEvent(wsMsg, nil)
}

func (x *Actions) EnvironmentVariableDeleted(envVar *models.EnvironmentVariable) {
	project := models.Project{}
	if x.db.Where("id = ?", envVar.ProjectId).First(&project).RecordNotFound() {
		log.InfoWithFields("envvar not found", log.Fields{
			"service": envVar,
		})
	}

	spew.Dump(fmt.Sprintf("projects/%s/environmentVariables/deleted", project.Slug))

	wsMsg := plugins.WebsocketMsg{
		Event:   fmt.Sprintf("projects/%s/environmentVariables/deleted", project.Slug),
		Payload: envVar,
	}
	x.events <- transistor.NewEvent(wsMsg, nil)
}

func (x *Actions) EnvironmentVariableUpdated(envVar *models.EnvironmentVariable) {
	project := models.Project{}
	if x.db.Where("id = ?", envVar.ProjectId).First(&project).RecordNotFound() {
		log.InfoWithFields("envvar not found", log.Fields{
			"envVar": envVar,
		})
	}

	spew.Dump(fmt.Sprintf("projects/%s/environmentVariables/updated", project.Slug))

	wsMsg := plugins.WebsocketMsg{
		Event:   fmt.Sprintf("projects/%s/environmentVariables/updated", project.Slug),
		Payload: envVar,
	}

	x.events <- transistor.NewEvent(wsMsg, nil)
}
