package dspy

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/supadev-ai/go-dspy/llm"
)

// Predictor is an LLM-backed module that uses a Signature to transform inputs to outputs.
type Predictor[I any, O any] struct {
	Signature Signature[I, O]
	Client    llm.Client
}

// NewPredictor creates a new Predictor with the given signature and LLM client.
func NewPredictor[I any, O any](sig Signature[I, O], client llm.Client) *Predictor[I, O] {
	return &Predictor[I, O]{
		Signature: sig,
		Client:    client,
	}
}

// Forward implements the Module interface.
func (p *Predictor[I, O]) Forward(ctx context.Context, input I) (O, error) {
	var output O

	prompt := p.buildPrompt(input)
	
	response, err := p.Client.Generate(ctx, prompt)
	if err != nil {
		return output, ErrModuleExecution("predictor.Forward", err)
	}

	// Parse the response into the output type
	parsed, err := p.parseResponse(response)
	if err != nil {
		return output, ErrModuleExecution("predictor.parseResponse", err)
	}

	return parsed, nil
}

// buildPrompt constructs a prompt from the signature and input.
func (p *Predictor[I, O]) buildPrompt(input I) string {
	var parts []string

	if p.Signature.Description != "" {
		parts = append(parts, p.Signature.Description)
	}

	// Extract input fields
	inputFields := p.extractFields(input)
	if len(inputFields) > 0 {
		parts = append(parts, "\nInput:")
		for key, value := range inputFields {
			parts = append(parts, fmt.Sprintf("%s: %v", key, value))
		}
	}

	// Add output instruction
	outputFields := p.getOutputFieldNames()
	if len(outputFields) > 0 {
		parts = append(parts, "\nOutput the following fields:")
		for _, field := range outputFields {
			parts = append(parts, fmt.Sprintf("- %s", field))
		}
	}

	return strings.Join(parts, "\n")
}

// extractFields extracts field names and values from a struct.
func (p *Predictor[I, O]) extractFields(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return result
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldVal := val.Field(i)
		result[field.Name] = fieldVal.Interface()
	}

	return result
}

// getOutputFieldNames returns the names of output struct fields.
func (p *Predictor[I, O]) getOutputFieldNames() []string {
	var output O
	typ := reflect.TypeOf(output)
	
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	var names []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.IsExported() {
			names = append(names, field.Name)
		}
	}

	return names
}

// parseResponse attempts to parse the LLM response into the output type.
// It tries JSON parsing first, then falls back to simple field extraction.
func (p *Predictor[I, O]) parseResponse(response string) (O, error) {
	var output O
	
	// Try JSON parsing first (if response looks like JSON)
	responseTrimmed := strings.TrimSpace(response)
	if strings.HasPrefix(responseTrimmed, "{") || strings.HasPrefix(responseTrimmed, "[") {
		if err := json.Unmarshal([]byte(responseTrimmed), &output); err == nil {
			return output, nil
		}
		// If JSON parsing fails, fall through to field extraction
	}

	// Fall back to field extraction from text
	val := reflect.ValueOf(&output).Elem()
	typ := reflect.TypeOf(output)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if val.IsNil() {
			val.Set(reflect.New(typ))
		}
		val = val.Elem()
	}

	if typ.Kind() == reflect.Struct {
		// Look for field names in the response (case-insensitive)
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if !field.IsExported() {
				continue
			}

			// Try multiple patterns: "FieldName:", "fieldname:", etc.
			fieldName := field.Name
			patterns := []string{
				fieldName + ":",
				strings.ToLower(fieldName) + ":",
				strings.Title(strings.ToLower(fieldName)) + ":",
			}

			for _, pattern := range patterns {
				idx := strings.Index(strings.ToLower(response), strings.ToLower(pattern))
				if idx != -1 {
					valueStart := idx + len(pattern)
					// Find the end of the value (newline or end of string)
					valueEnd := strings.Index(response[valueStart:], "\n")
					if valueEnd == -1 {
						valueEnd = len(response) - valueStart
					}
					valueStr := strings.TrimSpace(response[valueStart : valueStart+valueEnd])
					
					// Set the field value based on type
					fieldVal := val.Field(i)
					if fieldVal.CanSet() {
						switch field.Type.Kind() {
						case reflect.String:
							fieldVal.SetString(valueStr)
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							// Try to parse as int (simplified)
							// In production, use strconv
						}
					}
					break
				}
			}
		}
	} else if typ.Kind() == reflect.String {
		val.SetString(response)
	}

	return output, nil
}
