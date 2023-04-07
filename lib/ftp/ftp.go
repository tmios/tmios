package ftp

import (
	"bytes"
	"log"
	"time"

	"tmios/lib/errors"

	"github.com/jlaffaye/ftp"
)

var (
	ErrPathExists = errors.Conflict(600100, "该路径已存在文件名")
)

type Server struct {
	FtpUrl   string
	Username string
	Password string
	Timeout  time.Duration
}

type File struct {
	RelativePath string
	Filename     string
	Content      string
}

func NewClient(server *Server) (*ftp.ServerConn, error) {
	c, err := ftp.Dial(server.FtpUrl, ftp.DialWithTimeout(server.Timeout))
	if err != nil {
		return nil, err
	}

	if err := c.Login(server.Username, server.Password); err != nil {
		return nil, err
	}

	return c, nil
}

func getEntry(c *ftp.ServerConn, relativePath string) *ftp.Entry {
	entries, err := c.List("/")
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if entry.Name == relativePath {
			return entry
		}
	}

	return nil
}

func GetEntry(server *Server, relativePath string) (*ftp.Entry, error) {
	c, err := NewClient(server)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := c.Quit(); err != nil {
			log.Fatal(err)
		}
	}()

	entry := getEntry(c, relativePath)

	return entry, nil
}

func safeMkdirAll(c *ftp.ServerConn, relativePath string) error {
	entry := getEntry(c, relativePath)
	if entry != nil {
		if entry.Type == ftp.EntryTypeFolder {
			return nil
		} else {
			return ErrPathExists
		}
	}

	if err := c.MakeDir(relativePath); err != nil {
		return err
	}

	return nil
}

func SafeMkdirAll(server *Server, relativePath string) error {
	c, err := NewClient(server)
	if err != nil {
		return err
	}

	defer func() {
		if err := c.Quit(); err != nil {
			log.Fatal(err)
		}
	}()

	return safeMkdirAll(c, relativePath)
}

func putContent(c *ftp.ServerConn, relativePath, filename, content string) error {
	if err := c.ChangeDir("/"); err != nil {
		return err
	}

	if err := safeMkdirAll(c, relativePath); err != nil {
		return err
	}

	if err := c.ChangeDir(relativePath); err != nil {
		return err
	}

	data := bytes.NewBufferString(content)
	if err := c.Stor(filename, data); err != nil {
		return err
	}
	return nil
}

func PutFiles(server *Server, files []*File) error {
	c, err := NewClient(server)
	if err != nil {
		return err
	}

	defer func() {
		if err := c.Quit(); err != nil {
			log.Fatal(err)
		}
	}()

	for _, file := range files {
		if err := putContent(c, file.RelativePath, file.Filename, file.Content); err != nil {
			return err
		}
	}

	return nil
}
