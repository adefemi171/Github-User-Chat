package bson

import (
	"github.com/stretchr/codecs/constants"
	"gopkg.in/mgo.v2/bson"
)

// BsonCodec converts objects to and from BSON.
type BsonCodec struct{}

// Marshal converts an object to BSON.
func (b *BsonCodec) Marshal(object interface{}, options map[string]interface{}) ([]byte, error) {
	return bson.Marshal(object)
}

// Unmarshal converts JSON into an object.
func (b *BsonCodec) Unmarshal(data []byte, obj interface{}) error {
	return bson.Unmarshal(data, obj)
}

// ContentType returns the content type for this codec.
func (b *BsonCodec) ContentType() string {
	return constants.ContentTypeBSON
}

// FileExtension returns the file extension for this codec.
func (b *BsonCodec) FileExtension() string {
	return constants.FileExtensionBSON
}

// CanMarshalWithCallback returns whether this codec is capable of marshalling a response containing a callback.
func (b *BsonCodec) CanMarshalWithCallback() bool {
	return false
}
