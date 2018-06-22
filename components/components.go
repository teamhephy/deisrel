package components

import (
	"errors"
	"fmt"
	"os"
        "context"

	"github.com/google/go-github/github"
)

var diffURLFormat = "https://github.com/teamhephy/%s/compare/%s...master"

// ComponentVersion is a combination of the different types
type ComponentVersion struct {
	Name             string `json:"name"`
	ChartVersion     string `json:"chart"`
	ComponentVersion string `json:"component"`
	Diff             string `json:"diff"`
	Clean            bool   `json:"clean"`
}

type ByName []ComponentVersion

func (v ByName) Len() int           { return len(v) }
func (v ByName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByName) Less(i, j int) bool { return v[i].Name < v[j].Name }

// CheckVersions checks the versions of all components
func CheckVersions(chart map[string]interface{}, repositoryMap map[string]string,
	ghclient *github.Client) ([]ComponentVersion, error) {
	versions := []ComponentVersion{}

	for repo, name := range repositoryMap {
		version := ComponentVersion{Name: name}
		version.ChartVersion = getChartVersion(chart, name)

		componentVersion, err := getRespositoryVersion(ghclient, repo)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		version.ComponentVersion = componentVersion

		// Check if tag points to latest commit on master
		clean, err := checkTag(ghclient, repo, version.ComponentVersion)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		version.Clean = clean

		version.Diff = fmt.Sprintf(diffURLFormat, repo, version.ComponentVersion)

		versions = append(versions, version)
	}
	return versions, nil
}

func getChartVersion(chart map[string]interface{}, component string) string {
	if dependencies, ok := chart["dependencies"].([]interface{}); ok {
		for _, dependency := range dependencies {
			if name, ok := dependency.(map[interface{}]interface{})["name"]; ok {
				if name == component {
					if version, ok := dependency.(map[interface{}]interface{})["version"]; ok {
						return version.(string)
					}
				}
			}
		}
	}

	return "unknown"
}

func getRespositoryVersion(client *github.Client, repo string) (string, error) {
	tags, _, err := client.Repositories.ListTags(context.Background(), "teamhephy", repo, nil)
	if err != nil {
		return "unknown", err
	}

	if len(tags) < 1 {
		return "none", errors.New("No tags for component")
	}

	return *tags[0].Name, nil
}

func checkTag(client *github.Client, repo, tagName string) (bool, error) {
	master, _, err := client.Repositories.GetBranch(context.Background(), "teamhephy", repo, "master")
	if err != nil {
		return false, err
	}

	object, _, err := client.Git.GetRef(context.Background(), "teamhephy", repo, "refs/tags/"+tagName)
	if err != nil {
		return false, err
	}

	// If tag is a light tag, return the object iteself
	if *object.Object.Type != "tag" {
		return *master.Commit.SHA == *object.Object.SHA, nil
	}

	// If tag is an annotated tag, return the object it points to.
	tag, _, err := client.Git.GetTag(context.Background(), "teamhephy", repo, *object.Object.SHA)
	if err != nil {
		return false, err
	}

	return *master.Commit.SHA == *tag.Object.SHA, nil
}
