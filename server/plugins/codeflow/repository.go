package codeflow

import (
	"errors"
	"log"
	"math"
	"strings"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/codeflow/server/plugins"
	"gopkg.in/mgo.v2/bson"
)

// CurrentProject finds project by slug
func CurrentProject(r *rest.Request) (Project, error) {
	project := Project{}
	slug := r.PathParam("slug")

	projectCol := db.C("projects")
	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		return project, err
	}

	return project, nil
}

// CurrentUser returns current user identifed by JWT token
func CurrentUser(r *rest.Request) (User, error) {
	userCol := db.C("users")
	user := User{}

	if err := userCol.Find(bson.M{"username": r.Env["REMOTE_USER"]}).One(&user); err != nil {
		return user, err
	}

	return user, nil
}

// GetUserByEmail finds user by email and returns the User object or error
func GetUserByEmail(email string) (User, error) {
	userCol := db.C("users")
	user := User{}

	if err := userCol.Find(bson.M{"email": email}).One(&user); err != nil {
		return user, err
	} else {
		return user, nil
	}
}

// GetUserById finds user by id
func GetUserById(id bson.ObjectId) (User, error) {
	userCol := db.C("users")
	user := User{}

	if err := userCol.Find(bson.M{"_id": id}).One(&user); err != nil {
		return user, err
	} else {
		return user, nil
	}
}

// CreateUser creates a new user and sets Id field
func CreateUser(user *User) error {
	userCol := db.C("users")
	user.Id = bson.NewObjectId()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := userCol.Insert(&user); err != nil {
		return err
	}

	return nil
}

// UpdateUser finds existing user by id and updates the record
func UpdateUser(id bson.ObjectId, user *User) error {
	userCol := db.C("users")
	user.UpdatedAt = time.Now()

	if err := userCol.Update(bson.M{"_id": id}, bson.M{"$set": user}); err != nil {
		return err
	}

	return nil
}

// UpdateService finds existing service by id and updates the record
func UpdateService(id bson.ObjectId, service *Service) error {
	serviceCol := db.C("services")
	service.UpdatedAt = time.Now()

	if err := serviceCol.Update(bson.M{"_id": id}, bson.M{"$set": service}); err != nil {
		return err
	}

	return nil
}

// GetUserBookmarks returns list of user bookmarks
func GetUserBookmarks(userId bson.ObjectId) ([]Bookmark, error) {
	bookmarkCol := db.C("bookmarks")
	projectCol := db.C("projects")
	bookmarks := []Bookmark{}

	if err := bookmarkCol.Find(bson.M{"userId": userId}).All(&bookmarks); err != nil {
		return bookmarks, err
	}

	for idx, bookmark := range bookmarks {
		project := Project{}

		if err := projectCol.Find(bson.M{"_id": bookmark.ProjectId}).One(&project); err != nil {
			// Project not found so we need to delete the bookmark
			bookmarkCol.Remove(bson.M{"userId": userId, "projectId": bookmark.ProjectId})
			bookmarks = append(bookmarks[:idx], bookmarks[idx+1:]...)
			continue
		}

		bookmarks[idx].Name = project.Name
		bookmarks[idx].Slug = project.Slug
	}

	return bookmarks, nil
}

// CreateUserBookmark creates a new user bookmark
func CreateUserBookmark(userId bson.ObjectId, projectId bson.ObjectId) error {
	bookmarkCol := db.C("bookmarks")
	projectCol := db.C("projects")
	userCol := db.C("users")
	project := Project{}
	user := User{}
	bookmark := Bookmark{}

	if err := projectCol.Find(bson.M{"_id": projectId}).One(&project); err != nil {
		return err
	}

	if err := userCol.Find(bson.M{"_id": userId}).One(&user); err != nil {
		return err
	}

	bookmark.ProjectId = project.Id
	bookmark.UserId = user.Id

	if _, err := bookmarkCol.Upsert(bson.M{"userId": user.Id, "projectId": project.Id}, &bookmark); err != nil {
		return err
	}

	return nil
}

