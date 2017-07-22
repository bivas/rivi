package autoassign

type rule struct {
	FromRoles    []string `mapstructure:"roles"`
	Require      int      `mapstructure:"require"`
	Experimental struct {
		FromSuggestions bool `mapstructure:"suggest"`
	} `mapstructure:"experimental"`
}

func (r *rule) Defaults() {
	if r.Require == 0 {
		r.Require = 1
	}
}
