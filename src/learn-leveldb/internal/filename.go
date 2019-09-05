package internal

import (
	"fmt"
)

func makeFileName(dbname string, number uint64, suffix string) string {
	return fmt.Sprintf("%s/%06d.%s", dbname, number, suffix)
}

func TableFileName(dbname string, number uint64) string {
	return makeFileName(dbname, number, "ldb")
}

func DescriptorFileName(dbname string, number uint64) string {
	return fmt.Sprintf("%s/MANIFEST-%06d", dbname, number)
}

func CurrentFileName(dbname string) string {
	return dbname + "/CURRENT"
}
func TempFileName(dbname string, number uint64) string {
	return makeFileName(dbname, number, "dbtmp")
}