// DeleteUserBookmark removes user bookmark
func DeleteUserBookmark(userId bson.ObjectId, projectId bson.ObjectId) error {
	userCol := db.C("users")
	projectCol := db.C("projects")
	bookmarkCol := db.C("bookmarks")

	var project Project
	var user User

	if err := projectCol.Find(bson.M{"_id": projectId}).One(&project); err != nil {
		return err
	}

	if err := userCol.Find(bson.M{"_id": userId}).One(&user); err != nil {
		return err
	}

	if err := bookmarkCol.Remove(bson.M{"userId": user.Id, "projectId": project.Id}); err != nil {
		return err
	}

	return nil
}

// GetProject finds project by condition
func GetProject(condition bson.M) (Project, error) {
	projectCol := db.C("projects")
	var project Project

	if err := projectCol.Find(condition).One(&project); err != nil {
		return project, err
	} else {
		return project, nil
	}
}

// GetProjectByRepository finds project by repository name
func GetProjectByRepository(repository string) (Project, error) {
	if project, err := GetProject(bson.M{"repository": repository}); err != nil {
		return project, err
	} else {
		return project, nil
	}
}

// GetProjectById finds project by repository id
func GetProjectById(id bson.ObjectId) (Project, error) {
	if project, err := GetProject(bson.M{"_id": id}); err != nil {
		return project, err
	} else {
		return project, nil
	}
}

// UpdateProject finds existing project by id and updates the record
func UpdateProject(id bson.ObjectId, project *Project) error {
	projectCol := db.C("projects")
	project.UpdatedAt = time.Now()

	if err := projectCol.Update(bson.M{"_id": id}, bson.M{"$set": project}); err != nil {
		return err
	}

	return nil
}

// GetFeature finds feature by conditon
func GetFeature(condition bson.M) (Feature, error) {
	featureCol := db.C("features")
	var feature Feature

	if err := featureCol.Find(condition).One(&feature); err != nil {
		return feature, err
	} else {
		return feature, nil
	}
}

// GetFeatureById finds feature by id
func GetFeatureById(id bson.ObjectId) (Feature, error) {
	if feature, err := GetFeature(bson.M{"_id": id}); err != nil {
		return feature, err
	} else {
		return feature, nil
	}
}

// GetFeature finds feature by conditon
func GetFeatures(condition bson.M) ([]Feature, error) {
	featureCol := db.C("features")
	var features []Feature

	if err := featureCol.Find(condition).All(&features); err != nil {
		return features, err
	} else {
		return features, nil
	}
}

// GetFeatureByHash finds feture by hash
func GetFeatureByHash(hash string) (Feature, error) {
	if feature, err := GetFeature(bson.M{"hash": hash}); err != nil {
		return feature, err
	} else {
		return feature, nil
	}
}

// CreateFeature creates new feature
func CreateFeature(feature *Feature) error {
	featureCol := db.C("features")
	feature.Id = bson.NewObjectId()
	feature.CreatedAt = time.Now()
	feature.UpdatedAt = time.Now()

	if err := featureCol.Insert(&feature); err != nil {
		return err
	}

	return nil
}

// GetBuildByHash finds build by hash
func GetBuildByHash(hash string) (Build, error) {
	buildCol := db.C("builds")
	var build Build

	if err := buildCol.Find(bson.M{"featureHash": hash}).One(&build); err != nil {
		return build, err
	} else {
		return build, nil
	}
}

// GetBuildByHash finds build by hash
func GetBuildByHashAndType(hash string, typ string) (Build, error) {
	buildCol := db.C("builds")
	var build Build

	if err := buildCol.Find(bson.M{"featureHash": hash, "type": typ}).One(&build); err != nil {
		return build, err
	} else {
		return build, nil
	}
}

// GetBuildByHashAndState finds build by hash
func GetBuildByHashAndState(hash string, state plugins.State) (Build, error) {
	buildCol := db.C("builds")
	var build Build

	if err := buildCol.Find(bson.M{"featureHash": hash, "state": state}).One(&build); err != nil {
		return build, err
	} else {
		return build, nil
	}
}

// UpdateBuild finds existing build by id and updates the record
func UpdateBuild(id bson.ObjectId, build *Build) error {
	buildCol := db.C("builds")
	build.UpdatedAt = time.Now()

	if err := buildCol.Update(bson.M{"_id": id}, bson.M{"$set": build}); err != nil {
		return err
	}

	return nil
}

