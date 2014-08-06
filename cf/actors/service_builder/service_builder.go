package service_builder

import (
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/models"
)

type ServiceBuilder interface {
	GetServicesForBroker(string) ([]models.ServiceOffering, error)

	GetServiceVisibleToOrg(string, string) ([]models.ServiceOffering, error)
	GetServicesVisibleToOrg(string) ([]models.ServiceOffering, error)
	AttachPlansToService(models.ServiceOffering) (models.ServiceOffering, error)
}

type Builder struct {
	serviceRepo               api.ServiceRepository
	servicePlanRepo           api.ServicePlanRepository
	servicePlanVisibilityRepo api.ServicePlanVisibilityRepository
	orgRepo                   api.OrganizationRepository
}

func NewBuilder(service api.ServiceRepository, plan api.ServicePlanRepository, vis api.ServicePlanVisibilityRepository, org api.OrganizationRepository) Builder {
	return Builder{
		serviceRepo:               service,
		servicePlanRepo:           plan,
		servicePlanVisibilityRepo: vis,
		orgRepo:                   org,
	}
}

//func (builder Builder) attachOrgsToPlans(plans []models.ServicePlanFields) ([]models.ServicePlanFields, error) {
// 	visMap, err := builder.buildPlanToOrgsVisibilityMap()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for planIndex, _ := range plans {
// 		plan := &plans[planIndex]
// 		plan.OrgNames = visMap[plan.Guid]
// 	}

// 	return plans, nil
// }

// func (builder Builder) attachPlansToServices(services []models.ServiceOffering) ([]models.ServiceOffering, error) {
// 	var err error
// 	for serviceIndex, _ := range services {
// 		service := services[serviceIndex]
// 		services[serviceIndex], err = builder.AttachPlansToService(service)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return services, nil
// }

func (builder Builder) AttachPlansToService(service models.ServiceOffering) (models.ServiceOffering, error) {
	// 	plans, err := builder.servicePlanRepo.Search(map[string]string{"service_guid": service.Guid})
	// 	if err != nil {
	// 		return models.ServiceOffering{}, err
	// 	}
	// 	service.Plans, err = builder.attachOrgsToPlans(plans)
	// 	if err != nil {
	// 		return models.ServiceOffering{}, err
	// 	}
	// 	return service, nil
	return models.ServiceOffering{}, nil
}

func (builder Builder) GetServiceVisibleToOrg(serviceName string, orgName string) ([]models.ServiceOffering, error) {
	return nil, nil
}

func (builder Builder) GetServicesVisibleToOrg(orgName string) ([]models.ServiceOffering, error) {
	return nil, nil
}

// func (builder Builder) attachServicesToBrokers(brokers []models.ServiceBroker) ([]models.ServiceBroker, error) {
// 	for index, _ := range brokers {
// 		services, err := builder.serviceRepo.ListServicesFromBroker(brokers[index].Guid)
// 		if err != nil {
// 			return nil, err
// 		}
// 		brokers[index].Services, err = builder.attachPlansToServices(services)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return brokers, nil
// }

// func (builder Builder) attachSingleServiceToBrokers(brokers []models.ServiceBroker, serviceName string) ([]models.ServiceBroker, error) {
// 	service, err := builder.serviceRepo.FindServiceOfferingByLabel(serviceName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for index, _ := range brokers {
// 		if brokers[index].Guid == service.BrokerGuid {
// 			//check to see if its guid is contained in the list of the broker's service guid
// 			brokers[index].Services, err = builder.attachPlansToServices([]models.ServiceOffering{service})
// 			if err != nil {
// 				return nil, err
// 			}
// 			return []models.ServiceBroker{brokers[index]}, nil
// 		} else {
// 			continue
// 		}
// 	}
// 	return []models.ServiceBroker{}, nil
// }

// func (builder Builder) buildBrokersFromServicesToPlansMap(serviceToPlansMap map[string][]models.ServicePlanFields) ([]models.ServiceBroker, error) {
// 	brokersToServicesMap, err := builder.buildBrokersToServicesMap(serviceToPlansMap)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return builder.getAllBrokersFromServicesMap(brokersToServicesMap)
// }

// func (builder Builder) buildBrokersToServicesMap(serviceMap map[string][]models.ServicePlanFields) (map[string][]models.ServiceOffering, error) {
// 	brokersToServices := make(map[string][]models.ServiceOffering)

