package common

// Status is a custom enum type for states.
type Status int

// Declared enums for various states.
const (
	Pending Status = 0 << iota
	Active
	Archived
)

// String prints out a human readable version of the status.
func (s Status) String() string {
	switch {
	case s == Pending:
		return "Pending"
	case s == Active:
		return "Active"
	case s == Archived:
		return "Archived"
	}
	return "Unknown"
}
