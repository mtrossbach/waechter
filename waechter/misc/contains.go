package misc

func ContainsAll(s []string, e []string) bool {
	for _, ee := range e {
		if !Contains(s, ee) {
			return false
		}
	}
	return true
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func StringsMissingInList(baseList []string, searchList []string) []string {
	if len(baseList) == 0 {
		return searchList
	}
	if len(searchList) == 0 {
		return []string{}
	}
	var result []string
	for _, s := range searchList {
		found := false
		for _, b := range baseList {
			if b == s {
				found = true
				break
			}
		}
		if !found {
			result = append(result, s)
		}
	}
	return result
}
