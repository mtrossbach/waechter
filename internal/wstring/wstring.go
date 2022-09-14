package wstring

func StrDef(str1 string, str2 string) string {
	if len(str1) == 0 {
		return str2
	}
	return str1
}
