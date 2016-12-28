package codeflow

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/mgo.v2/bson"
)

type Bookmarks struct {
	Path string
}

func (x *Bookmarks) Register(api *rest.Api) []*rest.Route {
	var routes []*rest.Route
	routes = append(routes,
		rest.Get(x.Path, x.bookmarks),
		rest.Post(x.Path, x.createBookmarks),
		rest.Delete(x.Path, x.deleteBookmarks),
	)
	log.Printf("Started the codeflow bookmarks handler on %s\n", x.Path)
	return routes
}

func (x *Bookmarks) bookmarks(w rest.ResponseWriter, r *rest.Request) {
	user, _ := CurrentUser(r)
	bookmarks := []Bookmark{}

	bookmarkCol := db.C("bookmarks")

	if err := bookmarkCol.Find(bson.M{"userId": user.Id}).All(&bookmarks); err != nil {
		fmt.Println("Bookmarks error: " + err.Error())
	}

	projectCol := db.C("projects")
	for idx, bookmark := range bookmarks {
		project := Project{}
		if err := projectCol.Find(bson.M{"_id": bookmark.ProjectId}).One(&project); err != nil {

		}
		bookmarks[idx].Name = project.Name
		bookmarks[idx].Slug = project.Slug
	}

	w.WriteJson(bookmarks)
}

func (x *Bookmarks) createBookmarks(w rest.ResponseWriter, r *rest.Request) {
	user, _ := CurrentUser(r)
	bookmark := Bookmark{}
	project := Project{}

	err := r.DecodeJsonPayload(&bookmark)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bookmarkCol := db.C("bookmarks")
	projectCol := db.C("projects")
	if err := projectCol.Find(bson.M{"_id": bookmark.ProjectId}).One(&project); err != nil {

	}

	bookmark.ProjectId = project.Id
	bookmark.UserId = user.Id

	if info, err := bookmarkCol.Upsert(bson.M{"userId": user.Id, "projectId": project.Id}, &bookmark); err != nil {
		spew.Dump(info)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	x.bookmarks(w, r)
}

func (x *Bookmarks) deleteBookmarks(w rest.ResponseWriter, r *rest.Request) {
	user, _ := CurrentUser(r)
	bookmark := Bookmark{}

	err := r.DecodeJsonPayload(&bookmark)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bookmarkCol := db.C("bookmarks")
	bookmarkCol.Remove(bson.M{"userId": user.Id, "projectId": bookmark.ProjectId})

	x.bookmarks(w, r)
}
