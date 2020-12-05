package main

import (
    ConfigValues "JiraAlert/Config"
    "JiraAlert/Util"
    "bytes"
    "encoding/json"
    "github.com/andygrunwald/go-jira"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "log"
    "net/http"
    "strconv"
    "time"
)

//Global Values
var knownIssues []string
var cv ConfigValues.ConfigValues
var markImmediatelyAsKnown bool //If set true, the next run will skip the alerting and mark an issue immediately as known. Will be auto reset.

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

    cv = ConfigValues.ConfigValues{}
    cv.LoadAndValidateConfig()

    markImmediatelyAsKnown = cv.DoInitialPost

    log.Println("Initialize application")
    tp := jira.BasicAuthTransport{
        Username: cv.JiraUsername,
        Password: cv.JiraPassword,
    }

    client, err := jira.NewClient(tp.Client(), cv.JiraUrl)
    if err != nil {
        log.Fatal(err)
    }

    filter, _, err := client.Filter.Get(cv.JiraFilterId)

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
    err = http.ListenAndServe(":"+strconv.Itoa(cv.PrometheusPort), nil)

    if err != nil {
        log.Fatal(err)
    }

    <-finished //Wait forever ;)
}


func heartBeat(finished chan bool, client *jira.Client, filter *jira.Filter) {
    for range time.Tick(time.Second * time.Duration(cv.JiraCheckInterval)) {
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
            if !Util.Contains(knownIssues, issue.Key) {
                //Marks the issues as known, without writing an alert.
                if !markImmediatelyAsKnown {
                    alerts = append(alerts, issue)
                } else {
                    markImmediatelyAsKnown = false
                }
                knownIssues = append(knownIssues, issue.Key)
            }
        }

        if prevNumberOfKnownIssues != len(knownIssues){
            log.Println("Number of known issues: " + strconv.Itoa(len(knownIssues)))
        }

        //Alert for new issues
        for _, issue := range alerts {

                message := MatterHook{
                    Text: ":rotating_light:  **" + issue.Fields.Priority.Name + "** " + issue.Key + " " + issue.Fields.Summary,
                }

                messageJson, _ := json.Marshal(message)
                req, err := http.NewRequest("POST", cv.WebhookUrl, bytes.NewBuffer(messageJson))
                req.Header.Set("Content-Type", "application/json")

                matterMostClient := &http.Client{}

                stopWatch = time.Now()

                resp, err := matterMostClient.Do(req)

                stopWatchTimeElapsed = time.Since(stopWatch)
                mattermostRequestDuration.Observe(float64(stopWatchTimeElapsed.Seconds()))

                if err != nil {
                    panic(err)
                } else {
                    mattermostCallsMade.Inc()
                }
                defer resp.Body.Close()
        }
    }

    //Will never be reached!
    log.Println("The cake is a lie and by the way: You should never have been able to get here! How did you do it?")
    finished <- true
}
