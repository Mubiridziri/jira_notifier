package telegram

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"jira_notifier/helpers"
	"jira_notifier/models"
	"log"
	"os"
	strings2 "strings"
)

type Label struct {
	Text     string
	FontSize float64
	XPos     int
	YPos     int
}

func getCoverForNewIssue(issue models.Issue) string {
	labels := getLabelsByIssue(issue)
	return DrawImageByTemplate("resources/new_issue_cover.jpg", labels)

}

func getCoverForUpdatedIssue(issue models.Issue) string {
	labels := getLabelsByIssue(issue)
	return DrawImageByTemplate("resources/updated_issue_cover.jpg", labels)
}

func getLabelsByIssue(issue models.Issue) []Label {
	return []Label{
		{Text: issue.User.Name + "!", XPos: 290, YPos: 280, FontSize: 16},
		{Text: fmt.Sprintf("%v %v", issue.Tag, issue.Title), XPos: 98, YPos: 380, FontSize: 14},
		{Text: issue.Author, XPos: 247, YPos: 825, FontSize: 14},
		{Text: issue.Priority, XPos: 334, YPos: 910, FontSize: 14},
	}
}

func DrawImageByTemplate(templatePath string, labels []Label) string {
	path, err := os.Open(templatePath)

	if err != nil {
		panic(err)
	}

	defer path.Close()

	img, _, err := image.Decode(path)

	if err != nil {
		panic(err)
	}

	//REDRAW IMAGE
	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Src)

	//FONT
	fontFile, err := os.ReadFile("fonts/Helvetica-Bold.ttf")

	f, err := opentype.Parse(fontFile)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	if err != nil {
		log.Fatalf("failed to create new face: %v", err)
	}

	for _, label := range labels {
		face, err := opentype.NewFace(f, &opentype.FaceOptions{
			Size:    float64(img.Bounds().Dx() / 5),
			DPI:     label.FontSize,
			Hinting: font.HintingNone,
		})

		if err != nil {
			panic(err)
		}

		wrappedStrings := WrapString(label.Text, 28)
		YPos := label.YPos
		for _, text := range wrappedStrings {
			DrawLabel(dst, face, label.XPos, YPos, text)
			YPos = YPos + 60
		}

	}

	newFile, err := os.CreateTemp("resources", "cover_")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer newFile.Close()

	b := bufio.NewWriter(newFile)
	if err := jpeg.Encode(b, dst, &jpeg.Options{Quality: 100}); err != nil {
		log.Fatalf("failed to encode image: %v", err)
	}

	return newFile.Name()
}

func WrapString(text string, length int) []string {
	if len(text) <= length {
		return []string{text}
	}

	subString := ""
	var strings []string
	puncs := []string{
		" ",
		",",
		".",
		":",
		"-",
	}

	runes := bytes.Runes([]byte(text))
	l := len(runes)
	for index, symbolRune := range runes {
		currentSymbol := string(symbolRune)
		subString = subString + currentSymbol

		if (index+1)%length == 0 {
			if helpers.InStringArray(puncs, currentSymbol) != -1 {
				strings = append(strings, subString)
				subString = ""
				continue
			}
			foundIndex := helpers.GetLastFoundSymbolIndex(subString, " ")
			subStrRunes := bytes.Runes([]byte(subString))

			str1 := subStrRunes[0:foundIndex]
			str2 := subStrRunes[foundIndex : len(subStrRunes)-1]
			strings = append(strings, strings2.TrimLeft(string(str1), " "))
			subString = strings2.TrimLeft(string(str2), " ") + currentSymbol

		} else if (index + 1) == l {
			strings = append(strings, subString)
		}
	}

	return strings
}

func DrawLabel(dst *image.RGBA, face font.Face, x int, y int, str string) {
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.RGBA{24, 24, 24, 255}),
		Face: face,
		Dot:  fixed.Point26_6{fixed.I(x), fixed.I(y)},
	}
	d.DrawString(str)
}
