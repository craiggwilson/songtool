package theory

var std = New(DefaultConfig())

func New(cfg *Config) *Theory {
	return &Theory{Config: cfg}
}

func NewDefault() *Theory {
	return New(DefaultConfig())
}

type Theory struct {
	Config *Config
}
