package utils

import "strings"

func CleanResp(response string) string {
	raw := strings.TrimSpace(response)

	if strings.HasPrefix(raw, "```json") {
		raw = strings.TrimPrefix(raw, "```json")
	} else if strings.HasPrefix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```json")
	}
	if strings.HasSuffix(raw, "```") {
		raw = strings.TrimSuffix(raw, "```")
	}

	raw = strings.TrimSpace(raw)
	return raw
}
