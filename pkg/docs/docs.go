package docs

import (
	"backend-challenge/internal/configs"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type OpenAPI map[string]any

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

func getServers(cfg *configs.Setting) []Server {
	if cfg.App.Env == "development" || cfg.App.Env == "production" {
		return []Server{
			{URL: cfg.App.Url, Description: "Staging"},
		}
	}
	return []Server{
		{URL: fmt.Sprintf("http://localhost:%d", cfg.App.Port), Description: "Local development"},
	}
}

func MergeSpecs(cfg *configs.Setting, baseFile, modulesDir string) ([]byte, error) {
	baseData, err := os.ReadFile(baseFile)
	if err != nil {
		return nil, fmt.Errorf("read base file error: %w", err)
	}

	var base OpenAPI
	if err := json.Unmarshal(baseData, &base); err != nil {
		return nil, err
	}

	base["servers"] = getServers(cfg)

	if base["paths"] == nil {
		base["paths"] = make(map[string]any)
	}
	if base["components"] == nil {
		base["components"] = map[string]any{"schemas": make(map[string]any)}
	}

	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil, fmt.Errorf("read modules dir error: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(modulesDir, entry.Name()))
		if err != nil {
			return nil, err
		}

		var module OpenAPI
		if err := json.Unmarshal(data, &module); err != nil {
			return nil, err
		}

		if mPaths, ok := module["paths"].(map[string]any); ok {
			for k, v := range mPaths {
				base["paths"].(map[string]any)[k] = v
			}
		}

		if mComp, ok := module["components"].(map[string]any); ok {
			if mSchemas, ok := mComp["schemas"].(map[string]any); ok {
				bComp := base["components"].(map[string]any)
				if bComp["schemas"] == nil {
					bComp["schemas"] = make(map[string]any)
				}
				bSchemas := bComp["schemas"].(map[string]any)
				for k, v := range mSchemas {
					bSchemas[k] = v
				}
			}
		}
	}

	return json.Marshal(base)
}
