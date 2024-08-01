package gotil

import "encoding/json"

// Convert map / struct / slice to string
func JsonStringify(data any) string {
	b, _ := json.Marshal(data)
	return string(b)
}
