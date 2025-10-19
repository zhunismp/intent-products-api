package entities

type Status string

const (
	STAGING Status = "staging"
	VALID   Status = "valid"
	BOUGHT  Status = "bought"
)
