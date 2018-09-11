package utils

import "log"

func CheckError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func CheckOk(ok bool, msg string) {
	if !ok {
		log.Fatal(msg)
	}
}
