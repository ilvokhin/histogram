package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

type bin struct {
	Val   float64
	Count int
	Min   float64
	Max   float64
}

func merge(x, y bin) bin {
	count := x.Count + y.Count
	val := (x.Val*float64(x.Count) + y.Val*float64(y.Count)) / float64(count)
	return bin{
		val,
		count,
		min(x.Min, y.Min),
		max(x.Max, y.Max),
	}
}

const (
	defaultMaxBins = 10
	dotsMaxWidth   = 79 / 2
	dot            = "*"
)

type Histogram struct {
	maxBins int
	bins    []bin
}

func (h *Histogram) Init() {
	h.maxBins = defaultMaxBins
	if h.bins != nil {
		h.bins = nil
	}
}

func (h *Histogram) maybeInit() {
	if h.maxBins == 0 {
		h.Init()
	}
}

func New(maxBins int) *Histogram {
	return &Histogram{maxBins, nil}
}

func (h *Histogram) Update(val float64) {
	h.maybeInit()
	if h.maxBins <= 0 {
		return
	}
	pos, found := slices.BinarySearchFunc(h.bins, val,
		func(bin bin, val float64) int {
			return cmp.Compare(bin.Val, val)
		})
	if found {
		h.bins[pos].Count++
		return
	}
	if h.bins == nil {
		// Add space for one temporary bin.
		h.bins = make([]bin, 0, h.maxBins+1)
	}
	h.bins = slices.Insert(h.bins, pos, bin{val, 1, val, val})
	if len(h.bins) > h.maxBins {
		h.prune()
	}
}

func (h *Histogram) prune() {
	if len(h.bins) < 2 {
		panic("invariant failed: len(bins) < 2")
	}
	minDiff := h.bins[1].Val - h.bins[0].Val
	minPos := 0
	for pos := 1; pos < len(h.bins)-1; pos++ {
		diff := h.bins[pos+1].Val - h.bins[pos].Val
		if diff < minDiff {
			minDiff = diff
			minPos = pos
		}
	}
	h.bins[minPos] = merge(h.bins[minPos], h.bins[minPos+1])
	h.bins = slices.Delete(h.bins, minPos+1, minPos+2)
}

func (h *Histogram) Visualize(w io.Writer) error {
	if h.bins == nil {
		return nil
	}
	scale := h.scale()
	fmt.Fprintf(w, "# each %v represents a count of %v\n", dot, scale)
	for pos, _ := range h.bins {
		err := h.visualizeOne(w, pos, scale)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Histogram) scale() int {
	if h.bins == nil {
		panic("invariant failed: h.bins == nil")
	}
	maxCount := h.bins[0].Count
	for _, bin := range h.bins[1:] {
		maxCount = max(maxCount, bin.Count)
	}
	scale := 1
	if maxCount > dotsMaxWidth {
		scale = maxCount / dotsMaxWidth
	}
	return scale
}

func (h *Histogram) visualizeOne(w io.Writer, pos, scale int) error {
	beg := h.bins[pos].Min
	end, interval := h.intervalEnd(pos)
	count := h.bins[pos].Count
	dots := strings.Repeat("*", count/scale)
	_, err := fmt.Fprintf(w, "[%10.4f, %10.4f%s - %6d: %s\n",
		beg, end, interval, count, dots)
	return err
}

func (h *Histogram) intervalEnd(pos int) (float64, string) {
	if pos == len(h.bins)-1 {
		return h.bins[pos].Max, "]"
	}
	return h.bins[pos+1].Min, ")"
}

func main() {
	h := Histogram{}
	stdin := bufio.NewReader(os.Stdin)
	for {
		var val float64
		_, err := fmt.Fscan(stdin, &val)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		h.Update(val)
	}
	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()
	err := h.Visualize(stdout)
	if err != nil {
		log.Fatal(err)
	}
}
