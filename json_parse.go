package gotil

// parse json string to map / slice / struct
func JsonParse(obj any, data string) error {
	return ParseStruct(obj, data, "json")
}
