package dsn

type Options struct {
	SslMode    bool
	ParseTime  bool
	SearchPath string
}

type OptionSetter func(o *Options)

func ParseTime(enabled bool) OptionSetter {
	return func(o *Options) {
		o.ParseTime = enabled
	}
}

func SearchPath(path string) OptionSetter {
	return func(o *Options) {
		o.SearchPath = path
	}
}

func evaluateOptions(args []OptionSetter) Options {
	// Init default value for options
	opts := Options{
		SslMode:    false,
		ParseTime:  true,
		SearchPath: "",
	}

	for _, arg := range args {
		arg(&opts)
	}

	return opts
}
