package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func FIFO(pages []int, capacity int) (faults int, faultRate float64) {
	interrupts := make([]string, len(pages))
	obsolete := make([]int, len(pages))
	storage := make([][]int, len(pages))
	faults = 0
	faultRate = 0.0
	for i := range storage {
		storage[i] = make([]int, capacity)
	}
	for i, m := range pages {
		if i == 0 {
			storage[0][0] = pages[i]
			interrupts[i] = "×"
			obsolete[i] = 0
		} else {
			storage[i] = storage[i-1]
			if ifInStorage(storage[i], m) {
				interrupts[i] = "√"
				obsolete[i] = 0
				continue
			}
			faults += 1
			interrupts[i] = "×"
			obsolete[i] = storage[i][capacity-1]
			storage[i] = append([]int{m}, storage[i]...)
			storage[i] = storage[i][0 : len(storage[i])-1]
		}
	}
	fmt.Println()
	faultRate = float64(faults) / float64(len(pages))
	PrintStorage("FIFO", pages, storage, capacity, interrupts, obsolete)
	return faults, faultRate
}

func LRU(pages []int, capacity int) (faults int, faultRate float64) {
	interrupts := make([]string, len(pages))
	obsolete := make([]int, len(pages))
	storage := make([][]int, len(pages))
	faults = 0
	faultRate = 0.0
	for i := range storage {
		storage[i] = make([]int, capacity)
	}
	for i, m := range pages {
		if i == 0 {
			storage[0][0] = pages[i]
			interrupts[i] = "×"
			obsolete[i] = 0
		} else {
			storage[i] = storage[i-1]
			if ifInStorage(storage[i], m) {
				storage[i] = append([]int{m}, storage[i]...)
				storage[i] = storage[i][0 : len(storage[i])-1]
				interrupts[i] = "√"
				obsolete[i] = 0
				continue
			}
			faults += 1
			interrupts[i] = "×"
			obsolete[i] = storage[i][capacity-1]
			storage[i] = append([]int{m}, storage[i]...)
			storage[i] = storage[i][0 : len(storage[i])-1]
		}
	}
	fmt.Println()
	faultRate = float64(faults) / float64(len(pages))
	PrintStorage("LRU", pages, storage, capacity, interrupts, obsolete)
	return faults, faultRate

}

func ifInStorage(storage []int, pageNum int) bool {
	for _, m := range storage {
		if pageNum == m {
			return true
		}
	}
	return false
}

func ReadPages(filePath string) ([]int, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return nil, err
	}
	defer file.Close()

	// 创建一个切片来存储数字
	var pages []int

	// 创建一个读取器以逐行读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 将每行的数字拆分并添加到切片中
		parts := strings.Fields(line)
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				fmt.Printf("无法将字符串转换为整数: %v\n", err)
				continue
			}
			pages = append(pages, num)
		}
	}

	// 检查扫描过程中是否出现错误
	if err := scanner.Err(); err != nil {
		fmt.Println("扫描文件时发生错误:", err)
		return nil, err
	}
	return pages, nil
}

func PrintStorage(method string, pages []int, storage [][]int, capacity int, interrupts []string, obsolete []int) {
	fmt.Printf("%s\t", method)
	for _, m := range pages {
		fmt.Printf("%d\t", m)
	}
	for i := 0; i < capacity; i++ {
		fmt.Printf("\n页%d\t", i+1)
		for j := 0; j < len(pages); j++ {
			if storage[j][i] == 0 {
				fmt.Printf(" \t")
			} else {
				fmt.Printf("%d\t", storage[j][i])
			}
		}
	}
	fmt.Printf("\n\t")
	for i := 0; i < len(interrupts); i++ {
		fmt.Printf("%s\t", interrupts[i])
	}
	fmt.Printf("\n淘汰\t")
	for i := 0; i < len(obsolete); i++ {
		if obsolete[i] == 0 {
			fmt.Printf("-\t")
		} else {
			fmt.Printf("%d\t", obsolete[i])
		}

	}
}
func main() {
	pages, err := ReadPages("/Users/barrywu/CCNU-OS/Lab/Lab4/page")
	if err != nil {
		panic(err)
	}
	capacity := 3
	faults, faultRate := FIFO(pages, capacity)
	fmt.Printf("\n缺页的总次数:%d\t缺页中断率:%f", faults, faultRate)
	fmt.Println()
	faults, faultRate = LRU(pages, capacity)
	fmt.Printf("\n缺页的总次数:%d\t缺页中断率:%f", faults, faultRate)
}
