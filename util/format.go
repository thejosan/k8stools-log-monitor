package util

import (
	"time"
	"strconv"
	"strings"
	"regexp"
)

func Int64toInt(num int64) int{
	strInt64 := strconv.FormatInt(num, 10)
	int ,_ := strconv.Atoi(strInt64)
	return int
}

func SliceToString(Slice []string) string {
	if Slice != nil {
		str := "^"
		for _, v := range Slice {
			str = str + "," + v
		}
		str = strings.Replace(str, "^,", "", 1)
		return str
	} else {
		return ""
	}
}

func StringTOSlice(str string) []string {
	strarr := strings.Split(str, ",")
	return strarr
}

func StringTOSlicef(str string,strf string) []string {
	strarr := strings.Split(str, strf)
	return strarr
}

func Checkids(str string) string {
	exp := regexp.MustCompile(`\d+(,\d+)*`)
	result2 := exp.FindAllString(str, -1)
	result := SliceToString(result2)
	return result
}

func GetNowTime() string{
	nowTime := time.Now()
	t := nowTime.String()
	timeStr := t[:19]
	return timeStr
}

func Show_substr(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl+rl > l {
			break
		}
		sl += rl
		ss += string(r)
	}
	return ss
}