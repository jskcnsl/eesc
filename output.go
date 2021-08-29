package main

import (
	"fmt"
	"os"

	"github.com/jskcnsl/eesc/config"
	"github.com/spf13/cobra"
)

type output interface {
	Write(string)
	Writeln(string)
	Close()
}

var (
	_ output = &StdOutput{}
	o output = nil
)

func initOutput(*cobra.Command, []string) error {
	var err error
	if config.OutputFile != "" {
		fmt.Printf("store output to file: %s\n", config.OutputFile)
		o, err = NewFileOutput(config.OutputFile)
		if err != nil {
			return err
		}
	} else {
		o = NewStdOutput()
	}

	return nil
}

func closeOutput(*cobra.Command, []string) {
	if o != nil {
		o.Close()
	}
}

func NewStdOutput() output {
	return &StdOutput{}
}

type StdOutput struct {
}

func (o *StdOutput) Write(l string) {
	fmt.Print(l)
}

func (o *StdOutput) Writeln(l string) {
	fmt.Println(l)
}

func (o *StdOutput) Close() {
}

type FileOutput struct {
	f *os.File
}

func NewFileOutput(fn string) (output, error) {
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &FileOutput{
		f: f,
	}, err
}

func (o *FileOutput) Write(l string) {
	_, _ = o.f.WriteString(l)
}

func (o *FileOutput) Writeln(l string) {
	_, _ = o.f.WriteString(l + "\n")
}

func (o *FileOutput) Close() {
	o.f.Close()
}
