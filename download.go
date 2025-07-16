package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func updateDump(url string) (bool, error) {
	etagFile := "dump.etag"

	parts := strings.Split(url, "/")
	dumpFile := parts[len(parts)-1]

	if _, err := os.Stat(os.Getenv("PATHDUMP") + "/" + dumpFile); os.IsNotExist(err) {
		err := downloadDump(os.Getenv("URL"))
		if err != nil {
			return false, err
		}
		return true, nil
	}

	etag, err := getEtag(os.Getenv("URL"))
	if err != nil {
		return false, err
	}

	etagFromFile, err := readEtagFile(os.Getenv("PATHDUMP") + "/" + etagFile)

	if err != nil {
		return false, err
	}

	if etag != etagFromFile {
		err := downloadDump(os.Getenv("URL"))
		if err != nil {
			return false, err
		}
		fmt.Println("Обновлнено")
		return true, nil
	}

	fmt.Println("Обновлений не требуется")
	return false, nil

}

func downloadPathCreate(dumpPath string) error {
	if _, err := os.Stat(dumpPath); os.IsNotExist(err) {
		err := os.Mkdir(dumpPath, 0775)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeEtagFile(filepath, etag string) error {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Ошибка при создании файла etag:", err)
		return err
	}
	defer file.Close() // Закрыть файл после завершения работы

	_, err = file.WriteString(etag)
	if err != nil {
		fmt.Println("Ошибка при записи в файл etag:", err)
		return err
	}

	// fmt.Println("ETag успешно записан в файл.")

	return nil
}

func readEtagFile(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Ошибка при чтении файла:", err)
		return "", err
	}

	etag := string(data)
	// fmt.Println("Считанный ETag:", etag)
	return etag, nil

}

func getEtag(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	etag := resp.Header.Get("Etag")
	if etag == "" {
		return "", errors.New("not header etag")
	}

	return etag, nil
}

func downloadDump(url string) error {
	parts := strings.Split(url, "/")
	dumpFile := parts[len(parts)-1]

	etag, err := getEtag(url)

	if err != nil {
		return err
	}

	err = downloadPathCreate(os.Getenv("PATHDUMP"))

	if err != nil {
		return err
	}

	out, err := os.Create(os.Getenv("PATHDUMP") + "/" + dumpFile)
	if err != nil {
		return errors.New("not create z-i filepath")
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return errors.New("not download")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("status not ok")

	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = writeEtagFile(os.Getenv("PATHDUMP")+"/"+"dump.etag", etag)
	if err != nil {
		return err
	}

	return nil
}
