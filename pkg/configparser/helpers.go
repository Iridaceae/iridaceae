package configparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func toStrVal(i interface{}) string {
	switch t := i.(type) {
	case string:
		return t
	case int:
		return strconv.FormatInt(int64(t), 10)
	case bool:
		return strconv.FormatBool(t)
	case float64:
		return strconv.FormatFloat(t, 'f', 3, 64)
	case fmt.Stringer:
		return t.String()
	default:
		panic("cannot convert given input to string")
	}
}

func toIntVal(i interface{}) int {
	switch t := i.(type) {
	case string:
		n, ok := strconv.Atoi(t)
		if ok == nil {
			return n
		}
		panic("cannot convert string to int")
	case int, float64:
		return t.(int)
	case bool:
		if t {
			return 1
		}
		return 0
	default:
		panic("cannot convert given input to int")
	}
}

func toFloat64Val(i interface{}) float64 {
	switch t := i.(type) {
	case string:
		n, ok := strconv.ParseFloat(t, 64)
		if ok == nil {
			return n
		}
		panic("cannot convert string to float64")
	case int:
		return float64(t)
	case float64:
		return t
	default:
		panic("cannot convert given input to float64")
	}
}

func toBoolVal(i interface{}) bool {
	switch t := i.(type) {
	case string:
		lower := strings.ToLower(strings.TrimSpace(t))
		if lower == "true" || lower == "yes" || lower == "on" || lower == "enabled" || lower == "1" {
			return true
		}
		return false
	case int, float64:
		return t.(int) > 0
	case bool:
		return t
	default:
		panic("cannot convert given input to bool")
	}
}

func matchOptionsRegex(key string) (bool, error) {
	b, _ := regexp.MatchString(OptionsRegex, key)
	if b {
		return b, nil
	}
	return b, ErrInvalidOptionsMatch
}
