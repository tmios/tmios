package utils

import "time"

// Retry 重试,返回错误退出
func Retry(times time.Duration, fnc func() error) error {
	for {
		err := fnc()
		if err != nil {
			return err
		}
		time.Sleep(times)
	}
}
