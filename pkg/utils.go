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

func Replace[T interface{}](data []T, target T) bool {
	for i := range data {
		if data[i] == target {
			data[i] = target
			return true
		}
	}
	return false
}

func Replace2D[T interface{}](data [][]T, target T) bool {
	for i := range data {
		if Replace(data[i], target) {
			return true
		}
	}
	return false
}

func ReplacePtr[T interface{}](data []*T, target *T) bool {
	for _, current := range data {
		if current == nil {
			continue
		}
		if target == nil {
			var zero T
			if *current == zero {
				*current = zero
				return true
			}
			continue
		}
		if *current == *target {
			*current = *target
			return true
		}
	}
	return false
}

func Replace2DPtr[T interface{}](data [][]*T, target *T) bool {
	for _, row := range data {
		if ReplacePtr(row, target) {
			return true
		}
	}
	return false
}
