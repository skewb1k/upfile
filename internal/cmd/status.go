package cmd

import (
	"fmt"
	"path/filepath"

	"upfile/internal/service"
	storeFs "upfile/internal/store/fs"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
)

func status() *cobra.Command {
	return &cobra.Command{
		Use:   "status [dir]",
		Short: "Print status of files in dir (default: current dir)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}

			absDir, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("get abs path to dir: %w", err)
			}

			s := service.New(storeFs.New(getBaseDir()))

			res, err := s.Status(cmd.Context(), absDir)
			if err != nil {
				return err
			}

			for _, e := range res {
				var status string

				switch e.Status {
				case service.EntryStatusModified:
					status = yellow("Modified")
				case service.EntryStatusUpToDate:
					status = green("Up to date")
				case service.EntryStatusDeleted:
					status = red("Deleted")
				}

				cmd.Printf("%-40s %s\n", e.Fname, status)
			}

			return nil
		},
	}
}
