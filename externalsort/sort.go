//go:build !solution

package externalsort

import (
	"container/heap"
	"io"
	"os"
	"sort"
)

type Item struct {
	val string
	ind int
}

type MinHeap []Item

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].val <= h[j].val }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(Item))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func NewReader(r io.Reader) LineReader {
	return &LinesReader{r: r}
}

func NewWriter(w io.Writer) LineWriter {
	return &LinesWriter{w: w}
}
func Merge(w LineWriter, readers ...LineReader) error {
	mh := make(MinHeap, 0, 2)
	h := &mh
	for i, r := range readers {
		str, errRead := r.ReadLine()
		if errRead == nil || (str != "" && errRead == io.EOF) {
			heap.Push(h, Item{val: str, ind: i})
		}
	}
	for h.Len() > 0 {
		currMinItem := (heap.Pop(h)).(Item)
		errWrite := w.Write(currMinItem.val)
		if errWrite != nil {
			return errWrite
		}
		newstr, errRead := readers[currMinItem.ind].ReadLine()
		if errRead == nil || (newstr != "" && errRead == io.EOF) {
			heap.Push(h, Item{val: newstr, ind: currMinItem.ind})
		}
	}
	return nil
}

func Sort(w io.Writer, in ...string) error {
	ReadFileSet := make([]LineReader, 0, 2)
	for _, s := range in {
		file, err := os.OpenFile(s, os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		ReadFileSet = append(ReadFileSet, &LinesReader{r: file})
		FileContent := make([]string, 0, 2)
		for {
			str, err := ReadFileSet[len(ReadFileSet)-1].ReadLine()
			if str == "" && err == io.EOF {
				break
			}
			FileContent = append(FileContent, str)
		}
		sort.Slice(FileContent, func(i, j int) bool {
			return FileContent[i] < FileContent[j]
		})
		file.Seek(0, 0)
		for _, str := range FileContent {
			(&LinesWriter{w: file}).Write(str)
		}
		file.Seek(0, 0)
	}
	Output := &LinesWriter{w: w}
	Merge(Output, ReadFileSet...)
	return nil
}
