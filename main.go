package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/gommon/color"
	"github.com/tidwall/gjson"
)

var arguments = struct {
	Output      string
	Concurrency int
	URL         string
	Type        string
	Arch        bool
	Best        bool
	Frenchy     bool
	Verbose     bool
}{}

var checkPre = color.Yellow("[") + color.Green("✓") + color.Yellow("]")
var tildPre = color.Yellow("[") + color.Green("~") + color.Yellow("]")
var crossPre = color.Yellow("[") + color.Red("✗") + color.Yellow("]")

func extractIDS() []string {
	var ids []string
	var value gjson.Result

	cmd := exec.Command("youtube-dl", "-j", "--flat-playlist", "https://www.youtube.com/"+arguments.Type+"/"+arguments.URL)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed listing channel's IDs %s\n", err)
	}

	gjson.ForEachLine(string(out), func(line gjson.Result) bool {
		value = gjson.Get(line.String(), "id")
		ids = append(ids, value.String())
		return true
	})

	return ids
}

func downloadVideo(ID string, nbIDs, i int, worker *sync.WaitGroup) {
	defer worker.Done()

	start := time.Now()
	url := "https://youtube.com/watch?v=" + ID
	outputDirectory := arguments.Output + "/" +
		arguments.Type + "_" +
		arguments.URL + "_" +
		start.Format("02-01-2006")

	// Create appropriate directory for saving the files
	os.MkdirAll(outputDirectory, os.ModePerm)

	if arguments.Best == true {
		cmd := exec.Command("youtube-dl",
			"-f (\"bestvideo[width>=1920]\"/bestvideo)+bestaudio/best",
			"-ciw",
			"--prefer-ffmpeg",
			"--merge-output-format=mkv",
			"--write-sub",
			"--all-subs",
			"--convert-subs=srt",
			"--add-metadata",
			"--write-description",
			"--write-annotations",
			"--write-all-thumbnails",
			"--write-info-json",
			url)
		cmd.Dir = outputDirectory
		out, err := cmd.CombinedOutput()
		if arguments.Verbose == true {
			fmt.Println(string(out))
		}
		if err != nil {
			log.Fatalf(crossPre+
				color.Yellow("[")+
				color.Red(ID)+
				color.Yellow("] ")+
				color.Red("Failed downloading video: %s\n"), err)
		}
	} else if arguments.Frenchy == true {
		cmd := exec.Command("youtube-dl",
			"--config-location config.txt",
			"-f (\"bestvideo[width>=1920]\"/bestvideo)+bestaudio/best",
			"--ignore-errors",
			"--no-continue",
			"--no-overwrites",
			"--no-post-overwrites",
			"--prefer-ffmpeg",
			"--merge-output-format=mkv",
			"--write-sub",
			"--all-subs",
			"--convert-subs=srt",
			"--add-metadata",
			"--write-description",
			"--write-annotations",
			"--write-all-thumbnails",
			"--write-info-json",
			"--download-archive=archive.txt",
			url)
		out, err := cmd.CombinedOutput()
		if arguments.Verbose == true {
			fmt.Println(string(out))
		}
		if err != nil {
			log.Fatalf(crossPre+
				color.Yellow("[")+
				color.Red(ID)+
				color.Yellow("] ")+
				color.Red("Failed downloading video: %s\n"), err)
		}
	} else {
		cmd := exec.Command("youtube-dl",
			"-f best",
			"-ciw",
			"--prefer-ffmpeg",
			"--merge-output-format=mkv",
			"--write-sub",
			"--all-subs",
			"--convert-subs=srt",
			"--add-metadata",
			"--write-description",
			"--write-annotations",
			"--write-all-thumbnails",
			"--write-info-json",
			url)
		cmd.Dir = outputDirectory
		out, err := cmd.CombinedOutput()
		if arguments.Verbose == true {
			fmt.Println(string(out))
		}
		if err != nil {
			log.Fatalf(crossPre+
				color.Yellow("[")+
				color.Red(ID)+
				color.Yellow("] ")+
				color.Red("Failed downloading video: %s\n"), err)
		}
	}

	// Pad ID
	index := strconv.Itoa(i)
	for len(index) < len(strconv.Itoa(nbIDs)) {
		index = "0" + index
	}

	fmt.Println(checkPre +
		color.Yellow("[") +
		color.Green(ID) +
		color.Yellow("]") +
		color.Yellow("[") +
		color.Green(time.Since(start)) +
		color.Yellow("]") +
		color.Yellow("[") +
		color.Green(index) +
		color.Green("/") +
		color.Green(nbIDs) +
		color.Yellow("]") +
		color.Green(" Downloaded!"))
}

func main() {
	var worker sync.WaitGroup
	var count int

	// Parse arguments and fill the arguments structure
	parseArgs(os.Args)

	// Extract channel IDs
	ids := extractIDS()

	// Download videos
	i := 1
	for _, id := range ids {
		worker.Add(1)
		count++
		go downloadVideo(id, len(ids), i, &worker)
		if count == arguments.Concurrency {
			worker.Wait()
			count = 0
		}
		i++
	}
	worker.Wait()
}
