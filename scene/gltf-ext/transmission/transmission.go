package transmission

import (
	"encoding/json"

	"github.com/qmuntal/gltf"
)

const ExtensionName = "KHR_materials_transmission"

func Unmarshal(data []byte) (any, error) {
	matTransmission := new(MaterialsTransmission)
	err := json.Unmarshal(data, matTransmission)
	return matTransmission, err
}

type MaterialsTransmission struct {
	TransmissionFactor  *float32          `json:"transmissionFactor,omitempty"`
	TransmissionTexture *gltf.TextureInfo `json:"transmissionTexture,omitempty"`
}
