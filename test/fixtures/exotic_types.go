package fixtures

// ExoticTypes tests parsing of map, chan, interface{}, ellipsis, and array types.
type ExoticTypes struct {
	MapField       map[string]int
	ChanField      chan bool
	InterfaceField interface{}
	SliceField     []string
	PointerField   *int
	NestedMap      map[string]map[int]bool
}

// unexportedStruct should be skipped by the parser.
type unexportedStruct struct {
	ID int64
}

// EmptyStruct has no exported fields, should be skipped.
type EmptyStruct struct {
	hidden int
}

// InterfaceWithMethods is used to test interface parsing.
type ExoticRepository interface {
	DoSomething(data map[string]int) (interface{}, error)
	Variadic(items ...string) error
	WithMapReturn() (map[string]int, error)
}
