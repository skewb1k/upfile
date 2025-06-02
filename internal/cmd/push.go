package cmd

// import (
// 	"errors"
//
// 	indexFs "upfile/internal/index/fs"
// 	"upfile/internal/service"
//
// 	"github.com/spf13/cobra"
// )
//
// func push() *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "push <path>",
// 		Short: "Push file to the upstream",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			s := service.New(indexFs.New(getBaseDir()))
//
// 			err := s.Push(cmd.Context(), args[0])
// 			if err != nil {
// 				if errors.Is(err, service.ErrUpToDate) {
// 					cmd.Println("Already up to date.")
// 					return nil
// 				}
//
// 				return err
// 			}
//
// 			return nil
// 		},
// 	}
// }
