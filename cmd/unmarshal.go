package cmd

import (
	"github.com/gkwa/bouncingbeaver/app"
	"github.com/spf13/cobra"
)

var inputFile string

var unmarshalCmd = &cobra.Command{
	Use:   "unmarshal",
	Short: "Unmarshal DynamoDB data example",
	Long:  "Demonstrates unmarshaling DynamoDB AttributeValue format to Go structs",
	RunE: func(cmd *cobra.Command, args []string) error {
		processor := app.NewProcessor(verbose)
		return processor.ProcessData(inputFile)
	},
}

func init() {
	unmarshalCmd.Flags().StringVarP(&inputFile, "file", "f", "internal/dynamodb/testdata/sample_input.json", "input file (use '-' for stdin)")
	rootCmd.AddCommand(unmarshalCmd)
}
