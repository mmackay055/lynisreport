package main

import (
	"fmt"
	"lynisreport/lynis"
	"os"
	flag "github.com/spf13/pflag"
)

// Options
var repOpt string // option for report location
var logOpt string // option for log location
var fmtTimestampOpt bool
var fmtJsonOpt bool // option to output data as json
var fmtYamlOpt bool // option to output data as yaml
var fmtNewLineOpt bool // option to append newline at end of output

const (
	ERR_REPORTFILE int = 2
	ERR_LOGFILE    int = 3
	ERR_PROCCESS   int = 4
        ERR_WRITELOG   int = 5
        ERR_INVALIDOPT int = 6
)

func init() {
	flag.StringVarP(&repOpt,
		"reportfile",
		"r",
		"/var/log/lynis-report.dat",
		"Specify where to find the Lynis report file")
	flag.StringVarP(&logOpt,
		"logfile",
		"l",
		"",
		"Specify where to log output of report. Default is to standard output")
        flag.BoolVarP(&fmtTimestampOpt,
                "timestamp",
                "t",
                false,
                "Prepend timestamp info before data output")
        flag.BoolVarP(&fmtJsonOpt,
                "json",
                "j",
                true,
                "Output data in json(default output)")
        flag.BoolVarP(&fmtYamlOpt,
                "yaml",
                "y",
                false,
                "Output data in yaml(not yet implemented)")
        flag.BoolVarP(&fmtNewLineOpt,
                "newline",
                "n",
                false,
                "Append new line character to end of output")
}

func main() {
	// Parse command line args
	flag.Parse()

        // set data formatters
        var formatter lynis.OutputFormatter
        if fmtYamlOpt {
                fmt.Fprintf(os.Stderr, "error: yaml output is not yet implemented\n")
                os.Exit(ERR_INVALIDOPT)
        } else {
                formatter = &lynis.FormatJSON{}
        }

        // add optional timestamp
        if fmtTimestampOpt {
                lynis.SetNext(formatter, &lynis.FormatTimestamp{})
        }

        // add new line to end of output
        if fmtNewLineOpt {
                lynis.SetNext(formatter,&lynis.FormatNewLine{})
        }

	// open report file
	var input *os.File
	if len(repOpt) < 1 {
		input = os.Stdin
	} else {
		var err error
		input, err = os.Open(repOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr,"error: %s\n",
                                err)
			os.Exit(ERR_REPORTFILE)
		}
	}

	// open log file
	var output *os.File
	if len(logOpt) < 1 {
		output = os.Stdout
	} else {
		var err error
		output, err = os.OpenFile(logOpt, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			fmt.Fprintf(os.Stderr,"error: failed to open log file %s\n",
				logOpt)
			os.Exit(ERR_LOGFILE)
		}
	}

	_, data, err := lynis.Process(input, formatter)
	if err != nil {
                //TODO log error message to log file
		fmt.Fprintf(os.Stderr,"error: failed to parse Lynis Report %s\n", err)
		os.Exit(ERR_PROCCESS)
	}

        if bytes, err := output.Write(data); err != nil {
                fmt.Fprintf(os.Stderr,"error: failed writting report to log file. Wrote %d bytes expected %d. %s\n", bytes, len(data), err.Error())
                os.Exit(ERR_WRITELOG)
        }
}
