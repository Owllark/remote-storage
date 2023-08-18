package file_system

import (
	"errors"
	"fmt"
	"os"
	"remote-storage/common"
	"strings"
)

const pathSeparator = string(os.PathSeparator)

type FileAction func(info *common.FileInfo) error

type FileSystemBrowser struct {
	currentPath    string
	currentNesting []string
	prevPaths      []string
	root           common.FileInfo
}

func NewFileSystemBrowser(root common.FileInfo) *FileSystemBrowser {
	var res FileSystemBrowser
	res.root = root
	res.currentPath = ""
	res.currentNesting = res.currentNesting[:0]
	res.prevPaths = res.prevPaths[:0]
	return &res

}

func (fs *FileSystemBrowser) Reset(root common.FileInfo) {
	fs.root = root
	fs.currentPath = ""
	fs.currentNesting = fs.currentNesting[:0]
	fs.prevPaths = fs.prevPaths[:0]
}

func (fs *FileSystemBrowser) Cd(path string) error {
	var err error
	switch path {
	case "..":
		{
			if len(fs.currentNesting) != 0 {
				fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentNesting))
				fs.currentNesting = fs.currentNesting[:len(fs.currentNesting)-1]
				fs.currentPath = joinPath(fs.currentNesting)
			}
		}
	case "/":
		{
			fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentNesting))
			fs.currentNesting = fs.currentNesting[:0]
			fs.currentPath = joinPath(fs.currentNesting)
		}
	case "":
		{
			fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentNesting))
			fs.currentNesting = fs.currentNesting[:0]
			fs.currentPath = joinPath(fs.currentNesting)
		}
	case "-":
		{
			if len(fs.prevPaths) != 0 {
				fs.currentNesting = splitPath(fs.prevPaths[len(fs.prevPaths)-1])
				fs.prevPaths = fs.prevPaths[:len(fs.prevPaths)]
				fs.currentPath = joinPath(fs.currentNesting)
			}
		}
	default:
		{
			if fs.isDirectoryExists(splitPath(fs.currentPath + path)) {
				fs.prevPaths = append(fs.prevPaths, joinPath(fs.currentNesting))
				fs.currentNesting = append(fs.currentNesting, splitPath(path)...)
				fs.currentPath = joinPath(fs.currentNesting)
			} else {
				err = errors.New(fmt.Sprintf("directory %s not found", path))
			}

		}

	}
	return err
}

func (fs *FileSystemBrowser) Ls(dirPath string) ([]string, error) {
	var res []string
	var dir *common.FileInfo
	if (fs.currentPath + dirPath) == pathSeparator {
		dir = &fs.root
	} else {
		dir = fs.FindFileByPath(splitPath(fs.currentPath + dirPath))
	}
	if dir == nil {
		err := errors.New(fmt.Sprintf("%s not found", dirPath))
		return nil, err
	}
	if len(dir.Children) == 0 {
		res = append(res, "Empty")
	} else {
		topLine := "Name\tIs Directory\tSize\tModified"
		res = append(res, topLine)
		for _, file := range dir.Children {
			s := fmt.Sprintf("%s\t\t%t\t%d\t%s", file.Name, file.IsDir, file.Size, file.Modified)
			res = append(res, s)
		}
	}
	return res, nil
}

func joinPath(path []string) string {
	var res string
	for _, s := range path {
		res += s + string(os.PathSeparator)
	}
	return res
}

func splitPath(path string) []string {
	var res []string
	parts := strings.Split(path, string(os.PathSeparator))
	for _, part := range parts {
		if part != "" {
			res = append(res, part)
		}
	}

	return res
}

func (fs *FileSystemBrowser) isDirectoryExists(path []string) bool {
	curDir := &fs.root
	for _, dir := range path {
		fileFound := findFileInDirectory(curDir, dir)
		if fileFound == nil {
			return false
		}
		curDir = fileFound
	}
	return curDir.IsDir
}

func (fs *FileSystemBrowser) FindFileByPath(path []string) *common.FileInfo {
	curFile := &fs.root
	for _, dir := range path {
		fileFound := findFileInDirectory(curFile, dir)
		if fileFound == nil {
			return nil
		}
		curFile = fileFound
	}
	return curFile
}

func (fs *FileSystemBrowser) SetFiles(files common.FileInfo) {
	fs.root = files
}

func (fs *FileSystemBrowser) GetCurrentPath() string {
	return fs.currentPath
}

func findFileInDirectory(fileInfo *common.FileInfo, targetName string) *common.FileInfo {

	// Iterate over the children recursively to find the file
	for _, child := range fileInfo.Children {
		if child.Name == targetName {
			// If the child is a file and its name matches the target name, return it
			return &child
		}
	}

	return nil
}

func findFileByName(fileInfo common.FileInfo, targetName string) *common.FileInfo {
	// Check if the current file matches the target name
	if fileInfo.Name == targetName {
		return &fileInfo
	}

	// Iterate over the children recursively to find the file
	for _, child := range fileInfo.Children {
		if child.IsDir {
			// If the child is a directory, recursively search within it
			result := findFileByName(child, targetName)
			if result != nil {
				return result
			}
		} else if child.Name == targetName {
			// If the child is a file and its name matches the target name, return it
			return &child
		}
	}

	// File not found
	return nil
}

func (fs *FileSystemBrowser) ForEveryFile(path string, action FileAction) error {
	pathList := splitPath(path)
	firstFile := fs.FindFileByPath(pathList)
	err := action(firstFile)
	if err != nil {
		return err
	}

	for _, child := range firstFile.Children {
		err := action(&child)
		if err != nil {
			return err
		}
		if child.IsDir {
			err := fs.ForEveryFile(joinPath(pathList[1:]), action)
			if err != nil {
				return err
			}
		}
	}

	// File not found
	return nil
}
