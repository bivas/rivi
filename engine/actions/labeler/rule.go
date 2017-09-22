package labeler

type rule struct {
	Label  string `mapstructure:"label"`
	Remove string `mapstructure:"remove"`
}
