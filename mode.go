package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jskcnsl/eesc/ast"
	"github.com/jskcnsl/eesc/client"
	"github.com/jskcnsl/eesc/config"

	"github.com/spf13/cobra"
)

func prepareTargets() ([]string, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	query := []string{}
	if config.QueryExp != "" {
		query = append(query, config.QueryExp)
	}
	if config.ExpFile != "" {
		fmt.Printf("read expression from: %s\n", config.ExpFile)
		content, err := os.ReadFile(config.ExpFile)
		if err != nil {
			return nil, err
		}
		query = append(query, strings.Split(string(content), "\n")...)
	}

	result := make([]string, 0, len(query))
	for _, q := range query {
		tq := strings.TrimSpace(config.JoinExp + " " + q)
		if tq != "" {
			result = append(result, tq)
		}
	}

	return result, nil
}

func search(cmd *cobra.Command, args []string) error {
	targets, err := prepareTargets()
	if err != nil {
		return err
	}

	esc, err := client.NewEsClient(config.ServerAddress, config.IndexName)
	if err != nil {
		return err
	}

	for _, t := range targets {
		o.Writeln(strings.Repeat("#", 32))

		o.Writeln(t)
		a := ast.NewAst(t)
		if err := a.Resolve(); err != nil {
			o.Writeln(err.Error())
			continue
		}
		if config.Verbose {
			o.Writeln(a.Details())
		}

		o.Writeln(strings.Repeat("#", 32))
		o.Writeln(strings.Repeat("=", 32))
		res, err := esc.Search(cmd.Context(), a.Query(), config.Size, config.Offset)
		if err != nil {
			o.Writeln(err.Error())
		} else {
			o.Writeln(res)
		}
	}

	return nil
}

func count(cmd *cobra.Command, args []string) error {
	targets, err := prepareTargets()
	if err != nil {
		return err
	}

	esc, err := client.NewEsClient(config.ServerAddress, config.IndexName)
	if err != nil {
		return err
	}

	for _, t := range targets {
		o.Writeln(strings.Repeat("#", 32))

		o.Writeln(t)
		a := ast.NewAst(t)
		if err := a.Resolve(); err != nil {
			o.Writeln(err.Error())
			continue
		}
		if config.Verbose {
			o.Writeln(a.Details())
		}

		o.Writeln(strings.Repeat("#", 32))
		o.Writeln(strings.Repeat("=", 32))
		res, err := esc.Count(cmd.Context(), a.Query())
		if err != nil {
			o.Writeln(err.Error())
		} else {
			o.Writeln(res)
		}
	}

	return nil
}
