package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BrendanBurkhart/raytracer/internal/camera"
	"github.com/BrendanBurkhart/raytracer/internal/scene"
)

func main() {
	sceneData := os.Args[1:]

	sceneCount := 0

	for _, path := range sceneData {
		sceneCount += walkPath(path)
	}

	fmt.Printf("Sucessfully rendered %d scene(s)\n", sceneCount)
}

func walkPath(path string) (sceneCount int) {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error while walking %s: %v\n", path, err)
		return
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		var subpaths []os.FileInfo
		subpaths, err = ioutil.ReadDir(path)
		if err != nil {
			fmt.Printf("Error while walking %s: %v\n", path, err)
			return
		}

		for _, subpath := range subpaths {
			sceneCount += walkPath(fmt.Sprintf("%s/%s", path, subpath.Name()))
		}
	case mode.IsRegular():
		if filepath.Ext(path) != ".json" {
			return
		}

		err = renderScene(path)
		if err != nil {
			fmt.Printf("Error from %s: %v\n", path, err)
		} else {
			sceneCount++
		}
	}

	return
}

func renderScene(inputPath string) error {
	outputPath := fmt.Sprintf("%s.png", strings.TrimSuffix(inputPath, ".json"))

	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("unable to open data file: %v", err)
	}
	defer input.Close()

	data := &struct {
		Width  int           `json:"width"`
		Height int           `json:"height"`
		Camera camera.Camera `json:"camera"`
		Scene  scene.Scene   `json:"scene"`
	}{}

	if err = json.NewDecoder(input).Decode(data); err != nil {
		return fmt.Errorf("couldn't unmarshal scene data: %v", err)
	}

	if err = data.Scene.Initialize(); err != nil {
		return fmt.Errorf("couldn't initialize scene: %v", err)
	}

	if err != nil {
		return fmt.Errorf("unable to create scene: %v", err)
	}

	err = data.Camera.SetImageSize(data.Width, data.Height)
	if err != nil {
		return fmt.Errorf("error setting camera image size: %v", err)
	}

	fmt.Printf("Rendering scene (using %s lens) from: %s\n", data.Camera.GetLensName(), inputPath)

	if err = data.Camera.Render(&data.Scene, 15, 2<<10); err != nil {
		return fmt.Errorf("error while raytracing scene: %v", err)
	}

	output, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to open output file: %v", err)
	}
	defer output.Close()

	if err = data.Camera.Save(output); err != nil {
		return fmt.Errorf("unable to encode rendering: %v", err)
	}

	if err = output.Sync(); err != nil {
		return fmt.Errorf("unable to save rendering as PNG: %v", err)
	}

	return nil
}
