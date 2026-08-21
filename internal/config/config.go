package config

import ( //import

	"os"
)

type Config struct { //structure d'une configuration de connexion à la bdd
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	Port       string
}

func LoadConfig() (*Config, error) { //charge les données d'environement dans une configuration
	//chargement des variables d'environnement

	return &Config{
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", ""),
		DBUser:     getEnv("POSTGRES_USER", ""),
		DBPassword: getEnv("POSTGRES_PASSWORD", ""),
		DBName:     getEnv("POSTGRES_DB", ""),
		DBSSLMode:  getEnv("DB_SSL_MODE", ""),
		Port:       getEnv("PORT", ""),
	}, nil
}

func getEnv(key, defaultValue string) string { //récupère les données d'environnement
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
