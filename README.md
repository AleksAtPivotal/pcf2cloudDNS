# Pivotal Cloud Foundry to Cloud DNS

POC code to integrate PCF (Pivotal Cloud Foundry) with a Cloud DNS by using the Open Service Broker API.

## Installation

Rename the `broker.json.template`

```sh
cp ./configs/broker.json.template ./configs/broker.json
```

Edit the broker.json updating at least the part with comments

```json
    "AppsDomain": "apps.sonoma.cf-app.com", // CF Apps domain
    "auth": {
        "username": "admin", 
        "password": "password" 
    },
    "CloudFoundry": {
        "username": "admin", // CF API admin username
        "password": "XBEZR7xXF37MTsqtwBpvjU7Pi-quBSbW" // CF API admin password 
    },
    "dns": {
        "Domain": "pcf.systems", // route53 domain name
        "HostedZoneID": "Z369O88AFHE8O1", // aws route53 hosted zone ID
        "TTL": 300
    },
```

Login to Pivotal Cloud Foundry and push the service broker and determine the service broker url

```sh
cf push
cf apps
```

Add `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` variables to the servicebroker application and re-stage the application as it will be making AWS API Calls to modify route53 zones.

Register the service broker

```sh
export BROKERNAME=router
export BROKER_USERNAME=admin
export BROKER_PASSWORD=password
cf create-service-broker $BROKERNAME $BROKER_USERNAME $BROKER_PASSWORD https://router-osb.apps.sonoma.cf-app.com
```

where `$BROKER_USERNAME` and `$BROKER_PASSWORD` matches the variables below from the `broker.json` file.

```json
    "auth": {
        "username": "admin", 
        "password": "password" 
    },
```

Create a shared domain for PCF that matches the `dns.Domain` value from `broker.json`.

```sh
cf create-shared-domain pcf.systems
```

