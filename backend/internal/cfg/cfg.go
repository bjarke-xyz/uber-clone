package cfg

import "os"

type Cfg struct {
	Env                       string
	Port                      string
	OSRApiKey                 string
	DatabaseConnectionPoolUrl string
	FirebaseProjectId         string
}

func NewConfig() *Cfg {
	cfg := &Cfg{
		Env:                       os.Getenv("ENV"),
		Port:                      os.Getenv("PORT"),
		OSRApiKey:                 os.Getenv("OSR_API_KEY"),
		DatabaseConnectionPoolUrl: os.Getenv("DATABASE_CONNECTION_POOL_URL"),
		FirebaseProjectId:         os.Getenv("FIREBASE_PROJECT_ID"),
	}
	return cfg
}
