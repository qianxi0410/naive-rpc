package evangelion

import "github.com/qianxi0410/naive-rpc/codec"

const NAME = "EVANGELION"

func init() {
	codec.RegisterCodec(NAME, &ServerCodec{}, &ClientCodec{})
	codec.RegisterSessionBuilder(NAME, &DefaultSessionBuilder{})
}
