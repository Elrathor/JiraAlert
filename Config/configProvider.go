package ConfigProvider

import (
	"JiraAlert/Util"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"strconv"
)

type ConfigProvider struct {
	JiraUsername         string
	JiraUsernameKey      string
	JiraPassword         string
	JiraPasswordKey      string
	JiraFilterId         int
	JiraFilterIdKey      string
	JiraUrl              string
	JiraUrlKey           string
	JiraCheckInterval    int
	JiraCheckIntervalKey string
	WebhookUrl           string
	WebhookUrlKey        string
	PrometheusPort       int
	PrometheusPortKey    string
	DoInitialPost        bool
	logger               *zap.Logger
}

func NewConfigProvider(logger *zap.Logger) (element ConfigProvider) {
	element = ConfigProvider{
		logger: logger,
	}
	element.registerKeys()
	return element
}

func (cp *ConfigProvider) registerKeys() {
	cp.logger.Info("Registering config keys")
	cp.JiraFilterIdKey = "JIRA_FILTER_ID"
	cp.JiraUsernameKey = "JIRA_USERNAME"
	cp.JiraPasswordKey = "JIRA_PASSWORD"
	cp.JiraUrlKey = "JIRA_URL"
	cp.JiraCheckIntervalKey = "JIRA_CHECK_INTERVAL"
	cp.WebhookUrlKey = "WEBHOOK_ULR"
	cp.PrometheusPortKey = "PROMETHEUS_PORT"
}

func (cp *ConfigProvider) LoadAndValidateConfig() {
	var wasFound bool

	cp.logger.Info("Reading config")
	err := godotenv.Load()
	if err != nil {
		cp.logger.Fatal("Error loading .env file")
	}

	cp.logger.Info("Validating config values")
	cp.JiraUsername, wasFound = os.LookupEnv(cp.JiraUsernameKey)
	cp.validateString(cp.JiraUsernameKey, cp.JiraUsername, wasFound)

	cp.JiraPassword, wasFound = os.LookupEnv(cp.JiraPasswordKey)
	cp.validateString(cp.JiraPasswordKey, cp.JiraPassword, wasFound)

	cp.WebhookUrl, wasFound = os.LookupEnv(cp.WebhookUrlKey)
	cp.validateString(cp.WebhookUrlKey, cp.WebhookUrl, wasFound)

	cp.JiraUrl, wasFound = os.LookupEnv(cp.JiraUrlKey)
	cp.validateString(cp.JiraUrlKey, cp.JiraUrl, wasFound)

	filterIdString, wasFound := os.LookupEnv(cp.JiraFilterIdKey)
	cp.validateString(cp.JiraFilterIdKey, filterIdString, wasFound)

	cp.JiraFilterId, err = strconv.Atoi(filterIdString)
	if err != nil {
		cp.logger.Fatal("The given value has to be numeric", zap.String("key", cp.JiraFilterIdKey))
	}

	checkIntervalString, wasFound := os.LookupEnv(cp.JiraCheckIntervalKey)
	cp.validateString(cp.JiraCheckIntervalKey, checkIntervalString, wasFound)

	cp.JiraCheckInterval, err = strconv.Atoi(checkIntervalString)
	if err != nil {
		cp.logger.Fatal("The given value has to be numeric", zap.String("key", cp.JiraCheckIntervalKey))
	}

	prometheusPortString, wasFound := os.LookupEnv(cp.PrometheusPortKey)
	cp.validateString(cp.PrometheusPortKey, prometheusPortString, wasFound)

	cp.PrometheusPort, err = strconv.Atoi(prometheusPortString)
	if err != nil {
		cp.logger.Fatal("The given value has to be numeric", zap.String("key", cp.PrometheusPortKey))
	}

	cp.logger.Info("Reading command line arguments")
	args := os.Args[1:]
	if Util.Contains(args, "--NoInitialPost") {
		cp.DoInitialPost = false
		cp.logger.Info("Initial post to mattermost", zap.Bool("enabled", false))
	} else {
		cp.DoInitialPost = true
		cp.logger.Info("Initial post to mattermost", zap.Bool("enabled", true))
	}
}

func (cp *ConfigProvider) validateString(key string, value string, wasFound bool) {
	if !wasFound {
		cp.logger.Fatal("Expected config value not found in the .env file", zap.String("key", key))
	} else {
		cp.logger.Info("Read new config value", zap.String("key", key), zap.String("value", value))
	}
}
