package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/akatranlp/sentinel/provider"
)

type ProviderConfig struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	IconURL      string `json:"icon_url"`
	BaseURL      string `json:"base_url"`
	Enabled      bool   `json:"enabled"`
}

func InitProviders() ([]provider.Provider, error) {
	f, err := os.Open("examples/basic/providers.json")
	if err != nil {
		return nil, err
	}

	var pCfgs []ProviderConfig
	if err := json.NewDecoder(f).Decode(&pCfgs); err != nil {
		return nil, err
	}

	var providers []provider.Provider

	for _, cfg := range pCfgs {
		if !cfg.Enabled {
			continue
		}
		provider, err := provider.ProviderFactory(provider.FactoryParams{
			Name:         cfg.Name,
			Slug:         cfg.Slug,
			BaseURL:      cfg.BaseURL,
			IconURL:      cfg.IconURL,
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Type:         cfg.Type,
		})
		if err != nil {
			fmt.Printf("error while initializing provider %s of type %s. Error: %v\n", cfg.Name, cfg.Type, err)
			return nil, err
		}
		providers = append(providers, provider)
	}

	return providers, nil
}
