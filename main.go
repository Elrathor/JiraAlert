package main

import (
    "bytes"
    "encoding/json"
    "github.com/andygrunwald/go-jira"
    "github.com/joho/godotenv"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"
)

//Global Values
var knownIssues []string
var url string

//Prometheus Metrics
var (
    jiraRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
        Name: "jiraalert_jira_request_duration",
        Help: "The time it took to query the jira api",
    })
)

var (
    jiraCallsMade = promauto.NewCounter(prometheus.CounterOpts{
        Name: "jiraalert_jira_calls_made",
        Help: "The total number of requests made to the jira api",
    })
)

var (
    mattermostRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
        Name: "jiraalert_mattermost_request_duration",
        Help: "The time it took to send the mattermost webhook request",
    })
)

var (
    mattermostCallsMade = promauto.NewCounter(prometheus.CounterOpts{
        Name: "jiraalert_mattermost_calls_made",
        Help: "The total number of requests made to the mattermost webhook",
    })
)

// Struct for json marshaling
type MatterHook struct {
    Text string
}

func main() {
    // used to determine if a value was found in the config file
    var wasFound bool

    log.Println("Reading config")
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    log.Println("Validate config values")
    username, wasFound := os.LookupEnv("JIRA_USERNAME")
    if !wasFound {
        log.Fatal("JIRA_USERNAME has to be present in the .env file")
    } else {
        log.Println("Username: " + username)
    }

    password, wasFound := os.LookupEnv("JIRA_PASSWORD")
    if !wasFound {
        log.Fatal("JIRA_PASSWORD has to be present in the .env file")
    } else {
        log.Println("Password: <Secret>" )
    }

    url, wasFound = os.LookupEnv("WEBHOOK_ULR")
    if !wasFound {
        log.Fatal("WEBHOOK_ULR has to be present in the .env file")
    } else {
        log.Println("Webhook URL: " + url)
    }

    jiraUrl, wasFound := os.LookupEnv("JIRA_URL")
    if !wasFound {
        log.Fatal("JIRA_URL has to be present in the .env file")
    } else {
        log.Println("Jira URL: " + jiraUrl)
    }

    filterIdString, wasFound := os.LookupEnv("JIRA_FILTER_ID")
    if !wasFound {
        log.Fatal("JIRA_FILTER_ID has to be present in the .env file")
    }

    filterID, err := strconv.Atoi(filterIdString)
    if err != nil {
        log.Fatal("JIRA_FILTER_ID has to be a numeric value")
    } else {
        log.Println("Filter ID: " + filterIdString)
    }

    log.Println("Initialize application")
    tp := jira.BasicAuthTransport{
        Username: username,
        Password: password,
    }

    client, err := jira.NewClient(tp.Client(), jiraUrl)
    if err != nil {
        log.Fatal(err)
    }

    filter, _, err := client.Filter.Get(filterID)

    if err != nil {
        log.Fatal(err)
    } else {
        log.Println("Using filter: " + filter.Name)
    }

    log.Println("Initialize monitoring")
    http.Handle("/metrics", promhttp.Handler())

    log.Println("Start watcher")
    finished := make(chan bool)
    go heartBeat(finished, client, filter)

    log.Println("Starting monitoring")
    err = http.ListenAndServe(":2112", nil)

    if err != nil {
        log.Fatal(err)
    }

    <-finished //Wait forever ;)
}

// https://play.golang.org/p/Qg_uv_inCek
// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
    for _, v := range s {
        if v == str {
            return true
        }
    }

    return false
}

func heartBeat(finished chan bool, client *jira.Client, filter *jira.Filter) {
    for range time.Tick(time.Second * 20) {
        stopWatch := time.Now()

        issues, _, err := client.Issue.Search(filter.Jql, nil)

        stopWatchTimeElapsed := time.Since(stopWatch)
        jiraRequestDuration.Observe(stopWatchTimeElapsed.Seconds())

        if err != nil {
            panic(err)
        } else {
            jiraCallsMade.Inc()
        }

        var alerts []jira.Issue
        prevNumberOfKnownIssues := len(knownIssues)

        //Check if issue is already know if not set it to alert list and mark as know
        for _, issue := range issues {
            if !contains(knownIssues, issue.Key) && issue.Fields.Priority.ID != "4" { // 4 = low
                alerts = append(alerts, issue)
                knownIssues = append(knownIssues, issue.Key)
            }
        }

        if prevNumberOfKnownIssues != len(knownIssues){
            log.Println("Number of known issues: " + strconv.Itoa(len(knownIssues)))
        }

        //Alert for new issues
        for _, issue := range alerts {

            if issue.Fields.Priority.ID != "4" {

                message := MatterHook{
                    Text: ":rotating_light:  **" + issue.Fields.Priority.Name + "** " + issue.Key + " " + issue.Fields.Summary,
                }

                messageJson, _ := json.Marshal(message)
                req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageJson))
                req.Header.Set("Content-Type", "application/json")

                matterMostClient := &http.Client{}

                stopWatch = time.Now()

                resp, err := matterMostClient.Do(req)

                stopWatchTimeElapsed = time.Since(stopWatch)
                mattermostRequestDuration.Observe(float64(stopWatchTimeElapsed))

                if err != nil {
                    panic(err)
                } else {
                    mattermostCallsMade.Inc()
                }
                defer resp.Body.Close()
            }
        }
    }

    //Will never be reached!
    log.Println("The cake is a lie and by the way: You should never have been able to get here! How did you do it?")
    finished <- true
}
