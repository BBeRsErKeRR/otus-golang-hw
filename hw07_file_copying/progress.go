package main

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultProgressChar      = '■'
	defaultEmptyProgressChar = '□'
	defaultRefreshRate       = time.Millisecond * 50
	defaultColor             = "\033[36m"
	endColor                 = "\033[0m"
	barLength                = 20
)

func CreateNew(total int) *Progress {
	return CreateNew64(int64(total))
}

func CreateNew64(total int64) *Progress {
	pb := &Progress{
		Total:             total,
		RefreshRate:       defaultRefreshRate,
		ProgressChar:      defaultProgressChar,
		EmptyProgressChar: defaultEmptyProgressChar,
		Color:             defaultColor,
		finish:            make(chan struct{}),
	}
	return pb
}

type Progress struct {
	Total             int64
	RefreshRate       time.Duration
	ProgressChar      rune
	EmptyProgressChar rune
	Color             string
	// Width             int

	prefix     string
	postfix    string
	current    int64
	previous   int64
	startTime  time.Time
	changeTime time.Time
	mu         sync.Mutex
	finish     chan struct{}
}

func (pb *Progress) Postfix(postfix string) {
	pb.postfix = postfix
}

func (pb *Progress) Prefix(prefix string) {
	pb.prefix = prefix
}

func (pb *Progress) Increment() int {
	return pb.Add(1)
}

func (pb *Progress) Add(add int) int {
	return int(pb.Add64(int64(add)))
}

func (pb *Progress) Add64(add int64) int64 {
	return atomic.AddInt64(&pb.current, add)
}

func (pb *Progress) Get() int64 {
	c := atomic.LoadInt64(&pb.current)
	return c
}

func (pb *Progress) Set(current int) *Progress {
	return pb.Set64(int64(current))
}

func (pb *Progress) Set64(current int64) *Progress {
	atomic.StoreInt64(&pb.current, current)
	return pb
}

func (pb *Progress) print(total, current int64) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	color := pb.Color
	pc := pb.ProgressChar
	epc := pb.EmptyProgressChar

	var percentBox, barBox, elapsedTimeBox, totalBox, out string

	var percent float64
	if total > 0 {
		percent = float64(current) / (float64(total) / float64(100))
	} else {
		percent = float64(current) / float64(100)
	}
	percentBox = fmt.Sprintf("  (%.3f%%)", percent)

	progressLength := int(barLength * percent / 100)
	emptyProgressLength := barLength - progressLength
	barBox = strings.Repeat(string(pc), progressLength)
	if emptyProgressLength > 0 {
		barBox += strings.Repeat(string(epc), emptyProgressLength)
	}

	fromStart := time.Since(pb.startTime)
	elapsedTimeBox = fmt.Sprintf(" -> %.2fs", fromStart.Seconds())

	totalBox = fmt.Sprintf(" (%v/%v)  ", current, total)

	// (1/5) [*][*][][][][] -> 10% 10sec
	out = pb.prefix + totalBox + color + barBox + endColor + percentBox + elapsedTimeBox + pb.postfix

	fmt.Print("\r" + out)
}

func (pb *Progress) Start() {
	pb.startTime = time.Now()
	pb.Update()
	go pb.refresher()
}

func (pb *Progress) Finish() {
	close(pb.finish)
	pb.Update()
	fmt.Println("")
}

func (pb *Progress) Update() {
	c := atomic.LoadInt64(&pb.current)
	p := atomic.LoadInt64(&pb.previous)
	t := atomic.LoadInt64(&pb.Total)
	if p != c {
		pb.mu.Lock()
		pb.changeTime = time.Now()
		pb.mu.Unlock()
		atomic.StoreInt64(&pb.previous, c)
	}
	pb.print(t, c)
}

func (pb *Progress) refresher() {
	for {
		select {
		case <-pb.finish:
			return
		case <-time.After(pb.RefreshRate):
			pb.Update()
		}
	}
}

func (pb *Progress) Write(p []byte) (n int) {
	n = len(p)
	pb.Add(n)
	return
}

func (pb *Progress) Read(p []byte) (n int) {
	n = len(p)
	pb.Add(n)
	return
}

func (pb *Progress) NewProxyReader(r io.Reader) *Reader {
	return &Reader{r, pb, time.Microsecond * 0}
}

func (pb *Progress) NewProxyFreezeReader(r io.Reader, d time.Duration) *Reader {
	return &Reader{r, pb, d}
}

func (pb *Progress) NewProxyWriter(r io.Writer) *Writer {
	return &Writer{r, pb}
}
