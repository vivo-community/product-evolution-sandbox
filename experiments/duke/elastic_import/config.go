package widgets_import

type Config struct {
	Database database
	Elastic  elasticSearch `toml:"elastic"`
}

type elasticSearch struct {
	Url string
}

type database struct {
	Server   string
	Port     int
	Database string
	User     string
	Password string
}


