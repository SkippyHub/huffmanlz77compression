package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func treeToEcharts(node HuffmanTree, parent string) []opts.TreeData {
	var data []opts.TreeData
	switch n := node.(type) {
	case HuffmanLeaf:
		data = append(data, opts.TreeData{Name: fmt.Sprintf("%s:%d", string(n.Value), n.Frequency), Value: n.Frequency})
	case HuffmanNode:
		leftChildren := treeToEcharts(n.Left, "L")
		rightChildren := treeToEcharts(n.Right, "R")
		children := make([]*opts.TreeData, len(leftChildren)+len(rightChildren))
		for i, v := range leftChildren {
			children[i] = &v
		}
		for i, v := range rightChildren {
			children[i+len(leftChildren)] = &v
		}
		data = append(data, opts.TreeData{Name: fmt.Sprintf("%d", n.Frequency), Value: n.Frequency, Children: children})
	}
	return data
}

func GenerateEcharts(tree HuffmanTree, variableName string) {
	page := components.NewPage()
	treeChart := charts.NewTree()
	treeChart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Huffman Tree"}),
	)

	maxDepth := getMaxDepth(tree) // Get the maximum node depth
	treeChart.AddSeries("tree", treeToEcharts(tree, "root"))
	treeChart.SetSeriesOptions(
		charts.WithTreeOpts(opts.TreeChart{InitialTreeDepth: maxDepth}),
	)
	page.AddCharts(treeChart)

	// Create the file name
	fileName := fmt.Sprintf("huffman_%s_tree.html", variableName)

	// Create a file
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Render the page to the file
	page.Render(f)
}

// getMaxDepth returns the maximum depth of the Huffman tree.
func getMaxDepth(node HuffmanTree) int {
	switch n := node.(type) {
	case HuffmanLeaf:
		return 0
	case HuffmanNode:
		leftDepth := getMaxDepth(n.Left)
		rightDepth := getMaxDepth(n.Right)
		if leftDepth > rightDepth {
			return leftDepth + 1
		}
		return rightDepth + 1
	}
	return 0
}

// A HuffmanTree interface represents a tree node in Huffman encoding.
type HuffmanTree interface {
	Freq() int
}

// HuffmanNode represents a node in the Huffman tree.
type HuffmanNode struct {
	Frequency int         `json:"freq"`
	Value     rune        `json:"value"`
	Left      HuffmanTree `json:"left"`
	Right     HuffmanTree `json:"right"`
}

// HuffmanLeaf represents a leaf in the Huffman tree.
type HuffmanLeaf struct {
	Frequency int    `json:"freq"`
	Value     rune   `json:"value"`
	Char      string `json:"char"`
}

// Freq returns the frequency of occurrence of a Huffman node.
func (hn HuffmanNode) Freq() int {
	return hn.Frequency
}

// Freq returns the frequency of occurrence of a Huffman leaf.
func (hl HuffmanLeaf) Freq() int {
	return hl.Frequency
}

// A PriorityQueue implements heap.Interface and holds Huffman trees.
type PriorityQueue []HuffmanTree

func (pq PriorityQueue) Len() int {
	return len(pq)
}

// Less returns true if the frequency of tree at index i is less than tree at index j.
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq() < pq[j].Freq()
}

// swap swaps two Huffman trees in a priority queue.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push adds x as an element to a priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(HuffmanTree)
	*pq = append(*pq, item)
}

// Pop removes the minimum element (according to Less) from a priority queue and returns it.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// BuildTree builds the Huffman tree for Huffman encoding and decoding.

func BuildTree(leaves []HuffmanLeaf) HuffmanTree {
	var trees PriorityQueue
	for _, leaf := range leaves {
		trees = append(trees, leaf)
	}

	heap.Init(&trees)

	for trees.Len() > 1 {
		// two trees with least frequency
		a := heap.Pop(&trees).(HuffmanTree)
		b := heap.Pop(&trees).(HuffmanTree)

		// put into new node and re-insert into queue
		heap.Push(&trees, HuffmanNode{
			a.Freq() + b.Freq(),
			0,
			a,
			b,
		})
	}

	return heap.Pop(&trees).(HuffmanTree)
}

// BuildFrequencyTable builds a frequency table from a string.
func BuildFrequencyTable(s string) []HuffmanLeaf {
	frequencies := make(map[rune]int)
	for _, char := range s {
		frequencies[char]++
	}

	var leaves []HuffmanLeaf
	for char, freq := range frequencies {
		leaves = append(leaves, HuffmanLeaf{Frequency: freq, Value: char, Char: string(char)})
	}

	return leaves
}

// main function to test the Huffman encoding and decoding.
func main() {

	str := "abracadabra"
	upper := strings.ToUpper(str)
	lower := strings.ToLower(str)
	mixed := "AbRaCaDaBrA"

	fmt.Println("Uppercase:", upper)
	fmt.Println("Lowercase:", lower)
	fmt.Println("Mixed case:", mixed)

	// Build frequency table
	upperFrequencies := BuildFrequencyTable(upper)
	upperHuffmanTree := BuildTree(upperFrequencies)

	lowerFrequencies := BuildFrequencyTable(lower)
	lowerHuffmanTree := BuildTree(lowerFrequencies)

	mixedFrequencies := BuildFrequencyTable(mixed)
	mixedHuffmanTree := BuildTree(mixedFrequencies)

	GenerateEcharts(upperHuffmanTree, "upper")
	GenerateEcharts(lowerHuffmanTree, "lower")
	GenerateEcharts(mixedHuffmanTree, "mixed")

}
