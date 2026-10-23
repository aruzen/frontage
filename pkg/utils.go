package pkg

import (
	"fmt"
	"github.com/google/uuid"
	"math"
	"strconv"
)

import "log/slog"

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// LocalizeTag 責務がローカライズ用からResourceIDとしての役割に変化してしまっている
type LocalizeTag string

type Localized interface {
	LocalizeTag() LocalizeTag
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
		grid[i] = data[i*size.Height : (i+1)*size.Height]
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
		grid[i] = data[i*size.Height : (i+1)*size.Height]
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

// ToInt converts numeric values to int without truncation or overflow.
// It rejects fractions, non-finite values and values outside the int range.
func ToInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case int64:
		converted := int(val)
		if int64(converted) != val {
			return 0, fmt.Errorf("value %d is outside the int range", val)
		}
		return converted, nil
	case float64:
		return floatToInt(val)
	case float32:
		return floatToInt(float64(val))
	case uint:
		if val > ^uint(0)>>1 {
			return 0, fmt.Errorf("value %d is outside the int range", val)
		}
		return int(val), nil
	case uint64:
		if val > uint64(^uint(0)>>1) {
			return 0, fmt.Errorf("value %d is outside the int range", val)
		}
		return int(val), nil
	case nil:
		return 0, fmt.Errorf("value is nil")
	default:
		return 0, fmt.Errorf("expected number, got %T", v)
	}
}

func floatToInt(value float64) (int, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || math.Trunc(value) != value {
		return 0, fmt.Errorf("expected a finite integer, got %v", value)
	}
	limit := math.Ldexp(1, strconv.IntSize-1)
	if value < -limit || value >= limit {
		return 0, fmt.Errorf("value %v is outside the int range", value)
	}
	return int(value), nil
}

func ToUUID(v interface{}) (uuid.UUID, error) {
	if v == nil {
		return uuid.Nil, fmt.Errorf("value is nil")
	}
	s, ok := v.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("expected string value, got %T", v)
	}
	return uuid.Parse(s)
}

// PointToMap converts pkg.Point to a generic map for serialization.
func PointToMap(p Point) map[string]interface{} {
	return map[string]interface{}{"x": p.X, "y": p.Y}
}

// PointFromMap parses map into pkg.Point.
func PointFromMap(v interface{}) (Point, error) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return Point{}, fmt.Errorf("expected map for point, got %T", v)
	}
	x, err := ToInt(m["x"])
	if err != nil {
		return Point{}, fmt.Errorf("x: %w", err)
	}
	y, err := ToInt(m["y"])
	if err != nil {
		return Point{}, fmt.Errorf("y: %w", err)
	}
	return Point{X: x, Y: y}, nil
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
