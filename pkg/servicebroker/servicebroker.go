package servicebroker

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/alekssaul/pcf2cloudDNS/pkg/dns"

	"github.com/alekssaul/pcf2cloudDNS/pkg/cf"

	"github.com/alekssaul/pcf2cloudDNS/pkg/config"
	"github.com/pivotal-cf/brokerapi"
)

// InstanceCreator implements interface to create the service instance
type InstanceCreator interface {
	Create(instanceID string) error
	Destroy(instanceID string) error
	InstanceExists(instanceID string) (bool, error)
}

// InstanceBinder implements interface to bind to service instance
type InstanceBinder interface {
	Bind(instanceID string, bindingID string) (InstanceCredentials, error)
	Unbind(instanceID string, bindingID string) error
	InstanceExists(instanceID string) (bool, error)
}

// ServiceBroker - Main struct for ServiceBroker
type ServiceBroker struct {
	InstanceCreators map[string]InstanceCreator
	InstanceBinders  map[string]InstanceBinder
	Config           config.BrokerConfig
	Host             string
}

// InstanceCredentials - This service provider does not provide Instance Credentials
type InstanceCredentials struct {
}

// Services - Implements the Service Catalog
func (serviceBroker *ServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	planList := []brokerapi.ServicePlan{}
	for _, plan := range serviceBroker.Config.Plans {
		planList = append(planList, plan)
	}

	return []brokerapi.Service{
		brokerapi.Service{
			ID:          serviceBroker.Config.Catalog.ID,
			Name:        serviceBroker.Config.Catalog.Name,
			Description: serviceBroker.Config.Catalog.Description,
			Bindable:    *serviceBroker.Config.Catalog.Bindable,
			Plans:       planList,
			Metadata:    &serviceBroker.Config.Metadata,
			Tags: []string{
				"globalrouter",
				"dns",
			},
		},
	}, nil
}

// Bind - Binds to the service instane
func (serviceBroker *ServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	var binding brokerapi.Binding
	log.Println("Bind a new Service")

	err := cf.Bind(serviceBroker.Host, details.BindResource.AppGuid, details.BindResource.SpaceGuid, serviceBroker.Config.CloudFoundry.Username, serviceBroker.Config.CloudFoundry.Password, serviceBroker.Config.DNS.Domain)
	if err != nil {
		//return binding, err
		log.Printf("Error: %s", err)
		return binding, nil
	}

	return binding, nil
}

// Provision - Provisions the service instance
func (serviceBroker *ServiceBroker) Provision(ctx context.Context, instanceID string, serviceDetails brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	log.Println("Provisioning a new Service")

	var params config.ServiceParameters
	json.Unmarshal(serviceDetails.RawParameters, &params)
	if params.Host == "" {
		return spec, errors.New("\"host\" parameter must be passed")
	}

	spec = brokerapi.ProvisionedServiceSpec{}
	if serviceDetails.PlanID == "" {
		return spec, errors.New("plan_id required")
	}

	err = dns.ChangeAWSRecord("CREATE",
		params.Host+"."+serviceBroker.Config.DNS.Domain,
		serviceBroker.Config.DNS.HostedZoneID,
		"*."+serviceBroker.Config.AppsDomain,
		serviceBroker.Config.DNS.TTL,
		serviceBroker.Config.DNS.AWS.AccessKeyID,
		serviceBroker.Config.DNS.AWS.SecretAccessKey)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	serviceBroker.Host = params.Host

	return spec, nil
}

// Deprovision - Deprovisions the service instance
func (serviceBroker *ServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	log.Println("DeProvisioning Service")

	spec := brokerapi.DeprovisionServiceSpec{}
	if details.PlanID == "" {
		return spec, errors.New("plan_id required")
	}
	/*
		err := dns.ChangeAWSRecord("DELETE",
			// params.Host+"."+serviceBroker.Config.DNS.SharedDomain,
			"foo."+serviceBroker.Config.DNS.Domain,
			serviceBroker.Config.DNS.HostedZoneID,
			"*." + serviceBroker.Config.AppsDomain,
			serviceBroker.Config.DNS.TTL,
			serviceBroker.Config.DNS.AWS.AccessKeyID,
			serviceBroker.Config.DNS.AWS.SecretAccessKey)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	*/
	log.Printf("PlanID: %v\n", details.PlanID)
	log.Printf("ServiceID: %v\n", details.ServiceID)

	return spec, nil
}

// GetBinding - Not implemented
func (serviceBroker *ServiceBroker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	return brokerapi.GetBindingSpec{}, errors.New("GetBinding not implemented")
}

// GetInstance - Not implemented
func (serviceBroker *ServiceBroker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	return brokerapi.GetInstanceDetailsSpec{}, errors.New("GetInstance not implemented")
}

// LastBindingOperation - Not implemented
func (serviceBroker *ServiceBroker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, errors.New("LastBindingOperation not implemented")
}

// LastOperation - Always succeds for now
func (serviceBroker *ServiceBroker) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{
		State:       "succeeded",
		Description: "Always Success for now",
	}, nil
	/*
		return brokerapi.LastOperation{}, errors.New("LastOperation not implemented")
	*/
}

// Unbind - Not implemented
func (serviceBroker *ServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	return brokerapi.UnbindSpec{}, nil
}

// Update - Not implemented
func (serviceBroker *ServiceBroker) Update(cxt context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, errors.New("Update not implemented")
}
