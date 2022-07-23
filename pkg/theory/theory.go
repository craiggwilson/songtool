package theory

var defaultTheory = New(&defaultConfig)

func Default() *Theory {
	return defaultTheory
}

func New(cfg *Config) *Theory {
	return &Theory{Config: cfg}
}

type Theory struct {
	Config *Config
}
