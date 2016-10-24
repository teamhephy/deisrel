package actions

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

// ReleaseComponent creates a new GitHub release for a component.
func ReleaseComponent(client *github.Client, dest io.Writer) func(*cli.Context) error {
	return func(c *cli.Context) error {
		if c.NArg() != 2 {
			cli.ShowCommandHelp(c, "release")
			return cli.NewExitError("", 1)
		}
		component := c.Args().Get(0)
		newTag := c.Args().Get(1)
		if !isValidSemVerTag(newTag) {
			return cli.NewExitError("Invalid semantic version tag", 1)
		}
		dryRun := c.Bool("dry-run")
		sha := c.String("sha")

		if dryRun == true {
			fmt.Fprintln(dest, "Doing a dry run of the component release...")
		}

		// generate the changelog into a string
		buf := bytes.NewBufferString("")
		genChangelog := GenerateIndividualChangelog(client, buf)
		if err := genChangelog(c); err != nil {
			return err
		}
		changelog := buf.String()
		fmt.Println(changelog)

		// translate the component to a title
		title := getComponentTitle(component)

		// if it's not a dry run, prompt the user to confirm
		prompt := `
Please review the above changelog contents and ensure:
  1. All intended commits are mentioned
  2. The changes agree with the semver release tag (major, minor, or patch)

Create release for Deis %s %s?`
		if !dryRun && askForConfirmation(fmt.Sprintf(prompt, title, newTag)) {
			releaseName := fmt.Sprintf("Deis %s %s", title, newTag)

			release := github.RepositoryRelease{
				TargetCommitish: &sha,
				TagName:         &newTag,
				Name:            &releaseName,
				Body:            &changelog,
			}
			rel, _, err := client.Repositories.CreateRelease("deis", component, &release)
			if err != nil {
				return err
			}
			fmt.Fprintf(dest, "New release is available at %s\n", *rel.HTMLURL)
		}

		return nil
	}
}

func isValidSemVerTag(tag string) bool {
	regx := regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)
	return regx.MatchString(tag)
}

func getComponentTitle(component string) string {
	// use a lookup table for consistent titles
	titleMap := map[string]string{
		"charts":           "Workflow",
		"monitor":          "Monitoring",
		"nsq":              "NSQ",
		"workflow":         "Workflow Documentation",
		"workflow-cli":     "Workflow Client",
		"workflow-e2e":     "Workflow End-to-End Tests",
		"workflow-manager": "Workflow Manager",
	}
	title, ok := titleMap[component]
	if !ok {
		// otherwise just uppercase initial letters
		title = strings.Title(component)
	}
	return title
}

// askForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user.
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
