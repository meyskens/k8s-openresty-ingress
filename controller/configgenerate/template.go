package configgenerate

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
)

var hashMap = map[string][]byte{}

// WriteFilesFromTemplate writes the config files to disk using a specific template
func WriteFilesFromTemplate(in []ConfigFileValues, templatePath, outPath string) (bool, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	newHashMap := map[string][]byte{}

	out := map[string]string{}
	for _, file := range in {
		var b bytes.Buffer
		err := tmpl.Execute(&b, file)
		if err != nil {
			fmt.Println(err)
			return false, err
		}

		hash := sha256.New()
		newHashMap[file.Name] = hash.Sum(b.Bytes())
		out[file.Name] = b.String()
	}

	os.RemoveAll(outPath)
	os.MkdirAll(outPath, 0755)

	for name, content := range out {
		filePath := path.Join(outPath, name)
		ioutil.WriteFile(filePath, []byte(content), 0644)
	}

	changed := hasHashMapChanged(newHashMap, hashMap)
	hashMap = newHashMap

	return changed, nil
}

func hasHashMapChanged(new, old map[string][]byte) bool {
	if len(new) != len(old) {
		return true
	}

	for key, value := range new {
		oldValue, exists := old[key]
		if !exists {
			return true
		}

		if string(value) != string(oldValue) {
			return true
		}
	}

	return false
}
