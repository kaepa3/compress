package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Dst      string
	Src      string
	Name     string
	AddMonth bool
}

var config Config

func main() {
	// 設定読み込み
	readConfig()

	dest, err := os.Create(createDstName(config))
	if err != nil {
		panic(err)
	}

	zipWriter := zip.NewWriter(dest)
	defer zipWriter.Close()

	if err := addToZip(config.Src, zipWriter); err != nil {
		panic(err)
	}
}
func createDstName(conf Config) string {
	if conf.AddMonth {
		t := time.Now()
		return fmt.Sprintf("%s_%s_%d月.zip", conf.Dst, conf.Name, t.Month())
	}

	return fmt.Sprintf("%s_%s.zip", conf.Dst, conf.Name)
}
func readConfig() {
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		fmt.Println(err)
	}
}

func addToZip(filename string, zipWriter *zip.Writer) error {
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	stat, err := src.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(stat)
	if err != nil {
		return err
	}
	header.NonUTF8 = true

	writer, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, src)
	if err != nil {
		return err
	}

	return nil
}
