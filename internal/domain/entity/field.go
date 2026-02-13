package entity

type Field struct {
	Name       string
	Type       string
	IsPrimary  bool
	IsUnique   bool
	IsIgnored  bool
	TableName  string
}

func (f *Field) ShouldGenerate() bool {
	if f.IsIgnored {
		return false
	}
	if f.Type == "" {
		return false
	}
	return true
}
