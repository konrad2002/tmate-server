package misc

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

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

// ConvertToBSOND Recursively converts JSON-decoded BSONElement array into bson.D
func ConvertToBSOND(data interface{}) bson.D {
	doc := bson.D{}

	// Ensure we're working with an array of maps
	if elements, ok := data.([]interface{}); ok {
		for _, elem := range elements {
			if kv, ok := elem.(map[string]interface{}); ok {
				// Ensure it has "Key" and "Value"
				if key, keyOk := kv["Key"].(string); keyOk {
					doc = append(doc, bson.E{Key: key, Value: convertToBSONValue(kv["Value"])})
				}
			}
		}
	}

	return doc
}

// Helper function to process BSON values (including nested bson.D and bson.A)
func convertToBSONValue(value interface{}) interface{} {
	switch v := value.(type) {
	case []interface{}:
		// If it's an array, process it as bson.A
		arr := bson.A{}
		for _, item := range v {
			arr = append(arr, convertToBSONValue(item))
		}
		return arr
	case map[string]interface{}:
		// If it looks like a bson.D element, convert it properly
		if _, hasKey := v["Key"]; hasKey {
			return ConvertToBSOND([]interface{}{v}) // Convert as bson.D
		}
		// Otherwise, treat it as bson.M
		doc := bson.D{}
		for key, val := range v {
			doc = append(doc, bson.E{Key: key, Value: convertToBSONValue(val)})
		}
		return doc
	default:
		// Return primitive values as they are
		return v
	}
}
