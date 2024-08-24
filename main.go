package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	inDir             = "./input"
	outDir            = "./output"
	cardFrameFileName = inDir + "/frame.png"
	symbols           = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "A", "J", "Q", "K"}
	suits             = []string{"clubs", "diamonds", "hearts", "spades"}
)

func main() {
	fmt.Println("Generating card names to layer files...")

	cardNamesToLayerFiles, err := getCardNamesToLayerFiles()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print("Done!\n\n")
	fmt.Println("Generating card images...")

	err = generateCardDeckImages(cardNamesToLayerFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Done!")
}

func getCardNamesToLayerFiles() (map[string][]string, error) {
	cardNamesToLayerFiles := make(map[string][]string)

	if _, err := os.Stat(cardFrameFileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("card frame file (%s) does not exist", cardFrameFileName)
	}

	getColorFromSuit := func(suit string) string {
		switch suit {
		case "hearts", "diamonds":
			return "red"
		case "clubs", "spades":
			return "black"
		default:
			return ""
		}
	}

	// There are 4 types of files in inDir:
	// 1. frame:   (name: "frame.png", exists only one)
	// 2. content: (format: "symbol-suit.png.png", e.g. "4-clubs.png", "10-spades.png")
	// 3. symbol:  (format: "symbol-color.png", e.g. "4-black.png", "10-red.png")
	// 4. suit:    (format: "suit.png", e.g. "hearts.png", "spades.png")

	// 1. Search in inDir for each "symbol-suit.png" file.
	// 2. Search in inDir for corresponding "symbol-color.png" file.
	// 3. Search in inDir for corresponding "suit.png" file.
	// 4. If any of these files does not exist, continue loop. Otherwise, add to map.

	for _, symbol := range symbols {
		for _, suit := range suits {
			cardName := fmt.Sprintf("%s-%s", symbol, suit)

			contentFileName := fmt.Sprintf("%s/%s-%s.png", inDir, symbol, suit)
			if _, err := os.Stat(contentFileName); os.IsNotExist(err) {
				continue
			}

			symbolFileName := fmt.Sprintf("%s/%s-%s.png", inDir, symbol, getColorFromSuit(suit))
			if _, err := os.Stat(symbolFileName); os.IsNotExist(err) {
				continue
			}

			suitFileName := fmt.Sprintf("%s/%s.png", inDir, suit)
			if _, err := os.Stat(suitFileName); os.IsNotExist(err) {
				continue
			}

			cardNamesToLayerFiles[cardName] = []string{
				cardFrameFileName, contentFileName, symbolFileName, suitFileName,
			}
		}
	}

	return cardNamesToLayerFiles, nil
}

func generateCardDeckImages(cardNamesToLayerFiles map[string][]string) error {
	if len(cardNamesToLayerFiles) == 0 {
		return nil
	}

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err = os.Mkdir(outDir, 0o755)
		if err != nil {
			return err
		}
	}

	imageMagickOptions := "-gravity center -background None -layers Flatten"

	for cardName, layerFiles := range cardNamesToLayerFiles {
		input := strings.Join(layerFiles, " ")
		output := fmt.Sprintf("%s/%s.png", outDir, cardName)

		command := fmt.Sprintf("magick %s %s %s", input, imageMagickOptions, output)

		err := exec.Command("bash", "-c", command).Run()
		if err != nil {
			return err
		}
	}

	return nil
}
