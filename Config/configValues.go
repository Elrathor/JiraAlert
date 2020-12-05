package ConfigValues

import (
    "github.com/joho/godotenv"
    "log"
    "os"
    "strconv"
)

type ConfigValues struct {
    JiraUsername string
    JiraUsernameKey string
    JiraPassword string
    JiraPasswordKey string
    JiraFilterId int
    JiraFilterIdKey string
    JiraUrl string
    JiraUrlKey string
    JiraCheckInterval int
    JiraCheckIntervalKey string
    WebhookUrl string
    WebhookUrlKey string
}

func (cv *ConfigValues) registerKeys() {
    log.Println("Registering config keys")
    cv.JiraFilterIdKey = "JIRA_FILTER_ID"
    cv.JiraUsernameKey = "JIRA_USERNAME"
    cv.JiraPasswordKey = "JIRA_PASSWORD"
    cv.JiraUrlKey = "JIRA_URL"
    cv.JiraCheckIntervalKey = "JIRA_CHECK_INTERVAL"
    cv.WebhookUrlKey = "WEBHOOK_ULR"
}

func (cv *ConfigValues) LoadAndValidateConfig() {
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
}

func (cv ConfigValues) validateString(key string, value string, wasFound bool) {
    if !wasFound {
        log.Fatal(key + " has to be present in the .env file")
    } else {
        log.Println(key + ": " + value )
    }
}

