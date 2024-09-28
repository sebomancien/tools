package main

import (
	"log"
	"os"

	"github.com/sebomancien/bin2c/pkg/converter"
	"github.com/spf13/cobra"
)

var (
	input_file    string
	output_file   string
	array_name    string
	byte_per_line uint8
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "bin2c",
		Short: "bin2c is a simple CLI tool to convert binary files into C array objects",
		Long:  "bin2c is a simple CLI tool to convert binary files into C array objects",
		Run:   command,
	}

	rootCmd.Flags().StringVarP(&input_file, "input", "i", "", "input file")
	rootCmd.Flags().StringVarP(&output_file, "output", "o", "", "output file")
	rootCmd.Flags().StringVarP(&array_name, "name", "n", converter.DefaultArrayName, "name of the C array")
	rootCmd.Flags().Uint8VarP(&byte_per_line, "length", "l", converter.DefaultBytesPerLine, "number of bytes per line")
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

	config := converter.Config{
		ArrayName:   array_name,
		BytePerLine: byte_per_line,
	}

	converter.Convert(reader, writer, &config)
}
