package tolerantreader

import (
	"fmt"
	"github.com/oliveagle/jsonpath"
	"math"
	"reflect"
	"strconv"
	"time"
)

const ONLYDATE = "2006-01-02"
const PREFIX = "jsonpath"

func assign(target reflect.Value, val reflect.Value) error {
	if target.Kind() != val.Kind() {
		return fmt.Errorf("Target %s is not matching value %s.", target.Kind(), val.Kind())
	}
	target.Set(val)
	return nil
}

func Unmarshal(data map[string]interface{}, o interface{}) error {
	if reflect.TypeOf(o).Kind() != reflect.Ptr {
		return fmt.Errorf("Expected Pointer, got %s", reflect.TypeOf(o).Kind())
	}
	v := reflect.ValueOf(o).Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("Expected Ptr to Struct, got ptr to %s", v.Kind())
	}

	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		tag, found := typeOfS.Field(i).Tag.Lookup(PREFIX)
		if found {
			res, err := jsonpath.JsonPathLookup(data, tag)
			if err != nil {
				return err
			}

			switch v.Field(i).Kind() {
			case reflect.String:
				if reflect.ValueOf(res).Kind() == reflect.Float64 {

					f := reflect.ValueOf(res).Float()
					if f == math.Floor(f) {
						err := assign(v.Field(i), reflect.ValueOf(fmt.Sprintf("%.0f", res)))
						if err != nil {
							return err
						}
					} else {
						err := assign(v.Field(i), reflect.ValueOf(fmt.Sprintf("%f", res)))
						if err != nil {
							return err
						}
					}
				} else {
					err := assign(v.Field(i), reflect.ValueOf(res))
					if err != nil {
						return err
					}
				}
				break
			case reflect.Float64:
				if reflect.ValueOf(res).Kind() == reflect.String {
					s, err := strconv.ParseFloat(res.(string), 64)
					if err != nil {
						return err
					}
					err = assign(v.Field(i), reflect.ValueOf(s))
					if err != nil {
						return err
					}
				} else {
					err := assign(v.Field(i), reflect.ValueOf(res))
					if err != nil {
						return err
					}
				}
				break
			case reflect.Int:
				if reflect.ValueOf(res).Kind() == reflect.String {
					s, err := strconv.ParseInt(res.(string), 10, 64)
					if err != nil {
						return fmt.Errorf("Field %s : Failed to convert String to int", typeOfS.Field(i).Name)
					}
					v.Field(i).SetInt(int64(s))
					break
				} else if reflect.ValueOf(res).Kind() != reflect.Float64 {
					return fmt.Errorf("Field %s : Int must be numeric in Json, got %s",
						typeOfS.Field(i).Name,
						reflect.ValueOf(res).Kind())
				}
				v.Field(i).SetInt(int64(res.(float64)))
				break
			case reflect.Slice:
				if reflect.ValueOf(res).Kind() != reflect.Slice {
					return fmt.Errorf("Cannot handle Type %s as slice", reflect.ValueOf(res))
				}
				r2 := res.([]interface{})
				v.Field(i).Set(reflect.MakeSlice(v.Field(i).Type(), len(r2), len(r2)))
				if len(r2) < 1 {
					break
				}
				for j, w := range r2 {
					err := assign(v.Field(i).Index(j), reflect.ValueOf(w))
					if err != nil {
						return fmt.Errorf("%s in field '%s' at index %d", err, typeOfS.Field(i).Name, j)
					}
				}
				break
			case reflect.Struct:
				n := typeOfS.Field(i).Type.String()
				if n == "time.Time" {
					if reflect.TypeOf(res).Kind() != reflect.String {
						return fmt.Errorf("time must be 'string' in field %s", typeOfS.Field(i).Name)
					}
					var tval time.Time
					tval, err = time.Parse(time.RFC3339, res.(string))
					if err != nil {
						tval, err = time.Parse(ONLYDATE, res.(string))
						if err != nil {
							layout := "2006-01-02 15:04:05"
							tval, err = time.Parse(layout, res.(string))
							if err != nil {
								return fmt.Errorf("invalid time value in field %s: %s", typeOfS.Field(i).Name, err)
							}
						}
					}
					v.Field(i).Set(reflect.ValueOf(tval))
				}
				break
			}
		} else {
			fmt.Printf("Field: %s\tType: %s  NO TAG\n", typeOfS.Field(i).Name, v.Field(i).Type().String())
		}
	}
	return nil
}
