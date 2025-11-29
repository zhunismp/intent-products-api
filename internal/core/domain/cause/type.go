package cause

// CauseMap represents a mapping of cause IDs to their status
type CauseInfoMap map[string]*CauseInfo

type CauseInfo struct {
	Reason string
	Status bool
}