// CreateBuild creates new build
func CreateBuild(build *Build) error {
	buildCol := db.C("builds")
	build.Id = bson.NewObjectId()
	build.CreatedAt = time.Now()
	build.UpdatedAt = time.Now()

	if err := buildCol.Insert(&build); err != nil {
		return err
	}

	return nil
}

// CreateRelease creates new release
func CreateRelease(release *Release) error {
	releaseCol := db.C("releases")
	release.Id = bson.NewObjectId()
	release.State = plugins.Waiting
	release.CreatedAt = time.Now()
	release.UpdatedAt = time.Now()

	if err := releaseCol.Insert(&release); err != nil {
		return err
	}

	return nil
}

// PopulateRelease finds release by id and populates relations
func PopulateRelease(release *Release) error {
	var err error
	var headFeature Feature
	var tailFeature Feature
	var user User
	var workflows []Flow

	if !release.HeadFeatureId.Valid() {
		release.HeadFeatureId = release.HeadFeature.Id
	}

	if !release.TailFeatureId.Valid() {
		release.TailFeatureId = release.TailFeature.Id
	}

	if !release.UserId.Valid() {
		release.UserId = release.User.Id
	}

	if !release.HeadFeatureId.Valid() || !release.TailFeatureId.Valid() || !release.UserId.Valid() {
		return errors.New("One of feature ids is missing")
	}

	if headFeature, err = GetFeatureById(release.HeadFeatureId); err != nil {
		return err
	} else {
		release.HeadFeature = headFeature
	}

	if tailFeature, err = GetFeatureById(release.TailFeatureId); err != nil {
		return err
	} else {
		release.TailFeature = tailFeature
	}

	if user, err = GetUserById(release.UserId); err != nil {
		return err
	} else {
		release.User = user
	}

	if workflows, err = GetWorkflowsByReleaseId(release.Id); err != nil {
		return err
	} else {
		release.Workflow = workflows
	}

	return nil
}

// GetWorkflowsByReleaseId finds flows by release id
func GetWorkflowsByReleaseId(releaseId bson.ObjectId) ([]Flow, error) {
	workflowCol := db.C("workflows")
	var workflows []Flow

	if err := workflowCol.Find(bson.M{"releaseId": releaseId}).All(&workflows); err != nil {
		return workflows, err
	} else {
		return workflows, nil
	}
}

// GetRelease finds release by condition
func GetRelease(condition bson.M) (Release, error) {
	releaseCol := db.C("releases")
	var release Release

	if err := releaseCol.Find(condition).One(&release); err != nil {
		return release, err
	} else {
		return release, nil
	}
}

// GetReleaseById finds service by id
func GetReleaseById(id bson.ObjectId) (Release, error) {
	if release, err := GetRelease(bson.M{"_id": id}); err != nil {
		return release, err
	} else {
		return release, nil
	}
}

// GetReleasesByFeatureIdAndState finds release by id and state
func GetReleasesByFeatureIdAndState(featureId bson.ObjectId, state plugins.State) ([]Release, error) {
	releaseCol := db.C("releases")
	var releases []Release

	if err := releaseCol.Find(bson.M{"headFeatureId": featureId, "state": state}).All(&releases); err != nil {
		return releases, err
	} else {
		return releases, nil
	}
}

// GetService finds service by condition
func GetService(condition bson.M) (Service, error) {
	serviceCol := db.C("services")
	var service Service

	if err := serviceCol.Find(condition).One(&service); err != nil {
		return service, err
	} else {
		return service, nil
	}
}

// GetServiceById finds service by id
func GetServiceById(id bson.ObjectId) (Service, error) {
	if service, err := GetService(bson.M{"_id": id}); err != nil {
		return service, err
	} else {
		return service, nil
	}
}

// GetServices finds services by condition
func GetServices(condition bson.M) ([]Service, error) {
	serviceCol := db.C("services")
	var services []Service

	if err := serviceCol.Find(condition).All(&services); err != nil {
		return services, err
	} else {
		return services, nil
	}
}

