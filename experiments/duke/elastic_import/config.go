package widgets_import

type Config struct {
	Database  database      `toml:"database"`
	Elastic   elasticSearch `toml:"elastic"`
	Templates templatePaths `toml:"template"`
	Schemas   schemasPath   `toml:"schema"`
}

type elasticSearch struct {
	Url string `toml:"url"`
}

type database struct {
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type templatePaths struct {
	Layout  string `toml:"layout"`
	Include string `toml:"include"`
}

type schemasPath struct {
	Dir string `toml:"dir"`
}


