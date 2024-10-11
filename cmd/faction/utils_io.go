package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/voidshard/faction/pkg/structs"
)

type marshalable interface {
	MarshalYAML() ([]byte, error)
	UnmarshalYAML([]byte) error
}

func dumpYaml[T marshalable](data []T) ([]byte, error) {
	ydata := [][]byte{}
	for _, d := range data {
		b, err := d.MarshalYAML()
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

func readObjectsFromFile(desired structs.Object, files []string) ([]structs.Object, error) {
	found := []structs.Object{}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		for _, chunk := range bytes.Split(data, []byte("\n---\n")) {
			// TODO: I give up trying to make this work with generics or the interface
			switch desired.(type) {
			case *structs.World:
				obj := &structs.World{}
				err := obj.UnmarshalYAML(chunk)
				if err != nil {
					return nil, err
				}
				found = append(found, obj)
			case *structs.Faction:
				obj := &structs.Faction{}
				err := obj.UnmarshalYAML(chunk)
				if err != nil {
					return nil, err
				}
				found = append(found, obj)
			case *structs.Actor:
				obj := &structs.Actor{}
				err := obj.UnmarshalYAML(chunk)
				if err != nil {
					return nil, err
				}
				found = append(found, obj)
			}
		}
	}
	return found, nil
}
