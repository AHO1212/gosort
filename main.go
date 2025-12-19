package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	firstName = "ahmet_can"
	lastName  = "karayoluk"
	studentID = "231ADB260"
)

func main() {
	if len(os.Args) < 2 {
		printUsageAndExit("no mode provided")
	}

	mode := os.Args[1]
	args := os.Args[2:]

	switch mode {
	case "-r":
		handleRandomMode(args)
	case "-i":
		handleInputFileMode(args)
	case "-d":
		handleDirectoryMode(args)
	default:
		printUsageAndExit("unknown mode: " + mode)
	}
}

func printUsageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, "Error:", msg)
	}
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  gosort -r N")
	fmt.Fprintln(os.Stderr, "  gosort -i input.txt")
	fmt.Fprintln(os.Stderr, "  gosort -d incoming")
	os.Exit(1)
}

func processNumbers(nums []int, verbose bool) []int {
	if verbose {
		fmt.Printf("Original numbers: %v\n", nums)
	}

	chunks := chunkSlice(nums)

	if verbose {
		fmt.Println("Chunks before sorting:")
		printChunks(chunks)
	}

	sortChunksConcurrently(chunks)

	if verbose {
		fmt.Println("Chunks after sorting:")
		printChunks(chunks)
	}

	sorted := mergeSortedChunks(chunks)

	if verbose {
		fmt.Printf("Final sorted result: %v\n", sorted)
	}

	return sorted
}

func chunkSlice(nums []int) [][]int {
	n := len(nums)
	if n == 0 {
		return nil
	}

	chunksCount := int(math.Ceil(math.Sqrt(float64(n))))
	if chunksCount < 4 {
		chunksCount = 4
	}
	if chunksCount > n {
		chunksCount = n
	}

	baseSize := n / chunksCount
	remainder := n % chunksCount

	chunks := make([][]int, 0, chunksCount)
	start := 0
	for i := 0; i < chunksCount; i++ {
		size := baseSize
		if i < remainder {
			size++
		}
		end := start + size
		chunks = append(chunks, nums[start:end])
		start = end
	}

	return chunks
}

func printChunks(chunks [][]int) {
	for i, c := range chunks {
		fmt.Printf("  Chunk %d: %v\n", i, c)
	}
}

func sortChunksConcurrently(chunks [][]int) {
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for i := range chunks {
		go func(idx int) {
			defer wg.Done()
			sort.Ints(chunks[idx])
		}(i)
	}

	wg.Wait()
}

func mergeTwoSorted(a, b []int) []int {
	result := make([]int, 0, len(a)+len(b))
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}

	result = append(result, a[i:]...)
	result = append(result, b[j:]...)
	return result
}

func mergeSortedChunks(chunks [][]int) []int {
	if len(chunks) == 0 {
		return nil
	}
	result := chunks[0]
	for i := 1; i < len(chunks); i++ {
		result = mergeTwoSorted(result, chunks[i])
	}
	return result
}

func readIntsFromFile(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", path, err)
	}
	defer file.Close()

	var nums []int
	scanner := bufio.NewScanner(file)
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		val, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid integer on line %d: %v", lineNo, err)
		}
		nums = append(nums, val)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}

	return nums, nil
}

func writeIntsToFile(path string, nums []int) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create output file %s: %w", path, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, v := range nums {
		if _, err := fmt.Fprintf(writer, "%d\n", v); err != nil {
			return fmt.Errorf("failed writing to %s: %w", path, err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed flushing to %s: %w", path, err)
	}
	return nil
}

func handleRandomMode(args []string) {
	if len(args) != 1 {
		printUsageAndExit("incorrect -r usage")
	}

	n, err := strconv.Atoi(args[0])
	if err != nil || n < 10 {
		fmt.Fprintln(os.Stderr, "Error: N must be an integer >= 10")
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = rand.Intn(1000)
	}

	processNumbers(nums, true)
}

func handleInputFileMode(args []string) {
	if len(args) != 1 {
		printUsageAndExit("incorrect -i usage")
	}

	inputPath := args[0]
	nums, err := readIntsFromFile(inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input file:", err)
		os.Exit(1)
	}

	if len(nums) < 10 {
		fmt.Fprintln(os.Stderr, "Error: input file must contain at least 10 integers")
		os.Exit(1)
	}

	processNumbers(nums, true)
}

func handleDirectoryMode(args []string) {
	if len(args) != 1 {
		printUsageAndExit("incorrect -d usage")
	}

	dir := args[0]
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		fmt.Fprintln(os.Stderr, "Error:", dir, "is not a directory")
		os.Exit(1)
	}

	parent := filepath.Dir(dir)
	base := filepath.Base(dir)

	outDirName := fmt.Sprintf("%s_sorted_%s_%s_%s", base, firstName, lastName, studentID)
	outDirPath := filepath.Join(parent, outDirName)

	if err := os.MkdirAll(outDirPath, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Error creating output directory:", err)
		os.Exit(1)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading directory:", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".txt" {
			continue
		}

		wg.Add(1)
		go func(filename string) {
			defer wg.Done()

			inPath := filepath.Join(dir, filename)
			nums, err := readIntsFromFile(inPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Skipping %s: %v\n", filename, err)
				return
			}

			sorted := processNumbers(nums, false)

			outPath := filepath.Join(outDirPath, filename)
			if err := writeIntsToFile(outPath, sorted); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outPath, err)
			}
		}(entry.Name())
	}

	wg.Wait()
}
