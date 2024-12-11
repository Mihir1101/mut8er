package models

// MutationRule defines a struct for mutation logic
type MutationRule struct {
	Original string
	Mutant   string
}

type MutantDetails struct {
	OriginalLine string
	MutatedLine  string
	TestOutcome  string
	RuleApplied  MutationRule
}

type ContractMutationReport struct {
	FileName         string
	TotalMutants     int
	PassedMutants    int
	FailedMutants    int
	MutantDetails    []MutantDetails
}