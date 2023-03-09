package utils

import "net/http"

func IfErr(message any) map[string]any {
	output := make(map[string]any, 2)
	output["status"] = http.StatusOK
	output["message"] = message
	return output
}
