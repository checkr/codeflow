package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	graphql "github.com/neelance/graphql-go"
	uuid "github.com/satori/go.uuid"
)

type ServiceInput struct {
	ID             *string
	Name           string
	Command        string
	ServiceSpecId  string
	OneShot        bool
	Count          string
	ContainerPorts *[]models.ContainerPort
	ProjectId      string
}

type ServiceResolver struct {
	db      *gorm.DB
	Service models.Service
}

func (r *Resolver) Service(ctx context.Context, args *struct{ ID graphql.ID }) (*ServiceResolver, error) {
	service := models.Service{}
	if err := r.db.Where("id = ?", args.ID).First(&service).Error; err != nil {
		return nil, err
	}

	return &ServiceResolver{db: r.db, Service: service}, nil
}

func (r *Resolver) DeleteService(args *struct{ Service *ServiceInput }) (*ServiceResolver, error) {
	serviceId, err := uuid.FromString(*args.Service.ID)

	if err != nil {
		return &ServiceResolver{}, err
	}

	var service models.Service
	spew.Dump(args.Service)
	r.db.Where("id = ?", serviceId).Find(&service)
	r.db.Delete(&service)

	// delete all previous container ports
	var containerPorts []models.ContainerPort
	r.db.Where("service_id = ?", serviceId).Find(&containerPorts)

	// delete all container ports
	// replace with current

	for _, cp := range containerPorts {
		r.db.Delete(&cp)
	}

	r.actions.ServiceDeleted(&service)

	return &ServiceResolver{db: r.db, Service: service}, nil
}

func (r *Resolver) UpdateService(args *struct{ Service *ServiceInput }) (*ServiceResolver, error) {
	serviceId := uuid.FromStringOrNil(*args.Service.ID)
	serviceSpecId := uuid.FromStringOrNil(args.Service.ServiceSpecId)

	if serviceId == uuid.Nil || serviceSpecId == uuid.Nil {
		return nil, fmt.Errorf("Missing argument id")
	}

	var service models.Service
	if r.db.Where("id = ?", serviceId).Find(&service).RecordNotFound() {
		return nil, fmt.Errorf("Record not found with given argument id")
	}

	service.Command = args.Service.Command
	service.Name = args.Service.Name
	service.OneShot = args.Service.OneShot
	service.ServiceSpecId = serviceSpecId
	service.Count = args.Service.Count

	r.db.Save(&service)

	// delete all previous container ports
	var containerPorts []models.ContainerPort
	r.db.Where("service_id = ?", serviceId).Find(&containerPorts)

	// delete all container ports
	// replace with current

	for _, cp := range containerPorts {
		r.db.Delete(&cp)
	}

	if args.Service.ContainerPorts != nil {
		for _, cp := range *args.Service.ContainerPorts {
			containerPort := models.ContainerPort{
				ServiceId: service.ID,
				Port:      cp.Port,
				Protocol:  cp.Protocol,
			}
			r.db.Create(&containerPort)
		}
	}

	r.actions.ServiceUpdated(&service)

	return &ServiceResolver{db: r.db, Service: service}, nil
}

func (r *Resolver) CreateService(args *struct{ Service *ServiceInput }) (*ServiceResolver, error) {
	projectId, err := uuid.FromString(args.Service.ProjectId)
	if err != nil {
		return &ServiceResolver{}, err
	}
	serviceSpecId, err := uuid.FromString(args.Service.ServiceSpecId)
	if err != nil {
		return &ServiceResolver{}, err
	}
	spew.Dump(args.Service)
	service := models.Service{
		Name:          args.Service.Name,
		Command:       args.Service.Command,
		ServiceSpecId: serviceSpecId,
		OneShot:       args.Service.OneShot,
		Count:         args.Service.Count,
		ProjectId:     projectId,
		Created:       time.Now(),
	}

	r.db.Create(&service)
	spew.Dump(*args.Service.ContainerPorts)

	if args.Service.ContainerPorts != nil {
		for _, cp := range *args.Service.ContainerPorts {
			containerPort := models.ContainerPort{
				ServiceId: service.ID,
				Port:      cp.Port,
				Protocol:  cp.Protocol,
			}
			r.db.Create(&containerPort)
		}
	}

	r.actions.ServiceCreated(&service)

	return &ServiceResolver{db: r.db, Service: service}, nil
}

func (r *ServiceResolver) ID() graphql.ID {
	return graphql.ID(r.Service.Model.ID.String())
}

func (r *ServiceResolver) ServiceSpec(ctx context.Context) (*ServiceSpecResolver, error) {
	var serviceSpec models.ServiceSpec
	r.db.Model(r.Service).Related(&serviceSpec)
	return &ServiceSpecResolver{db: r.db, ServiceSpec: serviceSpec}, nil
}

func (r *ServiceResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	var project models.Project
	r.db.Model(r.Service).Related(&project)
	return &ProjectResolver{db: r.db, Project: project}, nil
}

func (r *ServiceResolver) Name() string {
	return r.Service.Name
}

func (r *ServiceResolver) Command() string {
	return r.Service.Command
}

func (r *ServiceResolver) Count() string {
	return r.Service.Count
}

func (r *ServiceResolver) Created() graphql.Time {
	return graphql.Time{Time: r.Service.Created}
}

func (r *ServiceResolver) OneShot() bool {
	return r.Service.OneShot
}

func (r *ServiceResolver) ContainerPorts(ctx context.Context) ([]*ContainerPortResolver, error) {
	var rows []models.ContainerPort
	var results []*ContainerPortResolver

	r.db.Where("service_id = ?", r.Service.ID).Order("created_at desc").Find(&rows)
	for _, cPort := range rows {
		spew.Dump(cPort)
		results = append(results, &ContainerPortResolver{ContainerPort: cPort})
	}
	spew.Dump(results)
	return results, nil
}

type ContainerPortResolver struct {
	ContainerPort models.ContainerPort
}

func (r *ContainerPortResolver) Port() string {
	return r.ContainerPort.Port
}

func (r *ContainerPortResolver) Protocol() string {
	return r.ContainerPort.Protocol
}
