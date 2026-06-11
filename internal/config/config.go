package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type CustomPreset struct {
	Name string    `json:"name"`
	Gain []float64 `json:"gain"`
}

type Config struct {
	CustomPresets [3]CustomPreset `json:"custom_presets,omitempty"`
	Language      string          `json:"language,omitempty"`
	LedKeepAlive  int             `json:"led_keep_alive,omitempty"` // 0=off, 30, 60 (minutes)
	Theme         string          `json:"theme,omitempty"`           // "light", "dark", "auto"
	DeviceIP      string          `json:"device_ip,omitempty"`       // last-known speaker IP, tried before mDNS
}

func (c *Config) initDefaults() {
	for i := range c.CustomPresets {
		if c.CustomPresets[i].Name == "" {
			c.CustomPresets[i].Name = fmt.Sprintf("Custom %d", i+1)
		}
		if c.CustomPresets[i].Gain == nil {
			c.CustomPresets[i].Gain = make([]float64, 7)
		}
	}
	if c.Language == "" {
		c.Language = "en"
	}
	if c.Theme == "" {
		c.Theme = "auto"
	}
}

func Dir() (string, error) {
	base, err := os.UserConfigDir() // ~/.config on Linux, %APPDATA% on Windows, ~/Library/... on mac
	if err != nil {
		return "", err
	}
	d := filepath.Join(base, "soundsticks")
	return d, os.MkdirAll(d, 0700)
}

func Load() (*Config, error) {
	d, err := Dir()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filepath.Join(d, "config.json"))
	if os.IsNotExist(err) {
		c := &Config{}
		c.initDefaults()
		return c, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var c Config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}
	c.initDefaults()
	return &c, nil
}

func (c *Config) Save() error {
	d, err := Dir()
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(d, "config.json"))
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
}
