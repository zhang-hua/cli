package actors

import (
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/models"
)

type ServiceActor struct {
	brokerRepo  api.ServiceBrokerRepository
	serviceRepo api.ServiceRepository
}

func NewServiceActor(broker api.ServiceBrokerRepository, service api.ServiceRepository) ServiceActor {
	return ServiceActor{
		brokerRepo:  broker,
		serviceRepo: service,
	}
}

func (actor ServiceActor) GetBrokersWithDependencies() ([]models.ServiceBroker, error) {
	brokers, err := actor.getServiceBrokers()
	if err != nil {
		return nil, err
	}

	for index, _ := range brokers {
		brokers[index].Services, err = actor.serviceRepo.ListServicesFromBroker(brokers[index].Guid)
	}
	if err != nil {
		return nil, err
	}

	return brokers, nil
}

func (actor ServiceActor) getServiceBrokers() (brokers []models.ServiceBroker, err error) {
	err = actor.brokerRepo.ListServiceBrokers(func(broker models.ServiceBroker) bool {
		brokers = append(brokers, broker)
		return true
	})
	return
}
