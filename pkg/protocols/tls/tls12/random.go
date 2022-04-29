package tls12

import (
	"fmt"
	"time"

	"github.com/elmasy-com/bytebuilder"
)

/*
	struct {
	    uint32 gmt_unix_time;
	    opaque random_bytes[28];
	} Random;
*/

type random struct {
	Time        time.Time
	RandomBytes []byte
}

func marshalRandom() []byte {

	buf := bytebuilder.NewEmpty()

	buf.WriteGMTUnixTime32(time.Now())

	buf.WriteRandom(28)

	return buf.Bytes()
}

func unmarshalRandom(bytes []byte) (random, error) {

	if bytes == nil {
		return random{}, fmt.Errorf("bytes is nil")
	}

	var (
		random random
		ok     bool
		buf    = bytebuilder.NewBuffer(bytes)
	)

	if random.Time, ok = buf.ReadGMTUnixTime32(); !ok {
		return random, fmt.Errorf("failed to read Random")
	}

	if random.RandomBytes = buf.ReadBytes(28); random.RandomBytes == nil {
		return random, fmt.Errorf("failed to read RandomBytes")
	}

	return random, nil
}
