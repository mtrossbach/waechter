package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"path"
	"regexp"
)

func File() string {
	return file
}

func Dir() string {
	return path.Dir(file)
}

func Print() {

	fmt.Printf("########################################\n")

	data, _ := yaml.Marshal(*instance)
	str := string(data)

	str = regexp.MustCompile(`token: (\S{3,})`).ReplaceAllString(str, "token: [set]")
	str = regexp.MustCompile(`password: (\S{3,})`).ReplaceAllString(str, "password: [set]")
	str = regexp.MustCompile(`pin: (\S{3,})`).ReplaceAllString(str, "pin: [set]")

	fmt.Printf("%s", str)
	fmt.Printf("########################################\n")
}
