package internal

type JobHandler interface {
	Send(*Message)
}
