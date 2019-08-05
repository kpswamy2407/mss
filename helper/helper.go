package helper

import (
	"crypto/md5"
	"fmt"
	"io"
)

//GenerateMd5 is used to generate the md5 hash for given string
func GenerateMd5(str string) (hash string) {
	h := md5.New()
	io.WriteString(h, str)
	md5 := h.Sum(nil)
	return fmt.Sprintf("%x", string(md5))

}
