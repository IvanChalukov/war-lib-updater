package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	tempDir, err := unzipWarFile(warFilePath)
	if err != nil {
		fmt.Printf("Error unzipping WAR file: %s\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	for _, lib := range libraries {
		if err := downloadAndReplaceLibrary(tempDir, lib); err != nil {
			fmt.Printf("Error processing library '%s': %s\n", lib, err)
			os.Exit(1)
		}
	}

	if err := rezipWarFile(tempDir, warFilePath); err != nil {
		fmt.Printf("Error rezipping WAR file: %s\n", err)
		os.Exit(1)
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

func downloadAndReplaceLibrary(tempDir, lib string) error {
	parts := strings.Split(lib, ":")
	if len(parts) != 3 {
		return fmt.Errorf("invalid library format")
	}
	groupID, artifactID, version := parts[0], parts[1], parts[2]

	mavenURL := fmt.Sprintf("https://repo1.maven.org/maven2/%s/%s/%s/%s-%s.jar",
		strings.ReplaceAll(groupID, ".", "/"), artifactID, version, artifactID, version)

	artifactName := fmt.Sprintf("%s-%s.jar", artifactID, version)
	jarPath, err := downloadFile(mavenURL, artifactName)
	if err != nil {
		return err
	}
	defer os.Remove(jarPath)

	libPath := filepath.Join(tempDir, "WEB-INF/lib", fmt.Sprintf("%s-%s.jar", artifactID, version))
	return replaceFile(jarPath, libPath)
}

func unzipWarFile(warFilePath string) (string, error) {
	tempDir, err := os.MkdirTemp("", "unzipped-war-")
	if err != nil {
		return "", err
	}

	r, err := zip.OpenReader(warFilePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		filePath := filepath.Join(tempDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := unzipFile(f, filePath); err != nil {
			return "", err
		}
	}
	return tempDir, nil
}

func unzipFile(f *zip.File, filePath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

func rezipWarFile(tempDir, warFilePath string) error {
	newWarFile, err := os.Create(warFilePath)
	if err != nil {
		return err
	}
	defer newWarFile.Close()

	zipWriter := zip.NewWriter(newWarFile)
	defer zipWriter.Close()

	return filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == tempDir {
			return nil
		}
		return addToZip(zipWriter, path, tempDir)
	})
}

func addToZip(zipWriter *zip.Writer, filePath, baseDir string) error {
	fileToZip, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name, _ = filepath.Rel(baseDir, filePath)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func replaceFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}
