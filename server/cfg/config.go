package cfg

type Config struct {
	gofigure      interface{} `envPrefix:"IFER" order:"flag,env"`
	Port          int         `env:"port" flag:"port" flagDesc:"Port to run the api server on"`
	Db            string      `env:"db" flag:"db" flagDesc:"Url to mongodb"`
	PasswordSalt  string      `env:"passwordSalt" flag:"passwordSalt" flagDesc:"Salt to use for the password"`
	ClientPath    string      `env:"clientPath" flag:"clientPath" flagDesc:"Path to where the client code is stored"`
	AllowedOrigin string      `env:"allowedOrigin" flag:"allowedOrigin" flagDesc:"The CORS allowed origin"`
	JwtSecret     string      `env:"jwtSecret" flag:"jwtSecret" flagDesc:"The secret key for the JWT token"`
}

func Load(getOverrides func(s interface{}) error) (Config, error) {
	config := getDefaults()

	err := getOverrides(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func getDefaults() Config {
	return Config{
		Port:          3030,
		Db:            "mongodb://localhost:27017/build-monitor-v2",
		PasswordSalt:  "you-really-need-to-change-this",
		ClientPath:    "../client/dist",
		AllowedOrigin: "*",
		JwtSecret:     "you-really-need-to-change-this-one-also",
	}
}
