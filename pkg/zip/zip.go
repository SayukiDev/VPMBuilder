package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Compress(inputPath string, outputPath string) (err error) {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()
	ss, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if ss.IsDir() {
		return compressDir(w, inputPath)
	}
	err = addFile(w, filepath.Base(inputPath), inputPath)
	return nil
}

func compressDir(writer *zip.Writer, path string) (err error) {
	return filepath.WalkDir(path, func(newP string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		realPath := newP
		if len(newP) == len(path) {
			return nil
		}

		newP = trimRootPath(newP, path)
		err = addFile(writer, newP, realPath)
		if err != nil {
			return err
		}
		return nil
	})
}

func trimRootPath(path string, rootPath string) string {
	return filepath.ToSlash(filepath.Clean(path)[len(rootPath)+1:])
}

func addFile(writer *zip.Writer, path string, realPath string) (err error) {
	fi, err := os.Stat(realPath)
	if err != nil {
		return err
	}
	hd, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}

	hd.Name = path
	if fi.IsDir() {
		hd.Name += "/"
		hd.Method = zip.Store
		_, err := writer.CreateHeader(hd)

		return err
	}
	hd.Method = zip.Deflate

	f, err := os.Open(realPath)
	if err != nil {
		return err
	}
	defer f.Close()
	fileWriter, err := writer.CreateHeader(hd)
	_, err = io.Copy(fileWriter, f)
	return err
}
