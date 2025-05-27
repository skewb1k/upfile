package service

import (
	"context"
)

func (s Service) Diff(ctx context.Context, path string) (string, error) {
	// dmp := diffmatchpatch.New()
	//
	// diffs := dmp.DiffMain(text1, text2, false)
	//
	// return dmp.DiffPrettyText(diffs), nil

	// fname := filepath.Base(path)
	// entryDir := filepath.Dir(path)
	//
	// if _, err := s.entries.GetEntryByDir(ctx, fname, entryDir); !errors.Is(err, entries.ErrNotFound) {
	// }
	//
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	return fmt.Errorf("failed to read file: %w", err)
	// }
	//
	// hash := computeHash(data)
	return "", nil
}
