package configgenerate

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
)

// WriteFilesFromTemplate writes the config files to disk using a specific template
func WriteFilesFromTemplate(in []ConfigFileValues, templatePath, outPath string) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	out := map[string]string{}
	for _, file := range in {
		var b bytes.Buffer
		err := tmpl.Execute(&b, file)
		if err != nil {
			fmt.Println(err)
			return err
		}

		out[file.Name] = b.String()
	}

	os.RemoveAll(outPath)
	os.MkdirAll(outPath, 0755)

	for name, content := range out {
		filePath := path.Join(outPath, name)
		ioutil.WriteFile(filePath, []byte(content), 0644)
	}

	return nil
}
