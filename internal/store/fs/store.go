package storeFs

type Store struct {
	BaseDir string
}

func New(baseDir string) *Store {
	return &Store{
		BaseDir: baseDir,
	}
}
