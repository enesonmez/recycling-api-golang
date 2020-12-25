package config

type Config struct {
	DBHost        string `json:"dbHost"`
	DBPort        int    `json:"dbPort"`
	DBUsername    string `json:"dbUsername"`
	DBPassword    string `json:"dbPassword"`
	DBName        string `json:"dbName"`
	Email         string `json:"email"`
	EmailPassword string `json:"emailPassword"`
	BaseURL       string `json:"baseURL"`
}
