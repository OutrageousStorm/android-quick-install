package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	device := flag.String("device", "", "Target device serial")
	parallel := flag.Int("parallel", 4, "Number of parallel installations")
	reinstall := flag.Bool("reinstall", false, "Reinstall if already present")
	flag.Parse()

	apks := flag.Args()
	if len(apks) == 0 {
		fmt.Println("Usage: adb-install [options] file1.apk file2.apk ...")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("📲 Quick Install — %d APKs, %d parallel\n\n", len(apks), *parallel)

	// Validate files exist
	for _, apk := range apks {
		if _, err := os.Stat(apk); err != nil {
			fmt.Printf("❌ %s: not found\n", apk)
			os.Exit(1)
		}
	}

	// Install with semaphore for concurrency control
	sem := make(chan struct{}, *parallel)
	var wg sync.WaitGroup
	results := make(chan string, len(apks))

	for _, apk := range apks {
		wg.Add(1)
		go func(apkPath string) {
			defer wg.Done()
			sem <- struct{}{}        // acquire
			defer func() { <-sem }() // release

			installAPK(apkPath, *device, *reinstall, results)
		}(apk)
	}

	wg.Wait()
	close(results)

	// Print results
	success := 0
	for msg := range results {
		fmt.Println(msg)
		if strings.Contains(msg, "✓") {
			success++
		}
	}

	fmt.Printf("\n✅ %d/%d installed\n", success, len(apks))
}

func installAPK(apkPath string, device string, reinstall bool, results chan<- string) {
	name := filepath.Base(apkPath)
	cmd := exec.Command("adb")
	
	if device != "" {
		cmd.Args = append(cmd.Args, "-s", device)
	}
	
	cmd.Args = append(cmd.Args, "install")
	if reinstall {
		cmd.Args = append(cmd.Args, "-r")
	}
	cmd.Args = append(cmd.Args, apkPath)

	output, err := cmd.CombinedOutput()
	
	if err == nil && strings.Contains(string(output), "Success") {
		results <- fmt.Sprintf("✓ %s", name)
	} else {
		results <- fmt.Sprintf("✗ %s", name)
	}
}
