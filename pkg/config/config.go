package config

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Providers             []Provider    `yaml:"providers"`
	EnableMonitorDeletion bool          `yaml:"enableMonitorDeletion"`
	MonitorNameTemplate   string        `yaml:"monitorNameTemplate"`
	ResyncPeriod          int           `yaml:"resyncPeriod,omitempty"`
	CreationDelay         time.Duration `yaml:"creationDelay,omitempty"`
}

// UnmarshalYAML interface to deserialize specific types
func (c *Config) UnmarshalYAML(data []byte) error {
	type Alias Config
	aux := struct {
		CreationDelay string `yaml:"creationDelay,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := yaml.Unmarshal(data, &aux); err != nil {
		return err
	}

	delay, err := time.ParseDuration(aux.CreationDelay)
	if err != nil {
		return err
	}
	c.CreationDelay = delay

	return nil
}

type Provider struct {
	Name              string      `yaml:"name"`
	ApiKey            string      `yaml:"apiKey"`
	ApiURL            string      `yaml:"apiURL"`
	AlertContacts     string      `yaml:"alertContacts"`
	AlertIntegrations string      `yaml:"alertIntegrations"`
	Username          string      `yaml:"username"`
	Password          string      `yaml:"password"`
	AccountEmail      string      `yaml:"accountEmail"`
	AppInsightsConfig AppInsights `yaml:"appInsightsConfig"`
}

type AppInsights struct {
	Name          string        `yaml:"name"`
	Location      string        `yaml:"location"`
	ResourceGroup string        `yaml:"resourceGroup"`
	Frequency     int32         `yaml:"frequency"`
	GeoLocation   []interface{} `yaml:"geoLocation"`
	EmailAction   EmailAction   `yaml:"emailAction"`
	WebhookAction WebhookAction `yaml:"webhookAction"`
}

type EmailAction struct {
	SendToServiceOwners bool     `yaml:"send_to_service_owners"`
	CustomEmails        []string `yaml:"custom_emails"`
}

type WebhookAction struct {
	ServiceURI string `yaml:"service_uri"`
}

func ReadConfig(filePath string) Config {
	var config Config
	// Read YML
	log.Info("Reading YAML Configuration", filePath)
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	// Unmarshall
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Panic(err)
	}

	return config
}

func GetControllerConfig() Config {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		configFilePath = "../../configs/testConfigs/test-config.yaml"
	}

	config := ReadConfig(configFilePath)

	return config
}
