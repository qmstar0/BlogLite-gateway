package config

type Proxy struct {
	Source string `toml:"source" mapstructure:"source"`
	Target string `toml:"target" mapstructure:"target"`
}

type Assets struct {
	Source      string `toml:"source" mapstructure:"source"`
	Dir         string `toml:"dir" mapstructure:"dir"`
	StripPrefix string `toml:"strip_prefix" mapstructure:"strip_prefix"`
}

type Config struct {
	Assets Assets `toml:"assets" mapstructure:"assets"`

	Port   []int   `toml:"port" mapstructure:"port"`
	Proxys []Proxy `toml:"proxy" mapstructure:"proxy"`
}
