package example

type Configuration struct {
	Values []ValueConfiguration `koanf:"values"`
}

type ValueConfiguration struct {
	Example string `koanf:"example"`
	Extra   int    `koanf:"extra"`
	Enable  bool   `konaf:"enable"`
}
