package codeflow

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/checkr/codeflow/server/plugins"
	"github.com/extemporalgenome/slug"
	"github.com/maxwellhealth/bongo"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	bongo.DocumentBase `bson:",inline"`
	Name               string `bson:"name" json:"name"`
	Username           string `bson:"username" json:"username"`
	Email              string `bson:"email" json:"email"`
}

type Project struct {
	bongo.DocumentBase    `bson:",inline"`
	Name                  string   `bson:"name" json:"name"`
	Slug                  string   `bson:"slug" json:"slug"`
	Repository            string   `bson:"repository" json:"repository"`
	Secret                string   `bson:"secret" json:"secret"`
	Pinged                bool     `bson:"pinged" json:"pinged"`
	GitSshUrl             string   `bson:"gitSshUrl" json:"gitSshUrl" validate:"required"`
	RsaPrivateKey         string   `bson:"rsaPrivateKey" json:"-"`
	RsaPublicKey          string   `bson:"rsaPublicKey" json:"rsaPublicKey"`
	Bokmarked             bool     `bson:"-" json:"bookmarked"`
	ContinuousIntegration bool     `bson:"continuousIntegration" json:"continuousIntegration"`
	ContinuousDelivery    bool     `bson:"continuousDelivery" json:"continuousDelivery"`
	Workflows             []string `bson:"workflows" json:"workflows"`
	LogsUrl               string   `bson:"-" json:"logsUrl"`
}

func (p *Project) AfterFind(*bongo.Collection) error {
	p.LogsUrl = strings.Replace(viper.GetString("plugins.codeflow.logs_url"), "##PROJECT-NAMESPACE##", fmt.Sprintf("production-%v", p.Slug), -1)
	return nil
}

func (p *Project) BeforeSave(collection *bongo.Collection) error {
	regex, _ := regexp.Compile("(?:git|ssh|git@[\\w\\.]+):((?:\\/\\/)?[\\w\\.@:\\/~_-]+)\\.git(?:\\/?|\\#[\\d\\w\\.\\-_]+?)$")
	repository := regex.FindStringSubmatch(p.GitSshUrl)[1]
	p.Name = repository
	p.Repository = repository
	p.Slug = slug.Slug(repository)

	return nil
}

func (p *Project) Validate(collection *bongo.Collection) []error {
	var err []error

	regex, _ := regexp.Compile("(?:git|ssh|git@[\\w\\.]+):((?:\\/\\/)?[\\w\\.@:\\/~_-]+)\\.git(?:\\/?|\\#[\\d\\w\\.\\-_]+?)$")
	if !regex.MatchString(p.GitSshUrl) {
		err = append(err, errors.New("Wrong Git SSH url"))
	}

	return err
}

type Bookmark struct {
	bongo.DocumentBase `bson:",inline"`
	ProjectId          bson.ObjectId `bson:"projectId" json:"projectId"`
	UserId             bson.ObjectId `bson:"userId" json:"-"`
	Name               string        `bson:"-" json:"name"`
	Slug               string        `bson:"-" json:"slug"`
}

func (b *Bookmark) AfterFind(*bongo.Collection) error {
	project := Project{}
	if err := db.Collection("projects").FindById(b.ProjectId, &project); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Projects::FindById::DocumentNotFoundError: _id: `%v`", b.ProjectId)
			db.Collection("bookmarks").DeleteDocument(b)
			return err
		} else {
			log.Printf("Projects::FindById::Error: %s", err.Error())
			return err
		}
	}

	b.Name = project.Name
	b.Slug = project.Slug

	return nil
}

type Service struct {
	bongo.DocumentBase `bson:",inline"`
	ProjectId          bson.ObjectId `bson:"projectId" json:"projectId"`
	State              plugins.State `bson:"state" json:"state"`
	Name               string        `bson:"name" json:"name"`
	Count              int           `bson:"count" json:"count"`
	Command            string        `bson:"command" json:"command"`
	Listeners          []Listener    `bson:"listeners" json:"listeners"`
}

type LoadBalancer struct {
	bongo.DocumentBase `bson:",inline"`
	Name               string         `bson:"name" json:"name"`
	ProjectId          bson.ObjectId  `bson:"projectId" json:"projectId"`
	ServiceId          bson.ObjectId  `bson:"serviceId" json:"serviceId"`
	Extension          string         `bson:"extension" json:"extension"`
	DNSName            string         `bson:"dnsName" json:"dnsName"`
	Type               string         `bson:"type" json:"type"`
	ListenerPairs      []ListenerPair `bson:"listenerPairs" json:"listenerPairs"`
	State              plugins.State  `bson:"state" json:"state"`
	StateMessage       string         `bson:"stateMessage" json:"stateMessage"`
}

type Feature struct {
	bongo.DocumentBase `bson:",inline"`
	ProjectId          bson.ObjectId `bson:"projectId" json:"projectId"`
	Message            string        `bson:"message" json:"message"`
	User               string        `bson:"user" json:"user"`
	Hash               string        `bson:"hash" json:"hash"`
	ParentHash         string        `bson:"parentHash" json:"parentHash"`
	ExternalLink       string        `bson:"externalLink" json:"externalLink"`
}

