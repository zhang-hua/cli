package models

type ServicePlanFields struct {
	Guid        string
	Name        string
	Free        bool
	Description string
	Public      bool
	Active      bool
	OrgNames    []string
}

type ServicePlan struct {
	ServicePlanFields
	ServiceOffering ServiceOfferingFields
}
