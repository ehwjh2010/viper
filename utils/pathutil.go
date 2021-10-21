package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

//EnsurePathExist 确认文件或文件夹是否存在
func EnsurePathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

//MakeDir 创建单一目录, 不支持创建多级目录
//@param path 路径
//@param exist_no_error 路径已存在时是否返回错误
func MakeDir(path string, existReturnError bool) (bool, error) {
	if IsEmptyStr(path) {
		return false, errors.New("invalid path")
	}

	exist, err := EnsurePathExist(path)
	if err != nil {
		return false, err
	}

	if exist {
		if existReturnError {
			return false, errors.New(fmt.Sprintf("%s had exist!!!", path))
		} else {
			return true, nil
		}
	} else {
		err := os.Mkdir(path, 0777)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

//RemovePath 完全删除文件夹或文件, 对于文件夹包括子文件以及子文件夹
//@param path 路径
//@param noExistReturnError 路径不存在时是否返回错误
func RemovePath(path string, noExistReturnError bool) (bool, error) {
	if IsEmptyStr(path) {
		return false, errors.New("invalid path")
	}

	exist, err := EnsurePathExist(path)
	if err != nil {
		return false, err
	}

	if !exist {
		if noExistReturnError {
			return false, errors.New("path no exist")
		} else {
			return true, nil
		}

	} else {
		err := os.RemoveAll(path)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}

//PathJoin 路径拼接
func PathJoin(paths ...string) string {
	p := filepath.Join(paths...)
	return p
}

//RelativePath2AbsPath 相对路径转化绝对路径
func RelativePath2AbsPath(relativePath string) (string, error) {
	absPath, err := filepath.Abs(relativePath)

	if err != nil {
		return "", err
	} else {
		return absPath, nil
	}
}

//PathSplit 路径分割, 返回目录以及文件名
func PathSplit(path string) (string, string) {
	dirName := filepath.Dir(path)
	fileName := filepath.Base(path)
	return dirName, fileName
}