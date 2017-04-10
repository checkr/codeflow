package github

import (
	"fmt"
	"time"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
)

const plugin_name = "github_webhook"

type Event interface {
	NewEvent() agent.Event
}

type Repository struct {
	Repository string `json:"full_name"`
	Private    bool   `json:"private"`
	Stars      int    `json:"stargazers_count"`
	Forks      int    `json:"forks_count"`
	Issues     int    `json:"open_issues_count"`
	SshUrl     string `json:"ssh_url"`
	HttpsUrl   string `json:"https_url"`
}

type Sender struct {
	User  string `json:"login"`
	Admin bool   `json:"site_admin"`
}

type CommitComment struct {
	Commit string `json:"commit_id"`
	Body   string `json:"body"`
}

type Deployment struct {
	Commit      string `json:"sha"`
	Task        string `json:"task"`
	Environment string `json:"environment"`
	Description string `json:"description"`
}

type Page struct {
	Name   string `json:"page_name"`
	Title  string `json:"title"`
	Action string `json:"action"`
}

type Issue struct {
	Number   int    `json:"number"`
	Title    string `json:"title"`
	Comments int    `json:"comments"`
}

type IssueComment struct {
	Body string `json:"body"`
}

type Team struct {
	Name string `json:"name"`
}

type PullRequest struct {
	Number       int    `json:"number"`
	State        string `json:"state"`
	Title        string `json:"title"`
	Comments     int    `json:"comments"`
	Commits      int    `json:"commits"`
	Additions    int    `json:"additions"`
	Deletions    int    `json:"deletions"`
	ChangedFiles int    `json:"changed_files"`
}

type PullRequestReviewComment struct {
	File    string `json:"path"`
	Comment string `json:"body"`
}

type Release struct {
	TagName string `json:"tag_name"`
}

type DeploymentStatus struct {
	State       string `json:"state"`
	Description string `json:"description"`
}

type HeadCommit struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type CommitCommentEvent struct {
	Comment    CommitComment `json:"comment"`
	Repository Repository    `json:"repository"`
	Sender     Sender        `json:"sender"`
}

