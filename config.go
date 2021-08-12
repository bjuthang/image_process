package main

type Config struct {
	Host        string  `json:"HOST"`
	Port        string  `json:"PORT"`
	ImgPath     string  `json:"IMG_PATH"`
	ExhaustPath string  `json:"EXHAUST_PATH"`
	Quality     float32 `json:"QUALITY"`
}

var (
	configPath  string
	config      Config
	jobs        int
	verboseMode bool
)

const (
	sampleConfig = `
{
  "HOST": "127.0.0.1",
  "PORT": "3333",
  "QUALITY": 80,
  "IMG_PATH": "./pics",
  "EXHAUST_PATH": "./exhaust",
}`
)
