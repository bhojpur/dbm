package dialect

// QuotePolicy describes quote handle policy
type QuotePolicy int

// All QuotePolicies
const (
	QuotePolicyAlways QuotePolicy = iota
	QuotePolicyNone
	QuotePolicyReserved
)
