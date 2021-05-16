package idetool

import "io/ioutil"

func LoadHceFile(hceFile string) (bts []byte, err error) {
	return ioutil.ReadFile(hceFile)
}
