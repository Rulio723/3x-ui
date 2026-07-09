package sub

import "testing"

func TestBuildRulioClashConfigInjectsProxies(t *testing.T) {
	proxies := []map[string]any{{
		"name": "test-node",
		"type": "direct",
	}}

	config, err := buildRulioClashConfig(proxies)
	if err != nil {
		t.Fatalf("buildRulioClashConfig() error = %v", err)
	}

	got, ok := config["proxies"].([]map[string]any)
	if !ok {
		t.Fatalf("config[proxies] type = %T, want []map[string]any", config["proxies"])
	}
	if len(got) != 1 || got[0]["name"] != "test-node" {
		t.Fatalf("config[proxies] = %#v, want injected proxy", got)
	}
	if _, ok := config["proxy-groups"]; !ok {
		t.Fatal("config missing proxy-groups from template")
	}
	if _, ok := config["rules"]; !ok {
		t.Fatal("config missing rules from template")
	}
}

func TestBuildRulioClashConfigWithCustomTemplate(t *testing.T) {
	proxies := []map[string]any{{
		"name": "test-node",
		"type": "direct",
	}}
	config, err := buildRulioClashConfigWithTemplate(proxies, `
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
		t.Fatalf("buildRulioClashConfigWithTemplate() error = %v", err)
	}
	got, ok := config["proxies"].([]map[string]any)
	if !ok {
		t.Fatalf("config[proxies] type = %T, want []map[string]any", config["proxies"])
	}
	if len(got) != 1 || got[0]["name"] != "test-node" {
		t.Fatalf("config[proxies] = %#v, want injected proxy", got)
	}
	if _, ok := config["proxy-providers"]; ok {
		t.Fatal("config must not keep proxy-providers")
	}
	if config["secret"] != "" {
		t.Fatalf("config[secret] = %v, want empty string", config["secret"])
	}
}
