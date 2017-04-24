package main

func failedIfError(err error) {
	println(err)

	if err != nil {
		panic(err)
	}
}

// InList for Check elem is exist in list
func InList(lst []string, elem string) bool {
	for _, x := range lst {
		if elem == x {
			return true
		}
	}

	return false
}
