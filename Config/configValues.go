package ConfigValues

import (
	"JiraAlert/Util"
	"github.com/joho/godotenv"
	"log"
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
}

func (cv *ConfigValues) registerKeys() {
	log.Println("Registering config keys")
	cv.JiraFilterIdKey = "JIRA_FILTER_ID"
	cv.JiraUsernameKey = "JIRA_USERNAME"
	cv.JiraPasswordKey = "JIRA_PASSWORD"
	cv.JiraUrlKey = "JIRA_URL"
	cv.JiraCheckIntervalKey = "JIRA_CHECK_INTERVAL"
	cv.WebhookUrlKey = "WEBHOOK_ULR"
	cv.PrometheusPortKey = "PROMETHEUS_PORT"
}

func (cv *ConfigValues) LoadAndValidateConfig() {
	cv.registerKeys()
	var wasFound bool

	log.Println("Reading config")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Validate config values")
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
		log.Fatal(cv.JiraFilterIdKey + " has to be a numeric value")
	}

	checkIntervalString, wasFound := os.LookupEnv(cv.JiraCheckIntervalKey)
	cv.validateString(cv.JiraCheckIntervalKey, checkIntervalString, wasFound)

	cv.JiraCheckInterval, err = strconv.Atoi(checkIntervalString)
	if err != nil {
		log.Fatal(cv.JiraCheckIntervalKey + " has to be a numeric value in seconds")
	}

	prometheusPortString, wasFound := os.LookupEnv(cv.PrometheusPortKey)
	cv.validateString(cv.PrometheusPortKey, prometheusPortString, wasFound)

	cv.PrometheusPort, err = strconv.Atoi(prometheusPortString)
	if err != nil {
		log.Fatal(cv.PrometheusPortKey + " has to be a numeric value in seconds")
	}

	log.Println("Reading command line arguments")
	args := os.Args[1:]
	if Util.Contains(args, "--NoInitialPost") {
		cv.DoInitialPost = false
		log.Println("Initial post to mattermost: disabled")
	} else {
		cv.DoInitialPost = true
		log.Println("Initial post to mattermost: enabled")
	}
}

func (cv *ConfigValues) validateString(key string, value string, wasFound bool) {
	if !wasFound {
		log.Fatal(key + " has to be present in the .env file")
	} else {
		log.Println(key + ": " + value)
	}
}