// 	for serviceGuid, plans := range serviceMap {
// 		service, err := builder.serviceRepo.GetServiceOfferingByGuid(serviceGuid)
// 		if err != nil {
// 			return nil, err
// 		}
// 		service.Plans = plans
// 		brokersToServices[service.BrokerGuid] = append(brokersToServices[service.BrokerGuid], service)
// 	}
// 	return brokersToServices, nil
// }

// func (builder Builder) buildOrgToPlansVisibilityMap(planToOrgsMap map[string][]string) map[string][]string {
// 	visMap := make(map[string][]string)
// 	for planGuid, orgNames := range planToOrgsMap {
// 		for _, orgName := range orgNames {
// 			visMap[orgName] = append(visMap[orgName], planGuid)
// 		}
// 	}

// 	return visMap
// }

// func (builder Builder) buildPlanToOrgsVisibilityMap() (map[string][]string, error) {
// 	orgLookup := make(map[string]string)
// 	builder.orgRepo.ListOrgs(func(org models.Organization) bool {
// 		orgLookup[org.Guid] = org.Name
// 		return true
// 	})

// 	visibilities, err := builder.servicePlanVisibilityRepo.List()
// 	if err != nil {
// 		return nil, err
// 	}

// 	visMap := make(map[string][]string)
// 	for _, vis := range visibilities {
// 		visMap[vis.ServicePlanGuid] = append(visMap[vis.ServicePlanGuid], orgLookup[vis.OrganizationGuid])
// 	}

// 	return visMap, nil
// }

// func (builder Builder) buildServicesToVisiblePlansMap(orgName string) (map[string][]models.ServicePlanFields, error) {
// 	allPlans, err := builder.servicePlanRepo.Search(nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	servicesToPlansMap := make(map[string][]models.ServicePlanFields)
// 	PlanToOrgsVisMap, err := builder.buildPlanToOrgsVisibilityMap()
// 	if err != nil {
// 		return nil, err
// 	}
// 	OrgToPlansVisMap := builder.buildOrgToPlansVisibilityMap(PlanToOrgsVisMap)
// 	filterOrgPlans := OrgToPlansVisMap[orgName]

// 	for _, plan := range allPlans {
// 		if builder.containsGuid(filterOrgPlans, plan.Guid) {
// 			plan.OrgNames = PlanToOrgsVisMap[plan.Guid]
// 			servicesToPlansMap[plan.ServiceOfferingGuid] = append(servicesToPlansMap[plan.ServiceOfferingGuid], plan)
// 		} else if plan.Public {
// 			servicesToPlansMap[plan.ServiceOfferingGuid] = append(servicesToPlansMap[plan.ServiceOfferingGuid], plan)
// 		}
// 	}

// 	return servicesToPlansMap, nil
// }

// func (builder Builder) buildSingleServiceToVisiblePlansMap(serviceLabel string, orgName string) (map[string][]models.ServicePlanFields, error) {
// 	serviceToVisiblePlansMap, err := builder.BuildServicesToVisiblePlansMap(orgName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	service, err := builder.serviceRepo.FindServiceOfferingByLabel(serviceLabel)
// 	if err != nil {
// 		return nil, err
// 	}
// 	serviceToFilter, ok := serviceToVisiblePlansMap[service.Guid]
// 	if !ok {
// 		// Service is not visible to Org
// 		return nil, nil
// 	}
// 	serviceMap := make(map[string][]models.ServicePlanFields)
// 	serviceMap[service.Guid] = serviceToFilter
// 	return serviceMap, nil
// }

// func (builder Builder) containsGuid(guidSlice []string, guid string) bool {
// 	for _, g := range guidSlice {
// 		if g == guid {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (builder Builder) filterBrokerList(brokers []models.ServiceBroker, brokerName string) []models.ServiceBroker {
// 	for brokerIndex, _ := range brokers {
// 		broker := &brokers[brokerIndex]
// 		if broker.Name == brokerName {
// 			return []models.ServiceBroker{brokers[brokerIndex]}
// 		}
// 	}

// 	// Could not find brokerFlag in visible brokers.
// 	return nil
// }

// func (builder Builder) getAllBrokersFromServicesMap(brokersToServices map[string][]models.ServiceOffering) ([]models.ServiceBroker, error) {
// 	var brokers []models.ServiceBroker

// 	for brokerGuid, services := range brokersToServices {
// 		if brokerGuid == "" {
// 			continue
// 		}
// 		broker, err := builder.brokerRepo.FindByGuid(brokerGuid)
// 		if err != nil {
// 			return nil, err
// 		}

// 		broker.Services = services
// 		brokers = append(brokers, broker)
// 	}

// 	return brokers, nil
// }
