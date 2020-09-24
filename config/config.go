// YOU CAN EDIT YOUR CUSTOM CONFIG HERE

package config

// Config holds all settings
var defaultConfig = []byte(`
environment: D
grpc_address: 10000
http_address: 9000


`)

type Config struct {
	Base `mapstructure:",squash"`
	// Custom here

}
