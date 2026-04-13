package config

func (c *Config) applyServerDefaults() {
	if len(c.Server.CORSAllowedOrigins) == 0 {
		c.Server.CORSAllowedOrigins = []string{"https://codefriend.ru", "http://localhost:3000"}
	}
	if len(c.Server.WebsocketAllowedHosts) == 0 {
		c.Server.WebsocketAllowedHosts = []string{"codefriend.ru", "localhost:40040", "127.0.0.1:40040"}
	}
}
