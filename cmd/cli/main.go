package main

import (
	"log"
	"os"

	"github.com/sebomancien/bin2c/pkg/converter"
	"github.com/spf13/cobra"
)

const (
	DEFAULT_INPUT_FILE     = "input.bin"
	DEFAULT_OUTPUT_FILE    = "output.c"
	DEFAULT_ARRAY_NAME     = "myArray"
	DEFAULT_BYTES_PER_LINE = 32
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

	rootCmd.Flags().StringVarP(&input_file, "input", "i", DEFAULT_INPUT_FILE, "input file")
	rootCmd.Flags().StringVarP(&output_file, "output", "o", DEFAULT_OUTPUT_FILE, "output file")
	rootCmd.Flags().StringVarP(&array_name, "name", "n", DEFAULT_ARRAY_NAME, "name of the C array")
	rootCmd.Flags().Uint8VarP(&byte_per_line, "length", "l", DEFAULT_BYTES_PER_LINE, "number of bytes per line")
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

	writer, err := os.Create(output_file)
	if err != nil {
		log.Fatal("An error occured while opening the output file:", err)
	}
	defer writer.Close()

	config := converter.Config{
		ArrayName:   array_name,
		BytePerLine: byte_per_line,
	}

	converter.Convert(reader, writer, &config)
}