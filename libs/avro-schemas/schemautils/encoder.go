package schemautils

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
)

type AvroEncoder struct {
	schemaID int
	codec    *goavro.Codec
}

func NewAvroEncoder(schemaRegistryURL, subject, schemaStr string) (*AvroEncoder, error) {
	client := srclient.CreateSchemaRegistryClient(schemaRegistryURL)

	schema, err := client.GetLatestSchema(subject)
	if err != nil {
		schema, err = client.CreateSchema(subject, schemaStr, srclient.Avro)
		if err != nil {
			return nil, err
		}
	}

	codec, err := goavro.NewCodec(schema.Schema())
	if err != nil {
		return nil, err
	}

	return &AvroEncoder{
		schemaID: schema.ID(),
		codec:    codec,
	}, nil
}

func (a *AvroEncoder) Encode(data map[string]interface{}) ([]byte, error) {
	avroBinary, err := a.codec.BinaryFromNative(nil, data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteByte(0)
	binary.Write(&buf, binary.BigEndian, int32(a.schemaID))
	buf.Write(avroBinary)

	return buf.Bytes(), nil
}

func (a *AvroEncoder) Decode(data []byte) (map[string]interface{}, error) {
    if len(data) < 5 {
        return nil, fmt.Errorf("invalid data length: %d", len(data))
    }

    payload := data[5:]

    decoded, _, err := a.codec.NativeFromBinary(payload)
    if err != nil {
        return nil, err
    }

    return decoded.(map[string]interface{}), nil
}
