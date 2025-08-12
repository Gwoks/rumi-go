package utils

import "time"

func GetJakartaLocationTime() *time.Location {
	jktTimeLoc, _ := time.LoadLocation("Asia/Jakarta")
	return jktTimeLoc
}
