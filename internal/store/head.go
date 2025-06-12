package store

import (
	"bufio"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/skewb1k/upfile/pkg/sha256"
)

func (s Store) CheckHead(ctx context.Context, filename string) (bool, error) {
	headMap, err := s.readAll()
	if err != nil {
		return false, err
	}

	_, ok := headMap[filename]
	return ok, nil
}

func (s Store) GetHead(ctx context.Context, filename string) (sha256.SHA256, error) {
	var hash sha256.SHA256
	headMap, err := s.readAll()
	if err != nil {
		return hash, err
	}

	h, ok := headMap[filename]
	if !ok {
		return hash, ErrNotFound
	}

	decoded, err := hex.DecodeString(h)
	if err != nil || len(decoded) != 32 {
		return hash, err
	}

	copy(hash[:], decoded)
	return hash, nil
}

func (s Store) SetHead(ctx context.Context, filename string, commitHash sha256.SHA256) error {
	headMap, err := s.readAll()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	headMap[filename] = hex.EncodeToString(commitHash[:])

	if err := os.MkdirAll(s.BaseDir, 0o700); err != nil {
		return err
	}

	f, err := os.Create(s.getHeadsPath())
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for k, v := range headMap {
		if strings.Contains(k, " ") {
			return errors.New("invalid filename with space")
		}
		fmt.Fprintf(w, "%s %s\n", k, v)
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

func (s Store) readAll() (map[string]string, error) {
	headMap := make(map[string]string)
	f, err := os.Open(s.getHeadsPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return headMap, nil
		}

		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		headMap[parts[0]] = parts[1]
	}

	return headMap, scanner.Err()
}

func (s Store) GetFilenames(ctx context.Context) ([]string, error) {
	headMap, err := s.readAll()
	if err != nil {
		return nil, err
	}

	names := make([]string, len(headMap))

	var i int

	for name := range headMap {
		names[i] = name
		i++
	}

	return names, nil
}
