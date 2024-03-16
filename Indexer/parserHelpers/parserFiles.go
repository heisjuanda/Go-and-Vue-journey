package parserhelpers

import (
	"bufio"
	"fmt"
	"log"
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
func readFileLineByLine(filePath string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	sem <- struct{}{}        // add worker
	defer func() { <-sem }() // remove worker

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %s", filePath, err)
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
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file %s: %s", filePath, err)
	}
}

// checks if folderPath is a folder or a file to read it
func ExploreFolder(folderPath string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()

	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Printf("Error reading folder %s: %s", folderPath, err)
		return
	}
	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())

		if file.IsDir() {
			wg.Add(1)
			ExploreFolder(filePath, wg, sem)
		} else {
			wg.Add(1)
			go readFileLineByLine(filePath, wg, sem)
		}
	}
}

// according to a fileName it will create the go routines and max number of workers to start reading the files
func ProcessFiles(fileNames string) {
	var wg sync.WaitGroup
	const maxWorkers = constants.MAX_WORKERS

	start := time.Now()
	// semaphore & maxLength of channel
	sem := make(chan struct{}, maxWorkers)
	wg.Add(1)
	go ExploreFolder(fileNames, &wg, sem)

	cpuProfile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer cpuProfile.Close()
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	memProfile, err := os.Create("mem_profile.prof")
	if err != nil {
		log.Fatal("Could not create memory profile: ", err)
	}
	defer memProfile.Close()
	defer pprof.WriteHeapProfile(memProfile)

	wg.Wait()
	close(sem)
	fmt.Println("start index: ", start)
	fmt.Println("end index: ", time.Since(start))
	fmt.Print("total emails: ", len(allEmails))
}
