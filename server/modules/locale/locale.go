package locale

import (
	"fmt"
	"github.com/Unknwon/com"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"server/modules/setting"
)

var (
	directories = make(directorySet)
)

func LocalFiles() map[string][]byte {
	localeNames, err := Dir("locale")

	localFiles := make(map[string][]byte)

	for _, name := range localeNames {
		localFiles[name], err = Locale(name)

		if err != nil {
			logrus.Fatalf("Failed to load %s locale file. %v", name, err)
		}
	}

	return localFiles
}

// Dir returns all files from static or custom directory.
func Dir(name string) ([]string, error) {
	if directories.Filled(name) {
		return directories.Get(name), nil
	}

	var (
		result []string
	)

	confDir := path.Join(setting.AppConfPath, name)

	if com.IsDir(confDir) {
		files, err := com.StatDir(confDir, true)

		if err != nil {
			return []string{}, fmt.Errorf("Failed to read custom directory. %v", err)
		}

		result = append(result, files...)
	}

	return directories.AddAndGet(name, result), nil
}

// Locale reads the content of a specific locale from static or custom path.
func Locale(name string) ([]byte, error) {
	return fileFromDir(path.Join("locale", name))
}

// fileFromDir is a helper to read files from static or custom path.
func fileFromDir(name string) ([]byte, error) {
	confPath := path.Join(setting.AppConfPath, name)

	if com.IsFile(confPath) {
		return ioutil.ReadFile(confPath)
	}

	return []byte{}, fmt.Errorf("Asset file does not exist: %s", name)
}
