package main

import (
	"container/heap"
	"fmt"
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
		children := append(treeToEcharts(n.Left, "L"), treeToEcharts(n.Right, "R")...)
		data = append(data, opts.TreeData{Name: fmt.Sprintf("%d", n.Frequency), Value: n.Frequency, Children: children})
	}
	return data
}

func GenerateEcharts(tree HuffmanTree) {
	page := components.NewPage()
	treeChart := charts.NewTree()
	treeChart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Huffman Tree"}))
	treeChart.AddSeries("tree", treeToEcharts(tree, "root"))
	page.AddCharts(treeChart)
	page.Render("huffman_tree.html")
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
	frequencies := BuildFrequencyTable(str)
	huffmanTree := BuildTree(frequencies)
	// symbols := map[rune]int{'a': 5, 'b': 9, 'c': 12, 'd': 13, 'e': 16, 'f': 45}
	// huffmanTree := BuildTree(symbols)

	// fmt.Println(huffmanTree)
	// huffmanTreeJson, err := json.MarshalIndent(huffmanTree, "", "  ")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(huffmanTreeJson))

	dot := GenerateDot(huffmanTree)
	fmt.Println(dot)
}

func printTree(node HuffmanTree, dot *strings.Builder, parent string) {
	switch n := node.(type) {
	case HuffmanLeaf:
		dot.WriteString(fmt.Sprintf("\"%s\" [label=\"%s:%d\"];\n", parent, string(n.Value), n.Frequency))
	case HuffmanNode:
		left := fmt.Sprintf("%sL", parent)
		right := fmt.Sprintf("%sR", parent)
		dot.WriteString(fmt.Sprintf("\"%s\" [label=\"%d\"];\n", parent, n.Frequency))
		dot.WriteString(fmt.Sprintf("\"%s\" -> \"%s\";\n", parent, left))
		dot.WriteString(fmt.Sprintf("\"%s\" -> \"%s\";\n", parent, right))
		printTree(n.Left, dot, left)
		printTree(n.Right, dot, right)
	}
}

func GenerateDot(tree HuffmanTree) string {
	var dot strings.Builder
	dot.WriteString("digraph HuffmanTree {\n")
	printTree(tree, &dot, "root")
	dot.WriteString("}\n")
	return dot.String()
}
