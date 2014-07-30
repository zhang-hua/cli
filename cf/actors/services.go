package actors

import (
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/models"
)

type ServiceActor interface {
	GetAllBrokersWithDependencies() ([]models.ServiceBroker, error)
	GetBrokerWithDependencies(string) ([]models.ServiceBroker, error)
	GetBrokerWithSingleService(string) ([]models.ServiceBroker, error)
	GetBrokersWithVisibilityFromASingleOrg(string) ([]models.ServiceBroker, error)
	FilterBrokers(brokerFlag string, serviceFlag string, orgFlag string) ([]models.ServiceBroker, error)
}

type ServiceHandler struct {
	brokerRepo                api.ServiceBrokerRepository
	serviceRepo               api.ServiceRepository
	servicePlanRepo           api.ServicePlanRepository
	servicePlanVisibilityRepo api.ServicePlanVisibilityRepository
	orgRepo                   api.OrganizationRepository
}

func NewServiceHandler(broker api.ServiceBrokerRepository, service api.ServiceRepository, plan api.ServicePlanRepository, vis api.ServicePlanVisibilityRepository, org api.OrganizationRepository) ServiceHandler {
	return ServiceHandler{
		brokerRepo:                broker,
		serviceRepo:               service,
		servicePlanRepo:           plan,
		servicePlanVisibilityRepo: vis,
		orgRepo:                   org,
	}
}

func (actor ServiceHandler) FilterBrokers(brokerFlag string, serviceFlag string, orgFlag string) ([]models.ServiceBroker, error) {
	// Filter
	return nil, nil
}

func (actor ServiceHandler) GetBrokersWithVisibilityFromASingleOrg(orgName string) ([]models.ServiceBroker, error) {
	serviceToOrgPlansMap, err := actor.createMapOfServicesToVisiblePlans(orgName)
	if err != nil {
		return nil, err
	}

	brokers, err := actor.getAllBrokersFromServicesMap(serviceToOrgPlansMap)
	if err != nil {
		return nil, err
	}

	return brokers, nil
}

func (actor ServiceHandler) GetBrokerWithSingleService(serviceLabel string) ([]models.ServiceBroker, error) {
	service, err := actor.serviceRepo.FindServiceOfferingByLabel(serviceLabel)
	if err != nil {
		return nil, err
	}

	broker, err := actor.brokerRepo.FindByGuid(service.BrokerGuid)
	if err != nil {
		return nil, err
	}

	broker.Services = []models.ServiceOffering{service}
	brokers := []models.ServiceBroker{broker}

	brokers, err = actor.getServicePlans(brokers)
	if err != nil {
		return nil, err
	}

	return actor.getOrgs(brokers)
}

func (actor ServiceHandler) GetBrokerWithDependencies(brokerName string) ([]models.ServiceBroker, error) {
	broker, err := actor.brokerRepo.FindByName(brokerName)
	if err != nil {
		return nil, err
	}
	brokers := []models.ServiceBroker{broker}
	brokers, err = actor.getServices(brokers)
	if err != nil {
		return nil, err
	}

	brokers, err = actor.getServicePlans(brokers)
	if err != nil {
		return nil, err
	}

	return actor.getOrgs(brokers)
}

func (actor ServiceHandler) GetAllBrokersWithDependencies() ([]models.ServiceBroker, error) {
	brokers, err := actor.getAllServiceBrokers()
	if err != nil {
		return nil, err
	}

	brokers, err = actor.getServices(brokers)
	if err != nil {
		return nil, err
	}

	brokers, err = actor.getServicePlans(brokers)
	if err != nil {
		return nil, err
	}
	return actor.getOrgs(brokers)
}

func (actor ServiceHandler) getAllServiceBrokers() (brokers []models.ServiceBroker, err error) {
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
	//Is there a cleaner way to do this?
	for brokerIndex, _ := range brokers {
		broker := &brokers[brokerIndex]
		for serviceIndex, _ := range broker.Services {
			service := &broker.Services[serviceIndex]
			service.Plans, err = actor.servicePlanRepo.Search(map[string]string{"service_guid": service.Guid})
			if err != nil {
				return nil, err
			}
		}
	}
	return brokers, nil
}