// GetSeGetServicesByProjectId finds services by project id
func GetServicesByProjectId(projectId bson.ObjectId) ([]Service, error) {
	if services, err := GetServices(bson.M{"projectId": projectId}); err != nil {
		return services, err
	} else {
		return services, nil
	}
}

// CreateWorkflow creates a new workflow and sets Id field
func CreateWorkflow(flow *Flow) error {
	workflowCol := db.C("workflows")
	flow.Id = bson.NewObjectId()
	flow.CreatedAt = time.Now()
	flow.UpdatedAt = time.Now()

	if err := workflowCol.Insert(&flow); err != nil {
		return err
	}

	return nil
}

// UpdateFlow finds existing flow by id and updates the record
func UpdateFlow(id bson.ObjectId, flow *Flow) error {
	workflowCol := db.C("workflows")
	flow.UpdatedAt = time.Now()

	if err := workflowCol.Update(bson.M{"_id": id}, bson.M{"$set": flow}); err != nil {
		return err
	}

	return nil
}

// UpdateRelease finds existing release by id and updates the record
func UpdateRelease(id bson.ObjectId, release *Release) error {
	releaseCol := db.C("releases")
	release.UpdatedAt = time.Now()

	if err := releaseCol.Update(bson.M{"_id": id}, bson.M{"$set": release}); err != nil {
		return err
	}

	return nil
}

// CreateOrUpdateBuildByFeatureHash creates or updates build
func CreateOrUpdateBuildByFeatureHash(hash string, build *Build) error {
	var err error
	var b Build

	if b, err = GetBuildByHash(hash); err != nil {
		if err = CreateBuild(build); err != nil {
			return err
		}
	} else {
		if err = UpdateBuild(b.Id, build); err != nil {
			return err
		}
	}

	return nil
}

// CreateSecret creates a new secret and sets Id field
func CreateSecret(secret *Secret) error {
	secretCol := db.C("secrets")
	secret.Id = bson.NewObjectId()
	secret.CreatedAt = time.Now()

	if err := secretCol.Insert(&secret); err != nil {
		return err
	}

	return nil
}

// GetSecretById finds secret by id
func GetSecretById(id bson.ObjectId) (Secret, error) {
	secretCol := db.C("secrets")
	var secret Secret

	if err := secretCol.Find(bson.M{"_id": id}).One(&secret); err != nil {
		return secret, err
	} else {
		return secret, nil
	}
}

// UpdateSecret finds existing secret and updates it
func UpdateSecret(id bson.ObjectId, secret *Secret) error {
	secretCol := db.C("secrets")

	if secret.Deleted == true {
		secret.DeletedAt = time.Now()
	}
	if err := secretCol.Update(bson.M{"_id": id}, bson.M{"$set": secret}); err != nil {
		return err
	}

	return nil
}

// GetSecretsByProjectIdAndType finds project secrets by type
func GetSecretsByProjectIdAndType(projectId bson.ObjectId, typ plugins.Type) ([]Secret, error) {
	secretCol := db.C("secrets")
	var secrets []Secret

	if err := secretCol.Find(bson.M{"projectId": projectId, "deleted": false, "type": typ}).All(&secrets); err != nil {
		return secrets, err
	} else {
		return secrets, nil
	}
}

// GetSecretsByProjectId finds project secrets
func GetSecretsByProjectId(projectId bson.ObjectId) ([]Secret, error) {
	secretCol := db.C("secrets")
	var secrets []Secret

	if err := secretCol.Find(bson.M{"projectId": projectId, "deleted": false}).All(&secrets); err != nil {
		return secrets, err
	} else {
		return secrets, nil
	}
}

// GetReleaseServices finds active services that can be released
func GetReleaseServices(projectId bson.ObjectId) ([]Service, error) {
	serviceCol := db.C("services")
	var services []Service

	if err := serviceCol.Find(bson.M{"projectId": projectId, "state": bson.M{"$in": []plugins.State{plugins.Waiting, plugins.Running}}}).All(&services); err != nil {
		return services, err
	} else {
		return services, nil
	}
}

