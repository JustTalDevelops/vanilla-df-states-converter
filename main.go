package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/df-mc/dragonfly/dragonfly/world"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"io/ioutil"
	"os"
)

type BlockState struct {
	Block world.BlockState `nbt:"block"`
	ID    int32            `nbt:"id"`
}

func main() {
	// Read arguments and error if there are too little
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Usage: vanilla-df-states-converter.exe vanilla-pallete.nbt output.txt")
		fmt.Println("Output is returned in base64 with no gzip.")
		return
	}

	// Mark the input and output file names
	inputFile := args[0]
	outputFile := args[1]

	// Load the contents of the input file
	b, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	// The vanilla palette that the vanilla states are unmarshalled into temporarily.
	var vanillaPalette struct {
		Blocks []BlockState `nbt:"blocks"`
	}

	// The vanilla states are gzip compressed, so before we can unmarshal them, we decompress them.
	gr, err := gzip.NewReader(bytes.NewBuffer(b))
	if err == nil {
		// Read the bytes from the IO reader
		b, err = ioutil.ReadAll(gr)
		if err != nil {
			panic(err)
		}

		defer gr.Close()
	}

	// Unmarshal the bytes into the vanilla palette with BigEndian encoding.
	err = nbt.UnmarshalEncoding(b, &vanillaPalette, nbt.BigEndian)
	if err != nil {
		panic(err)
	}

	// Create a new encoder with an empty byte buffer
	buff := bytes.NewBuffer([]byte{})
	e := nbt.NewEncoder(buff)

	// Encode every block state and add it to the buffer
	for _, s := range vanillaPalette.Blocks {
		err := e.Encode(&s)
		if err != nil {
			panic(err)
		}
	}

	// Write the output bytes to a file in base64
	err = ioutil.WriteFile(outputFile, []byte(base64.StdEncoding.EncodeToString(buff.Bytes())), 0777)
	if err != nil {
		panic(err)
	}
}
