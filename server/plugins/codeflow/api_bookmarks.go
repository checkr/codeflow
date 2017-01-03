package codeflow

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
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
	var err error
	var user User
	var bookmarks []Bookmark

	if user, err = CurrentUser(r); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if bookmarks, err = GetUserBookmarks(user.Id); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteJson(bookmarks)
	}
}

func (x *Bookmarks) createBookmarks(w rest.ResponseWriter, r *rest.Request) {
	user, _ := CurrentUser(r)
	var bookmark Bookmark

	if err := r.DecodeJsonPayload(&bookmark); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := CreateUserBookmark(user.Id, bookmark.ProjectId); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	x.bookmarks(w, r)
}

func (x *Bookmarks) deleteBookmarks(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var user User
	var bookmark Bookmark

	if user, err = CurrentUser(r); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = r.DecodeJsonPayload(&bookmark); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = DeleteUserBookmark(user.Id, bookmark.ProjectId); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	x.bookmarks(w, r)
}
