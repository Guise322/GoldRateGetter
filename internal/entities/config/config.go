package config

type Config struct {
	SendingHours []int
	PriceType    string
	Items        map[string]float64
	Email        Email
}

type Email struct {
	From     string
	Pass     string
	To       string
	SmtpHost string
	SmtpPort int
}
