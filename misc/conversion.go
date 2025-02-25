package misc

import "fmt"

func AnyToInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case float64: // Sometimes MongoDB stores numbers as float64
		return int(v), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}
