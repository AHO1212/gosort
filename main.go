// Ahmet Can Karayoluk
// 231ADB260

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
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

func main() {
	rFlag := flag.Int("r", -1, "")
	iFlag := flag.String("i", "", "")
	dFlag := flag.String("d", "", "")
	flag.Parse()

	mode := 0
	if *rFlag != -1 {
		mode++
	}
	if *iFlag != "" {
		mode++
	}
	if *dFlag != "" {
		mode++
	}
	if mode != 1 {
		log.Fatal("Usage:\n  gosort -r N\n  gosort -i input.txt\n  gosort -d incoming")
	}

	var err error
	switch {
	case *rFlag != -1:
		err = runRandom(*rFlag)
	case *iFlag != "":
		err = runInputFile(*iFlag)
	case *dFlag != "":
		err = runDirectory(*dFlag)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func runRandom(n int) error {
	if n < 10 {
		return errors.New("N must be >= 10")
	}

	numbers := generateRandomNumbers(n)

	fmt.Println("Original numbers:")
	fmt.Println(numbers)

	chunks := splitIntoChunks(numbers)
	fmt.Println("\nChunks before sorting:")
	printChunks(chunks)

	sorted := sortChunksConcurrently(chunks)
	fmt.Println("\nChunks after sorting:")
	printChunks(sorted)

	merged := mergeSortedChunks(sorted)
	fmt.Println("\nFinal sorted result:")
	fmt.Println(merged)

	return nil
}

func runInputFile(path string) error {
	numbers, err := readIntegersFromFile(path)
	if err != nil {
		return err
	}
	if len(numbers) < 10 {
		return errors.New("input file must contain at least 10 integers")
	}

	fmt.Println("Original numbers:")
	fmt.Println(numbers)

	chunks := splitIntoChunks(numbers)
	fmt.Println("\nChunks before sorting:")
	printChunks(chunks)

	sorted := sortChunksConcurrently(chunks)
	fmt.Println("\nChunks after sorting:")
	printChunks(sorted)

	merged := mergeSortedChunks(sorted)
	fmt.Println("\nFinal sorted result:")
	fmt.Println(merged)

	return nil
}

func runDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	outDir := filepath.Base(dir) + "_sorted_ahmet_can_karayoluk_231ADB260"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	count := 0

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.ToLower(filepath.Ext(e.Name())) != ".txt" {
			continue
		}

		inPath := filepath.Join(dir, e.Name())
		numbers, err := readIntegersFromFile(inPath)
		if err != nil {
			return err
		}
		if len(numbers) < 10 {
			return fmt.Errorf("file %s has fewer than 10 integers", e.Name())
		}

		chunks := splitIntoChunks(numbers)
		sorted := sortChunksConcurrently(chunks)
		merged := mergeSortedChunks(sorted)

		outPath := filepath.Join(outDir, e.Name())
		if err := writeIntegersToFile(outPath, merged); err != nil {
			return err
		}

		count++
	}

	if count == 0 {
		return errors.New("no .txt files processed")
	}

	fmt.Printf("Directory mode finished: %d files processed.\n", count)
	return nil
}

func splitIntoChunks(numbers []int) [][]int {
	n := len(numbers)
	chunkCount := int(math.Ceil(math.Sqrt(float64(n))))
	if chunkCount < 4 {
		chunkCount = 4
	}
	if chunkCount > n {
		chunkCount = n
	}

	base := n / chunkCount
	rem := n % chunkCount

	chunks := make([][]int, 0, chunkCount)
	start := 0
	for i := 0; i < chunkCount; i++ {
		size := base
		if i < rem {
			size++
		}
		end := start + size
		chunks = append(chunks, numbers[start:end])
		start = end
	}
	return chunks
}

func sortChunksConcurrently(chunks [][]int) [][]int {
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for i := range chunks {
		go func(i int) {
			defer wg.Done()
			sort.Ints(chunks[i])
		}(i)
	}

	wg.Wait()
	return chunks
}

func mergeSortedChunks(chunks [][]int) []int {
	if len(chunks) == 0 {
		return []int{}
	}
	result := chunks[0]

	for i := 1; i < len(chunks); i++ {
		result = mergeTwoSorted(result, chunks[i])
	}
	return result
}

func mergeTwoSorted(a, b []int) []int {
	i, j := 0, 0
	res := make([]int, 0, len(a)+len(b))

	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			res = append(res, a[i])
			i++
		} else {
			res = append(res, b[j])
			j++
		}
	}

	res = append(res, a[i:]...)
	res = append(res, b[j:]...)
	return res
}

func generateRandomNumbers(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rand.Intn(1000)
	}
	return nums
}

func readIntegersFromFile(path string) ([]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []int
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		v, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid integer in file: %s", line)
		}
		out = append(out, v)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func writeIntegersToFile(path string, nums []int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, n := range nums {
		fmt.Fprintln(w, n)
	}
	return w.Flush()
}

func printChunks(chunks [][]int) {
	for i, c := range chunks {
		fmt.Printf("Chunk %d: %v\n", i, c)
	}
}
