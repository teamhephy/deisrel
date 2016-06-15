package actions

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/arschles/assert"
	"github.com/arschles/sys"
	"github.com/deis/deisrel/git"
	"github.com/deis/deisrel/testutil"
	"github.com/google/go-github/github"
)

func TestDownloadFiles(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	org := "deis"
	repo := "charts"
	var opt github.RepositoryContentGetOptions
	helmChart := WorkflowE2EChart
	helmChart.Files = []string{"README.md"}
	filePath := filepath.Join(helmChart.Name, helmChart.Files[0])

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/contents/%s", org, repo, helmChart.Name), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		resp := `[
		  {
		    "type": "file",
		    "size": 625,
		    "name": "` + helmChart.Files[0] + `",
		    "path": "` + filePath + `",
				"content": "contents",
		    "sha": "",
		    "url": "",
		    "git_url": "",
		    "html_url": "",
		    "download_url": "https://raw.githubusercontent.com/` + org + `/` + repo + `/master/` + filePath + `",
		    "_links": {
		      "self": "",
		      "git": "",
		      "html": ""
		    }
		  },
		  {
		    "type": "dir",
		    "size": 0,
		    "name": "` + helmChart.Name + `",
		    "path": "` + helmChart.Name + `",
		    "sha": "",
		    "url": "",
		    "git_url": "",
		    "html_url": "",
		    "download_url": null,
		    "_links": {
		      "self": "",
		      "git": "",
		      "html": ""
		    }
		  }
		]`
		fmt.Fprintf(w, resp)
	})

	got, err := downloadFiles(ts.Client, org, repo, &opt, helmChart)

	assert.NoErr(t, err)
	assert.Equal(t, len(got), 1, "length of downloaded file slice")
	assert.Equal(t, got[0].Name, helmChart.Files[0], "file name")
}

func TestDownloadFilesNotExist(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	org := "deis"
	repo := "charts"
	var opt github.RepositoryContentGetOptions
	helmChart := WorkflowE2EChart
	helmChart.Files = []string{"NotExist.md"}

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/contents/%s", org, repo, helmChart.Name), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		resp := `[
		  {
		    "type": "dir",
		    "size": 0,
				"name": "` + helmChart.Name + `",
		    "path": "` + helmChart.Name + `",
		    "sha": "",
		    "url": "",
		    "git_url": "",
		    "html_url": "",
		    "download_url": null,
		    "_links": {
		      "self": "",
		      "git": "",
		      "html": ""
		    }
		  }
		]`
		fmt.Fprintf(w, resp)
	})

	got, err := downloadFiles(ts.Client, org, repo, &opt, helmChart)

	assert.ExistsErr(t, err, "file doesn't exist")
	assert.Nil(t, got, "downloaded files slice")
}

func TestStageFiles(t *testing.T) {
	fakeFS := sys.NewFakeFS()

	readCloser := ioutil.NopCloser(bytes.NewBufferString(""))
	fileName := "testFile"
	ghFileToStage := git.File{ReadCloser: readCloser, Name: fileName}
	stagingDir := "staging"

	fakeFS.MkdirAll(stagingDir, os.ModePerm)
	stageFiles(fakeFS, []git.File{ghFileToStage}, stagingDir)

	_, err := fakeFS.ReadFile(filepath.Join(stagingDir, fileName))
	assert.NoErr(t, err)
}

func TestCreateDir(t *testing.T) {
	fakeFS := sys.NewFakeFS()

	err := createDir(fakeFS, "foo")

	assert.NoErr(t, err)
	assert.Equal(t,
		fakeFS.Files,
		map[string]*bytes.Buffer{"foo": &bytes.Buffer{}},
		"file system contents")
}

func TestCreateDirAlreadyExists(t *testing.T) {
	fakeFS := sys.NewFakeFS()

	fakeFS.Create("foo")
	err := createDir(fakeFS, "foo")

	assert.NoErr(t, err)
	assert.Equal(t,
		fakeFS.Files,
		map[string]*bytes.Buffer{"foo": &bytes.Buffer{}},
		"file system contents")
}

func TestUpdateFilesWithMinorRelease(t *testing.T) {
	fakeFS := sys.NewFakeFS()
	fakeFP := sys.NewFakeFP()

	fileName := "foo/bar"
	fakeFS.Create(fileName)
	fakeFS.WriteFile(fileName, []byte("name: workflow-dev, version: v2.0.0"), os.ModePerm)
	var deisRelease = releaseName{
		Full:  "foobar",
		Short: "bar",
	}
	err := updateFilesWithRelease(fakeFP, fakeFS, deisRelease, fileName)

	assert.NoErr(t, err)
	actualFileContents, err := fakeFS.ReadFile(fileName)
	assert.NoErr(t, err)
	assert.Equal(t,
		actualFileContents,
		[]byte(fmt.Sprintf("name: workflow-%s, version: %s", deisRelease.Short, deisRelease.Full)),
		"updated file")
}

func TestUpdateFilesWithMajorRelease(t *testing.T) {
	fakeFS := sys.NewFakeFS()
	fakeFP := sys.NewFakeFP()

	type testFile struct {
		name             string
		originalContents []byte
		expectedContents []byte
	}

	var majorRelease = releaseName{
		Full:  "major",
		Short: "major",
	}

	testFiles := []testFile{
		testFile{
			name:             "foo/bar",
			originalContents: []byte("name: workflow-dev, version: v2.0.0"),
			expectedContents: []byte(fmt.Sprintf("name: workflow-%s, version: %s", majorRelease.Full, majorRelease.Full)),
		},
		testFile{
			name:             "README",
			originalContents: []byte("Workflow 2.0.0-dev"),
			expectedContents: []byte("Workflow 2.0.0"),
		},
	}

	for _, testFile := range testFiles {
		fakeFS.Create(testFile.name)
		fakeFS.WriteFile(testFile.name, testFile.originalContents, os.ModePerm)

		err := updateFilesWithRelease(fakeFP, fakeFS, majorRelease, testFile.name)

		assert.NoErr(t, err)
		actualFileContents, err := fakeFS.ReadFile(testFile.name)
		fmt.Println(string(actualFileContents))
		assert.NoErr(t, err)
		assert.Equal(t,
			actualFileContents,
			testFile.expectedContents,
			fmt.Sprintf("updated file with release %s", majorRelease))
	}
}

func TestUpdateFilesWithReleaseWithoutRelease(t *testing.T) {
	fakeFS := sys.NewFakeFS()
	fakeFP := sys.NewFakeFP()

	fileName := "foo/bar"
	fakeFS.Create(fileName)
	fakeFS.WriteFile(fileName, []byte("name: workflow-dev, version: v2.0.0"), os.ModePerm)
	fakeFP.WalkInvoked = false

	err := updateFilesWithRelease(fakeFP, fakeFS, deisRelease, fileName)
	assert.NoErr(t, err)
	assert.Equal(t, fakeFP.WalkInvoked, false, "walk invoked")

	actualFileContents, err := fakeFS.ReadFile(fileName)
	assert.NoErr(t, err)
	assert.Equal(t,
		actualFileContents,
		[]byte("name: workflow-dev, version: v2.0.0"),
		"updated file")
}
