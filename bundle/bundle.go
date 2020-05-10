package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("path?")
		return
	}
	pwd := os.Args[1]
	sprites := spriteBundle(pwd)
	animations := animationBundle(pwd)
	tiles := tileBundle(pwd)
	resources := resourceBundle(pwd)

	var data strings.Builder
	data.WriteString("resources{")
	data.WriteString(resources)
	data.WriteString("},sprites{")
	data.WriteString(sprites)
	data.WriteString("},animations{")
	data.WriteString(animations)
	data.WriteString("},tiles{")
	data.WriteString(tiles)
	data.WriteString("},")

	bundle, err := os.Create(filepath.Join(pwd, "public", "wad"))
	if err != nil {
		panic(err)
	}
	defer bundle.Close()
	_, err = bundle.WriteString(data.String())
	if err != nil {
		panic(err)
	}
}

func tileBundle(pwd string) string {
	return compress(filepath.Join(pwd, "raw", "tiles"))
}

func resourceBundle(pwd string) string {
	fmt.Println("resources")

	var dir []os.FileInfo
	var err error
	var data strings.Builder

	shaders := filepath.Join(pwd, "public", "shaders")
	fmt.Println("shaders", shaders)
	dir, err = ioutil.ReadDir(shaders)
	if err != nil {
		panic(err)
	}
	data.WriteString("shaders[")
	for _, info := range dir {
		name := info.Name()
		extension := filepath.Ext(name)
		base := strings.TrimSuffix(name, extension)
		fmt.Println(name)
		data.WriteString(base)
		data.WriteString(",")
	}
	data.WriteString("],")

	images := filepath.Join(pwd, "public", "images")
	fmt.Println("images", images)
	dir, err = ioutil.ReadDir(images)
	if err != nil {
		panic(err)
	}
	data.WriteString("images[")
	for _, info := range dir {
		name := info.Name()
		extension := filepath.Ext(name)
		base := strings.TrimSuffix(name, extension)
		fmt.Println(name)
		data.WriteString(base)
		data.WriteString(",")
	}
	data.WriteString("],")

	music := filepath.Join(pwd, "public", "music")
	fmt.Println("music", music)
	dir, err = ioutil.ReadDir(music)
	if err != nil {
		panic(err)
	}
	data.WriteString("music[")
	for _, info := range dir {
		name := info.Name()
		fmt.Println(name)
		data.WriteString(name)
		data.WriteString(",")
	}
	data.WriteString("],")

	sounds := filepath.Join(pwd, "public", "sounds")
	fmt.Println("sounds", sounds)
	dir, err = ioutil.ReadDir(sounds)
	if err != nil {
		panic(err)
	}
	data.WriteString("sounds[")
	for _, info := range dir {
		name := info.Name()
		fmt.Println(name)
		data.WriteString(name)
		data.WriteString(",")
	}
	data.WriteString("]")

	return data.String()
}

func animationBundle(pwd string) string {
	path := filepath.Join(pwd, "raw", "animations")
	fmt.Println("animations", path)

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	var data strings.Builder
	for _, info := range dir {
		name := info.Name()
		fmt.Println(name)
		data.WriteString(name)
		data.WriteString("{")
		data.WriteString(compress(filepath.Join(path, name)))
		data.WriteString("},")
	}

	return data.String()
}

func scaleSprites(pwd string) string {
	path := filepath.Join(pwd, "raw", "sprites")
	fmt.Println("sprites", path)

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	var data strings.Builder
	for _, info := range dir {
		name := info.Name()
		fmt.Println(name)
		from := filepath.Join(path, name)
		to := filepath.Join(pwd, "raw", "sprites-upscale", name)
		err := os.MkdirAll(to, os.ModePerm)
		if err != nil {
			panic(err)
		}
		upscaleSprite(from, to)
	}

	return data.String()
}

func spriteBundle(pwd string) string {
	scaleSprites(pwd)
	path := filepath.Join(pwd, "raw", "sprites-upscale")
	fmt.Println("sprites", path)

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	var data strings.Builder
	for _, info := range dir {
		name := info.Name()
		fmt.Println(name)
		data.WriteString(name)
		data.WriteString("{")
		str := spritePacker(name, filepath.Join(path, name), filepath.Join(pwd, "public", "images"))
		data.WriteString(str)
		data.WriteString("},")
	}

	return data.String()
}

func compress(path string) string {
	var data strings.Builder

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := strings.ReplaceAll(scan.Text(), " ", "")
		data.WriteString(line)
		if strings.HasSuffix(line, "]") || strings.HasSuffix(line, "}") {
			data.WriteString(",")
		}
	}

	return data.String()
}
