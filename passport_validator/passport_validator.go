package passport_validator

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	// Проверяем что серия паспорта состоит из 4 цифр
	passportSeriesRegexp = regexp.MustCompile(`^\d{4}$`)
	// Проверяем что номер паспорта состоит из 6 цифр
	passportNumberRegexp = regexp.MustCompile(`^\d{6}$`)
	// Проверяем что код подразделения состоит из 6 цифр возможно с - (дефисом)
	issuedCodeRegexp = regexp.MustCompile(`^\d{3}\-?\d{3}$`)
)

var (
	ErrEmptyLastName                    = errors.New("last name field is empty")
	ErrEmptyFirstName                   = errors.New("first name field is empty")
	ErrEmptyPassportSeries              = errors.New("passport series is empty")
	ErrInvalidPassportSeries            = errors.New("last 2 digits in passport series is before 1997 or after now year")
	ErrInvalidPassportSeriesNot4Digits  = errors.New("passport series is not 4 digits")
	ErrEmptyPassportNumber              = errors.New("passport number is empty")
	ErrInvalidPassportNumber            = errors.New("passport number is not 6 digits")
	ErrInvalidIssueDateBefore14Birthday = errors.New("issued date before fourteenth birthday")
	ErrEmptyIssueDate                   = errors.New("issued date is zero")
	ErrEmptyBirthday                    = errors.New("birthday is zero")
	ErrEmptyIssuedCode                  = errors.New("issued code is empty")
	ErrInvalidIssuedCode                = errors.New("wrong issued code format")
	ErrInvalidBirthday                  = errors.New("birthday is zero or in the future")
	ErrNonCyrillicCharacter             = errors.New("contains non-Cyrillic symbol")
	ErrIssueDatePassportInFuture        = errors.New("passport issued in future")
	ErrIssueDatePassportExpiredAt20     = errors.New("passport expired at 20")
	ErrIssueDatePassportExpiredAt45     = errors.New("passport expired at 45")
)

const (
	Age20PassportChange  = 20
	Age45PassportChange  = 45
	PassportDaysValidity = 91
)

func IsPassportLastNameValid(lastName string) error {
	if lastName == "" {
		return ErrEmptyLastName
	}
	return nameValidator(lastName)
}

func IsPassportFirstNameValid(firstName string) error {
	if firstName == "" {
		return ErrEmptyFirstName
	}
	return nameValidator(firstName)
}

func IsPassportMiddleNameValid(middleName string) error {
	// Не проверяем на пустоту, так как отчества может не быть
	return nameValidator(middleName)
}

func nameValidator(s string) error {
	allowedCharacters := map[rune]bool{
		'-':  true,
		' ':  true,
		'.':  true,
		',':  true,
		'I':  true,
		'V':  true,
		'\'': true,
		'(':  true,
		')':  true,
	}
	for _, c := range s {
		if !unicode.Is(unicode.Cyrillic, c) && !allowedCharacters[c] {
			return ErrNonCyrillicCharacter
		}
	}
	return nil
}

func IsPassportSeriesValid(series string, checkDate time.Time) error {
	if series == "" {
		return ErrEmptyPassportSeries
	}

	// Примитивная проверка: 4 цифры серии паспорта РФ
	if !passportSeriesRegexp.MatchString(series) {
		return ErrInvalidPassportSeriesNot4Digits
	}

	// Проверяем, что паспорт выдан после 1997 года, но не больше чем текущий год +5 лет.
	// Современный бланк появился в 1997 году. Две последних цифры паспорта должны быть в диапазоне 97-99/00-time.Now().Year()%100+5
	// Иногда квота на паспорта заканчивается поэтому печатают в счет будуших квот
	blankReleaseDate := time.Date(1997, time.January, 01, 00, 0, 0, 0, time.UTC)
	nowYear := checkDate.Year()

	issueYear := series[len(series)-2:]
	issueYearInt, err := strconv.Atoi(issueYear)
	if err != nil {
		return ErrInvalidPassportSeriesNot4Digits
	}

	if issueYearInt <= nowYear%100+5 {
		issueYearInt = issueYearInt + 2000
	} else {
		issueYearInt = issueYearInt + 1900
	}

	if issueYearInt < blankReleaseDate.Year() || issueYearInt > nowYear+5 {
		return ErrInvalidPassportSeries
	}

	return nil
}

