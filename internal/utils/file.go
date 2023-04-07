package utils

import (
	"os"
	"strings"
)

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {

		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// RemoveFile 删除所给路径文件
func RemoveFile(path string) error {
	return os.Remove(path)
}

// PathName 获取路径中的文件名
func PathName(path string) string {
	p := path
	if strings.Contains(p, ":") {
		index := strings.LastIndex(p, ":")
		if len(p) > index+1 {
			p = p[index+1:]
		}
	}
	if strings.Contains(p, ".") {
		index := strings.LastIndex(p, ".")
		if len(p) > index+1 {
			p = p[:index]

		}
	}
	if strings.Contains(p, "/") {
		index := strings.LastIndex(p, "/")
		if len(p) > index+1 {
			p = p[index+1:]
		}
	}
	return p
}

// PathNameSuffix 文件名包括后缀
func PathNameSuffix(path string) string {
	p := path
	if strings.Contains(p, ":") {
		index := strings.LastIndex(p, ":")
		if len(p) > index+1 {
			p = p[index+1:]
		}
	}

	if strings.Contains(p, "/") {
		index := strings.LastIndex(p, "/")
		if len(p) > index+1 {
			p = p[index+1:]
		}
	}
	return p
}

func DirAndFileCreate(path string) (*os.File, error) {
	if flag, err := PathExists(path); err != nil {
		return nil, err
	} else {
		if flag {
			err := os.Remove(path)
			if err != nil {
				return nil, err
			}
		} else {
			pathName := PathNameSuffix(path)
			str := strings.ReplaceAll(path, pathName, "")
			err := os.MkdirAll(str, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
		create, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		return create, nil
	}
}
