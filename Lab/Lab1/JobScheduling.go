package main

import (
	"fmt"
	"sort"
)

type Job struct {
	JobNumber              int
	ArrivalTime            int
	EstimatedTime          int
	StartTime              int
	EndTime                int
	TurnaroundTime         int
	WeightedTurnaroundTime int
}

type Time struct {
	Hour   int
	Minute int
}

func TimeAdd(initialTime int, addedTime int) int {
	hour := initialTime / 100
	minute := initialTime % 100
	new_minute := (addedTime + minute) % 60
	new_hour := (hour + (addedTime+minute)/60) % 24
	return 100*new_hour + new_minute
}

func PrintTitle() {
	fmt.Println("作业\t\t进入时间\t估计运行时间\t开始时间\t结束时间\t周转时间(分钟)\t带权周转时间\t")
}

func PrintJob(jobs []Job) {
	for i, m := range jobs {
		fmt.Printf("JOB%d\t", i+1)
		PrintTime(m.ArrivalTime)
		fmt.Printf("\t    %d \t ", m.EstimatedTime)
		PrintTime(m.StartTime)
		PrintTime(m.EndTime)
		fmt.Printf("\t\t%d\t", m.TurnaroundTime)
		fmt.Printf("\t%d\t", m.WeightedTurnaroundTime)
		fmt.Printf("\n")
	}
}

func PrintTime(originalTime int) {
	fmt.Printf("\t%02d:%02d\t", originalTime/100, originalTime%100)
}

func main() {
	//var jobs []Job
	//var num int
	//fmt.Printf("输入作业个数:")
	//fmt.Scan(&num)
	//for i := 0; i < num; i++ {
	//	var job Job
	//	fmt.Scan(&job.JobNumber, &job.ArrivalTime, &job.EstimatedTime)
	//	jobs = append(jobs, job)
	//}
	jobs := []Job{
		{1, 800, 50, 0, 0, 0, 0},
		{2, 815, 30, 0, 0, 0, 0},
		{3, 830, 25, 0, 0, 0, 0},
		{4, 835, 20, 0, 0, 0, 0},
		{5, 845, 15, 0, 0, 0, 0},
		{6, 900, 10, 0, 0, 0, 0},
		{7, 920, 5, 0, 0, 0, 0},
	}

	// 使用FIFO算法进行作业调度
	fmt.Println("FIFO 调度算法：")
	PrintTitle()
	scheduleFIFO(jobs)
	//
	//// 使用短作业优先（SJF）算法进行作业调度
	//fmt.Println("\n短作业优先（SJF）调度算法：")
	//sjfJobs := scheduleSJF(jobs)
	//printJobSequence(sjfJobs)
	//
	//// 使用最高响应比优先（HRRN）算法进行作业调度
	//fmt.Println("\n最高响应比优先（HRRN）调度算法：")
	//hrrnJobs := scheduleHRRN(jobs)
	//printJobSequence(hrrnJobs)
}

func scheduleFIFO(jobs []Job) {
	for i, m := range jobs {
		if i == 0 {
			m.StartTime = m.ArrivalTime

		}
		m.EndTime = TimeAdd(m.StartTime, m.EstimatedTime)
	}
	PrintJob(jobs)
}

func scheduleSJF(jobs []Job) []Job {
	// 对作业按照估计执行时间排序
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].EstimatedTime < jobs[j].EstimatedTime
	})
	return jobs
}

func scheduleHRRN(jobs []Job) []Job {
	// 对作业按照响应比排序
	sort.Slice(jobs, func(i, j int) bool {
		ratio1 := float64(jobs[i].EstimatedTime) / float64(jobs[i].ArrivalTime)
		ratio2 := float64(jobs[j].EstimatedTime) / float64(jobs[j].ArrivalTime)
		return ratio1 > ratio2
	})
	return jobs
}

func printJobSequence(jobs []Job) {
	fmt.Println("作业序列：")
	for _, job := range jobs {
		fmt.Printf("作业%d 进入内存的时间：%d\n", job.JobNumber, job.ArrivalTime)
	}
}
