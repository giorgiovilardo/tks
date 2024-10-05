package internal

type Config struct {
	Leagues []League `koanf:"leagues"`
}

type League struct {
	Name string `koanf:"name"`
	URL  string `koanf:"url"`
}
