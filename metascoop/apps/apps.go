package apps

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type AppInfo struct {
	GitURL string `yaml:"git"`

	AuthorName string `yaml:"author"`
	repoAuthor string

	FriendlyName string `yaml:"name"`
	keyName      string
}

func (a AppInfo) Name() string {
	if a.FriendlyName != "" {
		return a.FriendlyName
	}
	return a.keyName
}

func (a AppInfo) Author() string {
	if a.AuthorName != "" {
		return a.AuthorName
	}
	return a.repoAuthor
}

// ParseAppFile returns the list of apps from the app file
func ParseAppFile(filepath string) (list []AppInfo, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer f.Close()

	var apps map[string]AppInfo

	err = yaml.NewDecoder(f).Decode(&apps)
	if err != nil {
		return
	}

	for k, a := range apps {
		a.keyName = k

		u, uerr := url.ParseRequestURI(a.GitURL)
		if uerr != nil {
			err = fmt.Errorf("problem with given git URL %q for app with key=%q, name=%q: %w", a.GitURL, k, a.Name(), uerr)
			return
		}

		split := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(split) == 0 {
			return
		}
		a.repoAuthor = split[0]

		list = append(list, a)
	}

	return
}