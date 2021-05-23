package main

import (
	"github.com/jskcnsl/eesc/config"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "eesc",
		Short: "Easy Elasticsearch Client",
	}
	rootCmd.Flags().StringVarP(&config.ServerAddress, "server", "s", "", "Address of database server")
	rootCmd.Flags().StringVarP(&config.IndexName, "idx", "x", "", "Index to work around")
	rootCmd.Flags().BoolVarP(&config.Verbose, "verbose", "v", false, "show details")
	rootCmd.Flags().StringVarP(&config.QueryExp, "query", "q", "", "single query expression")
	rootCmd.Flags().StringVarP(&config.JoinExp, "join", "j", "", "query which will join to each expression")
	rootCmd.Flags().StringVarP(&config.ExpFile, "file", "f", "", "collection of exporession")
	rootCmd.Flags().StringVarP(&config.OutputFile, "output", "o", "", "output file")
	_ = rootCmd.MarkFlagRequired("server")
	_ = rootCmd.MarkFlagRequired("idx")

	searchCmd := &cobra.Command{
		Use:     "search",
		Short:   "search elasticsearch with query",
		RunE:    search,
		PreRunE: initOutput,
		PostRun: closeOutput,
	}
	searchCmd.Flags().IntVarP(&config.Size, "size", "l", 10, "size of result")
	searchCmd.Flags().IntVarP(&config.Offset, "offset", "A", 0, "search offset")

	rootCmd.AddCommand(
		searchCmd,
		&cobra.Command{
			Use:     "count",
			Short:   "get count result with query",
			RunE:    count,
			PreRunE: initOutput,
			PostRun: closeOutput,
		},
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
