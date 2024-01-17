package main

import (
	"flag"
	"fmt"
	"github.com/essentialkaos/go-badge"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Threshold struct {
	yellow int
	green  int
}

type Params struct {
	label     string
	threshold Threshold
	color     string
	value     string
	link      string
	style     string
}

func main() {

	source := flag.String("filename", "output.out", "File containing the tests output")
	label := flag.String("text", "Coverage", "Text on the left side of the badge")
	yellowThreshold := flag.Int("yellow", 30, "At what percentage does the badge becomes yellow instead of red")
	greenThreshold := flag.Int("green", 70, "At what percentage does the badge becomes green instead of yellow")
	color := flag.String("color", "", "Color of the badge - green/yellow/red")
	target := flag.String("target", "coverage.svg", "Target file")
	value := flag.String("value", "", "Text on the right side of the badge")
	link := flag.String("link", "", "Link the badge goes to")
	style := flag.String("style", "plastic", "Style of the badge")

	flag.Parse()

	params := &Params{
		*label,
		Threshold{*yellowThreshold, *greenThreshold},
		*color,
		*value,
		*link,
		*style,
	}

	err := generateBadge(*source, *target, params)

	if err != nil {
		panic(err)
	}
}

func generateBadge(source string, target string, params *Params) error {
	var err error
	g, err := badge.NewGenerator("Verdana.ttf", 11)
	if err != nil {
		panic(err)
	}
	funcs := map[string]func(string, string, string) []byte{
		"flat":    g.GenerateFlat,
		"square":  g.GenerateFlatSquare,
		"plastic": g.GeneratePlastic,
	}

	var coverage string
	if params.value != "" {
		coverage = params.value
	} else {
		coverage, err = retrieveTotalCoverage(source)
	}

	if err != nil {
		return err
	}

	badgeColor := setColor(coverage, params.threshold.yellow, params.threshold.green, params.color)

	fn, ok := funcs[params.style]
	if !ok {
		return fmt.Errorf("wrong style")
	}

	err = saveSvg(target, fn(params.label, coverage, badgeColor))

	if err != nil {
		return err
	}

	fmt.Println("\033[0;36mGoBadge: Coverage badge updated to " + coverage + " in " + target + "\033[0m")

	return nil
}

func setColor(coverage string, yellowThreshold int, greenThreshold int, color string) string {
	coverageNumber, _ := strconv.ParseFloat(strings.Replace(coverage, "%", "", 1), 4)
	if color != "" {
		return color
	}
	if coverageNumber >= float64(greenThreshold) {
		return badge.COLOR_BRIGHTGREEN
	}
	if coverageNumber >= float64(yellowThreshold) {
		return badge.COLOR_YELLOW
	}
	return badge.COLOR_RED
}

func retrieveTotalCoverage(filename string) (string, error) {
	// Read coverage file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("\033[1;31mGoBadge: Error while opening the coverage file\033[0m")
		return "", err
	}
	defer file.Close()

	// split content by words and grab the last one (total percentage)
	b, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("\033[1;31mGoBadge: Error while reading the coverage file\033[0m")
		return "", err
	}
	words := strings.Fields(string(b))
	last := words[len(words)-1]

	return last, nil
}

func saveSvg(target string, data []byte) error {

	file, err := os.Create(target)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
