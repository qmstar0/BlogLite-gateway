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

type SSL struct {
	SSLCertFilePath string `toml:"ssl_cert_fp" mapstructure:"ssl_cert_fp"`
	SSLKeyFilePath  string `toml:"ssl_key_fp" mapstructure:"ssl_key_fp"`
}

type Config struct {
	Debug  bool    `toml:"debug" mapstructure:"debug"`
	SSL    *SSL    `toml:"ssl" mapstructure:"ssl"`
	Assets *Assets `toml:"assets" mapstructure:"assets"`
	Proxys []Proxy `toml:"proxy" mapstructure:"proxy"`
}
