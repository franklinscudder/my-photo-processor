package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/fs"
	"os"

	"github.com/nfnt/resize"
)

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	return img, err
}

func RemoveIndex(s []fs.DirEntry, index int) []fs.DirEntry {
	return append(s[:index], s[index+1:]...)
}

func main() {
	var inDirFlag = flag.String("i", ".", "specify an input directory")
	var outDirFlag = flag.String("o", "./photos-out", "specify an output directory")
	flag.Parse()

	d, err := os.Open(*inDirFlag)

	if err != nil {
		fmt.Println("Error opening specified input dir!")
		return
	}
	defer d.Close()

	files, err := d.ReadDir(-1)

	if err != nil {
		fmt.Println("Error reading specified input dir!")
		return
	}

	var filteredFiles []fs.DirEntry
	for _, path := range files {
		strName := path.Name()
		over5Length := len(strName) >= 5
		if over5Length && (strName[len(strName)-4:] == ".jpg" || strName[len(strName)-5:] == ".jpeg") {
			filteredFiles = append(filteredFiles, path)
		}
	}

	if len(filteredFiles) == 0 {
		fmt.Printf("No jpegs found in %v", *inDirFlag)
		return
	}

	err = os.Mkdir(*outDirFlag, 0755)
	if err != nil {
		fmt.Printf("Error creating output dir: %v", err)
		return
	}

	for _, path := range filteredFiles {
		img, err := getImageFromFilePath(*inDirFlag + "/" + path.Name())

		if err != nil {
			fmt.Printf("Error opening file %v: %v", path.Name(), err)
		}

		bounds := img.Bounds()
		dx := bounds.Max.X - bounds.Min.X
		dy := bounds.Max.Y - bounds.Min.Y

		newImage := resize.Thumbnail(uint(dx/2), uint(dy/2), img, resize.Bilinear)

		f, err := os.Create(*outDirFlag + "/R_" + path.Name())
		if err != nil {
			fmt.Printf("Error creating output file %v: %v", *outDirFlag+"/R_"+path.Name(), err)
		}
		defer f.Close()
		if err = jpeg.Encode(f, newImage, nil); err != nil {
			fmt.Printf("failed to encode: %v", err)
		}
	}
}
