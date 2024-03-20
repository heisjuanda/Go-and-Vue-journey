package parserhelpers

import (
	"bufio"
	"fmt"
	"mail-indexer/constants"
	"mail-indexer/models"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync"
	"time"
)

var (
	allEmails []models.Email
	mu        sync.Mutex
)

// given a filePath, reads the content line per line
func readFileLineByLine(filePath string, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	file, fileError := os.Open(filePath)
	if fileError != nil {
		fmt.Printf("Error opening file %s: %s", filePath, fileError)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	email := models.Email{}
	hasSubEmails := false
	for scanner.Scan() {
		line := scanner.Text()
		ParseLineMessage(line, &email, &hasSubEmails)
	}
	createSubEmails(&email)
	appendEmail(&email)
	if scannerError := scanner.Err(); scannerError != nil {
		fmt.Printf("Error reading file %s: %s", filePath, scannerError)
	}
}

// checks if folderPath is a folder or a file to read it
func ExploreFolder(folderPath string, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()

	files, fileError := os.ReadDir(folderPath)
	if fileError != nil {
		fmt.Printf("Error reading folder %s: %s", folderPath, fileError)
		return
	}
	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())

		if file.IsDir() {
			wg.Add(1)
			ExploreFolder(filePath, wg, semaphore)
		} else {
			wg.Add(1)
			go readFileLineByLine(filePath, wg, semaphore)
		}
	}
}

// according to a fileName it will create the go routines and max number of workers to start reading the files
func ProcessFiles(fileNames string) {
	var wg sync.WaitGroup
	const maxWorkers = constants.MAX_WORKERS

	start := time.Now()
	semaphore := make(chan struct{}, maxWorkers)
	wg.Add(1)
	go ExploreFolder(fileNames, &wg, semaphore)

	cpuProfile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("could not create CPU profile: ", err)
	}
	defer cpuProfile.Close()
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	memProfile, err := os.Create("mem_profile.prof")
	if err != nil {
		fmt.Println("Could not create memory profile: ", err)
	}
	defer memProfile.Close()
	defer pprof.WriteHeapProfile(memProfile)

	wg.Wait()
	close(semaphore)
	fmt.Println("start index: ", start)
	fmt.Println("end index: ", time.Since(start))
	fmt.Print("total emails: ", len(allEmails))
}
