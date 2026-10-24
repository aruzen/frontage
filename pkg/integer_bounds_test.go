package pkg

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestToIntRejectsLossyConversions(t *testing.T) {
	limit := math.Ldexp(1, strconv.IntSize-1)
	for _, value := range []interface{}{
		1.5, -1.5, float32(1.5), math.NaN(), math.Inf(1), math.Inf(-1),
		float32(math.Inf(1)), limit, -2 * limit, uint64(math.MaxUint64), ^uint(0),
	} {
		if got, err := ToInt(value); err == nil {
			t.Errorf("ToInt(%v [%T]) = %d without an error", value, value, got)
		}
	}
	if strconv.IntSize == 32 {
		for _, value := range []int64{math.MaxInt32 + 1, math.MinInt32 - 1} {
			if _, err := ToInt(value); err == nil {
				t.Errorf("out of range int64 %d accepted", value)
			}
		}
	}
}

func TestToIntPreservesRepresentableIntegers(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	limit := math.Ldexp(1, strconv.IntSize-1)
	largestFloat := math.Trunc(math.Nextafter(limit, 0))
	for _, tc := range []struct {
		value interface{}
		want  int
	}{
		{int(0), 0}, {int64(-7), -7}, {float32(12), 12}, {float64(-12), -12},
		{uint(4), 4}, {uint64(5), 5}, {float64(0), 0},
		{int64(maxInt), maxInt}, {int64(minInt), minInt}, {uint64(maxInt), maxInt},
		{-limit, minInt}, {largestFloat, int(largestFloat)},
	} {
		got, err := ToInt(tc.value)
		if err != nil || got != tc.want {
			t.Errorf("ToInt(%v [%T]) = %d, %v; want %d", tc.value, tc.value, got, err, tc.want)
		}
	}
}

func TestPointFromMapRejectsFractionalCoordinates(t *testing.T) {
	for _, field := range []string{"x", "y"} {
		input := map[string]interface{}{"x": float64(1), "y": float64(2)}
		input[field] = 1.5
		if _, err := PointFromMap(input); err == nil || !strings.Contains(err.Error(), field) {
			t.Errorf("fractional %s did not produce a coordinate error: %v", field, err)
		}
	}
}
