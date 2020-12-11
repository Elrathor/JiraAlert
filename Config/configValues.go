package ConfigValues

import (
	"JiraAlert/Util"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"strconv"
)

type ConfigValues struct {
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

func NewConfigValues(logger *zap.Logger) (element ConfigValues) {
	element = ConfigValues{
		logger: logger,
	}
	element.registerKeys()
	return element
}

func (cv *ConfigValues) registerKeys() {
	cv.logger.Info("Registering config keys")
	cv.JiraFilterIdKey = "JIRA_FILTER_ID"
	cv.JiraUsernameKey = "JIRA_USERNAME"
	cv.JiraPasswordKey = "JIRA_PASSWORD"
	cv.JiraUrlKey = "JIRA_URL"
	cv.JiraCheckIntervalKey = "JIRA_CHECK_INTERVAL"
	cv.WebhookUrlKey = "WEBHOOK_ULR"
	cv.PrometheusPortKey = "PROMETHEUS_PORT"
}

func (cv *ConfigValues) LoadAndValidateConfig() {
	var wasFound bool

	cv.logger.Info("Reading config")
	err := godotenv.Load()
	if err != nil {
		cv.logger.Fatal("Error loading .env file")
	}

	cv.logger.Info("Validating config values")
	cv.JiraUsername, wasFound = os.LookupEnv(cv.JiraUsernameKey)
	cv.validateString(cv.JiraUsernameKey, cv.JiraUsername, wasFound)

	cv.JiraPassword, wasFound = os.LookupEnv(cv.JiraPasswordKey)
	cv.validateString(cv.JiraPasswordKey, cv.JiraPassword, wasFound)

	cv.WebhookUrl, wasFound = os.LookupEnv(cv.WebhookUrlKey)
	cv.validateString(cv.WebhookUrlKey, cv.WebhookUrl, wasFound)

	cv.JiraUrl, wasFound = os.LookupEnv(cv.JiraUrlKey)
	cv.validateString(cv.JiraUrlKey, cv.JiraUrl, wasFound)

	filterIdString, wasFound := os.LookupEnv(cv.JiraFilterIdKey)
	cv.validateString(cv.JiraFilterIdKey, filterIdString, wasFound)

	cv.JiraFilterId, err = strconv.Atoi(filterIdString)
	if err != nil {
		cv.logger.Fatal("The given value has to be numeric", zap.String("key", cv.JiraFilterIdKey))
	}

	checkIntervalString, wasFound := os.LookupEnv(cv.JiraCheckIntervalKey)
	cv.validateString(cv.JiraCheckIntervalKey, checkIntervalString, wasFound)

	cv.JiraCheckInterval, err = strconv.Atoi(checkIntervalString)
	if err != nil {
		cv.logger.Fatal("The given value has to be numeric", zap.String("key", cv.JiraCheckIntervalKey))
	}

	prometheusPortString, wasFound := os.LookupEnv(cv.PrometheusPortKey)
	cv.validateString(cv.PrometheusPortKey, prometheusPortString, wasFound)

	cv.PrometheusPort, err = strconv.Atoi(prometheusPortString)
	if err != nil {
		cv.logger.Fatal("The given value has to be numeric", zap.String("key", cv.PrometheusPortKey))
	}

	cv.logger.Info("Reading command line arguments")
	args := os.Args[1:]
	if Util.Contains(args, "--NoInitialPost") {
		cv.DoInitialPost = false
		cv.logger.Info("Initial post to mattermost", zap.Bool("enabled", false))
	} else {
		cv.DoInitialPost = true
		cv.logger.Info("Initial post to mattermost", zap.Bool("enabled", true))
	}
}

func (cv *ConfigValues) validateString(key string, value string, wasFound bool) {
	if !wasFound {
		cv.logger.Fatal("Expected config value not found in the .env file", zap.String("key", key))
	} else {
		cv.logger.Info("Read new config value", zap.String("key", key), zap.String("value", value))
	}
}
