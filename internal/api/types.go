package api

import (
	"encoding/json"
	"strconv"
)

// FlexInt accepts both JSON strings ("42") and bare numbers (42).
// HK/Linkplay API uses both formats inconsistently across endpoints.
type FlexInt int

func (n *FlexInt) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*n = FlexInt(i)
		return nil
	}
	var i int
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	*n = FlexInt(i)
	return nil
}
