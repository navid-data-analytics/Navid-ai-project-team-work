package zk

// Option implements the available options for the client
type Option func(*Client)

// OptionPathPrefix sets the path prefix used for non absolute reads.
// E.g. path prefix of "/env" would yield kafka brokers path of "/env/brokers/ids"
func OptionPathPrefix(path string) Option {
	return func(c *Client) {
		c.pathPrefix = path
	}
}

// OptionKafkaBrokersPath sets the kafka brokers read path for the zookeeper client.
// The path has to be absolute, any configured path prefix is ignored.
func OptionKafkaBrokersPath(absPath string) Option {
	return func(c *Client) {
		c.kafkaBrokersPath = absPath
	}
}
