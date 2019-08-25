package generator

import (
	"bytes"
	"encoding/json"
	"errors"
	"kit/config"
	"kit/fs"
	"kit/template"
	"strings"

	"github.com/ozgio/strutil"
	"github.com/spf13/afero"
)

// NewService generates a new service with the given name
func NewService(name string) error {
	appFs := fs.AppFs()

	b, err := afero.Exists(appFs, "kit.json")

	if err != nil {
		return err
	} else if !b {
		return errors.New("not in a kit project, you need to be in a project to run this command")
	}

	// we should remove the '_' because of this guide https://blog.golang.org/package-names
	folderName := strings.ReplaceAll(strutil.ToSnakeCase(name), "_", "")

	if err := fs.CreateFolder(folderName, appFs); err != nil {
		return err
	}

	data := map[string]string{
		"ProjectModule": folderName,
	}

	serviceFile, err := template.CompileGoFromPath("/assets/templates/service/service.go.gotmpl", data)
	if err != nil {
		return err
	}
	svcFs := afero.NewBasePathFs(appFs, folderName)
	err = fs.WriteFile("service.go", serviceFile, svcFs)
	if err != nil {
		return err
	}
	configData, err := fs.ReadFile("kit.json", appFs)
	if err != nil {
		return errors.New("could not read kit.json")
	}
	var kitConfig config.KitConfig
	err = json.NewDecoder(bytes.NewBufferString(configData)).Decode(&kitConfig)
	if err != nil {
		return err
	}
	kitConfig.Services = append(kitConfig.Services, name)
	newData, err := json.MarshalIndent(kitConfig, "", "\t")
	if err != nil {
		return err
	}
	return fs.WriteFile("kit.json", string(newData), appFs)
}
