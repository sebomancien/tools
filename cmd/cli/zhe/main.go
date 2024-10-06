package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sebomancien/tools/internal/zhe"
	"github.com/spf13/cobra"
)

var (
	input_file  string
	output_file string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "zhe",
		Short: "Zhe Program is a brute force search algorithm to find the best combinations of variables that fit with contraints",
		Long:  "Zhe Program is a brute force search algorithm to find the best combinations of variables that fit with contraints",
		Run:   command,
	}

	rootCmd.Flags().StringVarP(&input_file, "input", "i", "", "input file (.yaml)")
	rootCmd.Flags().StringVarP(&output_file, "output", "o", "", "output file (Default to console)")
	rootCmd.MarkFlagRequired("input")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func command(cmd *cobra.Command, args []string) {
	reader, err := os.Open(input_file)
	if err != nil {
		log.Fatal("An error occured while opening the input file:", err)
	}
	defer reader.Close()

	config, err := zhe.ReadYAML(reader)
	if err != nil {
		log.Fatal(err)
	}

	tui := zhe.NewTui()
	solver := zhe.NewSolver(config)

	go func() {
		//start := time.Now()
		result, err := solver.Solve(10)
		//elapsed := time.Since(start)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("Solved in %f\n", elapsed.Seconds())

		progress := solver.GetProgress()
		tui.Send(progress)

		var writer *os.File
		if output_file == "" {
			writer = os.Stdout
		} else {
			writer, err = os.Create(output_file)
			if err != nil {
				log.Fatal("An error occured while opening the output file:", err)
			}
			defer writer.Close()
		}

		err = result.WriteYAML(writer)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Start a routine to send progress updates
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range ticker.C {
			progress := solver.GetProgress()
			tui.Send(progress)
		}
	}()

	// Run the tea program
	_, err = tui.Run()
	if err != nil {
		fmt.Printf("Error starting program: %v\n", err)
		os.Exit(1)
	}
}
