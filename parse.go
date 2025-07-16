package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/idna"
	"golang.org/x/text/encoding/charmap"
)

type URLDomain struct {
	domain       string
	urls         []string
	urlNormalize []string
}

func urlToASCII(url string) (string, error) {
	encodedHost, err := idna.ToASCII(url)
	if err != nil {
		return "", err
	}
	return encodedHost, nil

}

func (ud *URLDomain) normalUrls() {
	if len(ud.urls) > 0 {
		// fmt.Println(ud.domain, ud.urls)
		for _, rawURL := range ud.urls {
			// Удаляем схему и приводим к нижнему регистру
			rawURL = removeScheme(rawURL)

			// Добавляем временную схему для парсинга
			if !strings.Contains(rawURL, "://") && !strings.HasPrefix(rawURL, "//") {
				rawURL = "http://" + rawURL
			}

			// Парсим URL
			u, err := url.Parse(rawURL)
			if err != nil {
				continue
			}

			if net.ParseIP(u.Hostname()) != nil {
				continue
			}

			// Обрабатываем IDN (кириллицу в домене)
			domain := u.Hostname()
			asciiDomain, err := idna.ToASCII(domain)
			if err != nil {
				continue
			}

			// Обновляем хост с учетом порта
			if u.Port() != "" {
				u.Host = asciiDomain + ":" + u.Port()
			} else {
				u.Host = asciiDomain
			}

			// Нормализуем путь - убираем только завершающие слеши, но не добавляем их
			u.Path = strings.TrimRight(u.Path, "/")
			// Полное кодирование пути (включая кириллицу)
			if u.Path != "" {
				u.Path = url.PathEscape(u.Path)
				u.Path = strings.ReplaceAll(u.Path, "%2F", "/") // Но сохраняем слеши
			}
			u.Path = strings.TrimRight(u.Path, "/")

			// Query-параметры: пробелы -> '+', остальное кодируем
			if u.RawQuery != "" {
				q := u.Query()
				u.RawQuery = q.Encode() // Стандартное кодирование QueryEscape
			}

			// Убираем стандартные порты
			if u.Port() == "80" || u.Port() == "443" {
				u.Host = u.Hostname()
			}

			// Собираем URL без схемы
			result := u.Host + u.Path
			if u.RawQuery != "" {
				result += "?" + u.RawQuery
			}
			if u.Fragment != "" {
				result += "#" + u.Fragment
			}
			if strings.Contains(result, ":") {
				continue
			}

			ud.urlNormalize = append(ud.urlNormalize, result)

		}
	}
}

// Проверяет нестандартные пути
func hasInvalidPath(path string) bool {
	// Регулярное выражение для проверки нестандартных конструкций
	invalidPattern := regexp.MustCompile(`:`)
	return invalidPattern.MatchString(path)
}

func removeScheme(rawURL string) string {
	schemes := []string{"http://", "https://", "ftp://", "//"}
	lowerURL := strings.ToLower(rawURL)
	for _, scheme := range schemes {
		if strings.HasPrefix(lowerURL, scheme) {
			return rawURL[len(scheme):]
		}
	}
	return rawURL
}

func parse() {
	domains := make(map[string]struct{})
	regexpDomains := make(map[string]struct{})
	fullURLs := make(map[string]struct{})
	// Открываем CSV файл
	// fullURL := make(map[string]struct{}{})
	file, err := os.Open(os.Getenv("PATHDUMP") + "/" + os.Getenv("DUMPFILE"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Создаем CSV reader
	decoder := charmap.Windows1251.NewDecoder()
	reader := csv.NewReader(decoder.Reader(file))

	//	reader := csv.NewReader(file)

	// Опционально: настраиваем параметры чтения
	reader.Comma = ';' // разделитель (по умолчанию ',')
	//reader.Comment = '#'        // символ комментария
	reader.FieldsPerRecord = -1 // количество полей (-1 - переменное)

	// Читаем данные
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Обрабатываем запись (record - это []string)
		if len(record) == 3 {

			domain, url := record[1], record[2]
			if domain == "" && url == "" {
				continue
			}

			urls := make([]string, 0)

			if strings.Contains(url, " | ") {
				urls = strings.Split(url, " | ")
			} else {
				if url != "" {
					urls = append(urls, url)
				}
			}

			data := URLDomain{domain, urls, make([]string, 0)}

			data.normalUrls()

			for _, urlNormalize := range data.urlNormalize {
				fullURLs[urlNormalize] = struct{}{}
			}

			if len(data.domain) < 4 {
				continue
			}

			asciidomain, err := urlToASCII(data.domain)

			if err != nil {
				continue
			}

			//fmt.Println(asciidomain)
			switch {
			case strings.HasPrefix(asciidomain, "www."):
				domains[asciidomain] = struct{}{}
				domains[asciidomain[4:]] = struct{}{}

			case strings.HasPrefix(asciidomain, "*."):
				regexpDomains[asciidomain] = struct{}{}
			default:
				domains[asciidomain] = struct{}{}
			}

		}

	}

	//
	zfile, err := os.Create(os.Getenv("PATHDUMP") + "/" + os.Getenv("ZINFOFILE"))

	if err != nil {
		log.Fatal(err)
	}

	defer zfile.Close()

	writer := bufio.NewWriter(zfile)

	for v, _ := range fullURLs {
		_, err := writer.WriteString("full:" + v + "\n")
		if err != nil {
			fmt.Println("Ошибка при записи строки:", err)
			log.Fatal(err)
		}
	}
	for v, _ := range domains {
		_, err := writer.WriteString("domain:" + v + "\n")
		if err != nil {
			fmt.Println("Ошибка при записи строки:", err)
			log.Fatal(err)
		}
	}
	for v, _ := range regexpDomains {
		// *.r7-casino-rus3.win
		// to
		// regexp:(.*\.)?zzslot13\.com
		regexpDomain := `regexp:(.*\.)?` + strings.ReplaceAll(v[2:], ".", `\.`)
		_, err := writer.WriteString(regexpDomain + "\n")
		if err != nil {
			fmt.Println("Ошибка при записи строки:", err)
			log.Fatal(err)
		}
	}
	if err := writer.Flush(); err != nil {
		fmt.Println("Ошибка при сбросе буфера:", err)
	}
}
