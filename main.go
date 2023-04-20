package main

import (
	"bytes"
	"compress/gzip"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"io"
	"os"
)

func main() {
	inputFile := "block_palette.nbt"
	outputFile := "block_states.nbt"

	b, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	var vanillaPalette struct {
		Blocks []struct {
			Name       string         `nbt:"name"`
			Properties map[string]any `nbt:"states"`
			Version    int32          `nbt:"version"`
			NameHash   int64          `nbt:"name_hash"`
		} `nbt:"blocks"`
	}

	gr, err := gzip.NewReader(bytes.NewBuffer(b))
	if err == nil {
		b, err = io.ReadAll(gr)
		if err != nil {
			panic(err)
		}
		_ = gr.Close()
	}

	err = nbt.UnmarshalEncoding(b, &vanillaPalette, nbt.BigEndian)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	e := nbt.NewEncoder(buf)

	for _, s := range vanillaPalette.Blocks {
		err := e.Encode(struct {
			Name       string         `nbt:"name"`
			Properties map[string]any `nbt:"states"`
			Version    int32          `nbt:"version"`
		}{
			Name:       s.Name,
			Properties: s.Properties,
			Version:    s.Version,
		})
		if err != nil {
			panic(err)
		}
	}

	err = os.WriteFile(outputFile, buf.Bytes(), 0777)
	if err != nil {
		panic(err)
	}
}
