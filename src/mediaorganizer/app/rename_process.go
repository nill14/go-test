package app

import (
	"fmt"
	"os"
	"strings"

	log "../logbook"

	"path/filepath"
)

//RenameArgs The arguments for rename command
type RenameArgs struct {
	// Help indicates to display the help message and exit
	Help bool

	//SourceDir where to search for the media files
	SourceDir string

	//TargerDir where to place the renamed file
	TargetDir string

	//TargerDir where to place the renamed files
	TargerDir string

	//Verbosity how verbose the output shall be
	Verbosity int

	//Quietness reduces output by warnings, etc
	Quietness int

	//DryRun Only log the rename information but do not touch anything
	DryRun bool

	//Interactive confirm each operation
	Interactive bool

	//Rename indicates rename, otherwise copy (or link)
	Rename bool

	// Extensions List of alternative extensions. Empty string indicates the default choice (e.g. [.mp4,mp3])
	Extensions []string

	//Tag An alternative EXIF tag to use as a title
	Tag string
}

var args RenameArgs

const chanSize = 500

// Rename rename media based on EXIF
func Rename(renameArgs RenameArgs) {
	args = renameArgs
	files := make(chan string, chanSize)
	exifs := make(chan *MediaFile, chanSize)
	renames := make(chan *renameObject, chanSize)

	logbook := log.NewLogBook(args.Verbosity, args.Quietness)

	go scanFiles(files, logbook)
	go getExifs(files, exifs, logbook)
	go collectData(exifs, renames, logbook)

	renameFiles(renames, logbook)

}

func scanFiles(files chan string, logbook *log.LogBook) {
	logbook.Log(log.INFO, fmt.Sprintf("Scanning root path %s", args.SourceDir))
	counter := 0

	err := filepath.Walk(args.SourceDir, visit(files, args.Extensions, &counter, logbook))
	if err != nil {
		logbook.Log(log.FATAL, fmt.Sprintf("Filepath walk failed with dir: %s", args.SourceDir))
	}

	logbook.Log(log.INFO, fmt.Sprintf("Scan finished with %d files found", counter))
	close(files)
}

func visit(files chan string, extensions []string, counter *int, logbook *log.LogBook) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logbook.Log(log.ERROR, fmt.Sprintf("File walk failed for file: %s - %s", path, err))
			return err
		}

		if info.IsDir() {
			logbook.Log(log.TRACE, fmt.Sprintf("Walking directory %s", path))
			return nil
		}

		if strings.HasPrefix(args.TargerDir, args.SourceDir) && strings.HasPrefix(path, args.TargerDir) {
			//when TargetDir is a subdirectory of a SourceDir
			//skip output from output (e.g. SourceDir/TargetDir/TargetDir) on subsequent run
			logbook.Log(log.DEBUG, fmt.Sprintf("Skipping file in the output directory - %s", path))
			return nil
		}

		for _, ext := range extensions {
			if ext == filepath.Ext(path) {
				logbook.Log(log.TRACE, fmt.Sprintf("Found file %s", path))
				(*counter)++
				files <- path
				return nil
			}
		}

		logbook.Log(log.TRACE, fmt.Sprintf("Skipping: %s - extension '%s' looking for: %s", path, filepath.Ext(path), extensions))
		return nil
	}
}

func getExifs(files chan string, exifs chan *MediaFile, logbook *log.LogBook) {
	counter := 0

	et := OpenExiftool()
	defer et.Stop()

	for file := range files {

		mf, err := et.ReadMediaFile(file)

		if err != nil {
			logbook.Log(log.ERROR, fmt.Sprintf("Calling exiftool failed for the file %s", file))
			break
		}
		logbook.Log(log.DEBUG, fmt.Sprintf("EXIF: %s - %d attributes", file, len(mf.Info)))
		counter++
		exifs <- mf
	}
	logbook.Log(log.INFO, fmt.Sprintf("Exiftool finished with %d exifs found", counter))
	close(exifs)
}

//	QuickTime:Title for mp4 or ID3:Title for mp3
func getTitle(mf *MediaFile, logbook *log.LogBook) string {
	if args.Tag != "" {
		title, ok := mf.Info[args.Tag]
		if ok {
			logbook.Log(log.TRACE, fmt.Sprintf("File: %s, Title: %s", mf.Filename, title))
			return title
		}
		return ""
	}

	title, ok := mf.Info["Title"]
	if ok {
		logbook.Log(log.TRACE, fmt.Sprintf("File: %s, Title: %s", mf.Filename, title))
		return title
	}

	title, ok = mf.Info["QuickTime:Title"]
	if ok {
		logbook.Log(log.TRACE, fmt.Sprintf("File: %s, QuickTime:Title: %s", mf.Filename, title))
		return title
	}

	title, ok = mf.Info["ID3:Title"]
	if ok {
		logbook.Log(log.TRACE, fmt.Sprintf("File: %s, ID3:Title: %s", mf.Filename, title))
		return title
	}

	//logged one level up
	return ""
}

