package config

type ServerConfig struct {
	Name              string        `mapstructure:"name"`
	UserServiceConfig UserSrcConfig `mapstructure:"user-srv"`
}

type UserSrcConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
