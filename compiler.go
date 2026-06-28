package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JSONToJSCompiler struct {
	VariableName string
}

func NewCompiler(varName string) *JSONToJSCompiler {
	if varName == "" {
		varName = "dynamicData"
	}
	return &JSONToJSCompiler{VariableName: varName}
}

func (c *JSONToJSCompiler) Compile(jsonStr string) (string, error) {
	var rawData interface{}
	err := json.Unmarshal([]byte(jsonStr), &rawData)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	jsObject, err := c.formatJS(rawData, 0)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("const %s = %s;\n\n", c.VariableName, jsObject))
	sb.WriteString(fmt.Sprintf("export default %s;\n", c.VariableName))

	return sb.String(), nil
}

func (c *JSONToJSCompiler) formatJS(data interface{}, indentLevel int) (string, error) {
	indent := strings.Repeat("  ", indentLevel)
	nextIndent := strings.Repeat("  ", indentLevel+1)

	switch v := data.(type) {
	case map[string]interface{}:
		if len(v) == 0 {
			return "{}", nil
		}
		var lines []string
		for key, val := range v {
			formattedVal, err := c.formatJS(val, indentLevel+1)
			if err != nil {
				return "", err
			}
			escapedKey := key
			if strings.ContainsAny(key, " -+/\\*") {
				escapedKey = fmt.Sprintf(`"%s"`, key)
			}
			lines = append(lines, fmt.Sprintf("%s%s: %s", nextIndent, escapedKey, formattedVal))
		}
		return fmt.Sprintf("{\n%s\n%s}", strings.Join(lines, ",\n"), indent), nil

	case []interface{}:
		if len(v) == 0 {
			return "[]", nil
		}
		var items []string
		for _, val := range v {
			formattedVal, err := c.formatJS(val, indentLevel+1)
			if err != nil {
				return "", err
			}
			items = append(items, nextIndent+formattedVal)
		}
		return fmt.Sprintf("[\n%s\n%s]", strings.Join(items, ",\n"), indent), nil

	case string:
		escaped := strings.ReplaceAll(v, `"`, `\"`)
		return fmt.Sprintf(`"%s"`, escaped), nil

	case float64:
		return fmt.Sprintf("%g", v), nil

	case bool:
		return fmt.Sprintf("%t", v), nil

	case nil:
		return "null", nil

	default:
		return "undefined", nil
	}
}
