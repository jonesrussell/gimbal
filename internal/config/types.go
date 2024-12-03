package config

// Config holds all game configuration settings
type Config struct {
	Screen struct {
		Title  string `json:"title"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"screen"`
	Game struct {
		Speed    float64 `json:"speed"`
		NumStars int     `json:"numStars"`
		Debug    bool    `json:"debug"`
	} `json:"game"`
	Player struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"player"`
}
