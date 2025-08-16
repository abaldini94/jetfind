package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

const APPNAME = "jetfind"

var Default *Config = &Config{
	Filter: FilterConfig{
		Type:       "fuzzy",
		Algo:       "jarowinkler",
		Threashold: 0.9,
	},
	Findignore: FindIgnoreConfig{
		Enable:       false,
		HiddenIgnore: false,
	},
	Tui: TuiConfig{
		HighlightedFile: HighlightedFileConfig{
			Foreground: "#7DD3FC",
		},
		QueryBox: QueryBoxConfig{
			TextForeground:   "#D1D5DB",
			TextBackground:   "#1A1B23",
			BorderForeground: "#6366F1",
		},
	},
}

type Config struct {
	Filter     FilterConfig     `yaml:"filter"`
	Findignore FindIgnoreConfig `yaml:"findignore"`
	Tui        TuiConfig        `yaml:"tui"`
}

type FilterConfig struct {
	Type       string  `yaml:"type"`
	Algo       string  `yaml:"algorithm"`
	Threashold float64 `yaml:"threashold"`
}

type FindIgnoreConfig struct {
	Enable       bool `yaml:"enable"`
	HiddenIgnore bool `yaml:"hidden_ignore"`
}

type TuiConfig struct {
	HighlightedFile HighlightedFileConfig `yaml:"highlighted_file"`
	QueryBox        QueryBoxConfig        `yaml:"query_box"`
}

type HighlightedFileConfig struct {
	Foreground string `yaml:"foreground"`
	Background string `yaml:"background"`
}

type QueryBoxConfig struct {
	TextForeground   string `yaml:"text_foreground"`
	TextBackground   string `yaml:"text_background"`
	BorderForeground string `yaml:"border_foreground"`
}

func GetConfigDir() string {
	return filepath.Join(xdg.ConfigHome, APPNAME)
}

func GetConfigFilePath() (string, error) {
	cfgPath, err := xdg.ConfigFile(filepath.Join(APPNAME, "config.yml"))
	if err != nil {
		return "", fmt.Errorf("error while getting the configuration file path")
	}
	return cfgPath, nil
}

func Load(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func LoadOrDefault() *Config {
	var cfg *Config

	cfgPath, err := GetConfigFilePath()
	if err != nil {
		cfg = Default
	} else {
		cfg, err = Load(cfgPath)
		if err != nil {
			cfg = Default
		}
	}

	if reflect.DeepEqual(cfg.Tui, TuiConfig{}) {
		cfg.Tui = Default.Tui
	} else {
		if reflect.DeepEqual(cfg.Tui.HighlightedFile, HighlightedFileConfig{}) {
			cfg.Tui.HighlightedFile = Default.Tui.HighlightedFile
		}
		if reflect.DeepEqual(cfg.Tui.QueryBox, QueryBoxConfig{}) {
			cfg.Tui.QueryBox = Default.Tui.QueryBox
		}
	}

	if reflect.DeepEqual(cfg.Filter, FilterConfig{}) {
		cfg.Filter = Default.Filter
	}

	if reflect.DeepEqual(cfg.Findignore, FindIgnoreConfig{}) {
		cfg.Findignore = Default.Findignore
	}
	return cfg
}
