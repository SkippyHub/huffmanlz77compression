package main

import (
	"container/heap"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"

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
	mixed := "AbRaCaDaBrAAAAAbbbb"

	fmt.Println("Uppercase:", upper, "Bytes:", []byte(upper))
	printASCIItoBitsAndMemory(upper)

	fmt.Println("Lowercase:", lower, "Bytes:", []byte(lower))
	printASCIItoBitsAndMemory(lower)

	fmt.Println("Mixed case:", mixed, "Bytes:", []byte(mixed))
	printASCIItoBitsAndMemory(mixed)

	// Build frequency table
	upperFrequencies := BuildFrequencyTable(upper)
	upperHuffmanTree := BuildTree(upperFrequencies)

	lowerFrequencies := BuildFrequencyTable(lower)
	lowerHuffmanTree := BuildTree(lowerFrequencies)

	mixedFrequencies := BuildFrequencyTable(mixed)
	mixedHuffmanTree := BuildTree(mixedFrequencies)

	upperEncoding := buildEncoding(upperHuffmanTree, "")
	lowerEncoding := buildEncoding(lowerHuffmanTree, "")
	mixedEncoding := buildEncoding(mixedHuffmanTree, "")

	fmt.Println("Uppercase encoding:", upperEncoding)
	upperEncoded := applyHuffmanEncoding(upper, upperEncoding)
	fmt.Println("Uppercase encoded:", string(upperEncoded))
	printStringBitsAndMemory(string(upperEncoded))

	fmt.Println("Lowercase encoding:", lowerEncoding)
	lowerEncoded := applyHuffmanEncoding(lower, lowerEncoding)
	fmt.Println("Lowercase encoded:", string(lowerEncoded))
	printStringBitsAndMemory(string(lowerEncoded))

	fmt.Println("Mixed case encoding:", mixedEncoding)
	mixedEncoded := applyHuffmanEncoding(mixed, mixedEncoding)
	fmt.Println("Mixed case encoded:", string(mixedEncoded))
	printStringBitsAndMemory(string(mixedEncoded))

	// Decode
	upperDecoded := applyHuffmanDecoding(upperEncoded, upperEncoding)
	fmt.Println("Uppercase decoded:", string(upperDecoded))

	lowerDecoded := applyHuffmanDecoding(lowerEncoded, lowerEncoding)
	fmt.Println("Lowercase decoded:", lowerDecoded)

	mixedDecoded := applyHuffmanDecoding(mixedEncoded, mixedEncoding)
	fmt.Println("Mixed case decoded:", mixedDecoded)

	// Generate Echarts
	GenerateEcharts(upperHuffmanTree, "upper")
	GenerateEcharts(lowerHuffmanTree, "lower")
	GenerateEcharts(mixedHuffmanTree, "mixed")

	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	// Apply shift string
	upperShifted := applyShiftString(upper)
	fmt.Println("uppercase shifted:", upperShifted)
	upperShiftedFrequencyTable := BuildFrequencyTable(upperShifted)
	upperShiftedHuffmanTree := BuildTree(upperShiftedFrequencyTable)

	upperShiftedEncoding := buildEncoding(upperShiftedHuffmanTree, "")
	fmt.Println("uppercase shifted encoding:", upperShiftedEncoding)

	upperShiftedEncoded := applyHuffmanEncoding(upperShifted, upperShiftedEncoding)
	fmt.Println("uppercase shifted encoded:", string(upperShiftedEncoded))
	printStringBitsAndMemory(string(upperShiftedEncoded))

	upperShiftedDecoded := applyHuffmanDecoding(upperShiftedEncoded, upperShiftedEncoding)
	fmt.Println("uppercase shifted decoded:", upperShiftedDecoded)

	// Remove shift string
	upperUnshifted := removeShiftString(upperShiftedDecoded)
	fmt.Println("uppercase unshifted:", upperUnshifted)

	lowerShifted := applyShiftString(lower)
	fmt.Println("lowercase shifted:", lowerShifted)
	lowerShiftedFrequencyTable := BuildFrequencyTable(lowerShifted)
	lowerShiftedHuffmanTree := BuildTree(lowerShiftedFrequencyTable)
	lowerShiftedEncoding := buildEncoding(lowerShiftedHuffmanTree, "")
	fmt.Println("lowercase shifted encoding:", lowerShiftedEncoding)

	lowerShiftedEncoded := applyHuffmanEncoding(lowerShifted, lowerShiftedEncoding)
	fmt.Println("lowercase shifted encoded:", string(lowerShiftedEncoded))
	printStringBitsAndMemory(string(lowerShiftedEncoded))

	lowerShiftedDecoded := applyHuffmanDecoding(lowerShiftedEncoded, lowerShiftedEncoding)
	fmt.Println("lowercase shifted decoded:", lowerShiftedDecoded)

	mixedShifted := applyShiftString(mixed)
	fmt.Println("mixedcase shifted:", mixedShifted)
	mixedShiftedFrequencyTable := BuildFrequencyTable(mixedShifted)
	mixedShiftedHuffmanTree := BuildTree(mixedShiftedFrequencyTable)
	mixedShiftedEncoding := buildEncoding(mixedShiftedHuffmanTree, "")
	fmt.Println("mixedcase shifted encoding:", mixedShiftedEncoding)

	mixedShiftedEncoded := applyHuffmanEncoding(mixedShifted, mixedShiftedEncoding)
	fmt.Println("mixedcase shifted encoded:", string(mixedShiftedEncoded))
	printStringBitsAndMemory(string(mixedShiftedEncoded))

	mixedShiftedDecoded := applyHuffmanDecoding(mixedShiftedEncoded, mixedShiftedEncoding)
	fmt.Println("mixedcase shifted decoded:", mixedShiftedDecoded)

	// Remove shift string
	mixedUnshifted := removeShiftString(mixedShiftedDecoded)
	fmt.Println("mixedcase unshifted:", mixedUnshifted)

	// Generate Echarts
	GenerateEcharts(upperShiftedHuffmanTree, "upper_shifted")
	GenerateEcharts(lowerShiftedHuffmanTree, "lower_shifted")
	GenerateEcharts(mixedShiftedHuffmanTree, "mixed_shifted")

	// // LZ77
	// fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	// The path to the file to compress
	filePath := "sample.html"

	// Read the file
	sampleHTMLbytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	// huffman encoding
	// apply shift string
	htmlshifted := applyShiftString(html.UnescapeString(string(sampleHTMLbytes)))

	// Build frequency table
	htmlfrequencies := BuildFrequencyTable(string(htmlshifted))
	// fmt.Println("HTML frequencies:", htmlfrequencies)

	// Build the Huffman tree
	htmltree := BuildTree(htmlfrequencies)

	// Build the encoding table
	htmlencoding := buildEncoding(htmltree, "")

	// Apply the encoding
	htmlencoded := applyHuffmanEncoding(string(htmlshifted), htmlencoding)

	// Print the encoded data
	// fmt.Println("Encoded data:", htmlencoded)

	GenerateEcharts(htmltree, "html")

	// Decode
	htmldecoded := applyHuffmanDecoding(htmlencoded, htmlencoding)

	fmt.Println("HTML decoded:", htmldecoded)

	// Remove shift string
	htmlunshifted := removeShiftString(htmldecoded)

	// Print the decoded data
	fmt.Println("Decoded data:", htmlunshifted)

	// compare the bytes lengthscompressed
	fmt.Println("Original size:", len(sampleHTMLbytes), "bytes")
	fmt.Println("Compressed size:", len(htmlencoded), "bytes")

}

