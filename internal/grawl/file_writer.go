package grawl

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

type FileWriter struct {
	sync.RWMutex
	filePath        string
	fileInitialized bool
}

func NewFileWriter(filePath string) *FileWriter {
	return &FileWriter{
		filePath:        filePath,
		fileInitialized: false,
	}
}

func (f *FileWriter) InitFile() {
	if f.fileInitialized {
		return
	}

	fmt.Printf("Saving file \"%s\".\n", f.filePath)

	file, err := os.Create(f.filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	headers := f.getCsvHeader()
	f.write(headers, file)
	f.fileInitialized = true
}

func (f *FileWriter) WriteResultLine(r *Result) {
	f.RLock()
	defer f.RUnlock()

	if !f.fileInitialized {
		panic("csv not initialized yet")
	}

	file, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	line := f.getCsvRow(r)
	f.write(line, file)
}

func (f *FileWriter) getCsvRow(r *Result) []string {

	errorText := ""
	if r.error != nil {
		errorText = fmt.Sprintf("%v", r.error)
	}

	return []string{
		r.responseAt.Format(DateFormat),
		strconv.Itoa(r.statusCode),
		r.status,
		r.url,

		r.foundOnUrl,
		r.contentType,
		strconv.FormatInt(r.GetDuration().Milliseconds(), 10),
		strconv.Itoa(r.depth),
		r.urlRedirectedFrom,

		r.urlHost,
		r.urlPath,
		r.urlParmeters,
		r.urlFragment,

		errorText,
	}
}

func (f *FileWriter) getCsvHeader() []string {
	return []string{
		"Response time",
		"Status code",
		"Status",
		"URL",

		"Found on URL",
		"Content type",
		"Duration (ms)",
		"Depth",
		"Redirected from",

		"Host",
		"Path",
		"Parameters",
		"Fragment",

		"Info / error",
	}
}

func (f *FileWriter) write(text []string, file io.Writer) {
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	if err := writer.Write(text); err != nil {
		panic(err)
	}
}
