package encrypt

import (
	"crypto/sha1"
	"fmt"
)

func SHA1(key string) string {
	data := []byte(key)
	return fmt.Sprintf("%x", sha1.Sum(data))
}
