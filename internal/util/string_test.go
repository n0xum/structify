package util

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"User", "user"},
		{"UserName", "user_name"},
		{"UserID", "user_id"},
		{"ID", "id"},
		{"APIKey", "api_key"},
		{"XMLParser", "xml_parser"},
		{"HTTPServer", "http_server"},
		{"AvatarURL", "avatar_url"},
		{"CustomerID", "customer_id"},
		{"SKU", "sku"},
		{"TotalAmount", "total_amount"},
		{"simpleword", "simpleword"},
		{"UserProfile", "user_profile"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToSnakeCase(tt.input)
			if got != tt.want {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
