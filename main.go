package main

import (
	"compress/gzip"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if os.Getenv("URL") == "" || os.Getenv("PATHDUMP") == "" {
		log.Fatal("create .env URL, PATHDUMP")
	}

}

func main() {

	action, err := updateDump(os.Getenv("URL"))

	if err != nil {
		log.Fatal(err)
	}
	if action {
		err = unGzip()

		if err != nil {
			log.Fatal(err)
		}
		parse()
		createDat()
		err = os.Remove(os.Getenv("PATHDUMP") + "/" + os.Getenv("DUMPFILE"))
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove(os.Getenv("PATHDUMP") + "/" + os.Getenv("ZINFOFILE"))
		if err != nil {
			log.Fatal(err)
		}

	}

}

func unGzip() error {
	// Открываем файл .gz
	gzFile, err := os.Open(os.Getenv("PATHDUMP") + "/" + os.Getenv("DUMPGZ"))
	if err != nil {
		return err
	}
	defer gzFile.Close()

	// Создаем новый файл для разархивированных данных
	outFile, err := os.Create(os.Getenv("PATHDUMP") + "/" + os.Getenv("DUMPFILE"))
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Создаем новый gzip.Reader
	reader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Копируем данные из gzip.Reader в выходной файл
	_, err = io.Copy(outFile, reader)
	if err != nil {
		return err
	}
	return nil

}
