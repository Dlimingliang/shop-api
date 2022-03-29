package config

type ServerConfig struct {
	Name              string        `mapstructure:"name"`
	UserServiceConfig UserSrcConfig `mapstructure:"user-srv"`
	JWTConfig         JWTAuthConfig `mapstructure:"jwt"`
}

type UserSrcConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTAuthConfig struct {
	SignKey string `mapstructure:"sign-key"`
}
