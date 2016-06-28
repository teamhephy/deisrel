package actions

import (
	"fmt"
	"log"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/git"
	"github.com/google/go-github/github"
)

func MoveMilestone(ghClient *github.Client) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		oldMilestone := c.Args().Get(0)
		newMilestone := c.Args().Get(1)
		if oldMilestone == "" || newMilestone == "" {
			log.Fatal("Usage: mv <old-milestone> <new-milestone>")
		}
		ok := true
		if !c.Bool(YesFlag) {
			var err error
			ok, err = prompt()
			if err != nil {
				log.Fatal(err)
			}
		}
		if ok {
			var wg sync.WaitGroup
			done := make(chan bool)
			errCh := make(chan error)
			defer close(errCh)
			for _, repo := range allGitRepoNames {
				wg.Add(1)
				go func(repo string) {
					defer wg.Done()
					if err := git.MoveMilestone(ghClient, repo, oldMilestone, newMilestone, c.Bool(IncludeClosed)); err != nil {
						errCh <- fmt.Errorf("Error moving %s issues from milestone %s to milestone %s: %s", repo, oldMilestone, newMilestone, err)
					}
				}(repo)
			}
			go func() {
				wg.Wait()
				close(done)
			}()
			errs := []error{}
			for {
				select {
				case <-done:
					if len(errs) > 0 {
						var errStr string
						for _, err := range errs {
							log.Println(err)
							errStr = fmt.Sprintf("%s%s\n", errStr, err)
						}
						return fmt.Errorf(errStr)
					}
					return nil
				case err := <-errCh:
					errs = append(errs, err)
				}
			}
		}
		return nil
	}
}
