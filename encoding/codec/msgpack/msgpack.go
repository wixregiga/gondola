// Package msgpack provides a codec implementation using msgpack.
//
// To enable it in your app, import it like:
//
//  import (
//	_ "gondola/encoding/codec/msgpack"
//  )
package msgpack

import (
	"gondola/encoding/codec"

	gocodec "github.com/ugorji/go/codec"
)

var (
	msgpackCodec = &codec.Codec{Encode: msgpackMarshal, Decode: msgpackUnmarshal, Binary: true}
	handle       = &gocodec.MsgpackHandle{}
)

func msgpackMarshal(in interface{}) ([]byte, error) {
	var b []byte
	enc := gocodec.NewEncoderBytes(&b, handle)
	err := enc.Encode(in)
	return b, err
}

func msgpackUnmarshal(data []byte, out interface{}) error {
	dec := gocodec.NewDecoderBytes(data, handle)
	return dec.Decode(out)
}

func init() {
	codec.Register("msgpack", msgpackCodec)
}
