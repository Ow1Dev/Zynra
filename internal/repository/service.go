package repository

type serviceEntity struct {
	Address     string
}

type ServiceRepository struct {
	services map[string]serviceEntity
}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{
		services: make(map[string]serviceEntity),
	}
}

func (r *ServiceRepository) AddService(name string, address string) {
	r.services[name] = serviceEntity{
		Address: address,
	}
}

func (r *ServiceRepository) GetService(name string) (serviceEntity, bool) {
	entity, exists := r.services[name]
	return entity, exists
}
