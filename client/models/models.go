package models

type AIFuncCall struct {
	Input string                 `json:"input"`
	Meta  map[string]interface{} `json:"meta"`
}

type AIFuncResponse struct {
	ID       string `json:"id"`
	Response string `json:"response"`
}
