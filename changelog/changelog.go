package changelog

import (
	"fmt"
	"text/template"
)

const (
	tplStr = `{{ if and .OldRelease .NewRelease }}
### {{.OldRelease}} -> {{.NewRelease}}
{{end}}

{{- if (gt (len .Releases) 0) }}
#### Releases

{{range .Releases}}- {{.}}
{{end}}

{{- end}}
{{- if gt (len .Features) 0 }}
#### Features

{{range .Features}}- {{.}}
{{end}}

{{- end}}
{{- if gt (len .Refactors) 0 }}
#### Refactors

{{range .Refactors}}- {{.}}
{{end}}

{{- end}}
{{- if gt (len .Fixes) 0 }}
#### Fixes

{{range .Fixes}}- {{.}}
{{end}}

{{- end}}
{{- if gt (len .Documentation) 0 }}
#### Documentation

{{ range .Documentation}}- {{.}}
{{end}}

{{- end}}
{{- if gt (len .Tests) 0 }}
#### Tests

{{ range .Tests}}- {{.}}
{{end}}

{{- end}}
{{- if gt (len .Maintenance) 0 }}
#### Maintenance

{{range .Maintenance}}- {{.}}
{{end}}
{{- end}}`
)

var (
	// Tpl is the standard changelog template. Execute it with a Values struct
	Tpl = template.Must(template.New("changelog").Parse(tplStr))
)

// ByName takes an array of Values structs and sorts in alphabetical order of name
type ByName []Values

func (v ByName) Len() int           { return len(v) }
func (v ByName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByName) Less(i, j int) bool { return v[i].RepoName < v[j].RepoName }

// Values represents the values that are required to render a changelog
type Values struct {
	RepoName      string
	OldRelease    string
	NewRelease    string
	Features      []string
	Fixes         []string
	Documentation []string
	Tests         []string
	Maintenance   []string
	Refactors     []string
	Releases      []string
}

// MergeValues merges all of the slices in vals together into a single Values struct which has OldRelease set to oldRel and NewRelease set to newRel
func MergeValues(oldRel, newRel string, vals []Values) *Values {
	ret := &Values{OldRelease: oldRel, NewRelease: newRel}
	for _, val := range vals {
		ret.Features = append(ret.Features, val.Features...)
		ret.Refactors = append(ret.Refactors, val.Refactors...)
		ret.Fixes = append(ret.Fixes, val.Fixes...)
		ret.Documentation = append(ret.Documentation, val.Documentation...)
		ret.Tests = append(ret.Tests, val.Tests...)
		ret.Maintenance = append(ret.Maintenance, val.Maintenance...)
		if val.OldRelease != val.NewRelease {
			ret.Releases = append(ret.Releases, fmt.Sprintf("%s %s -> %s", val.RepoName, val.OldRelease, val.NewRelease))
		}
	}
	return ret
}
