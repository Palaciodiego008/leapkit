package generate

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// actionsFolder is the folder where the actions are stored
	actionsFolder = "internal"

	//go:embed action.go.tmpl
	actionTemplate string
)

// Action generates a new action
func Action(name string) error {
	path := strings.Split(name, string(filepath.Separator))
	actionPackage := "internal"
	fileName := path[len(path)-1] // file name is the last part of the path
	if len(path) > 1 {
		actionPackage = path[len(path)-2] // package name is the second to last part of the path
	}

	folder := strings.Join(path[:len(path)-1], string(filepath.Separator)) // folder is everything but the last part of the path
	actionName := cases.Title(language.English).String(filepath.Base(name))
	// Create the folder
	if actionPackage != "internal" {
		folder = folder + string(filepath.Separator) // add the separator if the package is not internal
		if err := os.MkdirAll(filepath.Join(actionsFolder, folder), 0755); err != nil {
			return fmt.Errorf("error creating folder: %w", err)
		}
	}

	// Create action.go
	file, err := os.Create(filepath.Join(actionsFolder, folder, fileName+".go"))
	if err != nil {
		return err
	}

	defer file.Close()
	template := template.Must(template.New("handler").Parse(actionTemplate))
	err = template.Execute(file, map[string]string{
		"Package":  actionPackage,
		"FileName": fileName,
		"Folder":   folder,

		"ActionName": actionName,
	})

	if err != nil {
		return err
	}

	// Create action.html
	_, err = os.Create(filepath.Join(actionsFolder, folder, fileName+".html"))
	if err != nil {
		return err
	}

	fmt.Println("Action files created successfully✅")

	return nil
}
