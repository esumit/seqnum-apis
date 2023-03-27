package config

type SeqNumServerConfig struct {
	Port         string
	IPAddress    string
	WriteTimeout int
	ReadTimeout  int
	IdleTimeout  int
}

type SeqNumCollectionConfig struct {
	CollectionTime int
}
