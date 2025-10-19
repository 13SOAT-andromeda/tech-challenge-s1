package utils

import (
	"net/url"
)

func ParamsToMap(params url.Values) map[string]interface{} {
	paramsMap := make(map[string]interface{})
	for key, value := range params {
		paramsMap[key] = value[0]
	}

	return paramsMap
}
