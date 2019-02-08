package hash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
)

func ReaderSHA256(r io.Reader) string {
	
	h := sha256.New()
	_, err := io.Copy(h, r)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))	
}
