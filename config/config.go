package config

type Mode int

const (
	S1 Mode = iota + 1
	S15
	S2
	S25
	S3
	S35
	S4
	S5
	S55
	S6
	S7
	S8
	S85
	S9
	S10
	F1
	F2
	F3
	F4
	F5
	G1
	G2
	G3
	G31
	G32
	GG
	G5
	G51
	G52
	G6
	G61
	G7
	G8
	G81
	G9
	G91
	G10
	G101
	Z1
	Z2
	ZZ
)

type Config struct {
	RealClientMode Mode
	ServerPort     int
	ServerHost     string
}

var ProxyConfig *Config = &Config{}
