package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/andygrunwald/go-jira"
)

func main() {

  jiraUrl := os.Getenv("JIRA_URL")
  jiraUsername := os.Getenv("JIRA_USERNAME")
  jiraPassword := os.Getenv("JIRA_PASSWORD")
  jiraRelease := os.Getenv("JIRA_RELEASE")
  jiraJQL := fmt.Sprintf("fixVersion = '%s'", jiraRelease)

  jiraAuth := jira.BasicAuthTransport{
		Username: jiraUsername,
		Password: jiraPassword,
	}

	jiraClient, err := jira.NewClient(jiraAuth.Client(), jiraUrl)

  if err != nil {
    fmt.Println(err)
  }

  issues, _, err := jiraClient.Issue.Search(jiraJQL, &jira.SearchOptions{})
  var storyPoints int

  if err != nil {
    fmt.Println(err)
  }

  for _, issue := range issues {
    customFields, _, _ := jiraClient.Issue.GetCustomFields(issue.ID)
    storyPoint, _ := strconv.Atoi(customFields["customfield_10025"])
    storyPoints += storyPoint
    sprintStat := strings.Fields(strings.Split(customFields["customfield_10020"], "name:")[1])
    sprint := sprintStat[0] + " " + sprintStat[1]
    fmt.Println(sprint, issue.Fields.Type.Name, issue.Key, issue.Fields.Summary, storyPoint)
  }

  fmt.Printf("\nRelease stats\n=====================\nName: %s\nTasks: %d\nStory Points: %d\n",
            jiraRelease,
            len(issues),
            storyPoints,
          )

}
