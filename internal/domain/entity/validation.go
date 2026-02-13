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

func ToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := s[i-1]
			next := byte(' ')
			if i+1 < len(s) {
				next = s[i+1]
			}
			if prev >= 'a' && prev <= 'z' || next >= 'a' && next <= 'z' {
				result = append(result, '_')
			}
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
