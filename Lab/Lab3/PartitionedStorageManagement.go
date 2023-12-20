package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Read interface {
	ReadFile(filename string) error
}

// 内存块结构体
type MemoryBlock struct {
	StartAddress int
	Size         int
	Occupied     bool
}

// 进程结构体
type Process struct {
	ID   int
	Size int
}
type MemoryBlocks []MemoryBlock
type Processes []Process

func (m *MemoryBlocks) ReadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("无法读取CSV文件:", err)
		return err
	}
	for i, record := range records {
		if i == 0 {
			continue
		}
		startAddress, _ := strconv.Atoi(record[0])
		size, _ := strconv.Atoi(record[1])
		status, _ := strconv.ParseBool(record[2])
		memoryBlock := MemoryBlock{
			StartAddress: startAddress,
			Size:         size,
			Occupied:     status,
		}
		*m = append(*m, memoryBlock)
	}
	return nil
}

func (p *Processes) ReadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("无法读取CSV文件:", err)
		return err
	}
	for i, record := range records {
		if i == 0 {
			continue
		}
		id, _ := strconv.Atoi(record[0])
		size, _ := strconv.Atoi(record[1])
		process := Process{
			ID:   id,
			Size: size,
		}
		*p = append(*p, process)
	}
	return nil
}

func PrintFreeMemory(memoryBlocks MemoryBlocks, assignment MemoryBlocks) {

	var newMemoryBlocks []MemoryBlock

	for _, block := range memoryBlocks {
		if block.Size != 0 {
			newMemoryBlocks = append(newMemoryBlocks, block)
		}
	}

	memoryBlocks = newMemoryBlocks
	fmt.Println("空闲区表")
	fmt.Printf("起始\t大小\t标志\n")
	for _, memory := range memoryBlocks {
		fmt.Printf("%d\t%d\t%t\n", memory.StartAddress, memory.Size, memory.Occupied)
	}
	fmt.Println("已分配区表")
	fmt.Printf("起始\t大小\t标志\n")
	sort.Slice(assignment, func(i, j int) bool {
		return assignment[i].StartAddress < assignment[j].StartAddress
	})
	for _, memory := range assignment {
		fmt.Printf("%d\t%d\t%t\n", memory.StartAddress, memory.Size, memory.Occupied)
	}
}

// 首次适应
func firstFit(memoryBlocks MemoryBlocks, processes Processes) {
	var copiedBlocks = make([]MemoryBlock, len(memoryBlocks))
	var assigment MemoryBlocks
	copy(copiedBlocks, memoryBlocks)
	sort.Slice(copiedBlocks, func(i, j int) bool {
		return copiedBlocks[i].StartAddress < copiedBlocks[j].StartAddress
	})
	for i, _ := range processes {
		for j, _ := range copiedBlocks {
			if copiedBlocks[j].Size >= processes[i].Size {
				assigment = append(assigment, MemoryBlock{
					StartAddress: copiedBlocks[j].StartAddress,
					Size:         processes[i].Size,
					Occupied:     true,
				})
				copiedBlocks[j].StartAddress += processes[i].Size
				copiedBlocks[j].Size -= processes[i].Size
				break
			}
		}
	}
	sort.Slice(copiedBlocks, func(i, j int) bool {
		return copiedBlocks[i].StartAddress < copiedBlocks[j].StartAddress
	})
	PrintFreeMemory(copiedBlocks, assigment)
}

// 最佳适应
func bestFit(memoryBlocks MemoryBlocks, processes Processes) {
	var copiedBlocks = make([]MemoryBlock, len(memoryBlocks))
	copy(copiedBlocks, memoryBlocks)
	var assigment MemoryBlocks
	sort.Slice(copiedBlocks, func(i, j int) bool {
		return copiedBlocks[i].Size < copiedBlocks[j].Size
	})
	for i, _ := range processes {
		fit := 65535
		fitIndex := 0
		for j, _ := range copiedBlocks {
			fitDegree := copiedBlocks[j].Size - processes[i].Size
			if fitDegree >= 0 && fitDegree < fit {
				fitIndex = j
				fit = fitDegree
			}
		}
		if fit != 65535 {
			assigment = append(assigment, MemoryBlock{
				StartAddress: copiedBlocks[fitIndex].StartAddress,
				Size:         processes[i].Size,
				Occupied:     true,
			})
			copiedBlocks[fitIndex].StartAddress += processes[i].Size
			copiedBlocks[fitIndex].Size -= processes[i].Size
		}
	}
	sort.Slice(copiedBlocks, func(i, j int) bool {
		return copiedBlocks[i].Size < copiedBlocks[j].Size
	})
	PrintFreeMemory(copiedBlocks, assigment)
}

// 最差适应
func worstFit(memoryBlocks MemoryBlocks, processes Processes) {
	var copiedBlocks = make([]MemoryBlock, len(memoryBlocks))
	copy(copiedBlocks, memoryBlocks)
	var assigment MemoryBlocks
	for i, _ := range processes {
		for j, _ := range copiedBlocks {
			sort.Slice(copiedBlocks, func(i, j int) bool {
				return copiedBlocks[i].Size > copiedBlocks[j].Size
			})
			if copiedBlocks[j].Size >= processes[i].Size {
				assigment = append(assigment, MemoryBlock{
					StartAddress: copiedBlocks[j].StartAddress,
					Size:         processes[i].Size,
					Occupied:     true,
				})
				copiedBlocks[j].StartAddress += processes[i].Size
				copiedBlocks[j].Size -= processes[i].Size
				break
			}
		}
	}
	sort.Slice(copiedBlocks, func(i, j int) bool {
		return copiedBlocks[i].Size > copiedBlocks[j].Size
	})
	PrintFreeMemory(copiedBlocks, assigment)
}

func main() {
	var memoryBlocks MemoryBlocks
	err := memoryBlocks.ReadFile("/Users/barrywu/CCNU-OS/Lab/Lab3/free_partition.csv")
	if err != nil {
		panic(err)
	}

	var processes Processes
	err = processes.ReadFile("/Users/barrywu/CCNU-OS/Lab/Lab3/partition_request_sequence.csv")
	if err != nil {
		panic(err)
	}
	fmt.Println("-------------------最先适应-------------------")
	firstFit(memoryBlocks, processes)
	fmt.Println("-------------------最佳适应-------------------")
	bestFit(memoryBlocks, processes)
	fmt.Println("-------------------最坏适应-------------------")
	worstFit(memoryBlocks, processes)
}