func (s CommitCommentEvent) NewEvent() agent.Event {
	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
		"commit":     s.Comment.Commit,
		"comment":    s.Comment.Body,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type CreateEvent struct {
	Ref        string     `json:"ref"`
	RefType    string     `json:"ref_type"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s CreateEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
		"ref":        s.Ref,
		"refType":    s.RefType,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type DeleteEvent struct {
	Ref        string     `json:"ref"`
	RefType    string     `json:"ref_type"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s DeleteEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
		"ref":        s.Ref,
		"refType":    s.RefType,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type DeploymentEvent struct {
	Deployment Deployment `json:"deployment"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s DeploymentEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository":  s.Repository.Repository,
		"private":     fmt.Sprintf("%v", s.Repository.Private),
		"user":        s.Sender.User,
		"admin":       fmt.Sprintf("%v", s.Sender.Admin),
		"stars":       s.Repository.Stars,
		"forks":       s.Repository.Forks,
		"issues":      s.Repository.Issues,
		"commit":      s.Deployment.Commit,
		"task":        s.Deployment.Task,
		"environment": s.Deployment.Environment,
		"description": s.Deployment.Description,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type DeploymentStatusEvent struct {
	Deployment       Deployment       `json:"deployment"`
	DeploymentStatus DeploymentStatus `json:"deployment_status"`
	Repository       Repository       `json:"repository"`
	Sender           Sender           `json:"sender"`
}

func (s DeploymentStatusEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository":     s.Repository.Repository,
		"private":        fmt.Sprintf("%v", s.Repository.Private),
		"user":           s.Sender.User,
		"admin":          fmt.Sprintf("%v", s.Sender.Admin),
		"stars":          s.Repository.Stars,
		"forks":          s.Repository.Forks,
		"issues":         s.Repository.Issues,
		"commit":         s.Deployment.Commit,
		"task":           s.Deployment.Task,
		"environment":    s.Deployment.Environment,
		"description":    s.Deployment.Description,
		"depState":       s.DeploymentStatus.State,
		"depDescription": s.DeploymentStatus.Description,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type ForkEvent struct {
	Forkee     Repository `json:"forkee"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s ForkEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
		"fork":       s.Forkee.Repository,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type GollumEvent struct {
	Pages      []Page     `json:"pages"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

// REVIEW: Going to be lazy and not deal with the pages.
func (s GollumEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type IssueCommentEvent struct {
	Issue      Issue        `json:"issue"`
	Comment    IssueComment `json:"comment"`
	Repository Repository   `json:"repository"`
	Sender     Sender       `json:"sender"`
}

func (s IssueCommentEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"issue":      fmt.Sprintf("%v", s.Issue.Number),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
		"title":      s.Issue.Title,
		"comments":   s.Issue.Comments,
		"body":       s.Comment.Body,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type IssuesEvent struct {
	Action     string     `json:"action"`
	Issue      Issue      `json:"issue"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s IssuesEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"issue":      fmt.Sprintf("%v", s.Issue.Number),
		"action":     s.Action,
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
		"title":      s.Issue.Title,
		"comments":   s.Issue.Comments,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type MemberEvent struct {
	Member     Sender     `json:"member"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s MemberEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository":      s.Repository.Repository,
		"private":         fmt.Sprintf("%v", s.Repository.Private),
		"user":            s.Sender.User,
		"admin":           fmt.Sprintf("%v", s.Sender.Admin),
		"stars":           s.Repository.Stars,
		"forks":           s.Repository.Forks,
		"issues":          s.Repository.Issues,
		"newMember":       s.Member.User,
		"newMemberStatus": s.Member.Admin,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type MembershipEvent struct {
	Action string `json:"action"`
	Member Sender `json:"member"`
	Sender Sender `json:"sender"`
	Team   Team   `json:"team"`
}

func (s MembershipEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"user":            s.Sender.User,
		"admin":           fmt.Sprintf("%v", s.Sender.Admin),
		"action":          s.Action,
		"newMember":       s.Member.User,
		"newMemberStatus": s.Member.Admin,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type PageBuildEvent struct {
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s PageBuildEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type PublicEvent struct {
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s PublicEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type PullRequestEvent struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
	Sender      Sender      `json:"sender"`
}

func (s PullRequestEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"action":       s.Action,
		"repository":   s.Repository.Repository,
		"private":      fmt.Sprintf("%v", s.Repository.Private),
		"user":         s.Sender.User,
		"admin":        fmt.Sprintf("%v", s.Sender.Admin),
		"prNumber":     fmt.Sprintf("%v", s.PullRequest.Number),
		"stars":        s.Repository.Stars,
		"forks":        s.Repository.Forks,
		"issues":       s.Repository.Issues,
		"state":        s.PullRequest.State,
		"title":        s.PullRequest.Title,
		"comments":     s.PullRequest.Comments,
		"commits":      s.PullRequest.Commits,
		"additions":    s.PullRequest.Additions,
		"deletions":    s.PullRequest.Deletions,
		"changedFiles": s.PullRequest.ChangedFiles,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type PullRequestReviewCommentEvent struct {
	Comment     PullRequestReviewComment `json:"comment"`
	PullRequest PullRequest              `json:"pull_request"`
	Repository  Repository               `json:"repository"`
	Sender      Sender                   `json:"sender"`
}

func (s PullRequestReviewCommentEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository":   s.Repository.Repository,
		"private":      fmt.Sprintf("%v", s.Repository.Private),
		"user":         s.Sender.User,
		"admin":        fmt.Sprintf("%v", s.Sender.Admin),
		"prNumber":     fmt.Sprintf("%v", s.PullRequest.Number),
		"stars":        s.Repository.Stars,
		"forks":        s.Repository.Forks,
		"issues":       s.Repository.Issues,
		"state":        s.PullRequest.State,
		"title":        s.PullRequest.Title,
		"comments":     s.PullRequest.Comments,
		"commits":      s.PullRequest.Commits,
		"additions":    s.PullRequest.Additions,
		"deletions":    s.PullRequest.Deletions,
		"changedFiles": s.PullRequest.ChangedFiles,
		"commentFile":  s.Comment.File,
		"comment":      s.Comment.Comment,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type PushEvent struct {
	Ref        string     `json:"ref"`
	Before     string     `json:"before"`
	After      string     `json:"after"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
	HeadCommit HeadCommit `json:"head_commit"`
}

func (s PushEvent) NewEvent() agent.Event {
	data := plugins.GitCommit{
		Repository: s.Repository.Repository,
		User:       s.Sender.User,
		Message:    s.HeadCommit.Message,
		Ref:        s.Ref,
		ParentHash: s.Before,
		Hash:       s.After,
		Created:    s.HeadCommit.Timestamp,
	}

	return agent.NewEvent(data, nil)
}

type PingEvent struct {
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s PingEvent) NewEvent() agent.Event {
	data := plugins.GitPing{
		Repository: s.Repository.Repository,
		User:       s.Sender.User,
	}
	return agent.NewEvent(data, nil)
}

type ReleaseEvent struct {
	Release    Release    `json:"release"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s ReleaseEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
		"tagName":    s.Release.TagName,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type RepositoryEvent struct {
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s RepositoryEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type StatusEvent struct {
	Hash       string     `json:"sha"`
	State      string     `json:"state"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
	Context    string     `json:"context"`
}

func (s StatusEvent) NewEvent() agent.Event {
	data := plugins.GitStatus{
		Repository: s.Repository.Repository,
		User:       s.Sender.User,
		Hash:       s.Hash,
		State:      s.State,
		Context:    s.Context,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type TeamAddEvent struct {
	Team       Team       `json:"team"`
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s TeamAddEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
		"teamName":   s.Team.Name,
	}
	m := agent.NewEvent(data, nil)
	return m
}

type WatchEvent struct {
	Repository Repository `json:"repository"`
	Sender     Sender     `json:"sender"`
}

func (s WatchEvent) NewEvent() agent.Event {

	data := map[string]interface{}{
		"repository": s.Repository.Repository,
		"private":    fmt.Sprintf("%v", s.Repository.Private),
		"user":       s.Sender.User,
		"admin":      fmt.Sprintf("%v", s.Sender.Admin),
		"stars":      s.Repository.Stars,
		"forks":      s.Repository.Forks,
		"issues":     s.Repository.Issues,
	}
	m := agent.NewEvent(data, nil)
	return m
}
