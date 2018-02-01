package types

import (
	"os"
)

type GGGFile interface {
	Decode(file *os.File) bool
	Validate(file *os.File) bool
	Convert(fileType string) []byte
}
