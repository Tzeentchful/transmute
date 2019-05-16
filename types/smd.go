package types

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tzeentchful/transmute/utils"
)

type MeshDefinition struct {
	NameLen    uint32
	FaceOffset uint32
}
type SMDVertex struct {
	X          float32
	Y          float32
	Z          float32
	Unk1       [4]int16
	U          utils.Float16
	V          utils.Float16
	BoneIndex  [4]byte
	BoneWeight [4]byte
}

type SMD struct {
	Header struct {
		Version        byte
		NumIdx         uint32
		NumVert        uint32
		Unk1           byte
		NumMeshs       byte
		Unk2           byte
		TotalStringLen uint32
		BoundingBox    [6]float32
	}

	MeshDefinitions []MeshDefinition
	Names           []string

	IndexBuffer []uint16

	VertexBuffer []SMDVertex
}

func NewSMD() (this *SMD) {
	this = new(SMD)
	return
}

func (this *SMD) Decode(file io.Reader) bool {
	fmt.Printf("Decoding\n")
	// Read the header
	headerData := utils.ReadNextBytes(file, 40)
	err := binary.Read(bytes.NewBuffer(headerData), binary.LittleEndian, &this.Header)
	if err != nil {
		log.Fatal("Header read failed ", err)
	}

	fmt.Printf("Header: %v\n", this.Header)

	this.MeshDefinitions = make([]MeshDefinition, this.Header.NumMeshs)
	for i := 0; i < int(this.Header.NumMeshs); i++ {
		meshDefsData := utils.ReadNextBytes(file, 8)
		err = binary.Read(bytes.NewBuffer(meshDefsData), binary.LittleEndian, &this.MeshDefinitions[i])
		if err != nil {
			log.Fatal("Mesh def read failed ", err)
		}
	}

	this.Names = make([]string, this.Header.NumMeshs)
	for i := 0; i < int(this.Header.NumMeshs); i++ {
		meshNameData := utils.ReadNextBytes(file, int(this.MeshDefinitions[i].NameLen))
		stringData := make([]byte, this.MeshDefinitions[i].NameLen)
		err = binary.Read(bytes.NewBuffer(meshNameData), binary.LittleEndian, &stringData)

		if err != nil {
			log.Fatal("String read failed ", err)
		}

		this.Names[i], err = utils.DecodeUTF16(stringData)

		fmt.Printf("String: %#x : %s\n", stringData, this.Names[i])

		if err != nil {
			log.Fatal("String read failed ", err)
		}
	}

	this.IndexBuffer = make([]uint16, this.Header.NumIdx*3)

	for i := 0; i < int(this.Header.NumIdx*3); i++ {
		this.IndexBuffer[i] = binary.LittleEndian.Uint16(utils.ReadNextBytes(file, 2))
	}

	this.VertexBuffer = make([]SMDVertex, this.Header.NumVert)

	for i := 0; i < int(this.Header.NumVert); i++ {
		vetData := utils.ReadNextBytes(file, 32)
		err = binary.Read(bytes.NewBuffer(vetData), binary.LittleEndian, &this.VertexBuffer[i])
		if err != nil {
			log.Fatal("Error reading vert ", err)
		}
		if this.VertexBuffer[i].U < 0 {
			fmt.Printf("Found at %d\n", i)
			break
		}
	}

	return true
}

func (this *SMD) Convert(fileType string, file *os.File) {
	fmt.Printf("Converting\n")
	w := bufio.NewWriter(file)

	// Write verts
	for i := 0; i < len(this.VertexBuffer); i++ {
		vert := this.VertexBuffer[i]
		// fmt.Printf("Vert: %v\n", vert)
		if vert.U < 0 {
			fmt.Printf("Found at %d\n", i)
			break
		}

		w.WriteString(fmt.Sprintf("v %f %f %f\n", vert.X, -vert.Z, vert.Y))
	}

	// Write UVs
	for i := 0; i < len(this.VertexBuffer); i++ {
		vert := this.VertexBuffer[i]
		w.WriteString(fmt.Sprintf("vt %s %s\n", vert.U, vert.V))
	}

	currGroup := 0
	// Write Faces
	for i := 0; i < int(this.Header.NumIdx); i++ {
		if len(this.MeshDefinitions) > currGroup && i == int(this.MeshDefinitions[currGroup].FaceOffset) {
			w.WriteString(fmt.Sprintf("g %s\n", this.Names[currGroup]))
			currGroup++
		}
		w.WriteString(
			fmt.Sprintf("f %d/%d %d/%d %d/%d\n",
				this.IndexBuffer[i*3]+1, this.IndexBuffer[i*3]+1,
				this.IndexBuffer[i*3+1]+1, this.IndexBuffer[i*3+1]+1,
				this.IndexBuffer[i*3+2]+1, this.IndexBuffer[i*3+2]+1))
	}

	w.Flush()
}
