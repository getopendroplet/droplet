package config

// Remote holds details for communication with a remote registry
type Remote struct {
	Addr     string `yaml:"addr"`
	Protocol string `yaml:"protocol,omitempty"`
}
