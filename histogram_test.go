package main

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	h := Histogram{}
	h.Init()

	if h.maxBins != defaultMaxBins {
		t.Errorf("Expect %v as maximum number of bins, but got %v",
			defaultMaxBins, h.maxBins)
	}
	if h.bins != nil {
		t.Errorf("Expect bins to be nil, but got %v", h.bins)
	}
}

func TestNew(t *testing.T) {
	maxBins := 16
	h := New(maxBins)

	if h.maxBins != maxBins {
		t.Errorf("Expect %v as maximum number of bins, but got %v",
			maxBins, h.maxBins)
	}
	if h.bins != nil {
		t.Errorf("Expect bins to be nil, but got %v", h.bins)
	}
}

func TestUpdateZeroInit(t *testing.T) {
	h := Histogram{}
	h.Update(1.0)

	if h.maxBins != defaultMaxBins {
		t.Errorf("Expect %v as maximum number of bins, but got %v",
			defaultMaxBins, h.maxBins)
	}
	want := []bin{
		bin{Val: 1.0, Count: 1, Min: 1.0, Max: 1.0},
	}
	if !reflect.DeepEqual(h.bins, want) {
		t.Errorf(
			"Update failed: got: %v, want: %v", h.bins, want)
	}
}

func TestUpdateNegMaxBins(t *testing.T) {
	h := Histogram{-1, nil}
	h.Update(1.0)

	if h.maxBins != -1 {
		t.Errorf("Expect %v as maximum number of bins, but got %v",
			-1, h.maxBins)
	}
	if h.bins != nil {
		t.Errorf(
			"Update failed: got: %v, want: %v", h.bins, nil)
	}
}

func TestUpdateExact(t *testing.T) {
	h := Histogram{}
	for i := 0; i < 3; i++ {
		h.Update(1.0)
	}

	want := []bin{
		bin{Val: 1.0, Count: 3, Min: 1.0, Max: 1.0},
	}
	if !reflect.DeepEqual(h.bins, want) {
		t.Errorf(
			"Update failed: got: %v, want: %v", h.bins, want)
	}
}

func TestUpdateInsert(t *testing.T) {
	h := Histogram{}
	h.Update(1.0)
	h.Update(0.0)
	h.Update(2.0)

	want := []bin{
		bin{Val: 0.0, Count: 1, Min: 0.0, Max: 0.0},
		bin{Val: 1.0, Count: 1, Min: 1.0, Max: 1.0},
		bin{Val: 2.0, Count: 1, Min: 2.0, Max: 2.0},
	}
	if !reflect.DeepEqual(h.bins, want) {
		t.Errorf(
			"Update failed: got: %v, want: %v", h.bins, want)
	}
}

func TestUpdatePrune(t *testing.T) {
	h := Histogram{maxBins: 3, bins: nil}
	for i := 0; i < 4; i++ {
		h.Update(float64(i))
	}

	want := []bin{
		// Bins with values 0.0 and 1.0 are merged together.
		bin{Val: 0.5, Count: 2, Min: 0.0, Max: 1.0},
		bin{Val: 2.0, Count: 1, Min: 2.0, Max: 2.0},
		bin{Val: 3.0, Count: 1, Min: 3.0, Max: 3.0},
	}
	if !reflect.DeepEqual(h.bins, want) {
		t.Errorf(
			"Update failed: got: %v, want: %v", h.bins, want)
	}
}

func TestScale(t *testing.T) {
	h := Histogram{}
	for i := 0; i < 3*dotsMaxWidth; i++ {
		h.Update(0.0)
	}

	got := h.scale()
	want := 3
	if want != got {
		t.Errorf("Scale failed: got = %v, want = %v", got, want)
	}
}