func applyShiftString(s string) string {
	var shifted strings.Builder
	isShifted := false
	for _, c := range s {
		if unicode.IsUpper(c) && !isShifted {
			isShifted = true
			shifted.WriteRune('↑')
			shifted.WriteRune(unicode.ToLower(c))
		} else if unicode.IsLower(c) && isShifted {
			isShifted = false
			shifted.WriteRune('↓')
			shifted.WriteRune(c)
		} else if isShifted {
			shifted.WriteRune(unicode.ToLower(c))
		} else {
			shifted.WriteRune(c)
		}
	}
	return shifted.String()
}

func removeShiftString(s string) string {
	var unshifted strings.Builder
	isShifted := false
	for _, c := range s {
		if c == '↑' {
			isShifted = true
		} else if c == '↓' {
			isShifted = false
		} else if isShifted {
			unshifted.WriteRune(unicode.ToUpper(c))
		} else {
			unshifted.WriteRune(c)
		}
	}
	return unshifted.String()
}

// printASCIItoBitsAndMemory prints the binary representation of each character in the given string
// and calculates the memory used in bits.
func printASCIItoBitsAndMemory(s string) {
	bits := ""
	for _, c := range s {
		bits += fmt.Sprintf("%08b ", c)
	}
	fmt.Println("Bits:", bits)
	fmt.Println("Memory used:", len(s)*8, "bits")
}

