package section

type (
	Processor struct {
		WebServer ProcessorWebServer
	}

	ProcessorWebServer struct {
		ListenPort int `required:"true" default:"8082"`
	}
)
