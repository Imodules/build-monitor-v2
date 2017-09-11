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
	TcProjectPollInterval      string      `env:"tcProjectPollInterval" flag:"tcProjectPollInterval" flagDesc:"How often to poll TeamCity to refresh project list"`
	TcBuildPollInterval        string      `env:"tcBuildPollInterval" flag:"tcBuildPollInterval" flagDesc:"How often to poll the TeamCity for new running builds"`
	TcRunningBuildPollInterval string      `env:"tcRunningBuildPollInterval" flag:"tcRunningBuildPollInterval" flagDesc:"How often to poll TeamCity for running build status"`
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
		TcUrl:         "http://pstuart.no-ip.org:8111",
		TcProjectPollInterval:      "20m",
		TcBuildPollInterval:        "20s",
		TcRunningBuildPollInterval: "5s",
	}
}
