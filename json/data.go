
package json

type Data map[string]any

func (d Data) String() string {
	jsonBytes, err := jsonify(d)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}