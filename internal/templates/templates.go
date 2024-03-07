package templates

import (
	"embed"
	"html/template"
	"io/fs"
	"strings"
)

func GetTemplates(filesystem embed.FS) (*template.Template, error) {
	funcMap := template.FuncMap{}
	t := template.New("").Funcs(funcMap)
	err := fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".gohtml") {
			_, err = t.ParseFS(filesystem, path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return t, err
}
