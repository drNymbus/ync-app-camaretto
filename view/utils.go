package view

import (
	"log"

	"os"
	"io"
	// "bytes"
	"image"
	// "image/color"

	// "golang.org/x/text/language"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	// "github.com/hajimehoshi/ebiten/v2/text/v2"
)

// @desc: Given a file path returns the content of the file as a byte slice
func GetFileByte(filename string) []byte {
	var err error

	// Opens the file
	var file *os.File
	file, err = os.Open(filename)
	if err != nil { log.Fatal("[parameters.getFileByte] Open file:", err) }

	// Get file size, to make sure we read it correctly/entirely
	var stat os.FileInfo
	stat, err = file.Stat()
	if err != nil { log.Fatal("[parameters.getFileByte] Get file info:", err) }

	// Turn the file into a byte slice
	var n int
	var fileByte []byte = make([]byte, stat.Size())
	n, err = file.Read(fileByte)
	if err != nil && err != io.EOF { log.Fatal("[parameters.getFileByte] Read file:", err) }
	if n != int(stat.Size()) { log.Println("Warning | [parameters.getFile] File read partially:", err) }

	return fileByte
}

// @desc: Given a file path load the file as an ebiten.Image object
func GetImage(filename string) *ebiten.Image {
	var err error
	var img *ebiten.Image

	img, _, err = ebitenutil.NewImageFromFile(filename)
	if err != nil {
		var msg string = "[view.GetImage] Load image (" + filename + "):"
		log.Fatal(msg, err)
	}

	return img
}

// @desc: 
func InitIcon(filepath string) (image.Image, error) {
	var err error
	var file *os.File

	file, err = os.Open(filepath)
	if err != nil {
		return nil, err
	}

	image, _, err := image.Decode(file)
	return image, err
}