package model

import (
	"frontage/pkg"
	"image/color"
)

type material struct {
	name  string
	color color.Color
	tag   pkg.ItemTag
}

type Material = *material

var (
	// 火
	PYRO = &material{name: "PYRO", color: color.RGBA{R: 255, G: 85, B: 0, A: 255}, tag: "Material/PYRO"}
	// 水
	HYDRO = &material{name: "HYDRO", color: color.RGBA{R: 0, G: 170, B: 255, A: 255}, tag: "Material/HYDRO"}
	// 空
	AERO = &material{name: "AERO", color: color.RGBA{R: 170, G: 255, B: 255, A: 255}, tag: "Material/AERO"}
	// 自然
	NATURO = &material{name: "NATURO", color: color.RGBA{R: 0, G: 170, B: 0, A: 255}, tag: "Material/NATURO"}
	// 地
	GEO = &material{name: "GEO", color: color.RGBA{R: 102, G: 51, B: 0, A: 255}, tag: "Material/GEO"}
	// エーテル(光)
	AETHER = &material{name: "AETHER", color: color.RGBA{R: 255, G: 255, B: 170, A: 255}, tag: "Material/AETHER"}
	// 血(闇)
	BLOOD = &material{name: "BLOOD", color: color.RGBA{R: 170, G: 0, B: 0, A: 255}, tag: "Material/BLOOD"}
	// 過去
	PALEO = &material{name: "PALEO", color: color.RGBA{R: 128, G: 128, B: 128, A: 255}, tag: "Material/PALEO"}
	// 未来
	MELO = &material{name: "MELO", color: color.RGBA{R: 0, G: 255, B: 255, A: 255}, tag: "Material/MELO"}
	// 忠誠
	ALLEGRO = &material{name: "ALLEGRO", color: color.RGBA{R: 255, G: 170, B: 0, A: 255}, tag: "Material/ALLEGRO"}
	// 信仰
	FAITH = &material{name: "FAITH", color: color.RGBA{R: 255, G: 255, B: 255, A: 255}, tag: "Material/FAITH"}
)

func (m Material) Name() string       { return m.name }
func (m Material) Color() color.Color { return m.color }
func (m Material) Tag() pkg.ItemTag   { return m.tag }

type Materials map[Material]int

func (m Materials) Copy() Materials {
	dst := make(map[Material]int, len(m))
	for k, v := range m {
		dst[k] = v
	}
	return dst
}

func EnumerateMaterial() []Material {
	return []Material{PYRO, HYDRO, AERO, NATURO, GEO, AETHER, BLOOD, PALEO, MELO, ALLEGRO, FAITH}
}
