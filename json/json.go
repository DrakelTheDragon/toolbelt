package json

import (
	"encoding/json"
)

func jsonify(obj any) ([]byte, error) {
	jsonBytes, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return nil, err
	}

	jsonBytes = append(jsonBytes, '\n')

	return jsonBytes, nil
}

func unjsonify(b []byte, dst any) error {
	return json.Unmarshal(b, &dst)
}

func JSONify(obj any) (Data, error) {
	jsonBytes, err := jsonify(obj)
	if err != nil {
		return nil, err
	}

	var jsonData Data

	if err := unjsonify(jsonBytes, &jsonData); err != nil {
		return nil, err
	}

	return jsonData, nil
}

func UnJSONify(data Data, dst any) error {
	return unjsonify([]byte(data.String()), dst)
}
