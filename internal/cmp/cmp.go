package cmp

import (
	"github.com/sirupsen/logrus"
)

// Cmp 组件处理中心
type Cmp struct {
	srv []Srv
}

// NewCmp conf第一个，阻塞的http在最后
func NewCmp(srv ...Srv) *Cmp {
	var arr []Srv
	for _, v := range srv {
		arr = append(arr, v)
	}
	return &Cmp{arr}
}

func (cmp *Cmp) Run() error {
	for _, v := range cmp.srv {
		err := v.Run()
		if err != nil {
			logrus.Fatal(err)
			return err
		}
	}
	return nil
}
