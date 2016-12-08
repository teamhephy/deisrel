package actions

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"sort"
	"sync"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/deis/deisrel/changelog"
	"github.com/deis/deisrel/components"
)

// GenerateChangelog is the CLI action for creating an aggregated changelog from all of the Deis Workflow repos.
func GenerateChangelog(client *github.Client, dest io.Writer) func(*cli.Context) error {
	return func(c *cli.Context) error {
		paramsFile := c.Args().Get(0)
		repoMapFile := c.Args().Get(1)
		if paramsFile == "" || repoMapFile == "" {
			log.Fatal("Usage: changelog global <previous chart requirements.lock file> <repo map>")
		}
		versions := []components.ComponentVersion{}

		res := make(map[string]interface{})
		out, err := ioutil.ReadFile(paramsFile)
		if err != nil {
			return cli.NewExitError(err.Error(), 2)
		}
		if err := yaml.Unmarshal(out, &res); err != nil {
			return cli.NewExitError(err.Error(), 3)
		}

		mapping := make(map[string]string)
		out, err = ioutil.ReadFile(repoMapFile)
		if err != nil {
			return cli.NewExitError(err.Error(), 2)
		}
		if err := json.Unmarshal(out, &mapping); err != nil {
			return cli.NewExitError(err.Error(), 3)
		}

		versions, err = components.CheckVersions(res, mapping, client)
		if err != nil {
			return cli.NewExitError(err.Error(), 4)
		}

		vals, errs := generateChangelogVals(client, mapping, versions)
		if len(errs) > 0 {
			for _, err := range errs {
				log.Printf("Error: %s", err)
			}
		}
		sort.Sort(changelog.ByName(vals))
		if err := changelog.Tpl.Execute(dest, changelog.MergeValues("", "", vals)); err != nil {
			log.Fatalf("could not template changelog: %s", err)
		}
		return nil
	}
}

func generateChangelogVals(client *github.Client, repoMap map[string]string, versions []components.ComponentVersion) ([]changelog.Values, []error) {
	var wg sync.WaitGroup
	done := make(chan bool)
	valsCh := make(chan changelog.Values)
	errCh := make(chan error)
	defer close(errCh)
	for repo := range repoMap {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			component := repoMap[repo]
			componentVersion, err := findComponentVersionByName(versions, component)
			if err != nil {
				errCh <- err
				return
			}
			vals := &changelog.Values{RepoName: repo, OldRelease: componentVersion.ChartVersion, NewRelease: componentVersion.ComponentVersion}
			if _, err := changelog.SingleRepoVals(client, vals, componentVersion.ComponentVersion, repo, true); err != nil {
				errCh <- err
				return
			}
			valsCh <- *vals
		}(repo)
	}
	go func() {
		// wait for all fetches from github to be complete before returning
		wg.Wait()
		close(done)
	}()

	vals := []changelog.Values{}
	errs := []error{}
	for {
		select {
		case <-done:
			return vals, errs
		case val := <-valsCh:
			vals = append(vals, val)
		case err := <-errCh:
			errs = append(errs, err)
		}
	}
}

// findComponentVersionByName finds a particular ComponentVersion from an array of
// ComponentVersions based on Name; returns errComponentVersionNotFound if not found
func findComponentVersionByName(componentVersions []components.ComponentVersion, componentName string) (components.ComponentVersion, error) {
	for _, componentVersion := range componentVersions {
		if componentVersion.Name == componentName {
			return componentVersion, nil
		}
	}
	return components.ComponentVersion{}, errComponentVersionNotFound{componentName: componentName}
}
