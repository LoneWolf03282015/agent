package proxy

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func reproduceDockerAPIResponse(data []interface{}, requestPath string) interface{} {
	// VolumeList operation returns an object, not an array.
	if strings.HasPrefix(requestPath, "/volumes") {
		responseObject := make(map[string]interface{})
		responseObject["Volumes"] = data
		return responseObject
	}

	return data
}

func responseToJSONArray(response *http.Response, requestPath string) ([]interface{}, error) {
	responseObject, err := getResponseBodyAsGenericJSON(response)
	if err != nil {
		return nil, err
	}

	obj, ok := responseObject.(map[string]interface{})
	if ok && obj["message"] != nil {
		return nil, errors.New(obj["message"].(string))
	}

	var responseData []interface{}

	// VolumeList operation returns an object, not an array.
	// We need to extract the volume list from the "Volumes" property.
	// Note that the content of the "Volumes" property might be null if no volumes
	// are found, we replace it with an empty array in that case.
	if strings.HasPrefix(requestPath, "/volumes") {
		obj := responseObject.(map[string]interface{})
		volumeObj := obj["Volumes"]
		if volumeObj != nil {
			responseData = volumeObj.([]interface{})
		} else {
			responseData = make([]interface{}, 0)
		}
	} else {
		responseData = responseObject.([]interface{})
	}

	return responseData, nil
}

func getResponseBodyAsGenericJSON(response *http.Response) (interface{}, error) {
	var data interface{}

	err := json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
