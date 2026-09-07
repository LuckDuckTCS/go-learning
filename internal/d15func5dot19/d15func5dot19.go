package d15func5dot19

func triple(i int) (result int) {
	defer func() {
		switch p := recover(); p {
		case nil:

		case "Something":
			result = i * 3
		default:
			panic(p)
		}
	}()
	double(i)
	return -1
}

func double(i int) {
	panic("Something")
}
