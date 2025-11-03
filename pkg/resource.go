package pkg

import "image/color"

type Material struct {
	name        string
	color       color.Color
	translateID string
}

var (
	// 火
	PYRO = Material{name: "PYRO", color: color.RGBA{R: 255, G: 85, B: 0, A: 255}, translateID: "Material/PYRO"}
	// 水
	HYDRO = Material{name: "HYDRO", color: color.RGBA{R: 0, G: 170, B: 255, A: 255}, translateID: "Material/HYDRO"}
	// 空
	AERO = Material{name: "AERO", color: color.RGBA{R: 170, G: 255, B: 255, A: 255}, translateID: "Material/AERO"}
	// 自然
	NATURO = Material{name: "NATURO", color: color.RGBA{R: 0, G: 170, B: 0, A: 255}, translateID: "Material/NATURO"}
	// 地
	GEO = Material{name: "GEO", color: color.RGBA{R: 102, G: 51, B: 0, A: 255}, translateID: "Material/GEO"}
	// エーテル(光)
	AETHER = Material{name: "AETHER", color: color.RGBA{R: 255, G: 255, B: 170, A: 255}, translateID: "Material/AETHER"}
	// 血(闇)
	BLOOD = Material{name: "BLOOD", color: color.RGBA{R: 170, G: 0, B: 0, A: 255}, translateID: "Material/BLOOD"}
	// 過去
	PALEO = Material{name: "PALEO", color: color.RGBA{R: 128, G: 128, B: 128, A: 255}, translateID: "Material/PALEO"}
	// 未来
	MELO = Material{name: "MELO", color: color.RGBA{R: 0, G: 255, B: 255, A: 255}, translateID: "Material/MELO"}
	// 忠誠
	ALLEGRO = Material{name: "ALLEGRO", color: color.RGBA{R: 255, G: 170, B: 0, A: 255}, translateID: "Material/ALLEGRO"}
	// 信仰
	FAITH = Material{name: "FAITH", color: color.RGBA{R: 255, G: 255, B: 255, A: 255}, translateID: "Material/FAITH"}
)

func (m *Material) Name() string        { return m.name }
func (m *Material) Color() color.Color  { return m.color }
func (m *Material) TranslateID() string { return m.translateID }

type Materials map[Material]int

func (m Materials) Copy() Materials {
	dst := make(map[Material]int, len(m))
	for k, v := range m {
		dst[k] = v
	}
	return dst
}
