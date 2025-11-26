package translator

import (
	"frontage/pkg/engine/model"
	"frontage/pkg/network/data"
)

type MaterialsTranslator struct {
}

func (mt *MaterialsTranslator) ToModel(d data.Materials) (model.Materials, error) {
	result := make(model.Materials)
	for _, m := range model.EnumerateMaterial() {
		v, f := d[string(m.Tag())]
		if !f {
			continue
		}
		result[m] = v
	}
	return result, nil
}

func (mt *MaterialsTranslator) FromModel(m model.Materials) (data.Materials, error) {
	result := make(data.Materials)
	for k, v := range m {
		result[string(k.Tag())] = v
	}
	return result, nil
}
