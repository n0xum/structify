package entity

type MethodKind int

const (
	MethodCreate MethodKind = iota
	MethodGetByID
	MethodUpdate
	MethodDelete
	MethodList
	MethodFindBy
	MethodSmartQuery // Auto-generated SQL from patterns
	MethodCustomSQL
)

type RepositoryInterface struct {
	Name       string
	EntityName string
	Methods    []RepositoryMethod
	Package    string
}

type RepositoryMethod struct {
	Name             string
	Kind             MethodKind
	Params           []MethodParam
	ReturnsSingle    bool
	ReturnsError     bool
	HasEntityReturn  bool   // false when method returns only error (no entity)
	ScalarReturnType string // set when method returns a scalar (e.g. float64, int64)
	EntityName       string
	FindByFields     []string
	CustomSQL        string
	GeneratedSQL     string // Auto-generated SQL from pattern
	QueryPattern     string // Matched pattern identifier
}

type MethodParam struct {
	Name string
	Type string
}
