package changelog

import (
	"testing"

	"github.com/arschles/assert"
)

func TestAppendToValues(t *testing.T) {
	changelogMessage := "changelog message"

	type testCase struct {
		commitMessage  string
		vals           *Values
		skippedCommits []string
	}

	testCases := []testCase{
		testCase{
			commitMessage:  "Merge pull request",
			vals:           &Values{},
			skippedCommits: []string{},
		},
		testCase{
			commitMessage:  "skipped commit",
			vals:           &Values{},
			skippedCommits: []string{changelogMessage},
		},
		testCase{
			commitMessage:  "chore()",
			vals:           &Values{Maintenance: []string{changelogMessage}},
			skippedCommits: []string{},
		},
		testCase{
			commitMessage:  "docs()",
			vals:           &Values{Documentation: []string{changelogMessage}},
			skippedCommits: []string{},
		},
		testCase{
			commitMessage:  "feat()",
			vals:           &Values{Features: []string{changelogMessage}},
			skippedCommits: []string{},
		},
		testCase{
			commitMessage:  "fix()",
			vals:           &Values{Fixes: []string{changelogMessage}},
			skippedCommits: []string{},
		},
		testCase{
			commitMessage:  "ref()",
			vals:           &Values{Refactors: []string{changelogMessage}},
			skippedCommits: []string{},
		},
		testCase{
			commitMessage:  "test()",
			vals:           &Values{Tests: []string{changelogMessage}},
			skippedCommits: []string{},
		},
	}

	for _, testCase := range testCases {
		var skippedCommits []string
		vals := &Values{}
		skippedCommits = appendToValues(vals, testCase.commitMessage, changelogMessage, skippedCommits)

		assert.Equal(t, len(skippedCommits), len(testCase.skippedCommits), "skipped commits length")
		assert.Equal(t, vals, testCase.vals, "vals")
	}
}

func TestAppendToValuesMultipleSkippedCommits(t *testing.T) {
	var skippedCommits []string
	changelogMessage := "changelog message"

	commitMessages := []string{"Merge pull request", "skipped commit 1", "skipped commit 2"}

	vals := &Values{}
	for _, commitMessage := range commitMessages {
		skippedCommits = appendToValues(vals, commitMessage, changelogMessage, skippedCommits)
	}

	assert.Equal(t, len(skippedCommits), 2, "skipped commits length")
	assert.Equal(t, vals, &Values{}, "vals")
}
