package utils

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func ParseTime(timeStr string) (*time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}

	tt, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	if err != nil {
		return nil, err
	}
	return &tt, nil
}

func ParseTimeDurationToStr(timeStr string, duration time.Duration) (string, error) {
	timeStr = strings.ReplaceAll(timeStr, "T", " ")
	timeStr = strings.ReplaceAll(timeStr, "Z", "")
	parse, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return "", err
	}
	t := parse.Add(duration)
	return t.Format("2006-01-02T15:04:05"), nil
}

func ParseTimeToCTS(clientTime string, nextTime string) string {
	if clientTime != "" {
		ctimeArr := strings.Split(clientTime, "+")
		if len(ctimeArr) == 2 {
			ctimeIndex := strings.Index(ctimeArr[1], ":")
			ctimeAdd := ctimeArr[1][:ctimeIndex]
			atoi, err := strconv.Atoi(ctimeAdd)
			if err != nil {
				logrus.Error(err)
			} else {
				next, err := ParseTimeDurationToStr(nextTime, time.Hour*time.Duration(atoi))
				if err != nil {
					logrus.Error(err)
				}
				return next
			}
		}
	}
	return nextTime
}
