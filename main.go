package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"os"
	"strconv"
	"strings"

	"github.com/andygrunwald/go-jira"
)

func main() {

  jiraUrl := os.Getenv("JIRA_URL")
  jiraUsername := os.Getenv("JIRA_USERNAME")
  jiraPassword := os.Getenv("JIRA_PASSWORD")
  jiraIssue    := "DEVOPS-1880"
  jiraIssueLabel := "devops-automation"
  // jiraRelease := os.Getenv("JIRA_RELEASE")
  // jiraJQL := fmt.Sprintf("fixVersion = '%s'", jiraRelease)

  jiraAuth := jira.BasicAuthTransport{
		Username: jiraUsername,
		Password: jiraPassword,
	}

	jiraClient, err := jira.NewClient(jiraAuth.Client(), jiraUrl)

  if err != nil {
    fmt.Println(err)
  }

  setIssueLabel(jiraClient, jiraIssue, jiraIssueLabel)
  // getSprintStats(jiraClient, jiraJQL, jiraRelease)
  // createIssue(jiraClient, "DEVOPS", "Task", "Summary here", "Description here")

}

func createIssue(jiraClient *jira.Client, jiraProject string, issueType string, issueSummary string, issueDescription string) {
  jiraIssue := jira.Issue{
		Fields: &jira.IssueFields{
			Description: issueDescription,
			Type: jira.IssueType{
				Name: issueType,
			},
			Project: jira.Project{
				Key: jiraProject,
			},
			Summary: issueSummary,
		},
	}
	issue, _, err := jiraClient.Issue.Create(&jiraIssue)
	if err != nil {
		panic(err)
	}

	fmt.Printf(issue.Key)
}

func setIssueLabel(jiraClient *jira.Client, jiraIssue string, issueLabel string) {

  type Labels struct {
    Add string `json:"add"`
  }
  type Update struct {
    Labels []Labels `json:"labels"`
  }
  // type Body struct {
  //   Update Update `json:"update"`
  // }

  // labelsList := make([]string, 1)
  // labelsList[0] = issueLabel
  labelsStruct := Labels{Add: issueLabel}
  labelsStructList := make([]Labels, 1)
  labelsStructList[0] = labelsStruct
  updateStruct := Update{Labels: labelsStructList}
  // bodyStruct := Body{Update: updateStruct}
  sss, _ := json.Marshal(updateStruct) //bodyStruct)
  fmt.Println(string(sss))


  data := map[string]interface{}{
    "update": updateStruct,
  }

  response, err := jiraClient.Issue.UpdateIssue(jiraIssue, data)
  resp_body, _ := ioutil.ReadAll(response.Body)
  fmt.Println("RESP: ", string(resp_body), "ERR: ",err)
}

func getSprintStats(jiraClient *jira.Client, jiraJQL string, jiraRelease string) {
  issues, _, err := jiraClient.Issue.Search(jiraJQL, &jira.SearchOptions{})
  var storyPoints int

  if err != nil {
    fmt.Println(err)
  }

  for _, issue := range issues {
    customFields, _, _ := jiraClient.Issue.GetCustomFields(issue.ID)
    storyPoint, _ := strconv.Atoi(customFields["customfield_10025"])
    storyPoints += storyPoint
    sprintName := getSprintName(customFields["customfield_10020"], "name:", " startDate:")
    fmt.Println(sprintName, issue.Fields.Type.Name, issue.Key, issue.Fields.Summary, storyPoint)
  }

  fmt.Printf("\nRelease stats\n=====================\nName: %s\nTasks: %d\nStory Points: %d\n",
            jiraRelease,
            len(issues),
            storyPoints,
          )
}

func getSprintName(customField string, startString string, endString string) (result string) {
    s := strings.Index(customField, startString)
    if s == -1 {
        return result
    }
    newString := customField[s+len(startString):]
    e := strings.Index(newString, endString)
    if e == -1 {
        return result
    }
    result = newString[:e]
    return result
}

