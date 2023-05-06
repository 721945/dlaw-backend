package utils

import "fmt"

func FormatFileSize(size int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	var i int
	var f float64 = float64(size)
	for i = 0; i < len(units) && f >= 1024; i++ {
		f /= 1024
	}
	return fmt.Sprintf("%.2f %s", f, units[i])
}