// GetExtensionByProjectIdAndName finds extension by project id and name
func GetExtensionByProjectIdAndName(projectId bson.ObjectId, name string) (LoadBalancer, error) {
	var extension LoadBalancer
	extensionCol := db.C("extensions")

	if err := extensionCol.Find(bson.M{"projectId": projectId, "name": name}).All(&extension); err != nil {
		return extension, err
	} else {
		return extension, nil
	}
}

// UpdateExtension finds existing extension by id and updates the record
func UpdateExtension(id bson.ObjectId, extension *LoadBalancer) error {
	extensionCol := db.C("extensions")
	extension.UpdatedAt = time.Now()

	if err := extensionCol.Update(bson.M{"_id": id}, bson.M{"$set": extension}); err != nil {
		return err
	}

	return nil
}

// GetExtensions finds feature by conditon
func GetExtensions(projectId bson.ObjectId) ([]LoadBalancer, error) {
	extensionCol := db.C("extensions")
	var extensions []LoadBalancer

	if err := extensionCol.Find(bson.M{"projectId": projectId}).All(&extensions); err != nil {
		return extensions, err
	} else {
		return extensions, nil
	}
}

// Calc produces a Paging calculated over
// Page, Limit, Count and DefaultLimit values of given Paging
func (p Pagination) Calc() *Pagination {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.Limit < 1 {
		if p.DefaultLimit > 0 {
			p.Limit = p.DefaultLimit
		} else {
			p.Limit = 10
		}
	}

	if p.Count < p.Limit {
		p.TotalPages = 1
	} else {
		p.TotalPages = int(math.Ceil(float64(p.Count) / float64(p.Limit)))
	}

	if p.Page > p.TotalPages {
		p.Page = p.TotalPages
	}

	p.Offset = (p.Page - 1) * p.Limit

	return &p
}

func StringToLoadBalancerType(s string) plugins.Type {
	switch strings.ToLower(s) {
	case "internal":
		return plugins.Internal
	case "external":
		return plugins.External
	case "office":
		return plugins.Office
	}
	return plugins.Internal
}

// CreateProject creates a new project and sets Id field
func CreateProject(project *Project) error {
	projectCol := db.C("projects")
	project.Id = bson.NewObjectId()
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	if err := projectCol.Insert(&project); err != nil {
		return err
	}

	return nil
}

// CreateService creates a new service and sets Id field
func CreateService(service *Service) error {
	serviceCol := db.C("services")
	service.Id = bson.NewObjectId()
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()

	if err := serviceCol.Insert(&service); err != nil {
		return err
	}

	return nil
}

// CreateExtension creates a new extension and sets Id field
func CreateExtension(extension *LoadBalancer) error {
	extensionCol := db.C("extensions")
	extension.Id = bson.NewObjectId()
	extension.CreatedAt = time.Now()
	extension.UpdatedAt = time.Now()

	if err := extensionCol.Insert(&extension); err != nil {
		return err
	}

	return nil
}

func GetCurrentlyDeployedFeature(projectId bson.ObjectId, feature *Feature) error {
	var err error
	var project Project
	var release Release
	featureCol := db.C("features")
	releaseCol := db.C("releases")

	if project, err = GetProjectById(projectId); err != nil {
		log.Println(err.Error())
		return err
	}

	if err := releaseCol.Find(bson.M{"projectId": project.Id, "state": plugins.Complete}).Sort("-$natural").One(&release); err != nil {
		return err
	} else {
		if err := featureCol.Find(bson.M{"_id": release.HeadFeatureId}).One(&feature); err != nil {
			return err
		}
	}

	return nil
}

