package main

import (
	"container/list"
	_ "container/list"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

const (
	Ready     = 1
	Waiting   = 2
	Running   = 3
	TimeSlice = 2
)

type Process struct {
	ID            int
	Status        int // 1: 就绪, 2: 等待, 3: 运行
	ExecutionTime int
	Priority      int // 优先数
}

func ReadProcess(filename string) ([]Process, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("无法读取CSV文件:", err)
		return nil, err
	}
	var processes []Process
	for i, record := range records {
		if i == 0 {
			// 跳过标题行
			continue
		}

		if len(record) != 4 {
			fmt.Printf("行 %d 数据项数量不正确\n", i+1)
			continue
		}

		id, _ := strconv.Atoi(record[0])
		status, _ := strconv.Atoi(record[1])
		executionTime, _ := strconv.Atoi(record[2])
		priority, _ := strconv.Atoi(record[3])

		process := Process{
			ID:            id,
			Status:        status,
			ExecutionTime: executionTime,
			Priority:      priority,
		}

		processes = append(processes, process)
	}
	return processes, nil
}

func ExecutionSequence(processID int, averageTime int) {
	fmt.Printf("%d\t\t%d\n", processID, averageTime)
}

func PushInQueue(queue *list.List, processes []Process) {
	for _, process := range processes {
		if process.Status == Ready {
			queue.PushBack(process)
		}
	}

}

// FIFO
func FIFO(processes []Process) float64 {
	// 就绪队列
	readyQueue := list.New()

	// 初始化等待时间
	waitingTime := 0
	totalWaitingTime := 0
	copiedProcesses := make([]Process, len(processes))
	copy(copiedProcesses, processes)
	for _, process := range copiedProcesses {
		if process.Status == Ready {
			readyQueue.PushBack(process)
		}
	}
	for e := readyQueue.Front(); e != nil; e = e.Next() {
		process := e.Value.(Process)
		for i := range copiedProcesses {
			if copiedProcesses[i].ID == process.ID {
				copiedProcesses[i].Status = Running
			}
		}
		//fmt.Printf("%d ", process.ID)
		ExecutionSequence(process.ID, waitingTime)
		totalWaitingTime += waitingTime
		waitingTime = waitingTime + process.ExecutionTime
	}

	return float64(totalWaitingTime) / float64(len(processes))

}

// 时间片轮转调度算法，
type ProcessTime struct {
	ProcessID   int
	ProcessTime int
	WaitingTime int
}

func RoundRobin(processes []Process) float64 {
	var processTime = make([]ProcessTime, len(processes))
	copiedProcesses := make([]Process, len(processes))
	copy(copiedProcesses, processes)

	for i, _ := range processTime {
		processTime[i] = ProcessTime{
			ProcessID:   copiedProcesses[i].ID,
			ProcessTime: copiedProcesses[i].ExecutionTime,
			WaitingTime: 0,
		}
	}

	timeSlice := TimeSlice
	readyQueue := list.New()
	PushInQueue(readyQueue, copiedProcesses)
	// 初始化等待时间
	totalWaitingTime := 0
	fmt.Printf("执行顺序: ")
	for {
		e := readyQueue.Front()
		if e == nil {
			break
		}
		process := e.Value.(Process)
		fmt.Printf("%d ", process.ID)
		for i := range copiedProcesses {
			if copiedProcesses[i].ID == process.ID {
				copiedProcesses[i].Status = Running
				processTime[i].WaitingTime = totalWaitingTime
				if processTime[i].ProcessTime < timeSlice {
					totalWaitingTime += processTime[i].ProcessTime
				} else {
					totalWaitingTime += timeSlice
				}
				processTime[i].ProcessTime -= timeSlice
				if processTime[i].ProcessTime > 0 {
					readyQueue.MoveToBack(e)
				} else {
					readyQueue.Remove(e)
				}
			}
		}
	}
	sum := 0
	fmt.Printf("\n进程ID\t\t等待时间")

	for _, m := range processTime {
		sum += m.WaitingTime
		fmt.Printf("\n%d\t\t%d", m.ProcessID, m.WaitingTime)
	}

	return float64(sum) / float64(len(processes))
}

