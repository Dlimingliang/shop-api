package config

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Namespace string `mapstructure:"namespace"`
	DataId    string `mapstructure:"dataId"`
	Group     string `mapstructure:"group"`
}

type ServerConfig struct {
	Name         string             `yaml:"name"`
	JWTConfig    JWTAuthConfig      `yaml:"jwt"`
	RedisConfig  RedisServerConfig  `yaml:"redis"`
	ConsulConfig ConsulServerConfig `yaml:"consul"`
}

type JWTAuthConfig struct {
	SignKey string `yaml:"sign-key"`
}

type RedisServerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type ConsulServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
