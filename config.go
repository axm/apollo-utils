package apollo

import (
	"encoding/json"
	"github.com/axm/apollo-utils"
	"io/ioutil"
	"os"
	"fmt"
	"path"
)

type Config map[string]*json.RawMessage

func ReadFileFromRelativePath(relativePath string) ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get cwd: %w", err)
	}
	path := path.Join(cwd, relativePath)

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	return contents, nil
}

func ReadDatabaseConnection(config *Config) (*apollo.DatabaseConnection, error) {
	var dc apollo.DatabaseConnection
	bytes, err := (*config)["Database"].MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("unable to read database settings section: %w", err)
	}
	err = json.Unmarshal(bytes, &dc)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database settings: %w", err)
	}

	return &dc, nil
}

func ReadRabbitSettings(config *Config) (*apollo.RabbitConnection, *apollo.RabbitConsumerSettings, error) {
	buffer, err := (*config)["Rabbit"].MarshalJSON()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read rabbit settings section: %w", err)
	}

	rabbitMap := make(Config)
	err = json.Unmarshal(buffer, &rabbitMap)

	buffer, err = rabbitMap["Connection"].MarshalJSON()
	var rc apollo.RabbitConnection
	err = json.Unmarshal(buffer, &rc)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse rabbit connection settings: %w", err)
	}

	var cs apollo.RabbitConsumerSettings
	buffer, err = rabbitMap["Consumer"].MarshalJSON()
	err = json.Unmarshal(buffer, &cs)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse rabbit consumer settings: %w", err)
	}

	return &rc, &cs, nil
}
