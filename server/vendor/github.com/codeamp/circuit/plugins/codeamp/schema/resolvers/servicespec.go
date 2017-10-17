package resolvers

import (
	"context"
	"fmt"

	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	graphql "github.com/neelance/graphql-go"
	uuid "github.com/satori/go.uuid"
)

type ServiceSpecInput struct {
	ID                     *string
	Name                   *string
	CpuRequest             *string
	CpuLimit               *string
	MemoryRequest          *string
	MemoryLimit            *string
	TerminationGracePeriod *string
}

type ServiceSpecResolver struct {
	db          *gorm.DB
	ServiceSpec models.ServiceSpec
}

func (r *Resolver) ServiceSpec(ctx context.Context, args *struct{ ID graphql.ID }) (*ServiceSpecResolver, error) {
	serviceSpec := models.ServiceSpec{}
	if err := r.db.Where("id = ?", args.ID).First(&serviceSpec).Error; err != nil {
		return nil, err
	}

	return &ServiceSpecResolver{db: r.db, ServiceSpec: serviceSpec}, nil
}

func (r *Resolver) CreateServiceSpec(args *struct{ ServiceSpec *ServiceSpecInput }) (*ServiceSpecResolver, error) {
	spew.Dump(args.ServiceSpec)
	serviceSpec := models.ServiceSpec{
		Name:                   *args.ServiceSpec.Name,
		CpuRequest:             *args.ServiceSpec.CpuRequest,
		CpuLimit:               *args.ServiceSpec.CpuLimit,
		MemoryRequest:          *args.ServiceSpec.MemoryRequest,
		MemoryLimit:            *args.ServiceSpec.MemoryLimit,
		TerminationGracePeriod: *args.ServiceSpec.TerminationGracePeriod,
	}

	r.db.Create(&serviceSpec)

	r.actions.ServiceSpecCreated(&serviceSpec)

	return &ServiceSpecResolver{db: r.db, ServiceSpec: serviceSpec}, nil
}

func (r *Resolver) UpdateServiceSpec(args *struct{ ServiceSpec *ServiceSpecInput }) (*ServiceSpecResolver, error) {
	serviceSpec := models.ServiceSpec{}

	serviceSpecId, err := uuid.FromString(*args.ServiceSpec.ID)
	if err != nil {
		return nil, fmt.Errorf("UpdateServiceSpec: Missing argument id")
	}

	if r.db.Where("id=?", serviceSpecId).Find(&serviceSpec).RecordNotFound() {
		return nil, fmt.Errorf("ServiceSpec not found with given argument id")
	}

	spew.Dump(args.ServiceSpec)

	serviceSpec.Name = *args.ServiceSpec.Name
	serviceSpec.CpuLimit = *args.ServiceSpec.CpuLimit
	serviceSpec.CpuRequest = *args.ServiceSpec.CpuRequest
	serviceSpec.MemoryLimit = *args.ServiceSpec.MemoryLimit
	serviceSpec.MemoryRequest = *args.ServiceSpec.MemoryRequest
	serviceSpec.TerminationGracePeriod = *args.ServiceSpec.TerminationGracePeriod

	r.db.Save(&serviceSpec)
	r.actions.ServiceSpecUpdated(&serviceSpec)

	return &ServiceSpecResolver{db: r.db, ServiceSpec: serviceSpec}, nil
}

func (r *Resolver) DeleteServiceSpec(args *struct{ ServiceSpec *ServiceSpecInput }) (*ServiceSpecResolver, error) {
	serviceSpec := models.ServiceSpec{}

	serviceSpecId, err := uuid.FromString(*args.ServiceSpec.ID)
	if err != nil {
		return nil, fmt.Errorf("Missing argument id")
	}

	if r.db.Where("id=?", serviceSpecId).Find(&serviceSpec).RecordNotFound() {
		return nil, fmt.Errorf("ServiceSpec not found with given argument id")
	}

	r.db.Delete(serviceSpec)

	r.actions.ServiceSpecDeleted(&serviceSpec)

	return &ServiceSpecResolver{db: r.db, ServiceSpec: serviceSpec}, nil
}

func (r *ServiceSpecResolver) ID() graphql.ID {
	return graphql.ID(r.ServiceSpec.Model.ID.String())
}

func (r *ServiceSpecResolver) Name(ctx context.Context) string {
	return r.ServiceSpec.Name
}

func (r *ServiceSpecResolver) CpuRequest(ctx context.Context) string {
	return r.ServiceSpec.CpuRequest
}

func (r *ServiceSpecResolver) CpuLimit(ctx context.Context) string {
	return r.ServiceSpec.CpuLimit
}

func (r *ServiceSpecResolver) MemoryLimit(ctx context.Context) string {
	return r.ServiceSpec.MemoryLimit
}

func (r *ServiceSpecResolver) MemoryRequest(ctx context.Context) string {
	return r.ServiceSpec.MemoryRequest
}

func (r *ServiceSpecResolver) TerminationGracePeriod(ctx context.Context) string {
	return r.ServiceSpec.TerminationGracePeriod
}

func (r *ServiceSpecResolver) Created() graphql.Time {
	return graphql.Time{Time: r.ServiceSpec.Created}
}
