package config

import (
	"github.com/pivotal-cf/brokerapi"
)

// BrokerConfig holds the main configuration of the broker
type BrokerConfig struct {
	Port           string                    `json:"port"`
	Host           string                    `json:"host"`
	AppsDomain     string                    `json:"AppsDomain"`
	CloudFoundry   AuthConfiguration         `json:"CloudFoundry"`
	Authentication AuthConfiguration         `json:"auth"`
	Catalog        brokerapi.ServicePlan     `json:"service"`
	Metadata       brokerapi.ServiceMetadata `json:"service_metadata"`
	Plans          []brokerapi.ServicePlan   `json:"plans"`
	DNS            DNSConfig                 `json:"dns"`
}

// AuthConfiguration contains ServiceBroker http basic auth creds
type AuthConfiguration struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// ServiceParameters - contains payload we render from service broker parameters
type ServiceParameters struct {
	Host string `json:"host"`
}

// DNSConfig - Holds the DNS Configuration
type DNSConfig struct {
	Domain       string               `json:"Domain"`
	HostedZoneID string               `json:"HostedZoneID"`
	AWS          AWSAuthConfiguration `json:"AWS"`
	TTL          int64                `json:"TTL"`
}

// AWSAuthConfiguration contains authentication information for AWS
type AWSAuthConfiguration struct {
	AccessKeyID     string `json:"aws_access_key_id"`
	SecretAccessKey string `json:"aws_secret_access_key"`
}
