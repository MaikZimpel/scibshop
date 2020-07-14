package configuration

type Config struct {
	Security struct {
		ClientSecret string `yaml:"clientSecret" env:"SECURITY_CLIENT_SECRET" env-default:"secret"`
	} `yaml:"security"`
	Server struct {
		Port string `yaml:"port" env:"SERVER_PORT" env-default:"8082"`
	} `yaml:"server"`
}
