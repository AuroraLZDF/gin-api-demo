package utils

import (
	"regexp"
	"log"
	"gopkg.in/go-playground/validator.v9"
)

//var validate *validator.Validate
var validate = validator.New()

func _validate(param interface{}, validator string) bool {
	err := validate.Var(param, validator)
	if err != nil {
		log.Println("validate error: " + err.Error())
		return false
	}
	return true
}

func Required(str string) bool {
	return _validate(str, "required")
}

func IsEmail(email string) bool {
	return _validate(email, "email")
}

func IsMobile(mobile string) bool {
	regular := `^1\d{10}$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobile)
}

func Equals(str1 string, str2 string) bool {
	/*if err = Required(str1); err != nil {
		return err
	}

	if err = Required(str2); err != nil {
		return err
	}

	if str1 != str2 {
		err = errors.New("两次输入不一致")
	}

	return*/

	return _validate(str1, "eq="+str2)
}

func NotEqual(str1 string, str2 string) bool {
	return _validate(str1, "ne="+str2)
}

func IsURL(url string) bool {
	return _validate(url, "url")
}

func IsNumeric(num int) bool {
	return _validate(num, "numeric")
}

func Max(num int, max int) bool {
	return _validate(num, "max="+string(max))
}

func Min(num int, max int) bool {
	return _validate(num, "min="+string(max))
}

func OneOf(value interface{}, arr []rune) bool {
	return _validate(value, "oneof="+string(arr))
}

func IsUnique(arr []rune) bool {
	return _validate(arr, "unique")
}

func IsFile(str string) bool {
	return _validate(str, "file")
}

func IsDir(str string) bool {
	return _validate(str, "dir")
}

func Contains(str string, substr string) bool {
	return _validate(str, "contains="+substr)
}

func Ip(str string) bool {
	return _validate(str, "ip")
}

func IpV4(str string) bool {
	return _validate(str, "ipv4")
}

func IpV6(str string) bool {
	return _validate(str, "ipv6")
}

func TcpAddr(str string) bool {
	return _validate(str, "tcp_addr")
}

func TcpV4Addr(str string) bool {
	return _validate(str, "tcp4_addr")
}

func TcpV6Addr(str string) bool {
	return _validate(str, "tcp6_addr")
}

func UdpAddr(str string) bool {
	return _validate(str, "udp_addr")
}

func UdpV4Addr(str string) bool {
	return _validate(str, "udp4_addr")
}

func UdpV6Addr(str string) bool {
	return _validate(str, "udp5_addr")
}

func IpAddr(str string) bool {
	return _validate(str, "ip_addr")
}

func IpV4Addr(str string) bool {
	return _validate(str, "ip4_addr")
}

func IpV6Addr(str string) bool {
	return _validate(str, "ip6_addr")
}
