package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/deis/deisrel/actions"
	"github.com/deis/deisrel/components"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

var version = "0.0.0" // replaced when building

func main() {
	app := cli.NewApp()

	token, ok := os.LookupEnv("GITHUB_ACCESS_TOKEN")
	if !ok {
		fmt.Println("\"GITHUB_ACCESS_TOKEN\" must be set to a valid GitHub access token.")
		os.Exit(1)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ghclient := github.NewClient(oauth2.NewClient(oauth2.NoContext, ts))

	app.Name = "deisrel"
	app.Usage = "Manage deis workflow releases."
	app.UsageText = "deisrel [options] <chart versions> <repo map>"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Value: "pretty",
			Usage: "format to print output in. Valid options: [pretty, json]",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "release",
			Action:    actions.ReleaseComponent(ghclient, os.Stdout),
			Usage:     "create a new tagged release of a Deis component",
			ArgsUsage: "<component> <new-tag>",
			Flags: []cli.Flag{
				cli.BoolTFlag{
					Name:  "dry-run",
					Usage: "print the actions to be taken, but don't change anything (default: true)",
				},
				cli.StringFlag{
					Name:  "sha",
					Usage: "commit to be tagged as a new release. (default: HEAD of master)",
				},
			},
		},
		cli.Command{
			Name: "changelog",
			Subcommands: []cli.Command{
				cli.Command{
					Name:        "global",
					Action:      actions.GenerateChangelog(ghclient, os.Stdout),
					Usage:       "deisrel changelog global <repo map> <old-release> <new-release>",
					Description: "Aggregate changelog entries from all known repositories for a specified release",
				},
				cli.Command{
					Name:        "individual",
					Action:      actions.GenerateIndividualChangelog(ghclient, os.Stdout),
					Usage:       "deisrel changelog individual <repo-name> <new-release>",
					Description: "Generate a changelog entry for an changes on an individual repository, from a specified old release through a specified git SHA. The release will be called the specified new release in the changelog's title",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "base-tag",
							Usage: "old tag where changelog starts. Optional, defaults to most recent tag",
						},
						cli.StringFlag{
							Name:  "sha",
							Usage: "commit that will be tagged in the new release. Optional, defaults to most recent commit on master.",
						},
					},
				},
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		res := make(map[string]interface{})

		if c.NArg() < 2 {
			return cli.NewExitError("A params and a repo mapping file is required", 1)
		}

		out, err := ioutil.ReadFile(c.Args().Get(0))
		if err != nil {
			return cli.NewExitError(err.Error(), 2)
		}

		err = toml.Unmarshal(out, &res)
		if err != nil {
			return cli.NewExitError(err.Error(), 3)
		}

		mapping := make(map[string][]string)
		out, err = ioutil.ReadFile(c.Args().Get(1))
		if err != nil {
			return cli.NewExitError(err.Error(), 2)
		}
		err = json.Unmarshal(out, &mapping)
		if err != nil {
			return cli.NewExitError(err.Error(), 3)
		}

		versions, err := components.CheckVersions(res, mapping, ghclient)
		if err != nil {
			return cli.NewExitError(err.Error(), 4)
		}

		sort.Sort(components.ByName(versions))

		if c.String("output") == "json" {
			out, err := json.MarshalIndent(versions, "", "  ")
			if err != nil {
				return cli.NewExitError(err.Error(), 5)
			}
			fmt.Println(string(out))
		} else if c.String("output") == "pretty" {
			// Add padding for version name to all output lines up
			longestName := 0
			for _, version := range versions {
				name := len(version.Name)

				if name > longestName {
					longestName = name
				}
			}

			for _, version := range versions {
				cleanMsg := "clean"
				if !version.Clean {
					cleanMsg = "dirty"
				}

				padding := strings.Repeat(" ", longestName-len(version.Name))

				fmt.Printf("%s%s %s -> %s (%s)\n", version.Name, padding, version.ChartVersion, version.ComponentVersion, cleanMsg)
				if !version.Clean {
					fmt.Printf("\t%s has unrelased changes. See %s\n", version.Name, version.Diff)
				}
			}
		} else {
			return cli.NewExitError(fmt.Sprintf("Unrecognized output format: %s", c.String("output")), 1)
		}

		return nil
	}

	app.Run(os.Args)
}
