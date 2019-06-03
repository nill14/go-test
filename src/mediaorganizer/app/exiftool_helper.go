package app

import (
	"bufio"
	"bytes"
	"fmt"

	//	"github.com/barsanuphe/goexiftool"
	"os"
	"strings"

	exiftool "github.com/mostlygeek/go-exiftool"
)

// ExiftoolHelper keeps open connection to exiftool
type ExiftoolHelper struct {
	et *exiftool.Stayopen
}

//MediaFile keeps track of a single file EXIF data
type MediaFile struct {
	Filename string
	Info     map[string]string
}

func newMediaFile(filename string, b []byte) *MediaFile {
	mf := &MediaFile{Filename: filename, Info: make(map[string]string)}
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		res := strings.SplitN(scanner.Text(), ":", 2)
		if len(res) > 1 {
			key := strings.TrimSpace(res[0])
			value := strings.TrimSpace(res[1])
			mf.Info[key] = value
		}
	}
	return mf
}

// String displays all metadata
func (m *MediaFile) String() string {
	txt := m.Filename + ":\n"
	for k, v := range m.Info {
		txt += "\t" + k + " = " + v + "\n"
	}
	return txt
}

//OpenExiftool opens connection to exiftool
func OpenExiftool() *ExiftoolHelper {
	et, err := exiftool.NewStayOpen("exiftool")
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
	return &ExiftoolHelper{et}
}

//ReadMediaFile reads EXIF for a file
func (h *ExiftoolHelper) ReadMediaFile(filename string) (*MediaFile, error) {
	b, err := h.et.Extract(filename)
	mf := newMediaFile(filename, b)
	//		mf, err := goexiftool.NewMediaFile(file)
	return mf, err
}

// Stop closes connection to exiftool
func (h *ExiftoolHelper) Stop() {
	h.et.Stop()
}
