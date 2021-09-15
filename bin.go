package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type CompressStatus int
type PDFSetting string

const (
	Ready CompressStatus = iota
	Start
	Compressing
	Success
	Error
)

const (
	Prepress PDFSetting = "prepress"
	Ebook    PDFSetting = "ebook"
	Screen   PDFSetting = "screen"
)

type CommandWriter struct {
	totalPage  int
	handlePage int
	progress   CompressProgress
	state      CompressStatus
}

func NewCommandWriter(progress CompressProgress) *CommandWriter {
	return &CommandWriter{
		totalPage:  0,
		handlePage: 0,
		progress:   progress,
	}
}

func (w *CommandWriter) Write(p []byte) (n int, err error) {

	output := string(p)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Processing pages") {
			w.totalPage, _ = strconv.Atoi(line[strings.Index(line, "through ")+8 : len(line)-1])
			w.progress(float32(1)/float32(w.totalPage), Start, "")
		} else if strings.HasPrefix(line, "Page ") {
			w.handlePage, _ = strconv.Atoi(line[5:])
			w.progress(float32(w.handlePage)/float32(w.totalPage), Compressing, "")
		} else if strings.Contains(line, "error") && w.state != Error {
			reason := line[strings.Index(line, "error")+6:]
			w.progress(100, Error, reason)
			w.state = Error
		}
	}

	return len(p), nil
}

type CompressProgress func(progress float32, status CompressStatus, reason string)

func CompressPdf(inputFile string, outputFile string, setting PDFSetting, progress CompressProgress) {
	cmd := exec.Command("/usr/bin/gs", "-sDEVICE=pdfwrite", "-dCompatibilityLevel=1.4", fmt.Sprintf("-dPDFSETTINGS=/%s", string(setting)),
		"-dNOPAUSE", "-dBATCH", fmt.Sprintf("-sOutputFile=%s", outputFile), inputFile)

	commandWriter := NewCommandWriter(progress)
	cmd.Stdout = commandWriter
	cmd.Stderr = commandWriter

	err := cmd.Run()
	if err == nil {
		progress(1, Success, "")
	}
}
