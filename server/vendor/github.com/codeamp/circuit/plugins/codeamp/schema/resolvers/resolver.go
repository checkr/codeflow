package resolvers

import (
	"github.com/codeamp/circuit/plugins/codeamp/actions"
	"github.com/codeamp/transistor"
	"github.com/jinzhu/gorm"
)

type Resolver struct {
	db      *gorm.DB
	events  chan transistor.Event
	actions *actions.Actions
}

func NewResolver(events chan transistor.Event, db *gorm.DB, actions *actions.Actions) *Resolver {
	return &Resolver{
		events:  events,
		db:      db,
		actions: actions,
	}
}