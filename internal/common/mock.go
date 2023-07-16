package common

type Mock struct {
	Dir     string
	Replies map[string]string
}

func (m Mock) CacheDir() string {
	return m.Dir + "/cache/"
}

func (m Mock) WorkDir() string {
	return m.Dir + "/work/"
}
