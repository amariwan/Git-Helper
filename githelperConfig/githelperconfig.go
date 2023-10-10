package githelperConfig

// Importing packages
import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"aland-mariwan.de/githelper/helper"
)

type GithelperConfig []GithelperConfigElement

func unmarshalGithelperConfig(data []byte) (GithelperConfig, error) {
	var r GithelperConfig
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GithelperConfig) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *GithelperConfig) MarshalIndent() ([]byte, error) {
	return json.MarshalIndent(r, "", "\t")
}

type GithelperConfigElement struct {
	Name         string              `json:"name"`
	Folder       string              `json:"folder"`
	Repository   string              `json:"repository"`
	IgnorePrefix string              `json:"ignorePrefix"`
	Versions     []map[string]string `json:"versions"`
}

func ReadGithelperJson(startPath string) (gitHelperConfig GithelperConfig, err error) {
	path := helper.FindUpwards(startPath, "githelper.json")
	// Open our jsonFile
	jsonFile, err := os.Open(filepath.Join(path, "githelper.json"))
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	gitHelperConfig, err = unmarshalGithelperConfig(byteValue)
	return
}

func VersionReferenced(versions []map[string]string, actualBranch string) string {
	for _, version := range versions {
		for k2, v2 := range version {
			if helper.Match(k2, actualBranch) {
				if v2 == "=" {
					return actualBranch
				}
				return v2
			}
		}
	}

	return actualBranch
}
