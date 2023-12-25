package testlearn

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) { //TestCase1
	val := Add(1, 2)
	if val != 3 {
		t.Errorf("Expected 3, val is  %d;want 1", val)
	}
}

func TestIntMinTableDriven(t *testing.T) {
	var tests = []struct {
		a, b int // 输入
		want int // 输出
	}{
		{0, 1, 0},
		{1, 0, 0},
		{2, -2, 0},
		{0, -1, -1},
		{-1, 0, -1},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d,%d", tt.a, tt.b) // 测试case名称
		t.Run(testname, func(t *testing.T) {
			ans := IntMin(tt.a, tt.b)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}