func IsPassportNumberValid(number string) error {
	if number == "" {
		return ErrEmptyPassportNumber
	}

	// Примитивная проверка: 6 цифр серии паспорта РФ
	if !passportNumberRegexp.MatchString(number) {
		return ErrInvalidPassportNumber
	}

	return nil
}

func IsPassportIssueDateValid(issueDate, birthday time.Time, checkDate time.Time) error {
	if issueDate.IsZero() {
		return ErrEmptyIssueDate
	}

	if birthday.IsZero() {
		return ErrEmptyBirthday
	}

	ageAtIssue := issueDate.Year() - birthday.Year()
	// Проверка на возможность выдачи паспорта до 14 лет
	if ageAtIssue < 14 {
		return ErrInvalidIssueDateBefore14Birthday
	}

	// Проверка выдачи паспорта в будущем
	if issueDate.After(checkDate) {
		return ErrIssueDatePassportInFuture
	}

	// Проверка выдачи паспорта до 20 лет
	if ageAtIssue < 20 {
		expirationDate := birthday.AddDate(Age20PassportChange, 0, PassportDaysValidity)
		if checkDate.After(expirationDate) {
			return ErrIssueDatePassportExpiredAt20
		}
	}

	// Проверка выдачи паспорта до 45 лет
	if ageAtIssue < 45 {
		expirationDate := birthday.AddDate(Age45PassportChange, 0, PassportDaysValidity)
		if checkDate.After(expirationDate) {
			return ErrIssueDatePassportExpiredAt45
		}
	}

	return nil
}

func IsPassportIssuerCodeValid(issuedCode string) error {
	if issuedCode == "" {
		return ErrEmptyIssuedCode
	}

	// Проверяем что код подразделения состоит из 6 цифр разделенных - (дефисом), если он есть
	if !issuedCodeRegexp.MatchString(issuedCode) {
		return ErrInvalidIssuedCode
	}

	return nil
}

func IsPassportBirthdayValid(birthday time.Time, checkDate time.Time) error {
	if birthday.IsZero() {
		return ErrInvalidBirthday
	}
	if birthday.After(checkDate) {
		return ErrInvalidBirthday
	}
	// Если меньше 18 лет
	if birthday.After(checkDate.AddDate(-18, 0, 0)) {
		return ErrInvalidBirthday
	}

	return nil
}

// PassportPlaceOfBirthNormalize двойные пробелы и дефисы заменяются на одинарные; пробелы слева и справа от дефиса удаляются "US-  -SR"->"USSR". SHOCK-11301
func PassportPlaceOfBirthNormalize(placeOfBirth string) string {
	if placeOfBirth == "" {
		return ""
	}

	replaces := map[string]string{
		"  ": " ",
		"--": "-",
		" -": "-",
		"- ": "-",
	}

	findReplace := true
	for findReplace {
		findReplace = false
		for from, to := range replaces {
			placeOfBirthNew := strings.Replace(placeOfBirth, from, to, -1)
			if placeOfBirthNew != placeOfBirth {
				findReplace = true
				placeOfBirth = placeOfBirthNew
			}
		}
	}
	return placeOfBirth
}

// PassportIssuedByNormalize двойные пробелы, знаки препинания, разделители, кавычки заменяются на одинарные символы
func PassportIssuedByNormalize(issuedBy string) string {
	if issuedBy == "" {
		return ""
	}

	replaces := map[string]string{
		"  ": " ",
		"..": ".",
		",,": ",",
		`""`: `"`,
		`''`: `'`,
		//TODO что то еще?
	}

	findReplace := true
	for findReplace {
		findReplace = false
		for from, to := range replaces {
			issuedByNew := strings.Replace(issuedBy, from, to, -1)
			if issuedByNew != issuedBy {
				findReplace = true
				issuedBy = issuedByNew
			}
		}
	}
	return issuedBy
}
