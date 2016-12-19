package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/deis/deisrel/testutil"

	"github.com/google/go-github/github"
)

func TestGetComponentTitle(t *testing.T) {
	testCases := map[string]string{
		"charts":           "Workflow",
		"workflow":         "Workflow Documentation",
		"logger":           "Logger",
		"workflow-manager": "Workflow Manager",
		"controller":       "Controller",
		"slugbuilder":      "Slugbuilder",
	}
	for component, title := range testCases {
		if title != getComponentTitle(component) {
			t.Errorf("Expected getComponentTitle(%s) to return \"%s\"", component, title)
		}
	}
}

func TestIsValidSemVerTag(t *testing.T) {
	goodTags := []string{"v2.8.0", "v2.8.1", "v0.99.100928", "v12.0.01"}
	for _, tag := range goodTags {
		if !isValidSemVerTag(tag) {
			t.Errorf("isValidSemVerTag(\"%s\") returned false", tag)
		}
	}
	badTags := []string{"2.8.0", "v2.8", "tag v2.8.0", "x2.8.0", "supernaut"}
	for _, tag := range badTags {
		if isValidSemVerTag(tag) {
			t.Errorf("isValidSemVerTag(\"%s\") returned true", tag)
		}
	}
}

func TestCreateRelease(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	component := "controller"
	title := getComponentTitle(component)
	sha := "master"
	newTag := "v1.2.3"
	body := "changelog"
	expectedName := "Deis Controller v1.2.3"

	ts.Mux.HandleFunc("/repos/deis/controller/releases", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "POST" {
			t.Errorf("Request method: %v, want POST", got)
		}
		v := new(github.RepositoryRelease)
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			t.Errorf("decoding %+v into a github.RepositoryRelease failed; Error = %+v", r.Body, err)
		}

		fmt.Fprintf(w, `{
			"name": "` + expectedName + `",
			"tag_name": "` + newTag + `",
			"target_commitish": "` + sha + `",
			"body": "` + body + `",
			"created_at": "0001-01-01T00:00:00Z",
			"published_at": "0001-01-01T00:00:00Z"
		}`)
	})

	got, _, err := createRelease(ts.Client, component, title, sha, newTag, body)
	if err != nil {
		t.Errorf("createRelease returned error: %v", err)
	}

	want := &github.RepositoryRelease{
		Name: &expectedName,
		TagName: &newTag,
		TargetCommitish: &sha,
		Body: &body,
		CreatedAt: &github.Timestamp{time.Time{}},
		PublishedAt: &github.Timestamp{time.Time{}},
	}

	if !releaseEqual(got, want) {
		t.Errorf("createRelease returned %+v, want %+v", got, want)
	}
}

func releaseEqual(got, want *github.RepositoryRelease) bool {
	return reflect.DeepEqual(got.Name, want.Name) &&
		reflect.DeepEqual(got.TagName, want.TagName) &&
		reflect.DeepEqual(got.TargetCommitish, want.TargetCommitish) &&
		reflect.DeepEqual(got.Body, want.Body) &&
		got.CreatedAt.String() == want.CreatedAt.String() &&
		got.PublishedAt.String() == want.PublishedAt.String()
}
