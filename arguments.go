package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func parseArgs(args []string) {
	// Create new parser object
	parser := argparse.NewParser("YTChannelArchiver", "YouTube channel archiver, powered by YouTube-DL")

	// Create flags
	output := parser.String("o", "output", &argparse.Options{
		Required: false,
		Default:  "./Channels",
		Help:     "Output directory"})

	concurrency := parser.Int("j", "concurrency", &argparse.Options{
		Required: false,
		Default:  2,
		Help:     "Concurrent jobs for download"})

	url := parser.String("u", "url", &argparse.Options{
		Required: true,
		Help:     "Channel username or ID"})

	channelType := parser.String("t", "type", &argparse.Options{
		Required: true,
		Help:     "Channel type"})

	arch := parser.Flag("", "arch", &argparse.Options{
		Required: false,
		Default:  true,
		Help:     "-Archivist prefeered arguments"})

	best := parser.Flag("", "best", &argparse.Options{
		Required: false,
		Default:  false,
		Help:     "Best quality arguments"})

	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Default:  false,
		Help:     "Verbose output"})

	// Parse input
	err := parser.Parse(args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	// Fill arguments structure
	arguments.Output = *output
	arguments.Concurrency = *concurrency
	arguments.URL = *url
	arguments.Type = *channelType
	arguments.Arch = *arch
	arguments.Best = *best
	arguments.Verbose = *verbose
}
