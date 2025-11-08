package pkg

type Size struct {
	Width  int
	Height int
}

type Point struct {
	X int
	Y int
}

type LocalizeTag string

type Localized interface {
	Name() LocalizeTag
	Description() LocalizeTag
	Summery() []Localized
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

func Make2D[T interface{}](size Size, init T) [][]T {
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

func Copy2D[T interface{}](size Size, o [][]T) [][]T {
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

func ReplacePtr[T interface{}](data []*T, target *T, value *T) bool {
	for i, current := range data {
		if current == target {
			data[i] = value
			return true
		}
	}
	return false
}

func Replace2DPtr[T interface{}](data [][]*T, target *T, value *T) bool {
	for _, row := range data {
		if ReplacePtr(row, target, value) {
			return true
		}
	}
	return false
}
