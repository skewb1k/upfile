package cmd

import (
	"path/filepath"

	indexFs "upfile/internal/index/fs"
	"upfile/internal/service"
	userfileFs "upfile/internal/userfile/fs"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

func withService(f func(
	cmd *cobra.Command,
	args []string,
	s *service.Service,
) error,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return f(cmd, args, service.New(indexFs.New(filepath.Join(xdg.DataHome, Name)), userfileFs.New()))
	}
}
