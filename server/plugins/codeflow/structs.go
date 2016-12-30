package codeflow

import (
	"time"

	"github.com/checkr/codeflow/server/plugins"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Username  string        `bson:"username" json:"username"`
	Email     string        `bson:"email" json:"email"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type Project struct {
	Id            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name          string        `bson:"name" json:"name"`
	Slug          string        `bson:"slug" json:"slug"`
	Repository    string        `bson:"repository" json:"repository"`
	Secret        string        `bson:"secret" json:"secret"`
	Pinged        bool          `bson:"pinged" json:"pinged"`
	GitSshUrl     string        `bson:"gitSshUrl" json:"gitSshUrl" validate:"required"`
	RsaPrivateKey string        `bson:"rsaPrivateKey" json:"-"`
	RsaPublicKey  string        `bson:"rsaPublicKey" json:"rsaPublicKey"`
	Bokmarked     bool          `bson:"-" json:"bookmarked"`
	CreatedAt     time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type Bookmark struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"-"`
	ProjectId bson.ObjectId `bson:"projectId" json:"projectId"`
	UserId    bson.ObjectId `bson:"userId" json:"-"`
	Name      string        `bson:"-" json:"name"`
	Slug      string        `bson:"-" json:"slug"`
}

type Service struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ProjectId bson.ObjectId `bson:"projectId" json:"projectId"`
	State     plugins.State `bson:"state" json:"state"`
	Name      string        `bson:"name" json:"name"`
	Count     int           `bson:"count" json:"count"`
	Command   string        `bson:"command" json:"command"`
	Listeners []Listener    `bson:"listeners" json:"listeners"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type Listener struct {
	Port     int    `bson:"port" json:"port"`
	Protocol string `bson:"protocol" json:"protocol"`
}

type ListenerPair struct {
	Source      Listener `bson:"source" json:"source"`
	Destination Listener `bson:"destination" json:"destination"`
}

type LoadBalancer struct {
	Id            bson.ObjectId  `bson:"_id,omitempty" json:"id"`
	Name          string         `bson:"name" json:"name"`
	ProjectId     bson.ObjectId  `bson:"projectId" json:"projectId"`
	ServiceId     bson.ObjectId  `bson:"serviceId" json:"serviceId"`
	Extension     string         `bson:"extension" json:"extension"`
	DNSName       string         `bson:"dnsName" json:"dnsName"`
	Type          string         `bson:"type" json:"type"`
	ListenerPairs []ListenerPair `bson:"listenerPairs" json:"listenerPairs"`
	State         plugins.State  `bson:"state" json:"state"`
	StateMessage  string         `bson:"stateMessage" json:"stateMessage"`
	CreatedAt     time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time      `bson:"updatedAt" json:"updatedAt"`
}

type Feature struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ProjectId    bson.ObjectId `bson:"projectId" json:"projectId"`
	Message      string        `bson:"message" json:"message"`
	User         string        `bson:"user" json:"user"`
	Hash         string        `bson:"hash" json:"hash"`
	ParentHash   string        `bson:"parentHash" json:"parentHash"`
	ExternalLink string        `bson:"externalLink" json:"externalLink"`
	CreatedAt    time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type Release struct {
	Id            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ProjectId     bson.ObjectId `bson:"projectId" json:"projectId"`
	HeadFeatureId bson.ObjectId `bson:"headFeatureId" json:"-"`
	HeadFeature   Feature       `bson:"-" json:"headFeature"`
	TailFeatureId bson.ObjectId `bson:"tailFeatureId" json:"-"`
	TailFeature   Feature       `bson:"-" json:"tailFeature"`
	UserId        bson.ObjectId `bson:"userId" json:"-"`
	User          User          `bson:"-" json:"user"`
	State         plugins.State `bson:"state" json:"state"`
	Secrets       []Secret      `bson:"secrets" json:"-"`
	Workflow      []Flow        `bson:"-" json:"workflow"`
	CreatedAt     time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type Flow struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ReleaseId bson.ObjectId `bson:"releaseId,omitempty" json:"releaseId"`
	Type      string        `bson:"type" json:"type"`
	Name      string        `bson:"name" json:"name"`
	Message   string        `bson:"message" json:"message"`
	State     plugins.State `bson:"state" json:"state"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"upadatedAt" json:"updatedAt"`
}

type Secret struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	ProjectId bson.ObjectId `bson:"projectId" json:"-"`
	Deleted   bool          `bson:"deleted" json:"deleted,omitempty"`
	Key       string        `bson:"key" json:"key"`
	Value     string        `bson:"value" json:"value"`
	Type      plugins.Type  `bson:"type" json:"type"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt,omitempty"`
	DeletedAt time.Time     `bson:"deletedAt" json:"deletedAt,omitempty"`
}

type Build struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	FeatureHash string        `bson:"featureHash" json:"featureHash"`
	Type        string        `bson:"type" json:"type"`
	State       plugins.State `bson:"state" json:"state"`
	Image       string        `bson:"image" json:"image"`
	BuildLog    string        `bson:"buildLog" json:"buildLog"`
	BuildError  string        `bson:"buildError" json:"buildError"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type Pagination struct {
	Page                int `json:"page"`
	Limit, DefaultLimit int `json:"limit"`
	Offset              int `json:"offset"`
	TotalPages          int `json:"total_pages"`
	Count               int `json:"count"`
}

type PageResults struct {
	Records    interface{} `json:"records"`
	Pagination *Pagination `json:"pagination"`
}

type ProjectSettings struct {
	ProjectId      bson.ObjectId `json:"projectId"`
	GitSshUrl      string        `json:"gitSshUrl"`
	Secrets        []Secret      `json:"secrets"`
	DeletedSecrets []Secret      `json:"deletedSecrets"`
	UpdatedAt      time.Time     `json:"updatedAt"`
}

type ProjectChange struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ProjectId bson.ObjectId `bson:"projectId" json:"projectId"`
	ReleaseId bson.ObjectId `bson:"releaseId,omitempty" json:"releaseId"`
	Name      string        `bson:"name" json:"name"`
	Message   string        `bson:"message" json:"message"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type Statistics struct {
	Projects  int       `bson:"projects" json:"projects"`
	Releases  int       `bson:"deploys" json:"releases"`
	Features  int       `bson:"features" json:"features"`
	Users     int       `bson:"users" json:"users"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
