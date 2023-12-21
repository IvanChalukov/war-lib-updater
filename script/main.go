package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: <war-file-or-url> <groupId:artifactId:version> [<groupId:artifactId:version> ...]")
		os.Exit(1)
	}

	warFileOrUrl := os.Args[1]
	libraries := os.Args[2:]

	warFilePath, err := handleWarFile(warFileOrUrl)
	if err != nil {
		fmt.Printf("Error handling WAR file: %s\n", err)
		os.Exit(1)
	}

	// To do unzip file
	fmt.Println(warFilePath)

	// TODO update libs
	for _, lib := range libraries {
		fmt.Println(lib)
	}

}

func handleWarFile(warFileOrUrl string) (string, error) {
	fileName := "new.war"

	if strings.HasPrefix(warFileOrUrl, "http") {
		fmt.Println("Start downloading")
		return downloadFile(warFileOrUrl, fileName)
	}

	return copyFile(warFileOrUrl, fileName)
}

func downloadFile(url string, dst string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", nil
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return dst, err
}

func copyFile(src, dst string) (string, error) {
	input, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(dst, input, 0666)
	return dst, err
}
