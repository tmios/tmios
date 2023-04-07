package utils

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"tmios/lib/errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type PageReq struct {
	PageIndex int
	PageSize  int
	Query     map[string]interface{}
}

func ApiRet(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"Code":    0,
		"Detail":  "成功",
		"Content": data,
	})
}

func ApiErr(c *gin.Context, err error) {
	myErr, ok := err.(errors.Error)
	if ok {
		httpCode := myErr.Status
		c.JSON(httpCode, gin.H{
			"Code":    myErr.Code,
			"Detail":  myErr.Detail,
			"Content": myErr.Content,
		})
	} else {
		// TODO: this is only for develop
		fmt.Println("====================== internal:", err)

		log.WithField("Stack", string(debug.Stack())).
			WithError(err).Error("Internal error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"Code":    -1,
			"Detail":  fmt.Sprintf("%s", err),
			"Content": nil,
		})
	}
}

func Sort[T any](x []T, less func(a, b T) bool) {
	n := len(x)
	for {
		swapped := false
		for i := 1; i < n; i++ {
			if less(x[i], x[i-1]) {
				x[i-1], x[i] = x[i], x[i-1]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

func Uniq[T comparable](arr []T) []T {
	var (
		all []T
		m   = make(map[T]bool)
	)

	for _, a := range arr {
		m[a] = true
	}

	for k := range m {
		all = append(all, k)
	}

	return all
}

func Filter[T any](arr []T, predicate func(T) bool) []T {
	var ret []T
	for _, a := range arr {
		if predicate(a) {
			ret = append(ret, a)
		}
	}

	return ret
}

func Equal[T comparable](s1 []T, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for _, v1 := range s1 {
		eq := false
		for _, v2 := range s2 {
			if v1 == v2 {
				eq = true
				break
			}
		}

		if !eq {
			return false
		}
	}

	return true
}
