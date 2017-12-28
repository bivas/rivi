package types

import "strings"

type State int

const (
	Failure State = iota
	Pending
	Success
	Error
)

func GetState(value string) State {
	switch strings.ToLower(value) {
	case "failure":
		return Failure
	case "pending":
		return Pending
	case "success":
		return Success
	}
	return Error
}
