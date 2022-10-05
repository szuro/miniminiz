package server

type ServerConfig struct {
	IP        string `yaml:"listen_ip"`
	Port      string `yaml:"listen_port"`
	cacheSize int    `yaml:"cache_size"`
}

type Config struct {
	Server ServerConfig
	Hosts  CheckList
}
