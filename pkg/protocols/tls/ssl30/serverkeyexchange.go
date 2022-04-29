package ssl30

// TODO

type serverKeyExchange struct {
	Body []byte
}

func unmarshalServerKeyExchange(bytes []byte) (serverKeyExchange, error) {

	return serverKeyExchange{Body: bytes}, nil
}
