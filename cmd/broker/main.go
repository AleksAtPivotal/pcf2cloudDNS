package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/alekssaul/pcf2cloudDNS/pkg/config"
	"github.com/alekssaul/pcf2cloudDNS/pkg/servicebroker"

	"code.cloudfoundry.org/lager"
	"github.com/alekssaul/pcf2cloudDNS/pkg/utils"
	"github.com/pivotal-cf/brokerapi"
)

// Options contain the flags passed to the broker
type Options struct {
	BrokerConfigPath string
}

var options Options

func init() {
	defaultConfigPath := utils.GetPath([]string{"configs", "broker.json"})
	flag.StringVar(&options.BrokerConfigPath, "c", defaultConfigPath, "use '-c' option to specify the config file path")

	flag.Parse()
}

func main() {

	brokerLogger := lager.NewLogger("pcf2cloudDNS-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	brokerLogger.Info("Starting CF pcf2cloudDNS broker")

	brokerLogger.Info("Config File: " + options.BrokerConfigPath)

	// read our configuration
	configFile, err := os.Open(options.BrokerConfigPath)
	if err != nil {
		brokerLogger.Fatal("Loading config file", err, lager.Data{
			"broker-config-path": options.BrokerConfigPath,
		})
	}

	byteValue, _ := ioutil.ReadAll(configFile)
	var config config.BrokerConfig
	json.Unmarshal(byteValue, &config)

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: config.Authentication.Username,
		Password: config.Authentication.Password,
	}

	serviceBroker := &servicebroker.ServiceBroker{
		InstanceCreators: map[string]servicebroker.InstanceCreator{},
		InstanceBinders:  map[string]servicebroker.InstanceBinder{},
		Config:           config,
	}

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerCredentials)

	http.Handle("/", brokerAPI)

	brokerLogger.Fatal("http-listen", http.ListenAndServe(config.Host+":"+config.Port, nil))
}
