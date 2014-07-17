package actors

import (
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/models"
)

type ServiceActor interface {
	GetBrokersWithDependencies() ([]models.ServiceBroker, error)
}

type ServiceHandler struct {
	brokerRepo      api.ServiceBrokerRepository
	serviceRepo     api.ServiceRepository
	servicePlanRepo api.ServicePlanRepository
}

func NewServiceHandler(broker api.ServiceBrokerRepository, service api.ServiceRepository, plan api.ServicePlanRepository) ServiceHandler {
	return ServiceHandler{
		brokerRepo:      broker,
		serviceRepo:     service,
		servicePlanRepo: plan,
	}
}

func (actor ServiceHandler) GetBrokersWithDependencies() ([]models.ServiceBroker, error) {
	brokers, err := actor.getServiceBrokers()
	if err != nil {
		return nil, err
	}

	brokers, err = actor.getServices(brokers)
	if err != nil {
		return nil, err
	}

	return actor.getServicePlans(brokers)
}

func (actor ServiceHandler) getServiceBrokers() (brokers []models.ServiceBroker, err error) {
	err = actor.brokerRepo.ListServiceBrokers(func(broker models.ServiceBroker) bool {
		brokers = append(brokers, broker)
		return true
	})
	return
}

func (actor ServiceHandler) getServices(brokers []models.ServiceBroker) ([]models.ServiceBroker, error) {
	var err error
	for index, _ := range brokers {
		brokers[index].Services, err = actor.serviceRepo.ListServicesFromBroker(brokers[index].Guid)
		if err != nil {
			return nil, err
		}
	}
	return brokers, nil
}

func (actor ServiceHandler) getServicePlans(brokers []models.ServiceBroker) ([]models.ServiceBroker, error) {
	var err error
	for brokerIndex, _ := range brokers {
		broker := &brokers[brokerIndex]
		for serviceIndex, _ := range broker.Services {
			service := &broker.Services[serviceIndex]
			service.Plans, err = actor.servicePlanRepo.Search(map[string]string{"service-guid": service.Guid})
			if err != nil {
				return nil, err
			}
		}
	}
	return brokers, nil
}
