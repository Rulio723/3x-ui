package clashtemplate

import (
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestCleanRemovesProviderAndSecret(t *testing.T) {
	cleaned, err := Clean(`
mixed-port: 7890
secret: "123456"
proxies:
  - name: old
proxy-providers:
  airport:
    type: http
proxy-groups:
  - name: Auto
    type: url-test
    use:
      - airport
rules:
  - MATCH,Auto
`)
	if err != nil {
		t.Fatalf("Clean() error = %v", err)
	}

	var config map[string]any
	if err := yaml.Unmarshal([]byte(cleaned), &config); err != nil {
		t.Fatalf("cleaned YAML invalid: %v", err)
	}
	if _, ok := config["proxies"]; ok {
		t.Fatal("top-level proxies must be removed")
	}
	if _, ok := config["proxy-providers"]; ok {
		t.Fatal("proxy-providers must be removed")
	}
	if config["secret"] != "" {
		t.Fatalf("secret = %v, want empty string", config["secret"])
	}
	groups, ok := config["proxy-groups"].([]any)
	if !ok || len(groups) != 1 {
		t.Fatalf("proxy-groups = %#v", config["proxy-groups"])
	}
	group, ok := groups[0].(map[string]any)
	if !ok {
		t.Fatalf("proxy group type = %T", groups[0])
	}
	if _, ok := group["use"]; ok {
		t.Fatal("proxy group use must be removed")
	}
	if group["include-all"] != true {
		t.Fatalf("include-all = %v, want true", group["include-all"])
	}
}

func TestCleanEmptyTemplate(t *testing.T) {
	cleaned, err := Clean("  \n")
	if err != nil {
		t.Fatalf("Clean() error = %v", err)
	}
	if cleaned != "" {
		t.Fatalf("Clean() = %q, want empty", cleaned)
	}
}

func TestCleanRejectsInvalidYAML(t *testing.T) {
	_, err := Clean("proxy-groups: [")
	if err == nil {
		t.Fatal("Clean() error = nil, want error")
	}
}
