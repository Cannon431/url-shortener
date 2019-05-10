package dotenv_parser

import (
	"io/ioutil"
	"strings"
)

func Parse(fileName string) (map[string]string, error) {
	config := make(map[string]string)
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return config, err
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		segments := strings.Split(line, "=")
		key := strings.TrimSpace(segments[0])
		value := strings.TrimSpace(segments[1])

		config[key] = value
	}

	return config, nil
}
