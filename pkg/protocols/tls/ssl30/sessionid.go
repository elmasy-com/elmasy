package ssl30

import (
	"github.com/elmasy-com/bytebuilder"
)

/*
	opaque SessionID<0..32>;
*/

func marshalSessionID() []byte {

	buf := bytebuilder.NewEmpty()

	buf.WriteUint8(32)
	buf.WriteRandom(32)

	return buf.Bytes()
}
