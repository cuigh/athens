package pack

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func ModifyMod(modFile string, module, replace string) error {
	b, err := ioutil.ReadFile(modFile)
	if err != nil {
		return err
	}

	r := regexp.MustCompile(`module\s+` + replace)
	b = r.ReplaceAll(b, []byte("module "+module))
	return ioutil.WriteFile(modFile, b, os.ModeAppend)
}

func ModifyZip(zipFile string, module string) error {
	tempZipFile := zipFile + ".tmp"
	err := os.Rename(zipFile, tempZipFile)
	if err != nil {
		return err
	}

	err = copyZip(tempZipFile, zipFile, module)
	if err != nil {
		return err
	}

	return os.Remove(tempZipFile)
}

func copyZip(srcZip string, dstZip string, mod string) error {
	r, err := zip.OpenReader(srcZip)
	if err != nil {
		return err
	}
	defer r.Close()

	dstZipFile, err := os.Create(dstZip)
	if err != nil {
		return err
	}
	defer dstZipFile.Close()

	w := zip.NewWriter(dstZipFile)
	for _, file := range r.File {
		err = copyFile(w, file, mod)
		if err != nil {
			return err
		}
	}

	return w.Close()
}

func copyFile(w *zip.Writer, f *zip.File, mod string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	header := f.FileHeader
	header.Name = mod + header.Name[strings.Index(header.Name, "@"):]

	file, err := w.CreateHeader(&header)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, rc)
	return err
}
