package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

// 清理标点并标准化为小写
func cleanWord(word string) string {
	return strings.ToLower(strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r)
	}))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run wordfreq.go <文件路径>")
		return
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	wordCount := make(map[string]int)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := cleanWord(scanner.Text())
		if word != "" {
			wordCount[word]++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}

	// 将 map 转为切片并按频率排序
	type wordFreq struct {
		Word  string
		Count int
	}
	var freqs []wordFreq
	for w, c := range wordCount {
		freqs = append(freqs, wordFreq{w, c})
	}
	sort.Slice(freqs, func(i, j int) bool {
		return freqs[i].Count > freqs[j].Count
	})

	// 输出前 20 个高频词
	fmt.Println("Top 20 词频:")
	for i, wf := range freqs {
		if i >= 20 {
			break
		}
		fmt.Printf("%2d. %-15s %d\n", i+1, wf.Word, wf.Count)
	}
}
