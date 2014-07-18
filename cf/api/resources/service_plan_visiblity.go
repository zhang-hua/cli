package resources

import "github.com/cloudfoundry/cli/cf/models"

type ServicePlanVisibilityResource struct {
	Resource
	Entity ServicePlanVisibilityEntity
}

type ServicePlanVisibilityEntity struct {
	ServicePlanGuid  string
	OrganizationGuid string
}

func (resource ServicePlanVisibilityResource) ToFields() (fields models.ServicePlanVisibilityFields) {
	fields.Guid = resource.Metadata.Guid
	fields.ServicePlanResource = resources.Entity.ServicePlanGuid
	fields.OrganizationGuid = resources.Entity.OrganizationGuid
	return
}
