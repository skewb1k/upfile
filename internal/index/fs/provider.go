package indexFs

type IndexFsProvider struct {
	BaseDir string
}

func NewProvider(baseDir string) *IndexFsProvider {
	return &IndexFsProvider{
		BaseDir: baseDir,
	}
}
