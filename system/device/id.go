package device

import "strings"

type Id string

func NewId(prefix string, entity string) Id {
	return Id(prefix + "::" + entity)
}

func (d Id) Prefix() string {
	p := strings.Split(string(d), "::")
	if len(p) > 0 {
		return p[0]
	}
	return ""
}

func (d Id) Entity() string {
	p := strings.Split(string(d), "::")
	if len(p) > 1 {
		return p[1]
	}
	return ""
}
