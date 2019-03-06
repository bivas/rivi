package types

type Repository interface {
	GetCollaborators() []string
	IsCollaborator(name string) bool
	GetAvailableLabels() []string
	GetRulesFile() string
}