func printStringBitsAndMemory(s string) {
	bits := ""
	for i, c := range s {
		bits += string(c)
		if (i+1)%8 == 0 {
			bits += " "
		}
	}
	fmt.Println("Bits:", bits)
	fmt.Println("Memory used:", len(s), "bits")
}

// buildEncoding takes a HuffmanTree and a prefix string and returns a map that represents the encoding of each character in the tree.
// If the node is a leaf node, the character value is mapped to the prefix.
// If the node is an internal node, the left child is assigned a prefix of "0" and the right child is assigned a prefix of "1".
// The function recursively builds the encoding for each subtree and merges the results into a single map.
func buildEncoding(node HuffmanTree, prefix string) map[rune]string {
	encoding := make(map[rune]string)
	if leaf, ok := node.(HuffmanLeaf); ok {
		encoding[leaf.Value] = prefix
	} else if n, ok := node.(HuffmanNode); ok {
		leftEncoding := buildEncoding(n.Left, prefix+"0")
		for k, v := range leftEncoding {
			encoding[k] = v
		}
		rightEncoding := buildEncoding(n.Right, prefix+"1")
		for k, v := range rightEncoding {
			encoding[k] = v
		}
	}
	return encoding
}

func applyHuffmanEncoding(s string, encoding map[rune]string) []rune {
	var encoded []rune
	for _, c := range s {
		// fmt.Println("Encoding:", string(c), encoding[c], []byte(encoding[c]))
		encoded = append(encoded, []rune(encoding[c])...)
	}
	// fmt.Println("Encoded:", string(encoded))

	return encoded
}
func reverseMap(m map[rune]string) map[string]rune {
	reversed := make(map[string]rune)
	for k, v := range m {
		reversed[v] = k
	}
	return reversed
}

func applyHuffmanDecoding(s []rune, encoding map[rune]string) string {
	reversed := reverseMap(encoding)
	var decoded strings.Builder
	var code strings.Builder
	for _, c := range s {
		code.WriteRune(c)
		if val, ok := reversed[code.String()]; ok {
			decoded.WriteRune(val)
			code.Reset()
		}
	}
	return decoded.String()
}

// LZ77
type LZ77Token struct {
	Distance int
	Length   int
	Next     byte
}

func LZ77Compress(input []byte, windowSize int) []LZ77Token {
	var result []LZ77Token
	for i := 0; i < len(input); {
		length, distance := longestMatch(input, i, windowSize)
		nextChar := byte(0)
		if i+length < len(input) {
			nextChar = input[i+length]
		}
		result = append(result, LZ77Token{Distance: distance, Length: length, Next: nextChar})
		i += length + 1
		if i >= len(input) {
			break
		}
	}
	return result
}

func longestMatch(data []byte, current int, windowSize int) (length, distance int) {
	start := max(0, current-windowSize)
	for i := start; i < current; i++ {
		l := 0
		for l < current-i && i+l < len(data) && current+l < len(data) && data[i+l] == data[current+l] {
			l++
		}
		if l > length {
			length = l
			distance = current - i
		}
	}
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func computeAndCompareCompressionRate(original []byte, compressed []LZ77Token, targetRate float64) {
	originalSize := len(original)
	compressedSize := len(compressed) * 3 // each LZ77Token consists of 3 parts

	compressionRate := float64(compressedSize) / float64(originalSize)

	fmt.Printf("Original size: %d bytes\n", originalSize)
	fmt.Printf("Compressed size: %d bytes\n", compressedSize)
	fmt.Printf("Compression rate: %.2f\n", compressionRate)

	if compressionRate < targetRate {
		fmt.Println("Compression rate is less than the target rate.")
	} else if compressionRate == targetRate {
		fmt.Println("Compression rate is equal to the target rate.")
	} else {
		fmt.Println("Compression rate is greater than the target rate.")
	}
}
