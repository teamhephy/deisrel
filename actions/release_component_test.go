package actions

import "testing"

func TestGetComponentTitle(t *testing.T) {
	testCases := map[string]string{
		"charts":           "Workflow",
		"workflow":         "Workflow Documentation",
		"logger":           "Logger",
		"workflow-manager": "Workflow Manager",
		"controller":       "Controller",
		"slugbuilder":      "Slugbuilder",
	}
	for component, title := range testCases {
		if title != getComponentTitle(component) {
			t.Errorf("Expected getComponentTitle(%s) to return \"%s\"", component, title)
		}
	}
}

func TestIsValidSemVerTag(t *testing.T) {
	goodTags := []string{"v2.8.0", "v2.8.1", "v0.99.100928", "v12.0.01"}
	for _, tag := range goodTags {
		if !isValidSemVerTag(tag) {
			t.Errorf("isValidSemVerTag(\"%s\") returned false", tag)
		}
	}
	badTags := []string{"2.8.0", "v2.8", "tag v2.8.0", "x2.8.0", "supernaut"}
	for _, tag := range badTags {
		if isValidSemVerTag(tag) {
			t.Errorf("isValidSemVerTag(\"%s\") returned true", tag)
		}
	}
}