// 优先数调度算法
func PriorityScheduling(processes []Process) float64 {
	var processTime = make([]ProcessTime, len(processes))
	copiedProcesses := make([]Process, len(processes))
	copy(copiedProcesses, processes)
	for i, _ := range processTime {
		processTime[i] = ProcessTime{
			ProcessID:   copiedProcesses[i].ID,
			ProcessTime: copiedProcesses[i].ExecutionTime,
			WaitingTime: 0,
		}
	}
	readyQueue := list.New()
	totalWaitingTime := 0

	sort.Slice(copiedProcesses, func(i, j int) bool {
		return copiedProcesses[i].Priority < copiedProcesses[j].Priority
	})
	for _, m := range copiedProcesses {
		readyQueue.PushBack(m)
	}
	fmt.Printf("执行顺序: ")
	for {
		e := readyQueue.Front()
		if e == nil {
			break
		}
		process := e.Value.(Process)
		fmt.Printf("%d ", process.ID)
		for i, _ := range copiedProcesses {
			if processes[i].ID == process.ID {
				processes[i].Status = Running
				processTime[i].WaitingTime = totalWaitingTime
				totalWaitingTime += processes[i].ExecutionTime
				readyQueue.Remove(e)
			}
		}
	}
	sum := 0
	fmt.Printf("\n进程ID\t\t等待时间")
	for _, m := range processTime {
		sum += m.WaitingTime
		fmt.Printf("\n%d\t\t%d", m.ProcessID, m.WaitingTime)
	}
	return float64(sum) / float64(len(processes))
}

func sortByPriority(processes []Process) float64 {
	var processTime = make([]ProcessTime, len(processes))
	copiedProcesses := make([]Process, len(processes))
	copy(copiedProcesses, processes)
	maxPriority := 0
	for i, _ := range processTime {
		processTime[i] = ProcessTime{
			ProcessID:   copiedProcesses[i].ID,
			ProcessTime: copiedProcesses[i].ExecutionTime,
			WaitingTime: 0,
		}
		if processes[i].Priority > maxPriority {
			maxPriority = processes[i].Priority
		}
	}
	readyQueue := make([]list.List, maxPriority)
	timeSlice := make([]int, maxPriority)
	num := 1
	for i := 0; i < maxPriority; i++ {
		timeSlice[i] = num
		num += 1
	}

	for _, m := range processes {
		readyQueue[m.Priority-1].PushBack(m)
	}
	totalWaiting := 0
	fmt.Printf("执行顺序: ")
	for i, _ := range readyQueue {
		for {
			e := readyQueue[i].Front()
			if e == nil {
				break
			}
			process := e.Value.(Process)
			fmt.Printf("%d ", process.ID)
			for j, n := range copiedProcesses {
				if n.ID == process.ID {
					if i != maxPriority-1 {
						copiedProcesses[j].Status = Running
						processTime[j].WaitingTime = totalWaiting
						if processTime[j].ProcessTime <= timeSlice[i] {
							totalWaiting += processTime[j].ProcessTime
						} else {
							totalWaiting += timeSlice[i]
						}
						processTime[j].ProcessTime -= timeSlice[i]
						if processTime[j].ProcessTime > 0 {
							readyQueue[i+1].PushBack(copiedProcesses[j])
							readyQueue[i].Remove(e)
						} else {
							readyQueue[i].Remove(e)
						}
					} else {
						processTime[j].WaitingTime = totalWaiting
						if processTime[j].ProcessTime < timeSlice[i] {
							totalWaiting += processTime[j].ProcessTime
						} else {
							totalWaiting += timeSlice[i]
						}
						processTime[j].ProcessTime -= timeSlice[i]
						if processTime[j].ProcessTime > 0 {
							readyQueue[i].MoveToBack(e)
						} else {
							readyQueue[i].Remove(e)
						}
					}
				}
			}
		}
	}

	sum := 0
	fmt.Printf("\n进程ID\t\t等待时间")
	for _, m := range processTime {
		sum += m.WaitingTime
		fmt.Printf("\n%d\t\t%d", m.ProcessID, m.WaitingTime)
	}
	return float64(sum) / float64(len(processes))
}

func main() {
	processes, err := ReadProcess("/Users/barrywu/CCNU-OS/Lab/Lab2/process.csv")
	if err != nil {
		panic(err)
	}
	for _, m := range processes {
		fmt.Println(m)
	}

	fmt.Println("------------FIFO------------")
	fmt.Printf("执行序列\t等待时间\n")
	avg := FIFO(processes)
	fmt.Printf("平均等待时间:%.2f\n", avg)
	fmt.Println("------------时间片轮转------------")
	avg = RoundRobin(processes)
	fmt.Printf("\n平均等待时间:%.2f\n", avg)
	fmt.Println("------------优先数调度------------")
	avg = PriorityScheduling(processes)
	fmt.Printf("\n平均等待时间:%.2f\n", avg)
	fmt.Println("------------分级调度算法------------")
	avg = sortByPriority(processes)
	fmt.Printf("\n平均等待时间:%.2f\n", avg)
}
