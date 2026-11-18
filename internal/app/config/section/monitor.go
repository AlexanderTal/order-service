package section

type Monitor struct {
	LogLevel    string `required:"true" envconfig:"LOG_LEVEL"`
	Environment string `required:"true"`
}
