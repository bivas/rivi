package slack

type rule struct {
	Message string `mapstructure:"message-template"`
	Channel string `mapstructure:"channel"`
	Notify  string `mapstructure:"notify"`
}

func (r *rule) Defaults() {
}
