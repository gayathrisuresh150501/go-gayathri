package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type Student struct {
	Name    string
	Age     int
	Classes []string
}

func structToJSONDefinition(data interface{}) (string, error) {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("input is not a struct")
	}

	result := make(map[string]interface{})
	fields := make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		fieldInfo := map[string]interface{}{
			"type": field.Type.String(),
		}

		// Handle maps
		if field.Type.Kind() == reflect.Map {
			fieldInfo["key_type"] = field.Type.Key().String()
			fieldInfo["value_type"] = field.Type.Elem().String()
		}

		// Handle slices/arrays
		if field.Type.Kind() == reflect.Slice {
			elemType := field.Type.Elem()
			fieldInfo["element_type"] = elemType.String()
			// If the slice contains structs, process them
			if elemType.Kind() == reflect.Struct {
				embeddedFields, err := structToJSONDefinition(reflect.New(elemType).Elem().Interface())
				if err != nil {
					return "", err
				}
				fieldInfo["element_fields"] = json.RawMessage(embeddedFields)
			}
		}

		// Handle embedded structs
		if field.Type.Kind() == reflect.Struct {
			embeddedStruct := reflect.New(field.Type).Elem().Interface()
			embeddedFields, err := structToJSONDefinition(embeddedStruct)
			if err != nil {
				return "", err
			}
			fieldInfo["fields"] = json.RawMessage(embeddedFields)
		}

		// Add JSON tag if present
		if tag := field.Tag.Get("json"); tag != "" {
			fieldInfo["json_tag"] = tag
		}

		fields[field.Name] = fieldInfo
	}
	result[t.Name()] = map[string]interface{}{
		"fields": fields,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func main() {
	// Generate JSON representation of the struct
	jsonDef, err := structToJSONDefinition(Student{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Write JSON to a file
	file, err := os.Create("output.json")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(jsonDef)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Println("JSON definition written to output.json")
}
