package utils

import (
	"io"
	"os"
	"path/filepath"
)

type FileItem struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Ext      string     `json:"ext"`
	IsDir    bool       `json:"isDir"`
	Size     int64      `json:"size"`
	Time     int64      `json:"time"`
	Children []FileItem `json:"children"`
}

func FileInfo(filename string) (FileItem, error) {
	info := FileItem{}
	file, err := os.Open(filename)
	if err != nil {
		return FileItem{}, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	info.Name = fileInfo.Name()
	info.IsDir = fileInfo.IsDir()
	info.Path = filename
	if !info.IsDir {
		info.Ext = filepath.Ext(info.Name)
	}
	info.Size, _ = file.Seek(0, io.SeekEnd)
	info.Time = fileInfo.ModTime().Unix()
	info.Children = []FileItem{}
	return info, nil
}

func FileCatalog(filename string, fn func(item *FileItem) bool) ([]FileItem, error) {
	files, err := os.ReadDir(filename)
	if err != nil {
		return nil, err
	}
	list := []FileItem{}
	for _, file := range files {
		info, _ := file.Info()
		name := info.Name()
		path := filename + string(filepath.Separator) + name
		isDir := info.IsDir()
		ext := ""
		if !isDir {
			ext = filepath.Ext(name)
		}
		val := FileItem{
			Name:     name,
			Path:     path,
			Ext:      ext,
			IsDir:    isDir,
			Size:     info.Size(),
			Time:     info.ModTime().Unix(),
			Children: []FileItem{},
		}
		bl := fn(&val)
		if bl && isDir {
			child, _ := FileCatalog(path, fn)
			val.Children = child
		}
		list = append(list, val)
	}
	return list, nil
}
