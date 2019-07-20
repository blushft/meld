package utility

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Attr(v interface{}) map[string]string {
	m := make(map[string]string)
	var vv reflect.Type
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		vv = reflect.TypeOf(v).Elem()
	} else {
		vv = reflect.TypeOf(v)
	}

	m["type"] = vv.String()
	m["value"] = vv.String()
	m["kind"] = vv.Kind().String()
	m["pkg"] = vv.PkgPath()
	m["name"] = vv.Name()

	if vv.Kind() == reflect.Func {
		m["variadic"] = strconv.FormatBool(vv.IsVariadic())
		m["numIn"] = strconv.Itoa(vv.NumIn())
		m["numOut"] = strconv.Itoa(vv.NumOut())
		for i := 0; i < vv.NumIn(); i++ {
			in := vv.In(i)
			kind := in.Kind().String()
			if in.Kind() == reflect.Ptr {
				kind = "*" + in.Elem().Kind().String()
			}
			if in.Kind() == reflect.Interface {
				kind = in.PkgPath() + "#" + in.Name()
			}
			m["arg"+strconv.Itoa(i)] = kind
		}
		for i := 0; i < vv.NumOut(); i++ {
			out := vv.Out(i)

			kind := out.Kind().String()
			if out.Kind() == reflect.Ptr {
				kind = "*" + out.Elem().Kind().String()
			}
			if out.Kind() == reflect.Interface {
				kind = out.PkgPath() + "#" + out.Name()
			}

			m["ret"+strconv.Itoa(i)] = kind
		}
	}
	return m
}

func BindInterface(v interface{}, dst interface{}) error {
	switch d := dst.(type) {
	case *string:
		switch vt := v.(type) {
		case string:
			*d = vt
		case []string:
			*d = strings.Join(vt, ",")
		case bool:
			*d = strconv.FormatBool(vt)
		case int64:
			*d = strconv.FormatInt(vt, 10)
		case float64:
			*d = strconv.FormatFloat(vt, 'f', -1, 64)
		case int:
			*d = strconv.FormatInt(int64(vt), 10)
		default:
			panic(fmt.Sprintf("Unhandled type %+v for string conversion", reflect.TypeOf(vt)))
		}

	case *[]string:
		switch vt := v.(type) {
		case []interface{}:
			*d = []string{}
			for _, v := range vt {
				*d = append(*d, v.(string))
			}
		case []string:
			*d = vt
		case string:
			*d = strings.Split(vt, ",")
		default:
			panic(fmt.Sprintf("Unhandled type %+v for []string conversion", reflect.TypeOf(vt)))
		}

	case *bool:
		switch vt := v.(type) {
		case bool:
			*d = vt
		case string:
			vt = strings.ToLower(vt)
			if vt == "true" || vt == "on" || vt == "yes" || vt == "1" || vt == "t" {
				*d = true
			} else {
				*d = false
			}
		case int64:
			*d = vt > 0
		case int:
			*d = vt > 0
		case uint64:
			*d = vt > 0
		case uint:
			*d = vt > 0
		case float64:
			*d = vt > 0
		default:
			panic(fmt.Sprintf("Unhandled type %+v for bool conversion", reflect.TypeOf(vt)))
		}

	case *float64:
		switch vt := v.(type) {
		case int64:
			*d = float64(vt)
		case int32:
			*d = float64(vt)
		case int16:
			*d = float64(vt)
		case int8:
			*d = float64(vt)
		case int:
			*d = float64(vt)
		case uint64:
			*d = float64(vt)
		case uint32:
			*d = float64(vt)
		case uint16:
			*d = float64(vt)
		case uint8:
			*d = float64(vt)
		case uint:
			*d = float64(vt)
		case float64:
			*d = vt
		case float32:
			*d = float64(vt)
		case string:
			x, _ := strconv.ParseFloat(vt, 64)
			*d = float64(x)
		default:
			panic(fmt.Sprintf("Unhandled type %+v for float64 conversion", reflect.TypeOf(vt)))
		}

	case *int:
		switch vt := v.(type) {
		case int64:
			*d = int(vt)
		case int32:
			*d = int(vt)
		case int16:
			*d = int(vt)
		case int8:
			*d = int(vt)
		case int:
			*d = vt
		case uint64:
			*d = int(vt)
		case uint32:
			*d = int(vt)
		case uint16:
			*d = int(vt)
		case uint8:
			*d = int(vt)
		case uint:
			*d = int(vt)
		case float64:
			*d = int(vt)
		case float32:
			*d = int(vt)
		case string:
			*d, _ = strconv.Atoi(vt)
		default:
			panic(fmt.Sprintf("Unhandled type %+v for int conversion", reflect.TypeOf(vt)))
		}

	case *time.Time:
		switch vt := v.(type) {
		case time.Time:
			*d = vt
		case int:
			*d = time.Unix(int64(vt), 0)
		case int64:
			*d = time.Unix(vt, 0)
		default:
			panic(fmt.Sprintf("Unhandled type %+v for time.Time conversion", reflect.TypeOf(vt)))
		}

	case *url.Values:
		switch vt := v.(type) {
		case string:
			*d, _ = url.ParseQuery(vt)
		default:
			panic(fmt.Sprintf("Unhandled type %+v for url.Values conversion", reflect.TypeOf(vt)))
		}

	default:
		panic(fmt.Sprintf("Unhandled dst type %+v", reflect.TypeOf(dst)))
	}

	return nil
}
