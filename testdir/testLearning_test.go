package main

import "testing"

func TestAdd(t *testing.T) {//TestCase1
	val := Add(1, 2)
	if val != 3 {
		t.Errorf("Expected 3, val is  %d;want 1", val)
	}
}
