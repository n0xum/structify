package util

import "strings"

// ToSnakeCase converts a PascalCase or camelCase string to snake_case,
// correctly grouping consecutive uppercase letters as a single word:
//   - "UserName"   → "user_name"
//   - "ID"         → "id"
//   - "UserID"     → "user_id"
//   - "AvatarURL"  → "avatar_url"
//   - "XMLParser"  → "xml_parser"
//   - "HTTPServer" → "http_server"
//   - "SKU"        → "sku"
func ToSnakeCase(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	var result []rune
	for i, r := range runes {
		if r >= 'A' && r <= 'Z' && i > 0 {
			prev := runes[i-1]
			// Insert _ if previous is lowercase (e.g. "userName" → "user_Name")
			// OR if previous is uppercase and next is lowercase
			// (e.g. "HTMLContent" → "html_content": boundary before "C", not before "T","M","L")
			if prev >= 'a' && prev <= 'z' {
				result = append(result, '_')
			} else if i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z' {
				result = append(result, '_')
			}
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
