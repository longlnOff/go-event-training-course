package message
import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)
var (
	marshalerJSON = cqrs.JSONMarshaler{
		GenerateName: cqrs.StructName,
	}
)
