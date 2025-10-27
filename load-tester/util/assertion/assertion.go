package assertion

import (
	"fmt"
	"reflect"
)

// AssertMapValue checks if a value in a map is of the expected type and non-empty (for certain types).
//
// It returns the value and an error if any check fails.
func AssertMapValue[M ~map[K]V, K comparable, V any, T any](m M, key K) (T, error) {
	var zero T
	value, ok := m[key]
	if !ok {
		return zero, fmt.Errorf("key %v was not found in map", key)
	}

	expectedType := reflect.TypeOf((*T)(nil)).Elem()
	valueType := reflect.TypeOf(value)
	valueValue := reflect.ValueOf(value)

	// Special case for map[string]any
	if expectedType == reflect.TypeOf(map[string]any{}) {
		if valueType.Kind() == reflect.Map &&
			valueType.Key().Kind() == reflect.String &&
			valueType.Elem().Kind() == reflect.Interface {
			return valueValue.Interface().(T), nil
		}

		return zero, fmt.Errorf("failed to assert %v as map[string]any, got type %v", key, valueType)
	}

	// For other types, check if the value can be converted to T
	if !valueType.ConvertibleTo(expectedType) {
		return zero, fmt.Errorf("failed to assert %v as %v, got type %v", key, expectedType, valueType)
	}

	// Convert the value
	convertedValue := valueValue.Convert(expectedType).Interface().(T)

	// Check if the value is empty, but only for types that have a notion of emptiness
	if isEmptyCheckRequired(expectedType) {
		if isEmpty(reflect.ValueOf(convertedValue)) {
			return zero, fmt.Errorf("value is empty")
		}
	}

	return convertedValue, nil
}

// isEmptyCheckRequired determines if a type should be checked for emptiness
func isEmptyCheckRequired(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return true
	default:
		return false
	}
}

// isEmpty checks if a value is empty
func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	default:
		return false
	}
}
