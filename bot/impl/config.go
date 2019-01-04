package impl

type ConfigDatabase struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	Token    string         `yaml:"token"`
	Database ConfigDatabase `yaml:"database"`
}