var replacer = strings.NewReplacer(
	"?", "_",
	":", "_",
	"&", "_",
	"!", "_",
	"’", "_",
	"\"", "_",
	";", "_",
	"\"", "",
	"/", "-",
	"“", "",
	"”", "",
	"—", "-",
	".", "",
	"\u200b", "")

func escapeChars(str string) string {
	//replace characters not friendly to some file systems (e.g. FAT)
	return replacer.Replace(str)

}

type renameObject struct {
	oldPath string
	newPath string
}

func collectData(exifs chan *MediaFile, renames chan *renameObject, logbook *log.LogBook) {
	counter := 0
	for exif := range exifs {
		title := getTitle(exif, logbook)
		sourceFile := exif.Filename

		if title == "" {
			logbook.Log(log.WARN, fmt.Sprintf("Skipping undefined Title attribute - %s (%d attributes)", exif.Filename, len(exif.Info)))
			logbook.Log(log.DEBUG, fmt.Sprint(exif.String()))
			continue
		}

		//		name := filepath.Base(sourceFile)
		relativeDir := strings.TrimPrefix(filepath.Dir(sourceFile), args.SourceDir)
		ext := filepath.Ext(sourceFile)
		escaped := escapeChars(title)
		newPath := filepath.Join(args.TargerDir, relativeDir, escaped+ext)

		if newPath == sourceFile {
			logbook.Log(log.DEBUG, fmt.Sprintf("Skipping already renamed %s", sourceFile))
			continue
		}

		if escaped == "" {
			logbook.Log(log.DEBUG, fmt.Sprintf("Skipping empty new name: %s", sourceFile))
			continue
		}

		counter++
		renames <- &renameObject{oldPath: sourceFile, newPath: newPath}
	}
	logbook.Log(log.INFO, fmt.Sprintf("Collection finished with %d renames found", counter))
	close(renames)
}

func renameFiles(renames chan *renameObject, logbook *log.LogBook) {
	counter := 0
	for rename := range renames {

		logbook.Log(log.INFO, fmt.Sprintf("%60s -> %-60s", rename.oldPath, rename.newPath))

		counter++

		ok := !args.Interactive || confirm(fmt.Sprintf("Rename '%s' to '%s'? [Yn]:", rename.oldPath, rename.newPath))
		if !args.DryRun && ok {
			copyFile(rename, logbook)
		}

	}
	logbook.Log(log.INFO, fmt.Sprintf("Rename finished with %d files", counter))
}

func confirm(message string) bool {
	var answer string
	for {
		fmt.Print(message)
		fmt.Scan(&answer)
		answer = strings.ToLower(answer)
		if "y" == answer || "yes" == answer {
			return true
		} else if "n" == answer || "no" == answer {
			return false
		}
	}
}

func copyFile(rename *renameObject, logbook *log.LogBook) {
	var err error

	newDir := filepath.Dir(rename.newPath)
	if _, err = os.Stat(newDir); os.IsNotExist(err) {
		//create non existing directories if necessary
		err = os.MkdirAll(newDir, 0664|0111)
		if err != nil {
			logbook.Log(log.ERROR, fmt.Sprintf("Mkdir %s failed", newDir))
		}
	}

	if !args.Rename {
		if _, err = os.Stat(rename.newPath); err == nil {
			//remove target file if already exists
			err = os.Remove(rename.newPath)
			if err != nil {
				logbook.Log(log.ERROR, fmt.Sprintf("Remove %s failed - %s", rename.newPath, err.Error()))
			}
		}

		//keep it simple and fast for now
		err = os.Link(rename.oldPath, rename.newPath)
		if err != nil {
			logbook.Log(log.ERROR, fmt.Sprintf("Link %s failed - %s", rename.newPath, err.Error()))
		}

	} else {
		err = os.Rename(rename.oldPath, rename.newPath)
		if err != nil {
			logbook.Log(log.ERROR, fmt.Sprintf("Rename %s failed - %s", rename.newPath, err.Error()))
		}
	}
}
