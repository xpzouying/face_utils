package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var (
	flagFaceSvr = flag.String("face", "127.0.0.1:8080", "Face server address. Default is 127.0.0.1:8080")
	flagSrcDir  = flag.String("src", "", "source directory")
	flagDesDir  = flag.String("des", "", "Detect result directory")

	// image ext / type
	imgExt = []string{".jpg", ".bmp"}
)

// make detect request
// @input: name: detect file name
func makeDetectRequest(name string) (*http.Request, error) {

	// buffer to write
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// create img form file
	part, _ := writer.CreateFormFile("img", "img")
	// read file content
	data, err := ioutil.ReadFile(name)
	failedIfError(err)
	_, err = part.Write(data)
	failedIfError(err)
	writer.Close()

	detectURL := fmt.Sprintf("http://%s/detect", *flagFaceSvr)
	req, err := http.NewRequest("POST", detectURL, body)
	if err != nil {
		return nil, err
	}

	// MUST setting post header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

// get all files in src directory
// @return: files in srcDir
func getAllSrcFiles(srcDir string) ([]string, error) {
	var files []string

	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		// wask though path and process function

		// get file type ext
		ext := filepath.Ext(path)

		// if not support image type
		if !InList(imgExt, ext) {
			return nil
		}

		// append image file path to files
		files = append(files, path)

		return nil
	})

	return files, nil
}

// detect task
// @input: single image path
func detectTask(imgPath string) (DetectResult, error) {

	req, err := makeDetectRequest(imgPath)
	if err != nil {
		return DetectResult{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return DetectResult{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DetectResult{}, err
	}

	var detectRes DetectResult
	err = json.Unmarshal(body, &detectRes)
	if err != nil {
		return DetectResult{}, err
	}

	return detectRes, nil
}

// detect all images in path list
func detectAllImgs(imgsPath []string) ([]DetectResult, error) {
	var allRes []DetectResult

	for _, path := range imgsPath {
		res, err := detectTask(path)
		if err != nil {
			return nil, err
		}

		allRes = append(allRes, res)
	}

	return allRes, nil
}

func saveAllDetectFaces(detResList []DetectResult, desDir string) error {

	// if directory not exists, then create it
	if _, err := os.Stat(desDir); os.IsNotExist(err) {
		fmt.Println("Directory not exists. Create the dir: ", desDir)
		err = os.Mkdir(desDir, os.ModeDir|os.ModePerm)
		failedIfError(err)
	}

	// count faces
	var faceCount int

	// each detect result
	for _, detRes := range detResList {

		// each faceinfo
		for _, f := range detRes.FacesInfo {

			fmt.Println("String of base64: ", f.Image)

			imgBytes, err := base64.StdEncoding.DecodeString(f.Image)
			failedIfError(err)

			saveFileName := fmt.Sprintf("%s/face_%d.jpg", desDir, faceCount)
			fmt.Printf("Save filename: %s", saveFileName)
			err = ioutil.WriteFile(saveFileName, imgBytes, os.ModePerm)
			failedIfError(err)

			faceCount++
		}

	}

	return nil
}

func main() {
	flag.Parse()

	// get all images in src directory
	imgsPath, _ := getAllSrcFiles(*flagSrcDir)
	fmt.Println(imgsPath)

	detResList, err := detectAllImgs(imgsPath)
	failedIfError(err)

	err = saveAllDetectFaces(detResList, *flagDesDir)
	failedIfError(err)
}