func (actor ServiceHandler) getOrgs(brokers []models.ServiceBroker) ([]models.ServiceBroker, error) {
	visMap, err := actor.buildPlanToOrgsVisibilityMap()
	if err != nil {
		return nil, err
	}
	//Is there a cleaner way to do this?
	for brokerIndex, _ := range brokers {
		broker := &brokers[brokerIndex]
		for serviceIndex, _ := range broker.Services {
			service := &broker.Services[serviceIndex]
			for planIndex, _ := range service.Plans {
				plan := &service.Plans[planIndex]
				plan.OrgNames = visMap[plan.Guid]
			}
		}
	}
	return brokers, nil
}

func (actor ServiceHandler) buildPlanToOrgsVisibilityMap() (map[string][]string, error) {
	orgLookup := make(map[string]string)
	actor.orgRepo.ListOrgs(func(org models.Organization) bool {
		orgLookup[org.Guid] = org.Name
		return true
	})

	visibilities, err := actor.servicePlanVisibilityRepo.List()
	if err != nil {
		return nil, err
	}

	visMap := make(map[string][]string)
	for _, vis := range visibilities {
		visMap[vis.ServicePlanGuid] = append(visMap[vis.ServicePlanGuid], orgLookup[vis.OrganizationGuid])
	}

	return visMap, nil
}

func (actor ServiceHandler) buildOrgToPlansVisibilityMap(planToOrgsMap map[string][]string) map[string][]string {
	visMap := make(map[string][]string)
	for planGuid, orgNames := range planToOrgsMap {
		for _, orgName := range orgNames {
			visMap[orgName] = append(visMap[orgName], planGuid)
		}
	}

	return visMap
}

func (actor ServiceHandler) createMapOfServicesToVisiblePlans(orgName string) (map[string][]models.ServicePlanFields, error) {
	allPlans, err := actor.servicePlanRepo.Search(nil)
	if err != nil {
		return nil, err
	}

	servicesToPlansMap := make(map[string][]models.ServicePlanFields)
	PlanToOrgsVisMap, err := actor.buildPlanToOrgsVisibilityMap()
	if err != nil {
		return nil, err
	}
	OrgToPlansVisMap := actor.buildOrgToPlansVisibilityMap(PlanToOrgsVisMap)
	filterOrgPlans := OrgToPlansVisMap[orgName]

	for _, plan := range allPlans {
		if actor.containsGuid(filterOrgPlans, plan.Guid) {
			plan.OrgNames = PlanToOrgsVisMap[plan.Guid]
			servicesToPlansMap[plan.ServiceOfferingGuid] = append(servicesToPlansMap[plan.ServiceOfferingGuid], plan)
		} else if plan.Public {
			servicesToPlansMap[plan.ServiceOfferingGuid] = append(servicesToPlansMap[plan.ServiceOfferingGuid], plan)
		}
	}

	return servicesToPlansMap, nil
}

func (actor ServiceHandler) getAllBrokersFromServicesMap(serviceMap map[string][]models.ServicePlanFields) ([]models.ServiceBroker, error) {
	var brokers []models.ServiceBroker
	brokersToServices := make(map[string][]models.ServiceOffering)

	for serviceGuid, plans := range serviceMap {
		service, err := actor.serviceRepo.GetServiceOfferingByGuid(serviceGuid)
		if err != nil {
			return nil, err
		}
		service.Plans = plans
		brokersToServices[service.BrokerGuid] = append(brokersToServices[service.BrokerGuid], service)
	}

	for brokerGuid, services := range brokersToServices {
		if brokerGuid == "" {
			continue
		}
		broker, err := actor.brokerRepo.FindByGuid(brokerGuid)
		if err != nil {
			return nil, err
		}

		broker.Services = services
		brokers = append(brokers, broker)
	}

	return brokers, nil
}

func (actor ServiceHandler) containsGuid(guidSlice []string, guid string) bool {
	for _, g := range guidSlice {
		if g == guid {
			return true
		}
	}
	return false
}
