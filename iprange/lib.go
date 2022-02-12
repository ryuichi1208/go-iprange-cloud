package iprange

import (
	"fmt"
	"os"
	"path/filepath"
)

type OutputConfig struct {
	DirFullPath  string
	FileBaseName string
	FileFormat   string
}

func (c *OutputConfig) AudioFormat() string {
	return c.FileFormat
}

func (c *OutputConfig) AbsPath() string {
	name := fmt.Sprintf("%s.%s", c.FileBaseName, c.FileFormat)
	return filepath.Join(c.DirFullPath, name)
}

func (c *OutputConfig) IsExist() bool {
	_, err := os.Stat(c.AbsPath())
	return err == nil
}
