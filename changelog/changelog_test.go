package changelog

import (
	"bytes"
	"strings"
	"testing"

	"github.com/arschles/assert"
)

const (
	oldRelease = "old"
	newRelease = "new"
)

var (
	features      = []string{"feature1"}
	fixes         = []string{"fix1"}
	documentation = []string{"doc1"}
	tests         = []string{"test1"}
	maintenance   = []string{"maint1"}
	refactor      = []string{"ref1"}
)

func TestTemplate(t *testing.T) {
	type testCase struct {
		vals    Values
		missing string
	}
	testCases := []testCase{
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         fixes,
				Refactors:     refactor,
				Documentation: documentation,
				Tests:         tests,
				Maintenance:   maintenance,
			},
			missing: "",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      nil,
				Fixes:         fixes,
				Refactors:     refactor,
				Documentation: documentation,
				Tests:         tests,
				Maintenance:   maintenance,
			},
			missing: "#### Features",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         nil,
				Refactors:     refactor,
				Documentation: documentation,
				Tests:         tests,
				Maintenance:   maintenance,
			},
			missing: "#### Fixes",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         fixes,
				Refactors:     refactor,
				Documentation: nil,
				Tests:         tests,
				Maintenance:   maintenance,
			},
			missing: "#### Documentation",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         fixes,
				Refactors:     refactor,
				Documentation: documentation,
				Tests:         nil,
				Maintenance:   maintenance,
			},
			missing: "#### Tests",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         fixes,
				Refactors:     refactor,
				Documentation: documentation,
				Tests:         tests,
				Maintenance:   nil,
			},
			missing: "#### Maintenance",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         fixes,
				Refactors:     nil,
				Documentation: documentation,
				Tests:         tests,
				Maintenance:   maintenance,
			},
			missing: "#### Refactors",
		},
	}

	for i, testCase := range testCases {
		var buf bytes.Buffer
		if err := Tpl.Execute(&buf, testCase.vals); err != nil {
			t.Errorf("Error executing template %d (%s)", i, err)
			continue
		}
		if len(testCase.missing) > 0 {
			outStr := string(buf.Bytes())
			if strings.Contains(outStr, testCase.missing) {
				t.Errorf("Expected [%s] to be missing from the rendered template %d, but found it", testCase.missing, i)
			}
		}
	}
}

func TestMergeValues(t *testing.T) {
	val1 := Values{RepoName: "repo-a", OldRelease: "v1.2.3", NewRelease: "v1.2.4", Features: []string{"feat1"}}
	val2 := Values{RepoName: "repo-b", OldRelease: "v4.5.6", NewRelease: "v4.6.0", Fixes: []string{"fix1"}, Features: []string{"feat2"}}
	val3 := Values{OldRelease: "v1.2.3", NewRelease: "v1.2.3"} // no change; should not be added to res.Releases
	res := MergeValues("old", "new", []Values{val1, val2, val3})
	assert.Equal(t, res.OldRelease, "old", "old release")
	assert.Equal(t, res.NewRelease, "new", "new release")
	assert.Equal(t, len(res.Features), 2, "length of features slice")
	assert.Equal(t, len(res.Fixes), 1, "length of fixes slice")
	assert.Equal(t, res.Features, []string{"feat1", "feat2"}, "features slice")
	assert.Equal(t, res.Fixes, []string{"fix1"}, "fixes slice")
	assert.Equal(t, res.Releases, []string{"repo-a v1.2.3 -> v1.2.4", "repo-b v4.5.6 -> v4.6.0"}, "releases slice")
}
