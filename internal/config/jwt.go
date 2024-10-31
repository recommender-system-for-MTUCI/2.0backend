package config

type JWT struct {
	AccessTime  int
	RefreshTime int
	PublicKey   string
	PrivateKey  string
}
