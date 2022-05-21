package config

const (
	DefaultSearchConfigName = "config"
	DefaultSearchConfigType = "yml" // yaml
)

var (
	DefaultSearchPath = []string{".", "./config", "./configs"}
)

type Option func(opts *Options)

type Options struct {
	AbsPath []string // abs path -> high level priority

	// default: config
	ConfigName string
	// default: yml,yaml
	ConfigType string
	// default: ., ./config, ./configs;
	// priority: ./ -> ./config -> ./configs
	SearchPath []string
	// the profiles active -> []string{"dev","share"}
	// target: config.yml/config_dev.yml/config_share.yml
	// default: []string
	ProfilesActive []string
	// multi config merge depth, if necessary
	MergeDepth uint8
}

func (opts *Options) validate() {
	if isBlankString(opts.ConfigName) {
		opts.ConfigName = DefaultSearchConfigName
	}
	if isBlankString(opts.ConfigType) {
		opts.ConfigType = DefaultSearchConfigType
	}
	if isEmptyStringSlice(opts.SearchPath) {
		opts.SearchPath = DefaultSearchPath
	}
	if opts.MergeDepth == 0 {
		opts.MergeDepth = defaultMergeDepth
	}
}

func newOptions() *Options {
	return &Options{
		SearchPath:     make([]string, 0),
		ProfilesActive: make([]string, 0),
	}
}

func WithAbsPath(absPath ...string) Option {
	return func(opts *Options) {
		opts.AbsPath = absPath
	}
}

func WithConfigName(configName string) Option {
	return func(opts *Options) {
		opts.ConfigName = configName
	}
}

func WithConfigType(configType string) Option {
	return func(opts *Options) {
		opts.ConfigType = configType
	}
}

func WithSearchPath(searchPath ...string) Option {
	return func(opts *Options) {
		opts.SearchPath = searchPath
	}
}

func WithProfilesActive(profilesActive ...string) Option {
	return func(opts *Options) {
		opts.ProfilesActive = profilesActive
	}
}

func WithMergeDepth(mergeDepth uint8) Option {
	return func(opts *Options) {
		opts.MergeDepth = mergeDepth
	}
}
