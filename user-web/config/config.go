package config

type ServerConfig struct {
	Name              string            `mapstructure:"name"`
	UserServiceConfig UserSrcConfig     `mapstructure:"user-srv"`
	JWTConfig         JWTAuthConfig     `mapstructure:"jwt"`
	RedisConfig       RedisServerConfig `mapstructure:"redis"`
}

type UserSrcConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTAuthConfig struct {
	SignKey string `mapstructure:"sign-key"`
}

type RedisServerConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}
