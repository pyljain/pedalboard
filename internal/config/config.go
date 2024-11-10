package config

type Config struct {
	Authentication Authentication `yaml:"authentication"`
	Routes         []Route        `yaml:"routes"`
	Port           int            `yaml:"port"`
}

type Authentication struct {
	APIKey string `yaml:"api_key"`
}

type Route struct {
	Path      string     `yaml:"path"`
	Responses []Response `yaml:"responses"`
	Latency   Latency    `yaml:"latency"`
}

type Response struct {
	Status      int    `yaml:"status"`
	Body        string `yaml:"body"`
	Probability int    `yaml:"probability"`
}

type Latency struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}
