package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/voidshard/faction/pkg/kind"
	v1 "github.com/voidshard/faction/pkg/structs/v1"

	"gopkg.in/yaml.v3"
)

func dumpYaml(data []v1.Object) ([]byte, error) {
	ydata := [][]byte{}
	for _, d := range data {
		b, err := yaml.Marshal(d)
		if err != nil {
			return nil, err
		}
		ydata = append(ydata, bytes.TrimSpace(b))
	}

	return bytes.Join(ydata, []byte("\n---\n")), nil
}

func calculateFileHash(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func readObjectsFromFile(files []string) ([]v1.Object, error) {
	found := []v1.Object{}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		for _, chunk := range bytes.Split(data, []byte("\n---\n")) {
			obj, err := kind.New("", chunk)
			if err != nil {
				return nil, err
			}
			found = append(found, obj)
		}
	}
	return found, nil
}
