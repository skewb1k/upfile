package cmd

// import (
// 	"path/filepath"
//
// 	indexFs "upfile/internal/index/fs"
// 	"upfile/internal/service"
//
// 	"github.com/spf13/cobra"
// )
//
// func show() *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "show <filename>",
// 		Short: "Show upstream version of file",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			s := service.New(indexFs.New(getBaseDir()))
//
// 			upstreamContent, err := s.Show(cmd.Context(), filepath.Base(args[0]))
// 			if err != nil {
// 				return err
// 			}
//
// 			cmd.Print(upstreamContent)
//
// 			return nil
// 		},
// 	}
// }
