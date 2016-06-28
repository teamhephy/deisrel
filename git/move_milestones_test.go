package git

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/deis/deisrel/testutil"
)

var (
	repo         = "controller"
	issueNumber  = 1
	oldMilestone = "v2.1"
	newMilestone = "v2.2"
)

func TestMoveMilestoneWithOldMilestoneNotFound(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/deis/%s/milestones", repo), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Request method: %v, want GET", r.Method)
		}
		// Return a list of milestones
		fmt.Fprintf(w, `[ { "title": "%s" } ]`, newMilestone)
	})

	expected := "git.errMilestoneNotFound"
	if err := MoveMilestone(ts.Client, repo, oldMilestone, newMilestone, false); err == nil {
		t.Error("Did not receive expected error message")
	} else if errType := reflect.TypeOf(err).String(); errType != expected {
		t.Errorf("Expected a %s, but got a %s", expected, errType)
	}
}

func TestMoveMilestoneWithNewMilestoneNotFound(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/deis/%s/milestones", repo), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Request method: %v, want GET", r.Method)
		}
		// Return a list of milestones
		fmt.Fprintf(w, `[ { "title": "%s" } ]`, oldMilestone)
	})

	expected := "git.errMilestoneNotFound"
	if err := MoveMilestone(ts.Client, repo, oldMilestone, newMilestone, false); err == nil {
		t.Error("Did not receive expected error message")
	} else if errType := reflect.TypeOf(err).String(); errType != expected {
		t.Errorf("Expected a %s, but got a %s", expected, errType)
	}
}

func TestMoveMilestone(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/deis/%s/milestones", repo), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Request method: %v, want GET", r.Method)
		}
		// Return a list of milestones
		fmt.Fprintf(w, `
		[
		  {
			  "number": 1,
				"title": "%s"
		  },
		  {
				"number": 2,
		    "title": "%s" 
		  }
		]`, oldMilestone, newMilestone)
	})

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/deis/%s/issues", repo), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Request method: %v, want GET", r.Method)
		}
		// Return a list of issues
		fmt.Fprintf(w, `
		[
		  {
		    "number": %d
		  }
		]`, issueNumber)
	})

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/deis/%s/issues/%d", repo, issueNumber), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Request method: %v, want PATCH", r.Method)
		}
		// Return a list of issues
		fmt.Fprint(w, "{}")
	})

	if err := MoveMilestone(ts.Client, repo, oldMilestone, newMilestone, false); err != nil {
		t.Error(err)
	}
}
