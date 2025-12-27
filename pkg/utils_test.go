package pkg

import "testing"

func TestMake2D_Rectangular(t *testing.T) {
	size := Size{Width: 2, Height: 3}
	grid := Make2D(size, 7)
	if len(grid) != size.Width {
		t.Fatalf("width mismatch: got %d want %d", len(grid), size.Width)
	}
	for x := range grid {
		if len(grid[x]) != size.Height {
			t.Fatalf("height mismatch at x=%d: got %d want %d", x, len(grid[x]), size.Height)
		}
		for y := range grid[x] {
			if grid[x][y] != 7 {
				t.Fatalf("unexpected init value at (%d,%d): %d", x, y, grid[x][y])
			}
		}
	}
}

func TestCopy2D_Rectangular(t *testing.T) {
	size := Size{Width: 3, Height: 2}
	original := Make2D(size, 0)
	original[1][0] = 9
	original[2][1] = 5

	copyGrid := Copy2D(size, original)
	if len(copyGrid) != size.Width {
		t.Fatalf("width mismatch: got %d want %d", len(copyGrid), size.Width)
	}
	for x := range copyGrid {
		if len(copyGrid[x]) != size.Height {
			t.Fatalf("height mismatch at x=%d: got %d want %d", x, len(copyGrid[x]), size.Height)
		}
	}
	if copyGrid[1][0] != 9 || copyGrid[2][1] != 5 {
		t.Fatalf("copy values mismatch: got %d/%d", copyGrid[1][0], copyGrid[2][1])
	}

	copyGrid[1][0] = 4
	if original[1][0] == copyGrid[1][0] {
		t.Fatalf("copy should not share backing array")
	}
}
