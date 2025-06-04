package utils

func HandleErr(err error, desc string) {
	if err != nil {
		panic(desc + ": " + err.Error())
	}
}
