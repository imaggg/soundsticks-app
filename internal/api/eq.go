package api

import (
	"encoding/json"
	"fmt"
)

// builtinPresets — firmware-exact fs+gain for each preset.
// Must match verbatim: device halves gain in getEQList storage,
// but setActiveEQ requires the original full values to recognise the preset.
var builtinPresets = map[int]EQPayload{
	1: {Fs: []float64{125, 250, 500, 1000, 2000, 4000, 8000}, Gain: []float64{0, 0, 0, 0, 0, 0, 0}},
	2: {Fs: []float64{80, 250, 500, 1000, 9000, 4000, 9000}, Gain: []float64{-5, 0, 0, 4, -2, 4, 0}},
	3: {Fs: []float64{60, 200, 500, 1500, 2000, 7000, 8000}, Gain: []float64{-2, 7, 0, 3, 0, 2, 0}},
	4: {Fs: []float64{200, 250, 500, 1000, 1500, 10000, 6000}, Gain: []float64{-4, 0, 0, 0, -3, 0, -3}},
}

type EQPayload struct {
	Fs   []float64 `json:"fs"`
	Gain []float64 `json:"gain"`
}

type EQEntry struct {
	Band      FlexInt   `json:"band"`
	EQID      FlexInt   `json:"eq_id"`
	EQName    string    `json:"eq_name"`
	EQPayload EQPayload `json:"eq_payload"`
}

type EQListResponse struct {
	ActiveEQID FlexInt         `json:"active_eq_id"`
	EQList     []EQEntry       `json:"eq_list"`
	ErrorCode  json.RawMessage `json:"error_code"`
}

type PlayerStatus struct {
	Status string  `json:"status"`
	Vol    FlexInt `json:"vol"`
	Mute   FlexInt `json:"mute"`
	EQ     FlexInt `json:"eq"`
}

func (c *Client) GetEQList() (*EQListResponse, error) {
	var resp EQListResponse
	return &resp, c.get("getEQList", &resp)
}

func (c *Client) SetActiveEQ(eqID int) error {
	p, ok := builtinPresets[eqID]
	if !ok {
		return fmt.Errorf("unknown preset %d", eqID)
	}
	return c.post("setActiveEQ", map[string]any{
		"active_eq_id": fmt.Sprintf("%d", eqID),
		"band":         7,
		"eq_payload":   p,
	})
}

func (c *Client) SetCustomEQ(fs []float64, gain []float64) error {
	scaled := make([]float64, len(gain))
	for i, g := range gain {
		// 125Hz (index 0) has a larger negative range on firmware: -12 UI → -9 API.
		// All other values (and positive 125Hz) use the standard UI/2 factor.
		if i == 0 && g < 0 {
			scaled[i] = g * 9.0 / 12.0
		} else {
			scaled[i] = g / 2
		}
	}
	return c.post("setActiveEQ", map[string]any{
		"active_eq_id": "0",
		"band":         7,
		"eq_payload":   EQPayload{Fs: fs, Gain: scaled},
	})
}

func (c *Client) GetDeviceInfo() (map[string]any, error) {
	var info map[string]any
	return info, c.get("getDeviceInfo", &info)
}

func (c *Client) GetPlayerStatus() (*PlayerStatus, error) {
	var s PlayerStatus
	return &s, c.get("getPlayerStatus", &s)
}

func (c *Client) SetPlayerVol(vol int) error {
	return c.rawGet(fmt.Sprintf("setPlayerCmd:vol:%d", vol))
}

func (c *Client) SetPlayerMute(mute bool) error {
	v := 0
	if mute {
		v = 1
	}
	return c.rawGet(fmt.Sprintf("setPlayerCmd:mute:%d", v))
}
