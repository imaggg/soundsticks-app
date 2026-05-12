package api

import "encoding/json"

type Pattern struct {
	ID    FlexInt `json:"id"`
	Level FlexInt `json:"level"`
}

type PatternSupport struct {
	ID          FlexInt `json:"id"`
	Level       FlexInt `json:"level"`
	SupportWand bool    `json:"support_wand"`
}

type LightInfoData struct {
	Enable        FlexInt          `json:"enable"`
	Brightness    FlexInt          `json:"brightness"`
	DynamicLevel  FlexInt          `json:"dynamic_level"`
	ActivePattern Pattern          `json:"active_pattern"`
	LightElement  FlexInt          `json:"light_element"`
	SupportList   []PatternSupport `json:"support_list"`
}

type LightInfoResponse struct {
	LightInfo LightInfoData   `json:"light_info"`
	ErrorCode json.RawMessage `json:"error_code"`
}

func (c *Client) GetLightInfo() (*LightInfoData, error) {
	var resp LightInfoResponse
	return &resp.LightInfo, c.get("getLightInfo", &resp)
}

// SetLightInfo sends a partial light update.
// Payload is sent flat (no light_info wrapper), values as strings — device format.
// Example: {"brightness":"80"} or {"active_pattern":{"id":"1","level":"50"}}
func (c *Client) SetLightInfo(raw map[string]json.RawMessage) error {
	return c.post("setLightInfo", raw)
}
