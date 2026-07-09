package clashtemplate

import (
	"fmt"
	"strings"

	yaml "github.com/goccy/go-yaml"
)

func Clean(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", nil
	}
	config, err := Parse(raw)
	if err != nil {
		return "", err
	}
	CleanConfig(config)
	out, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)) + "\n", nil
}

func Parse(raw string) (map[string]any, error) {
	var config map[string]any
	if err := yaml.Unmarshal([]byte(raw), &config); err != nil {
		return nil, err
	}
	if len(config) == 0 {
		return nil, fmt.Errorf("empty Clash/Mihomo YAML template")
	}
	return config, nil
}

func CleanConfig(config map[string]any) {
	delete(config, "proxies")
	delete(config, "proxy-providers")
	delete(config, "proxyProviders")
	config["secret"] = ""
	cleanProxyGroups(config["proxy-groups"])
}

func cleanProxyGroups(groups any) {
	list, ok := groups.([]any)
	if !ok {
		return
	}
	for _, item := range list {
		group, ok := item.(map[string]any)
		if !ok {
			continue
		}
		_, hadUse := group["use"]
		delete(group, "use")
		if !hadUse {
			continue
		}
		group["include-all"] = true
	}
}
