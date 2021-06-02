package utils

import (
	"encoding/json"
	"fmt"
)

//Wrapper for success response and failure response
func ResponseSuccess(v interface{}) []byte {
	rs := struct {
		Data interface{} `json:"data"`
	}{Data: v}
	if res, err := json.Marshal(rs); err != nil {
		return ResponseFailure(err.Error())
	} else {
		return res
	}
}

func ResponseFailure(s string) []byte {
	errorS := fmt.Sprintf(`{"error":"%s"}`, s)
	return []byte(errorS)
}
