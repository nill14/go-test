package cmd

import (
	"fmt"
	"os"
	"strings"

	"../app"

	"github.com/spf13/cobra"
)

var renameArgs app.RenameArgs

var rootCmd = &cobra.Command{
	Use:   "mediaorganizer",
	Short: "Mediaorganizer is a tool to rename files based on EXIF data",
	Long:  `Mediaorganizer renames files by their EXIF title`,
	Run: func(cmd *cobra.Command, args []string) {
		renameArgs.SourceDir = strings.TrimSpace(renameArgs.SourceDir)
		renameArgs.TargerDir = strings.TrimSpace(renameArgs.TargetDir)
		if renameArgs.SourceDir == renameArgs.TargerDir {
			renameArgs.Rename = true
		}
		app.Rename(renameArgs)
	},
}

func initArgs() {
	renameArgs = app.RenameArgs{}
	flag := rootCmd.PersistentFlags()
	flag.BoolVarP(&renameArgs.Help, "help", "h", false, "Print help message.")
	flag.CountVarP(&renameArgs.Verbosity, "verbose", "v", "Print some extra information. (-v -vv -vvv) is supported.")
	flag.CountVarP(&renameArgs.Quietness, "quiet", "q", "Surpress warnings")
	flag.BoolVarP(&renameArgs.DryRun, "dry-run", "n", false, "Print the information but do not rename anything.")
	flag.BoolVarP(&renameArgs.Interactive, "interactive", "i", false, "Confirm each rename operation.")
	flag.StringVar(&renameArgs.Tag, "tag", "", "An alternative EXIF tag instead of title (e.g. QuickTime:Title, ID3:Title)")
	flag.StringSliceVar(&renameArgs.Extensions, "extensions", []string{".mp4", ".mp3"},
		"File extensions as a filter (defaults to .mp4,.mp3")

	directory, _ := os.Getwd()
	flag.StringVarP(&renameArgs.SourceDir, "directory", "d", directory, "Directory with media or '.' by default.")
	flag.StringVarP(&renameArgs.TargetDir, "target", "o", directory, "The output directory")
	flag.BoolVarP(&renameArgs.Rename, "rename", "r", false, `Use rename instead of copy 
	(when the source directory is the same as target rename is implied)`)

	// flag.Parse()
	// SourceDir = flag.Arg(0)

}

// Execute runs command line parser
func Execute() {
	initArgs()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
