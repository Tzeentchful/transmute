package main

import (
	"log"
	"os"

	"github.com/tzeentchful/transmute/types"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	filePath = kingpin.Arg("file", "Path to SMD file").Required().String()
	objPath  = kingpin.Arg("obj", "Path to OBJ file").Required().String()
)

var compressedMagic = [3]byte{0x43, 0x4D, 0x50}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	file, err := os.Open(*filePath)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	smd := types.NewSMD()
	smd.Decode(file)

	objFile, _ := os.Create(*objPath)
	defer objFile.Close()
	smd.Convert("obj", objFile)

}
