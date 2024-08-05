package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"
)

// Cache structure
type Cache struct {
	maxEntries int
	entries    map[string]string
	keys       []string
}

// NewCache creates a new cache with a given maximum number of entries
func NewCache(maxEntries int) *Cache {
	return &Cache{
		maxEntries: maxEntries,
		entries:    make(map[string]string),
		keys:       make([]string, 0, maxEntries),
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (string, bool) {
	val, exists := c.entries[key]
	return val, exists
}

// Set adds a key-value pair to the cache
func (c *Cache) Set(key, value string) {
	if _, exists := c.entries[key]; !exists {
		if len(c.entries) >= c.maxEntries {
			oldestKey := c.keys[0]
			c.keys = c.keys[1:]
			delete(c.entries, oldestKey)
		}
		c.keys = append(c.keys, key)
		c.entries[key] = value
	}
}

// Converts a binary string to its hexadecimal representation
func binToHex(binStr string) (string, error) {
	binBytes := make([]byte, (len(binStr)+7)/8)
	for i := 0; i < len(binStr); i += 8 {
		var binByte byte
		for j := 0; j < 8 && i+j < len(binStr); j++ {
			binByte = binByte<<1 | (binStr[i+j] - '0')
		}
		binBytes[i/8] = binByte
	}
	return strings.ToUpper(hex.EncodeToString(binBytes)), nil
}

// Converts a hexadecimal string to its binary representation
func hexToBin(hexStr string) (string, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	binStr := ""
	for _, b := range bytes {
		binStr += fmt.Sprintf("%08b", b)
	}
	return binStr, nil
}

// Converts mat.in to mat.in.x using caching
func convertWithCache(inputFile, outputFile string, cache *Cache) error {
	input, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)

	for scanner.Scan() {
		line := scanner.Text()
		if cachedValue, found := cache.Get(line); found {
			writer.WriteString(cachedValue + "\n")
		} else {
			parts := strings.Split(line, ":")
			matrixSize := parts[0]
			binaryStr := parts[1]

			hexStr, err := binToHex(binaryStr)
			if err != nil {
				return err
			}
			newLine := fmt.Sprintf("%s:%s", matrixSize, hexStr)
			cache.Set(line, newLine)
			writer.WriteString(newLine + "\n")
		}
	}
	writer.Flush()
	return scanner.Err()
}

// Converts mat.in to mat.in.x without caching
func convertWithoutCache(inputFile, outputFile string) error {
	input, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		matrixSize := parts[0]
		binaryStr := parts[1]

		hexStr, err := binToHex(binaryStr)
		if err != nil {
			return err
		}
		newLine := fmt.Sprintf("%s:%s", matrixSize, hexStr)
		writer.WriteString(newLine + "\n")
	}
	writer.Flush()
	return scanner.Err()
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: <mode> <input_file> <output_file> <cache_size>")
		return
	}

	mode := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]
	cacheSize := 5000 // Default cache size
	if len(os.Args) >= 5 {
		fmt.Sscanf(os.Args[4], "%d", &cacheSize)
	}

	cache := NewCache(cacheSize)

	switch mode {
	case "compress-cached":
		start := time.Now()
		if err := convertWithCache(inputFile, outputFile, cache); err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Printf("Cached conversion took %.2f seconds\n", time.Since(start).Seconds())
	case "compress-noncached":
		start := time.Now()
		if err := convertWithoutCache(inputFile, outputFile); err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Printf("Non-cached conversion took %.2f seconds\n", time.Since(start).Seconds())
	case "decompress-cached":
		start := time.Now()
		if err := convertWithCache(inputFile, outputFile, cache); err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Printf("Cached decompression took %.2f seconds\n", time.Since(start).Seconds())
	case "decompress-noncached":
		start := time.Now()
		if err := convertWithoutCache(inputFile, outputFile); err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Printf("Non-cached decompression took %.2f seconds\n", time.Since(start).Seconds())
	default:
		fmt.Println("Unknown mode. Use 'compress-cached', 'compress-noncached', 'decompress-cached', or 'decompress-noncached'.")
	}
}
