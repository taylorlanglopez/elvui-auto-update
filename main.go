package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var links []string

// This is my interface path for example: "C:\\Program Files (x86)\\World of Warcraft\\_retail_\\Interface\\"
// Requires \\ for backslash escaping
var retailPath string = "C:\\Program Files (x86)\\World of Warcraft\\_retail_\\Interface\\"

// This will get called for each HTML element found
func processElement(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element
	href, exists := element.Attr("href")
	if exists {
		links = append(links, href)
	}
}

func extFileName(s string, key byte) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == key {
			return s[i+1:]
		}
	}
	return ""
}

func ext(url string, key byte) string {
	for i := len(url) - 1; i >= 0; i-- {
		c := url[i]
		switch {
		case c == key:
			return url[i:]
		case '0' <= c && c <= '9':
		case 'A' <= c && c <= 'Z':
		case 'a' <= c && c <= 'z':
		default:
			return ""
		}
	}
	return ""
}

func getNewElvUIZip(arr []string) string {
	// Match this -> "elvui-11.27.zip"
	rExpObj, err := regexp.Compile("(elvui)-.*[.].*[.zip]")
	if err != nil {
		panic(err)
	}
	for _, value := range arr {
		if rExpObj.Match([]byte(value)) {
			return value
		}
	}
	return ""
}

func findZip(arr []string) string {
	if len(arr) < 1 {
		return ""
	}

	retSlice := []string{}
	for _, v := range arr {
		p := ext(v, '.')
		if p == ".zip" {
			retSlice = append(retSlice, v)
		}
	}
	fmt.Println("Possible .zip's -> ", retSlice)
	return getNewElvUIZip(retSlice)
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func checkDiffVersion(s string) bool {
	currentElvUITOC := retailPath + "AddOns\\ElvUI\\ElvUI.toc"
	outFile, err := os.Open(currentElvUITOC)

	if err != nil {
		fmt.Println("File not found at -> ", currentElvUITOC)
		return false
	}

	scanLines := bufio.NewScanner(outFile)
	var currentVersion string

	for scanLines.Scan() {
		if strings.Contains(scanLines.Text(), "Version") {
			currentVersion = scanLines.Text()
			break
		}
	}

	stringSlice := strings.Split(currentVersion, " ")
	verNum := stringSlice[len(stringSlice)-1]
	stringSlice = strings.Split(s, "-")
	installedNum := stringSlice[len(stringSlice)-1]
	finalInstalledNum := strings.Replace(installedNum, ".zip", "", 1)

	if finalInstalledNum == verNum {
		return false
	}

	return true
}

func main() {
	// Make HTTP request
	baseURL := "https://www.tukui.org"
	response, err := http.Get("https://www.tukui.org/download.php?ui=elvui")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Retrieved response from tukui.org.")
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find all links and process them with the function
	// defined earlier
	document.Find("a").Each(processElement)
	zipFragment := findZip(links)
	if zipFragment == "" {
		fmt.Println("ElvUI Zip pattern was not found, exiting.")
		return
	}
	fmt.Println("Zip found successfully -> ", zipFragment)
	fmt.Println("Checking version diff before downloading new ElvUI package...")
	if checkDiffVersion(zipFragment) == false {
		fmt.Println("ElvUI is already up to date, exiting")
		return
	}
	fmt.Println("ElvUI is correctly out of date, continuing")
	possibleDL := baseURL + zipFragment
	outputPath := retailPath
	resp, err := http.Get(possibleDL)
	if resp.StatusCode != 200 {
		fmt.Println("DL responded with an error, here is the link the script attempted to get ->", possibleDL)
	}
	if err != nil {
		panic(err)
	} else {
		fmt.Println("File downloaded from URL -> ", resp.Request.URL.String())
	}
	filename := extFileName(resp.Request.URL.String(), '/')
	finalPath := outputPath + filename
	out, e := os.Create(finalPath)
	if e != nil {
		panic(e)
	}
	fmt.Println("File created -> ", finalPath)
	_, err = io.Copy(out, resp.Body)
	Unzip(finalPath, outputPath+"AddOns\\")
	fmt.Println("File paths used during unzip process -> ", finalPath, outputPath+"AddOns\\")
	fmt.Println("Cleanig up zip at -> ", finalPath)
	delErr := os.Remove(finalPath)
	if delErr != nil {
		fmt.Println("Could not delete file at this requires manual cleanup -> ", finalPath)
	}
}
