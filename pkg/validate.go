package pkg

import (
	"regexp"
	"time"
)

const REGEX_PHONE = `^\d{9}$`
const REGEX_EMAIL = `^[\w\.-]+@[\w\.-]+\.\w{2,}$`
const REGEX_RUC = `^\d{11}$`

func ValidatePhone(phone string) bool {
	return regexp.MustCompile(REGEX_PHONE).MatchString(phone)
}

func ValidateEmail(email string) bool {
	return regexp.MustCompile(REGEX_EMAIL).MatchString(email)
}

func ValidateDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func ValidateRUC(ruc string) bool {
	return regexp.MustCompile(REGEX_RUC).MatchString(ruc)
}
