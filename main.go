package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dsnet/compress/brotli"
	"github.com/tzeentchful/transmute/types"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	filePath = kingpin.Arg("file", "Path to SMD file").Required().String()
	objPath  = kingpin.Arg("obj", "Path to OBJ file").Required().String()
)

var compressedMagic = []byte{0x43, 0x4D, 0x50}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	file, err := os.Open(*filePath)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	smd := types.NewSMD()

	var start = make([]byte, 3)
	file.Read(start)

	if bytes.Compare(start, compressedMagic) == 0 {
		var rawLength = make([]byte, 4)
		file.Read(rawLength)
		var length = binary.LittleEndian.Uint32(rawLength)

		fmt.Printf("decode length: 0x%#x %d\n", rawLength, length)

		brotliReader, err := brotli.NewReader(bufio.NewReader(file), nil)

		if err != nil {
			log.Fatalf("unexpected NewReader error: %v", err)
		}

		smd.Decode(brotliReader)
	} else {
		file.Seek(0, io.SeekStart)
		smd.Decode(file)
	}

	objFile, _ := os.Create(*objPath)
	defer objFile.Close()
	smd.Convert("obj", objFile)

}
