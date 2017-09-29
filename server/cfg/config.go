package cfg

type Config struct {
	gofigure                   interface{} `envPrefix:"BM" order:"flag,env"`
	Port                       int         `env:"port" flag:"port" flagDesc:"Port to run the api server on"`
	Db                         string      `env:"db" flag:"db" flagDesc:"Url to mongodb"`
	PasswordSalt               string      `env:"passwordSalt" flag:"passwordSalt" flagDesc:"Salt to use for the password"`
	ClientPath                 string      `env:"clientPath" flag:"clientPath" flagDesc:"Path to where the client code is stored"`
	AllowedOrigin              string      `env:"allowedOrigin" flag:"allowedOrigin" flagDesc:"The CORS allowed origin"`
	JwtSecret                  string      `env:"jwtSecret" flag:"jwtSecret" flagDesc:"The secret key for the JWT token"`
	TcUrl                      string      `env:"tcUrl" flag:"tcUrl" flagDesc:"The main url for the TeamCity REST API"`
	TcPollInterval             string      `env:"tcPollInterval" flag:"tcPollInterval" flagDesc:"How often to poll TeamCity for builds"`
	TcRunningBuildPollInterval string      `env:"tcRunningBuildPollInterval" flag:"tcRunningBuildPollInterval" flagDesc:"How often to poll TeamCity when we have running builds"`
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
		Port:                       3030,
		Db:                         "mongodb://localhost:27017/build-monitor-v2",
		PasswordSalt:               "you-really-need-to-change-this",
		ClientPath:                 "../client/dist",
		AllowedOrigin:              "*",
		JwtSecret:                  "you-really-need-to-change-this-one-also",
		TcUrl:                      "http://localhost:3031",
		TcPollInterval:             "20s",
		TcRunningBuildPollInterval: "5s",
	}
}
