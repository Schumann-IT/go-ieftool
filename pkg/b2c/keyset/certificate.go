package keyset

type Certificate struct {
	Body     string
	Password string
}

func NewCertificate(cert, pw string) *Certificate {
	return &Certificate{
		Body:     cert,
		Password: pw,
	}
}
