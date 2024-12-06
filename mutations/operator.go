package mutations

// MutationRule defines a struct for mutation logic
type MutationRule struct {
	Original string
	Mutant   string
}

// Predefined mutation rules
var MutationRules = []MutationRule{
	{"+", "-"},
	{"-", "+"},
	{">", "<"},
	{"<", ">"},
	{"*", "/"},
	{"/", "*"},
}
