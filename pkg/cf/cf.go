package cf

import (
	"log"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/go-cfenv"
)

// Bind - Bind will add a route to the App published on CF
func Bind(routehost string, appGUID string, spaceGUID string, Username string, Password string) (err error) {
	searchdomain := "pcf.systems"

	log.Println("Running CF Bind workflow")
	appEnv, _ := cfenv.Current()

	c := &cfclient.Config{
		ApiAddress:        appEnv.CFAPI,
		Username:          Username,
		Password:          Password,
		SkipSslValidation: true,
	}

	client, err := cfclient.NewClient(c)
	if err != nil {
		return err
	}

	SharedDomains, err := client.ListSharedDomains()
	if err != nil {
		return err
	}

	var sharedDomain cfclient.SharedDomain
	for _, domain := range SharedDomains {
		if domain.Name == searchdomain {
			sharedDomain = domain
		}
	}

	spaceroute, err := client.CreateRoute(cfclient.RouteRequest{
		DomainGuid: sharedDomain.Guid,
		SpaceGuid:  spaceGUID,
		Host:       routehost,
	})
	if err != nil {
		return err
	}

	_, err = client.MappingAppAndRoute(cfclient.RouteMappingRequest{
		RouteGUID: spaceroute.Guid,
		AppGUID:   appGUID,
		AppPort:   8080,
	})
	if err != nil {
		return err
	}

	return nil
}
