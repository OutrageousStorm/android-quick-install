package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func runAdb(args ...string) error {
	cmd := exec.Command("adb", args...)
	return cmd.Run()
}

func installAPK(path string) error {
	fmt.Printf("Installing %s...", filepath.Base(path))
	err := runAdb("install", "-r", "-g", path)
	if err != nil {
		fmt.Println(" ✗")
		return err
	}
	fmt.Println(" ✓")
	return nil
}

func main() {
	dirPtr := flag.String("dir", ".", "Directory containing APKs")
	parallelPtr := flag.Int("parallel", 2, "Number of parallel installs")
	flag.Parse()

	files, err := filepath.Glob(filepath.Join(*dirPtr, "*.apk"))
	if err != nil || len(files) == 0 {
		fmt.Println("❌ No APK files found in", *dirPtr)
		os.Exit(1)
	}

	fmt.Printf("⚡ Found %d APKs, installing with %d parallel workers

", len(files), *parallelPtr)

	semaphore := make(chan struct{}, *parallelPtr)
	var wg sync.WaitGroup
	success, failed := 0, 0
	var mu sync.Mutex

	for _, f := range files {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			semaphore <- struct{}{}        // acquire
			defer func() { <-semaphore }() // release

			if installAPK(path) == nil {
				mu.Lock()
				success++
				mu.Unlock()
			} else {
				mu.Lock()
				failed++
				mu.Unlock()
			}
		}(f)
	}

	wg.Wait()
	fmt.Printf("
✅ Success: %d  ❌ Failed: %d
", success, failed)
}
