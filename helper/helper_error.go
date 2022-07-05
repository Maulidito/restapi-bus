package helper

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}

}

func PanicIfErrorString(err string) {
	if err != "" {
		panic(err)
	}
}
