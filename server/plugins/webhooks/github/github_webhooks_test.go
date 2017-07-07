package github

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/checkr/codeflow/server/agent"
)

func testGithubWebhookRequest(event string, jsonString string, t *testing.T) {
	gh := &GithubWebhook{
		Path:   "/github",
		events: make(chan agent.Event, 1),
	}
	defer close(gh.events)
	req := httptest.NewRequest("POST", "/github", strings.NewReader(jsonString))
	req.Header.Add("X-Github-Event", event)
	w := httptest.NewRecorder()
	gh.eventHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("POST "+event+" returned HTTP status code %v.\nExpected %v", w.Code, http.StatusOK)
	}
}

func TestCommitCommentEvent(t *testing.T) {
	testGithubWebhookRequest("commit_comment", CommitCommentEventJSON(), t)
}

func TestDeleteEvent(t *testing.T) {
	testGithubWebhookRequest("delete", DeleteEventJSON(), t)
}

func TestDeploymentEvent(t *testing.T) {
	testGithubWebhookRequest("deployment", DeploymentEventJSON(), t)
}

func TestDeploymentStatusEvent(t *testing.T) {
	testGithubWebhookRequest("deployment_status", DeploymentStatusEventJSON(), t)
}

func TestForkEvent(t *testing.T) {
	testGithubWebhookRequest("fork", ForkEventJSON(), t)
}

func TestGollumEvent(t *testing.T) {
	testGithubWebhookRequest("gollum", GollumEventJSON(), t)
}

func TestIssueCommentEvent(t *testing.T) {
	testGithubWebhookRequest("issue_comment", IssueCommentEventJSON(), t)
}

func TestIssuesEvent(t *testing.T) {
	testGithubWebhookRequest("issues", IssuesEventJSON(), t)
}

func TestMemberEvent(t *testing.T) {
	testGithubWebhookRequest("member", MemberEventJSON(), t)
}

func TestMembershipEvent(t *testing.T) {
	testGithubWebhookRequest("membership", MembershipEventJSON(), t)
}

func TestPageBuildEvent(t *testing.T) {
	testGithubWebhookRequest("page_build", PageBuildEventJSON(), t)
}

func TestPublicEvent(t *testing.T) {
	testGithubWebhookRequest("public", PublicEventJSON(), t)
}

func TestPullRequestReviewCommentEvent(t *testing.T) {
	testGithubWebhookRequest("pull_request_review_comment", PullRequestReviewCommentEventJSON(), t)
}

func TestPushEvent(t *testing.T) {
	testGithubWebhookRequest("push", PushEventJSON(), t)
}

func TestReleaseEvent(t *testing.T) {
	testGithubWebhookRequest("release", ReleaseEventJSON(), t)
}

func TestRepositoryEvent(t *testing.T) {
	testGithubWebhookRequest("repository", RepositoryEventJSON(), t)
}

func TestStatusEvent(t *testing.T) {
	testGithubWebhookRequest("status", StatusEventJSON(), t)
}

func TestTeamAddEvent(t *testing.T) {
	testGithubWebhookRequest("team_add", TeamAddEventJSON(), t)
}

func TestWatchEvent(t *testing.T) {
	testGithubWebhookRequest("watch", WatchEventJSON(), t)
}
