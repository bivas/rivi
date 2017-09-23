package types

type JobHandler interface {
	Send(*Message)
}
