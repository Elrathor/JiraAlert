# JiraAlert
This tool offers a simple but monitored way to check a given filter for issues and send a webhook message to your mattermost
server if a new one appears in the list.

**Example usecase:** Monitoring and alerting for high and medium priority tickets within a Jira Servicedesk queue.

If there's a feature missing for you, just let me know via a [feature request](https://github.com/Elrathor/JiraAlert/issues/new/choose)

## Assumptions
Since this tool is currently still under development, some assumptions had to be made:
1. You use a current, on premise version of Jira which still supports password based api login
2. You have a user able to use the Jira api
3. You want to be notified via a message in a Mattermost channel
4. You are allowed to query your Jira api every ~20s

Most of the assumptions above are subject to change, when those values will become configurable. 

## Configuration
The following values have to be present inside a .env file inside the application directory.

|Key|Value|
|---|---|
|JIRA_USERNAME|The username the application will use to query Jira (required)|
|JIRA_PASSWORD|The corresponding password the application will use to query Jira (required)|
|JIRA_FILTER_ID|The numeric ID of the filter the application will query (required)|
|JIRA_URL|The base url of the Jira instance that should be queried (required)|
|JIRA_CHECK_INTERVAL|The numeric interval in seconds in which Jira will be queried (required)|
|WEBHOOK_ULR|The mattermost webhook that should be notified when an alert happens (required)|
|PROMETHEUS_PORT|The numeric port on which the Prometheus endpoint will be available (use 2112 if unsure) (required)|


## Installation
At the current point in time no prebuild packages are offered, but you can build the application from source.
1. Checkout the project via git
1. Install a current version of Go (latest version testet: go1.15.5)
2. Navigate your terminal into the project directory 
3. Install all dependencies via `go get`
4. Build the application via `go build`
5. Start the new executable generated by the go compiler :)

## Monitoring
By default, the application offers an `/metrics` endpoint at port 2112. (http://localhost:2112/metrics) This endpoint can
be used by a prometheus server to keep taps on this application.

## Roadmap
The following (in no particular order) things _should_ be done to improve this application:
- [x] Add more configurable values
- [ ] Make the application able to run inside a container (Docker)
- [ ] Auto recovery when connection to Mattermost or Jira failed
- [ ] Add Slack support
- [ ] Add support for OAuth / Token based authentication for jira
- [ ] publish issue Counter via Prometheus to display it on Graphana dashboards
- [ ] Split the application in smaller chunks
- [ ] Buy Elrathor some cookies
