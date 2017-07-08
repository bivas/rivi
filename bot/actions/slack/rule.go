package slack

type rule struct {
	Message string `mapstructure:"message-template"`
	Channel string `mapstructure:"channel"`
}

func (r *rule) Defaults() {
}
