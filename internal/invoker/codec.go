package invoker

// rawCodec is a grpc.Codec implementation that passes []byte payloads
// through without any marshalling. This allows grpcannon to send
// pre-encoded protobuf bytes directly.
type rawCodec struct{}

func (rawCodec) Marshal(v interface{}) ([]byte, error) {
	if b, ok := v.([]byte); ok {
		return b, nil
	}
	// Nil payload is valid for methods that take google.protobuf.Empty.
	return nil, nil
}

func (rawCodec) Unmarshal(data []byte, v interface{}) error {
	if p, ok := v.(*[]byte); ok {
		*p = data
	}
	return nil
}

func (rawCodec) Name() string { return "proto" }
