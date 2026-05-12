package api

// MomentConfig holds soundscape settings for setSmartButtonConfig.
type MomentConfig struct {
	SoundscapeID int // 1=Forest 2=Rain 3=Ocean 4=City
	Volume       int // percent_of_volume 0–100
	SleepTimer   int // seconds: 0=unlimited, 900=15m, 1800=30m, 2700=45m, 3600=60m
}

// ControlSoundscapeV2 starts (action_id=6) or stops (action_id=7) a soundscape.
func (c *Client) ControlSoundscapeV2(actionID, soundscapeID int) error {
	return c.post("controlSoundscapeV2", map[string]interface{}{
		"action_id":     actionID,
		"autoResume":    true,
		"fadeOut":       "true",
		"soundscape_id": soundscapeID,
	})
}

// SetSoundscapeElement sets the volume of a single sound element within a soundscape.
// sliderIndex is 0–2 (left-to-right UI order); device element_id is reversed (2-sliderIndex).
func (c *Client) SetSoundscapeElement(soundscapeID, sliderIndex, value int) error {
	return c.post("setSoundscapeV2Config", map[string]interface{}{
		"soundscape_id": soundscapeID,
		"element_id":    2 - sliderIndex, // device: 0=third slider, 1=second, 2=first
		"element_value": value,
	})
}

// SoundscapeElement is one volume control within a soundscape mode.
type SoundscapeElement struct {
	ID    FlexInt `json:"id"`
	Value FlexInt `json:"value"`
}

// SoundscapeV2Entry is a single mode entry returned by getSoundscapeV2Config.
type SoundscapeV2Entry struct {
	SoundscapeID FlexInt             `json:"soundscape_id"`
	Elements     []SoundscapeElement `json:"element_list"`
}

// GetSoundscapeV2Config fetches per-element volumes for all soundscape modes.
// Element IDs are device-reversed: id 2=first slider, id 1=second, id 0=third.
func (c *Client) GetSoundscapeV2Config() ([]SoundscapeV2Entry, error) {
	var resp struct {
		List []SoundscapeV2Entry `json:"soundscapeV2_list"`
	}
	return resp.List, c.get("getSoundscapeV2Config", &resp)
}

// SmartButtonState holds the parsed result of getSmartButtonConfig.
type SmartButtonState struct {
	SoundscapeID int
	SleepTimer   int
}

// GetSmartButtonConfig reads the current soundscape mode and timer from the device.
func (c *Client) GetSmartButtonConfig() (*SmartButtonState, error) {
	var raw struct {
		SmartConfig struct {
			Soundscape struct {
				ActiveSoundscapeID FlexInt `json:"active_soundscape_id"`
			} `json:"soundscape"`
			Timer struct {
				SleepTimer FlexInt `json:"sleep_timer"`
			} `json:"timer"`
		} `json:"smart_config"`
	}
	if err := c.get("getSmartButtonConfig", &raw); err != nil {
		return nil, err
	}
	return &SmartButtonState{
		SoundscapeID: int(raw.SmartConfig.Soundscape.ActiveSoundscapeID),
		SleepTimer:   int(raw.SmartConfig.Timer.SleepTimer),
	}, nil
}

// SetSmartButtonConfig pushes soundscape mode, volume, and sleep timer to the device.
func (c *Client) SetSmartButtonConfig(cfg MomentConfig) error {
	timerStatus := "disabled"
	if cfg.SleepTimer > 0 {
		timerStatus = "enabled"
	}
	return c.post("setSmartButtonConfig", map[string]interface{}{
		"atmos": map[string]interface{}{
			"atmos_level": 0,
			"status":      "off",
		},
		"music": map[string]interface{}{
			"album_cover": "",
			"music_id":    "",
		},
		"soundscape": map[string]interface{}{
			"active_soundscape_id": cfg.SoundscapeID,
			"supported_list":       []int{1, 2, 3, 4, 0},
		},
		"soundscapeV2": map[string]interface{}{
			"mix_with_music":    "disabled",
			"percent_of_volume": cfg.Volume,
		},
		"timer": map[string]interface{}{
			"sleep_timer":  cfg.SleepTimer,
			"timer_status": timerStatus,
		},
	})
}
