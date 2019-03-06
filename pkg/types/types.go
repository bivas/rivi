package types

import "os"

var (
	RulesConfigFileName = ".rivi.yaml"
)

func init() {
	envOverride := os.Getenv("RIVI_CONFIG_FILE")
	if envOverride != "" {
		RulesConfigFileName = envOverride
	}
}
