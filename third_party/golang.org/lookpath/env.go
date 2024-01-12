package lookpath

type Env interface {
	Getenv(key string) string
	LookupEnv(key string) (string, bool)
}