func GetFeaturesWithPagination(projectId bson.ObjectId, currPage int, limit int) (PageResults, error) {
	var err error
	var itemsCount int
	features := []Feature{}
	var pagination *Pagination
	var currentlyDeployedFeature Feature
	var pageResults PageResults

	featureCol := db.C("features")

	itemsCount, err = featureCol.Find(bson.M{"projectId": projectId}).Count()
	if err != nil {
		return pageResults, err
	}

	if err = GetCurrentlyDeployedFeature(projectId, &currentlyDeployedFeature); err != nil {
		pagination = Pagination{
			Page:  currPage,
			Limit: limit,
			Count: itemsCount,
		}.Calc()

		if err = featureCol.Find(bson.M{"projectId": projectId}).Limit(limit).Skip(pagination.Offset).Sort("-$natural").All(&features); err != nil {
			return pageResults, err
		}
	} else {
		itemsCount, err := featureCol.Find(bson.M{"projectId": projectId, "_id": bson.M{"$gt": currentlyDeployedFeature.Id}}).Count()
		if err != nil {
			return pageResults, err
		}

		pagination = Pagination{
			Page:  currPage,
			Limit: limit,
			Count: itemsCount,
		}.Calc()

		if err = featureCol.Find(bson.M{"projectId": projectId, "_id": bson.M{"$gt": currentlyDeployedFeature.Id}}).Limit(pagination.Limit).Skip(pagination.Offset).Sort("-$natural").All(&features); err != nil {
			return pageResults, err
		}
	}

	pageResults = PageResults{
		Records:    &features,
		Pagination: pagination,
	}

	return pageResults, nil
}

func GetReleasesWithPagination(projectId bson.ObjectId, currPage int, limit int) (PageResults, error) {
	var err error
	var itemsCount int
	releases := []Release{}
	var pagination *Pagination
	var pageResults PageResults

	releaseCol := db.C("releases")

	itemsCount, err = releaseCol.Find(bson.M{"projectId": projectId}).Count()
	if err != nil {

	}

	pagination = Pagination{
		Page:  currPage,
		Limit: limit,
		Count: itemsCount,
	}.Calc()

	if err := releaseCol.Find(bson.M{"projectId": projectId}).Limit(pagination.Limit).Skip(pagination.Offset).Sort("-$natural").All(&releases); err != nil {

	}

	for idx, rel := range releases {
		if err = PopulateRelease(&rel); err != nil {
			log.Println(err.Error())
			return pageResults, err
		}
		releases[idx] = rel
	}

	pageResults = PageResults{
		Records:    releases,
		Pagination: pagination,
	}

	return pageResults, nil
}

func GetCurrentRelease(projectId bson.ObjectId) (Release, error) {
	var err error
	var release Release
	releaseCol := db.C("releases")

	if err := releaseCol.Find(bson.M{"projectId": projectId, "state": plugins.Complete}).Sort("-$natural").One(&release); err != nil {
		return release, err
	}

	if err = PopulateRelease(&release); err != nil {
		return release, err
	}

	return release, nil
}

func GetProjectsWithPagination(currPage int, limit int) (PageResults, error) {
	var err error
	var itemsCount int
	projects := []Project{}
	var pagination *Pagination
	var pageResults PageResults

	projectCol := db.C("projects")

	itemsCount, err = projectCol.Find(bson.M{}).Count()
	if err != nil {

	}

	pagination = Pagination{
		Page:  currPage,
		Limit: limit,
		Count: itemsCount,
	}.Calc()

	if err := projectCol.Find(bson.M{}).Limit(pagination.Limit).Skip(pagination.Offset).Sort("-$natural").All(&projects); err != nil {

	}

	pageResults = PageResults{
		Records:    &projects,
		Pagination: pagination,
	}

	return pageResults, nil
}

func CollectStats(save bool) (Statistics, error) {
	stats := Statistics{}

	// Project stats
	projectCol := db.C("projects")
	if projectCount, err := projectCol.Find(bson.M{}).Count(); err != nil {
		return stats, err
	} else {
		stats.Projects = projectCount
	}

	// Feature stats
	featureCol := db.C("features")
	if featureCount, err := featureCol.Find(bson.M{}).Count(); err != nil {
		return stats, err
	} else {
		stats.Features = featureCount
	}

	// Release stats
	releaseCol := db.C("releases")
	if releaseCount, err := releaseCol.Find(bson.M{}).Count(); err != nil {
		return stats, err
	} else {
		stats.Releases = releaseCount
	}

	// User stats
	userCol := db.C("users")
	if userCount, err := userCol.Find(bson.M{}).Count(); err != nil {
		return stats, err
	} else {
		stats.Users = userCount
	}
	return stats, nil
}
