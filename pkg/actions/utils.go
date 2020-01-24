package actions

import "fmt"

func makeConfigMap(config map[interface{}]interface{}) map[string]string {
	mapString := make(map[string]string)

	for key, value := range config {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		mapString[strKey] = strValue
	}

	return mapString
}
