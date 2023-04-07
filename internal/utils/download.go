package utils

import (
	"errors"
	"io"
	"os"
	"sync"
)

const (
	EOF = 1024 + iota
	NORMALLY_EXIST
)

var codeToStr = map[int]string{
	EOF:            io.EOF.Error(),
	NORMALLY_EXIST: "normally close",
}

func CodeToStr(code int) string {
	v, _ := codeToStr[code]
	return v
}

var NormallyCloseError = errors.New(CodeToStr(NORMALLY_EXIST))
var EndOfError = errors.New(CodeToStr(EOF))

type Downloader struct {
	waiter sync.WaitGroup
	mu     sync.Mutex

	exitChans []chan int
}

func (d *Downloader) WaitForDone() {
	d.waiter.Wait()
}

func (d *Downloader) CacelAllAndWaitForDone() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.SendToExitChans()

	d.waiter.Wait()

	d.freeAllExitChan()
}

func (d *Downloader) NewExitChan() chan int {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.newExitChan()
}

func (d *Downloader) FreeExitChan(c chan int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.freeExitChan(c)
}

func (d *Downloader) FreeAllExitChan() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.freeAllExitChan()
}

func (d *Downloader) SendToExitChans() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, _ := range d.exitChans {
		go func(i int) {
			d.exitChans[i] <- 1
		}(i)
	}
}

func (d *Downloader) newExitChan() chan int {
	if nil == d.exitChans {
		d.exitChans = make([]chan int, 0)
	}

	var chn = make(chan int, 1)

	d.exitChans = append(d.exitChans, chn)

	return chn
}

func (d *Downloader) freeExitChan(c chan int) {
	for i, chn := range d.exitChans {
		if chn == c {
			d.exitChans[i] = nil
			d.exitChans = append(d.exitChans[:i], d.exitChans[i+1:]...)
			break
		}
	}
}

func (d *Downloader) freeAllExitChan() {
	for i, _ := range d.exitChans {
		d.exitChans[i] = nil
	}
	d.exitChans = nil
}

func (d *Downloader) Download(dst string, f func(buffer []byte) error) error {
	d.waiter.Add(1)
	defer d.waiter.Done()

	exitChan := d.NewExitChan()
	defer func() {
		d.FreeExitChan(exitChan)
	}()

	var file, err = os.Open(dst)
	if nil != err {
		return err
	}
	defer file.Close()
	var maxSize = 1024
	var buffer = make([]byte, maxSize)
	var n int

	for {
		select {
		case <-exitChan:
			return NormallyCloseError
		default:
			n, err = file.Read(buffer)
			if nil != err {
				if io.EOF != err {
					return err
				}
			}
			if 0 < n {
				err = f(buffer[:n])
				if nil != err {
					return nil
				}
			}
		}
		if n < maxSize {
			return EndOfError
		}
	}
}
