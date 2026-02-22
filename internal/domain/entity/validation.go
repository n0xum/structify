package entity

import "strings"

func ValidateFieldName(name string) error {
	if name == "" {
		return ErrInvalidFieldName
	}
	if strings.ContainsAny(name, " \t\n\r") {
		return ErrInvalidFieldName
	}
	return nil
}

func ValidateFieldType(goType string) error {
	baseType := goType
	if idx := strings.Index(baseType, "["); idx != -1 {
		baseType = baseType[:idx]
	}
	if idx := strings.Index(baseType, "."); idx != -1 {
		baseType = baseType[:idx]
	}
	baseType = strings.TrimPrefix(baseType, "*")
	baseType = strings.TrimPrefix(baseType, "[]")

	if baseType == "" {
		return ErrInvalidFieldName
	}
	return nil
}