type Release struct {
	bongo.DocumentBase `bson:",inline"`
	ProjectId          bson.ObjectId `bson:"projectId" json:"projectId"`
	HeadFeatureId      bson.ObjectId `bson:"headFeatureId" json:"-"`
	HeadFeature        Feature       `bson:"-" json:"headFeature"`
	TailFeatureId      bson.ObjectId `bson:"tailFeatureId" json:"-"`
	TailFeature        Feature       `bson:"-" json:"tailFeature"`
	UserId             bson.ObjectId `bson:"userId" json:"-"`
	User               User          `bson:"-" json:"user"`
	State              plugins.State `bson:"state" json:"state"`
	StateMessage       string        `bson:"stateMessage" json:"stateMessage"`
	Secrets            []Secret      `bson:"secrets" json:"-"`
	Services           []Service     `bson:"services" json:"-"`
	Workflow           []Flow        `bson:"-" json:"workflow"`
}

func (r *Release) AfterFind(c *bongo.Collection) error {
	headFeature := Feature{}
	if err := db.Collection("features").FindById(r.HeadFeatureId, &headFeature); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Release::AfterFind::Features::FindById::DocumentNotFoundError: _id: `%v`", r.HeadFeatureId)
			return err
		} else {
			log.Printf("Release::AfterFind::Features::FindById::Error: `%v`", err.Error())
			return err
		}
	}
	r.HeadFeature = headFeature

	tailFeature := Feature{}
	if err := db.Collection("features").FindById(r.TailFeatureId, &tailFeature); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Release::AfterFind::Features::FindById::DocumentNotFoundError: _id: `%v`", r.TailFeatureId)
			return err
		} else {
			log.Printf("Release::AfterFind::Features::FindById::Error: `%v`", err.Error())
			return err
		}
	}
	r.TailFeature = tailFeature

	user := User{}
	if err := db.Collection("users").FindById(r.UserId, &user); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Release::AfterFind::Users::FindById::DocumentNotFoundError: _id: `%v`", r.UserId)
			return err
		} else {
			log.Printf("Release::AfterFind::Users::FindById::Error: `%v`", err.Error())
			return err
		}
	}
	r.User = user

	workflows := []Flow{}
	flow := Flow{}

	results := db.Collection("workflows").Find(bson.M{"releaseId": r.Id})
	for results.Next(&flow) {
		workflows = append(workflows, flow)
	}
	r.Workflow = workflows

	return nil
}

type Flow struct {
	bongo.DocumentBase `bson:",inline"`
	ReleaseId          bson.ObjectId `bson:"releaseId,omitempty" json:"releaseId"`
	Type               string        `bson:"type" json:"type"`
	Name               string        `bson:"name" json:"name"`
	Message            string        `bson:"message" json:"message"`
	State              plugins.State `bson:"state" json:"state"`
}

type ExternalFlowStatus struct {
	bongo.DocumentBase `bson:",inline"`
	ProjectId          bson.ObjectId `bson:"projectId" json:"projectId"`
	Hash               string        `bson:"hash" json:"hash"`
	Context            string        `bson:"context" json:"context"`
	Message            string        `bson:"message" json:"message"`
	State              plugins.State `bson:"state" json:"state"`
	OriginalState      string        `bson:"originalState" json:"originalState"`
}

type Secret struct {
	bongo.DocumentBase `bson:",inline"`
	ProjectId          bson.ObjectId `bson:"projectId" json:"-"`
	Key                string        `bson:"key" json:"key"`
	Value              string        `bson:"value" json:"value"`
	Type               plugins.Type  `bson:"type" json:"type"`
	Deleted            bool          `bson:"deleted" json:"deleted"`
}

type Build struct {
	bongo.DocumentBase `bson:",inline"`
	FeatureHash        string        `bson:"featureHash" json:"featureHash"`
	Type               string        `bson:"type" json:"type"`
	State              plugins.State `bson:"state" json:"state"`
	Image              string        `bson:"image" json:"image"`
	BuildLog           string        `bson:"buildLog" json:"buildLog"`
	BuildError         string        `bson:"buildError" json:"buildError"`
}

type ProjectSettings struct {
	ProjectId             bson.ObjectId `json:"projectId"`
	GitSshUrl             string        `json:"gitSshUrl"`
	Secrets               []Secret      `json:"secrets"`
	DeletedSecrets        []Secret      `json:"deletedSecrets"`
	ContinuousIntegration bool          `json:"continuousIntegration"`
	ContinuousDelivery    bool          `json:"continuousDelivery"`
}

type ProjectChange struct {
	bongo.DocumentBase `bson:",inline"`
	ProjectId          bson.ObjectId `bson:"projectId" json:"projectId"`
	ReleaseId          bson.ObjectId `bson:"releaseId,omitempty" json:"releaseId"`
	Name               string        `bson:"name" json:"name"`
	Message            string        `bson:"message" json:"message"`
}

type Statistics struct {
	Projects int `bson:"projects" json:"projects"`
	Releases int `bson:"deploys" json:"releases"`
	Features int `bson:"features" json:"features"`
	Users    int `bson:"users" json:"users"`
}

type PageResults struct {
	Records    interface{}          `json:"records"`
	Pagination bongo.PaginationInfo `json:"pagination"`
}

type Listener struct {
	Port     int    `bson:"port" json:"port"`
	Protocol string `bson:"protocol" json:"protocol"`
}

type ListenerPair struct {
	Source      Listener `bson:"source" json:"source"`
	Destination Listener `bson:"destination" json:"destination"`
}
