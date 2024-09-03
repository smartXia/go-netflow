package utils

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

func HumanBytes(n int64) string {
	return humanize.Bytes(uint64(n))
}
func UnmarshalConfigFromFile(filePath string, config interface{}) error {
	data, err := GetFileContent(filePath)
	if err != nil {
		return err
	}
	if strings.HasSuffix(filePath, ".json") {
		return json.Unmarshal(data, config)
	}
	if strings.HasSuffix(filePath, ".yml") || strings.HasSuffix(filePath, ".yaml") {
		return yaml.Unmarshal(data, config)
	}
	return fmt.Errorf("unknown extension of file: %s", filePath)
}
func GetFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
