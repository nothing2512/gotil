package gotil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// obj : *struct
// data : JSON
// tag : tagname (optional / "")
func ParseStruct(obj any, data any, tag string) {
	_v := reflect.ValueOf(obj)
	_t := reflect.New(_v.Type().Elem()).Elem().Type()

	if _t.Kind() == reflect.Map {
		b, _ := json.Marshal(data)
		_ = json.Unmarshal(b, _v.Elem().Addr().Interface())
		return
	}

	if _t.Kind() == reflect.Slice {
		var temp []any
		b, _ := json.Marshal(data)
		_ = json.Unmarshal(b, &temp)
		var slice []any
		for _, x := range temp {
			__t := reflect.New(_v.Type().Elem().Elem())
			_n := __t.Interface()
			if reflect.TypeOf(x).Kind() == reflect.Ptr || reflect.TypeOf(x).Kind() == reflect.Map {
				ParseStruct(_n, x, tag)
				slice = append(slice, _n)
			} else {
				slice = append(slice, x)
			}
		}
		b, _ = json.Marshal(slice)
		_ = json.Unmarshal(b, _v.Elem().Addr().Interface())
		return
	}

	var temp JSON
	b, _ := json.Marshal(data)
	_ = json.Unmarshal(b, &temp)

	if _t.Kind() == reflect.Ptr {
		_t = _t.Elem()
		_v = reflect.ValueOf(obj).Elem()
	}

	i := 0
	for i < _t.NumField() {
		serialize := false
		_field := _t.Field(i)
		key := _field.Name
		if _field.Tag.Get("json") != "" {
			key = strings.Split(_field.Tag.Get("json"), ",")[0]
		}

		for _, tag := range strings.Split(_field.Tag.Get(tag), ";") {
			if !strings.Contains(tag, ":") {
				key = tag
			}
			if tag == "serializer:json" {
				serialize = true
			}
			if _k := strings.ReplaceAll(tag, "column:", ""); _k != tag {
				key = _k
			}
		}

		vField := _v.Elem().Field(i)
		if vField.Kind() == reflect.Struct {
			bd, _ := json.Marshal(temp)
			_ = json.Unmarshal(bd, vField.Addr().Interface())
		}

		for k, v := range temp {
			if v == nil {
				continue
			}
			if k == key && vField.CanSet() {
				if serialize {
					_ = json.Unmarshal([]byte(fmt.Sprintf("%v", v)), vField.Addr().Interface())
				} else {
					if vField.Kind() == reflect.Int {
						if reflect.TypeOf(v).Kind() == reflect.Float64 {
							v = int(v.(float64))
						}
						v, _ = strconv.Atoi(fmt.Sprintf("%v", v))
					}

					if vField.Kind() == reflect.String {
						v = fmt.Sprintf("%v", v)
					}

					_i := 0
					_s := ""
					if vField.Type() == reflect.TypeOf(&_i) {
						xv, _ := strconv.Atoi(fmt.Sprintf("%v", v))
						vField.Set(reflect.ValueOf(&xv))
					} else if vField.Type() == reflect.TypeOf(&_s) {
						xv := fmt.Sprintf("%v", v)
						vField.Set(reflect.ValueOf(&xv))
					} else if vField.Type() == reflect.TypeOf(time.Time{}) {
						for _, y := range []string{
							time.RFC822,
							time.StampMilli,
							"2006-01-02T15:04:05",
							"2006-01-02 15:04:05",
							"2006-01-02",
							time.RFC3339,
						} {
							x, err := time.Parse(y, fmt.Sprintf("%v", v))
							if err == nil {
								vField.Set(reflect.ValueOf(x))
								break
							}
						}
					} else if vField.Type() == reflect.TypeOf(&time.Time{}) {
						for _, y := range []string{
							time.RFC822,
							time.StampMilli,
							"2006-01-02T15:04:05",
							"2006-01-02T15:04:05.505Z",
							"2006-01-02 15:04:05",
							"2006-01-02",
							time.RFC3339,
						} {
							x, err := time.Parse(y, fmt.Sprintf("%v", v))
							if err == nil {
								vField.Set(reflect.ValueOf(&x))
								break
							}
						}
					} else if v != nil {
						vField.Set(reflect.ValueOf(v))
					}
				}
				break
			}
		}
		i++
	}
}
