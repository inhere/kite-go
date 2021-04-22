package pacutil

import (
	"encoding/base64"
	"io/ioutil"
)

func DecodeGfwList(gwfile string) (dst []byte, err error) {
	in, err := ioutil.ReadFile(gwfile)
	if err != nil {
		return
	}

	dst = make([]byte, base64.StdEncoding.DecodedLen(len(in)))

	_, err = base64.StdEncoding.Decode(dst, in)
	return
}
