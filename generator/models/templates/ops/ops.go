package ops

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Debug    string `yaml:"debug"`
}

type Ops struct {
	Workdir      string   `yaml:"workdir"`
	AppHTTPAlias string   `yaml:"app_http_alias"`
	MainDatabase Database `yaml:"database"`
}
