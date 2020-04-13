package demo

func F(str string) string {
	rtn1 := F1(str)
	rtn2 := F2(str)
	if len(rtn1) > len(rtn2) {
		return rtn1
	} else {
		return rtn2
	}
}

func F2(str string) string {
	if len(str) < 2 {
		return ""
	}
	p, q := 0, 0
	for i := 0; i < len(str); i++ {
		for j := 0; i-j >= 0 && i+j+1 < len(str); j++ {
			if str[i-j] != str[i+j+1] {
				break
			}
			if 2*j+2 > q-p {
				p = i - j
				q = i + j + 1
			}
		}
	}
	return str[p : q+1]
}

func F1(str string) string {
	p, q := 0, 0
	for i := 0; i < len(str); i++ {
		for j := 0; i-j >= 0 && i+j < len(str); j++ {
			if str[i-j] != str[i+j] {
				break
			}
			if j*2 > q-p {
				p = i - j
				q = i + j
			}
		}
	}
	return str[p : q+1]
}
