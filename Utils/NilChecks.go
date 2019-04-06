package utils

/*ErrorNilCheck returns error messsage if error exists*/
func ErrorNilCheck(error error) {
	if error != nil {
		println(error.Error())
		return
	}
	return
}

/*ErrorInvalidCheck returns error messsage if error exists*/
func ErrorInvalidCheck(error error) {
	if error != nil {
		println(error.Error())
		return
	}
	return
}
