package trigger

var (
	allowedMethods = []string{"GET", "POST"}
)

type rule struct {
	Endpoint string            `mapstructure:"endpoint"`
	Method   string            `mapstructure:"method"`
	Headers  map[string]string `mapstructure:"headers"`
	Body     string            `mapstructure:"body"`
}

func (r *rule) Defaults() {
	allowedMethod := false
	for _, allowed := range allowedMethods {
		if allowed == r.Method {
			allowedMethod = true
			break
		}
	}
	if !allowedMethod {
		r.Method = "POST"
	}
}
