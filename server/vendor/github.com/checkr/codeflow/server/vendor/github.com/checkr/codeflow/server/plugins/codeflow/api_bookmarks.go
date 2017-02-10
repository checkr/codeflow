package codeflow

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/maxwellhealth/bongo"
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
	user := User{}
	bookmarks := []Bookmark{}
	bookmark := Bookmark{}

	if err := CurrentUser(r, &user); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := db.Collection("bookmarks").Find(bson.M{"userId": user.Id})
	for results.Next(&bookmark) {
		if dnfError, ok := results.Error.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Bookmarks::Next::DocumentNotFoundError: %s", dnfError.Error())
			continue
		}
		bookmarks = append(bookmarks, bookmark)
	}

	w.WriteJson(bookmarks)
}

func (x *Bookmarks) createBookmarks(w rest.ResponseWriter, r *rest.Request) {
	user := User{}
	bookmark := Bookmark{}

	if err := CurrentUser(r, &user); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := r.DecodeJsonPayload(&bookmark); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bookmark.UserId = user.Id

	if err := db.Collection("bookmarks").Save(&bookmark); err != nil {
		log.Printf("Bookmarks::Save::Error: %v", err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	x.bookmarks(w, r)
}

func (x *Bookmarks) deleteBookmarks(w rest.ResponseWriter, r *rest.Request) {
	user := User{}
	bookmark := Bookmark{}

	if err := CurrentUser(r, &user); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := r.DecodeJsonPayload(&bookmark); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Collection("bookmarks").DeleteOne(bson.M{"projectId": bookmark.ProjectId, "userId": user.Id}); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	x.bookmarks(w, r)
}
