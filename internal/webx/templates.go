package webx

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

func GetTemplates() (*template.Template, error) {
	funcMap := template.FuncMap{}
	t := template.New("").Funcs(funcMap)
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".gohtml") {
			_, err = t.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return t, err
}
