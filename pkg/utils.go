package pkg

import "log/slog"

type Size struct {
	Width  int
	Height int
}

type Point struct {
	X int
	Y int
}

type ItemTag string

type HaveItemTag interface {
	Tag() ItemTag
}

type PacketHeader struct {
	Tag  uint16
	Size uint32
}

func SizeToPoint(s Size) Point {
	return Point{
		X: s.Width,
		Y: s.Height,
	}
}

func PointToSize(s Point) Size {
	return Size{
		Width:  s.X,
		Height: s.Y,
	}
}

func Make2D[T any](size Size, init T) [][]T {
	data := make([]T, size.Height*size.Width)
	grid := make([][]T, size.Width)
	for i := range grid {
		grid[i] = data[i*size.Width : (i+1)*size.Width]
		for j := range grid[i] {
			grid[i][j] = init
		}
	}
	return grid
}

func Copy2D[T any](size Size, o [][]T) [][]T {
	data := make([]T, size.Height*size.Width)
	grid := make([][]T, size.Width)
	for i := range grid {
		grid[i] = data[i*size.Width : (i+1)*size.Width]
		for j := range grid[i] {
			grid[i][j] = o[i][j]
		}
	}
	return grid
}

func Replace[T comparable](data []T, target T, value T) bool {
	for i := range data {
		if data[i] == target {
			data[i] = value
			return true
		}
	}
	return false
}

func Replace2D[T comparable](data [][]T, target T, value T) bool {
	for i := range data {
		if Replace(data[i], target, value) {
			return true
		}
	}
	return false
}

func ReplacePtr[T any](data []*T, target *T, value *T) bool {
	for i, current := range data {
		if current == target {
			data[i] = value
			return true
		}
	}
	return false
}

func Replace2DPtr[T any](data [][]*T, target *T, value *T) bool {
	for _, row := range data {
		if ReplacePtr(row, target, value) {
			return true
		}
	}
	return false
}

func Overwrite[T any](dist, src []T) {
	if len(dist) != len(src) {
		slog.Warn("len(dist) != len(src)", "src", len(src), "dist", len(dist))
		for i := 0; i < min(len(dist), len(src)); i++ {
			(dist)[i] = src[i]
		}
	} else {
		copy(dist, src)
	}
}

func Overwrite2D[T any](dist [][]T, src [][]T) {
	if len(dist) != len(src) {
		slog.Warn("len(dist) != len(src)", "src", len(src), "dist", len(dist))
		for i := 0; i < min(len(dist), len(src)); i++ {
			Overwrite(dist[i], src[i])
		}
	} else {
		for i := range dist {
			Overwrite(dist[i], src[i])
		}
	}
}

func CopyMap[K comparable, V any](v map[K]V) map[K]V {
	m := make(map[K]V)
	for key, value := range v {
		m[key] = value
	}
	return m
}
