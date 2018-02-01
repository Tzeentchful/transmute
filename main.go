package main

import (
	"fmt"
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

type Header struct {
	Version        byte
	NumIdx         uint32
	NumVert        uint32
	Unk1           byte
	NumMeshs       byte
	Unk2           byte
	TotalStringLen uint32
	BoundingBox    [6]float32
}

type MeshDef struct {
	NameLen    uint32
	FaceOffset uint32
}

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

	fmt.Printf("Parsed data:\n%+v\n", smd)

	objFile, _ := os.Create(*objPath)
	defer objFile.Close()
	smd.Convert("obj", objFile)

}
