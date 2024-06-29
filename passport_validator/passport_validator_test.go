package passport_validator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PassportLastName(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		lastName string
		wantErr  error
	}{
		"valid last name": {
			lastName: "Маск",
		},
		"valid last name with V": {
			lastName: "МаскV",
		},
		"valid last name with ,": {
			lastName: "Ма,скV",
		},
		"valid last name with '": {
			lastName: "Д'Артаньян",
		},
		"valid last name with .": {
			lastName: "Д.Артаньян",
		},
		"valid last name with (V)": {
			lastName: "Д'Артаньян(V)",
		},
		"valid double last name": {
			lastName: "Маск-Петрович",
		},
		"with whitespace last name": {
			lastName: "Ривейро И Ламасарес",
		},

		"blank last name": {
			lastName: "",
			wantErr:  ErrEmptyLastName,
		},
		"not cyrillic last name": {
			lastName: "musk",
			wantErr:  ErrNonCyrillicCharacter,
		},
		"пишем значение с арабской цифрой «4»": {
			lastName: "Александрови4",
			wantErr:  ErrNonCyrillicCharacter,
		},
		"недопустимая буква латинского алфавита": {
			lastName: "Иванов XV",
			wantErr:  ErrNonCyrillicCharacter,
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := IsPassportLastNameValid(tt.lastName)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			require.NoError(t, err)

		})
	}
}

func Test_PassportFirstName(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		firstName string
		wantErr   error
	}{
		"valid first name with '": {
			firstName: "Д'Артаньян",
		},
		"valid first name with .": {
			firstName: "Д.Артаньян",
		},
		"valid first name with (V)": {
			firstName: "Д'Артаньян(V)",
		},
		"with whitespace first name": {
			firstName: "Ривейро И Ламасарес",
		},
		"valid first name": {
			firstName: "илон",
		},
		"valid double first name": {
			firstName: "илон-максим",
		},
		"valid with whitespace first name": {
			firstName: "ил он",
		},
		"blank first name": {
			firstName: "",
			wantErr:   ErrEmptyFirstName,
		},
		"not cyrillic first name": {
			firstName: "Elon",
			wantErr:   ErrNonCyrillicCharacter,
		},
		"пишем значение с арабской цифрой «4»": {
			firstName: "Александрови4",
			wantErr:   ErrNonCyrillicCharacter,
		},
		"недопустимая буква латинского алфавита": {
			firstName: "Иванов XV",
			wantErr:   ErrNonCyrillicCharacter,
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := IsPassportFirstNameValid(tt.firstName)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			require.NoError(t, err)

		})
	}
}

func Test_PassportMiddleName(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		middleName string
		wantErr    error
	}{
		"valid middle name with '": {
			middleName: "Д'Артаньян",
		},
		"valid middle name with .": {
			middleName: "Д.Артаньян",
		},
		"valid middle name": {
			middleName: "Иванович",
		},
		"valid double middle name": {
			middleName: "Иванович-Маскович",
		},
		"valid with whitespace middle name": {
			middleName: "Иван ович",
		},
		"blank middle name": {
			middleName: "",
		},
		"not cyrillic middle name": {
			middleName: "Ivanovich",
			wantErr:    ErrNonCyrillicCharacter,
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := IsPassportMiddleNameValid(tt.middleName)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			require.NoError(t, err)

		})
	}
}

func Test_PassportSeries(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		series    string
		checkDate time.Time
		wantErr   error
	}{
		"valid series 2017": {
			series:    "4617",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
		"valid series 2000": {
			series:    "4600",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
		"valid series 1999": {
			series:    "4699",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
		"valid series 2026": {
			series:    "4626",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
		"valid series 2029 (corner case with quote)": {
			series:    "4629",
			checkDate: time.Date(2024, 02, 22, 0, 0, 0, 0, time.UTC),
		},
		"blank series": {
			series:    "",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrEmptyPassportSeries,
		},
		"5 digests on series": {
			series:    "46179",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidPassportSeriesNot4Digits,
		},
		"3 digests on series": {
			series:    "461",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidPassportSeriesNot4Digits,
		},
		"not digests on series": {
			series:    "lol",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidPassportSeriesNot4Digits,
		},
		"series before 1997": {
			series:    "4696",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidPassportSeries,
		},
		"series with V": {
			series:    "408V",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidPassportSeriesNot4Digits,
		},
		"series with V2": {
			series:    "4V23",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidPassportSeriesNot4Digits,
		},
		"series with V3": {
			series:    "V4323",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidPassportSeriesNot4Digits,
		},
		"series with Л": {
			series:    "408Л",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidPassportSeriesNot4Digits,
		},
		"series with Л2": {
			series:    "4Л08",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidPassportSeriesNot4Digits,
		},
		"series with Л3": {
			series:    "Л4308",
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidPassportSeriesNot4Digits,
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := IsPassportSeriesValid(tt.series, tt.checkDate)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			require.NoError(t, err)

		})
	}
}

func Test_PassportNumber(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		number  string
		wantErr error
	}{
		"valid number": {
			number: "657482",
		},
		"blank number": {
			number:  "",
			wantErr: ErrEmptyPassportNumber,
		},
		"not valid number not digests": {
			number:  "Ivanov",
			wantErr: ErrInvalidPassportNumber,
		},
		"not valid number 7 digests": {
			number:  "6574822",
			wantErr: ErrInvalidPassportNumber,
		},
		"not valid number 3 digests": {
			number:  "652",
			wantErr: ErrInvalidPassportNumber,
		},
		"not valid with V": {
			number:  "34215V",
			wantErr: ErrInvalidPassportNumber,
		},
		"not valid with V2": {
			number:  "34V215",
			wantErr: ErrInvalidPassportNumber,
		},
		"not valid with V3": {
			number:  "V657482",
			wantErr: ErrInvalidPassportNumber,
		},
		"not valid with Л": {
			number:  "34215Л",
			wantErr: ErrInvalidPassportNumber,
		},
		"not valid with Л2": {
			number:  "34Л215",
			wantErr: ErrInvalidPassportNumber,
		},
		"not valid with Л3": {
			number:  "Л657482",
			wantErr: ErrInvalidPassportNumber,
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := IsPassportNumberValid(tt.number)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			require.NoError(t, err)

		})
	}
}

func Test_PassportIssuerCode(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		issuerCode string
		wantErr    error
	}{
		"valid issuer code with dash": {
			issuerCode: "500-159",
		},
		"valid issuer code": {
			issuerCode: "500159",
		},
		"invalid issuer code blank": {
			issuerCode: "",
			wantErr:    ErrEmptyIssuedCode,
		},
		"invalid issuer code [3] is not dash(-)": {
			issuerCode: "50-0159",
			wantErr:    ErrInvalidIssuedCode,
		},
		"invalid issuer code [0:4] is not digest": {
			issuerCode: "qw2-159",
			wantErr:    ErrInvalidIssuedCode,
		},
		"invalid issuer code [4:7] is not digest": {
			issuerCode: "500-qwe",
			wantErr:    ErrInvalidIssuedCode,
		},
		"invalid issuer code 2 dash": {
			issuerCode: "500--159",
			wantErr:    ErrInvalidIssuedCode,
		},
		"invalid issuer code less 6 digest with -": {
			issuerCode: "1234-12",
			wantErr:    ErrInvalidIssuedCode,
		},
		"invalid issuer code less 6 digest": {
			issuerCode: "12345",
			wantErr:    ErrInvalidIssuedCode,
		},
		"invalid issuer code with Д and -h": {
			issuerCode: "123-12Д",
			wantErr:    ErrInvalidIssuedCode,
		},
		"invalid issuer code with Д": {
			issuerCode: "12345Д",
			wantErr:    ErrInvalidIssuedCode,
		},
		"invalid issuer code with V and -": {
			issuerCode: "123-12V",
			wantErr:    ErrInvalidIssuedCode,
		},
		"invalid issuer code with V": {
			issuerCode: "12345V",
			wantErr:    ErrInvalidIssuedCode,
		},
		"invalid issuer code more 6 digest with -": {
			issuerCode: "123-1233",
			wantErr:    ErrInvalidIssuedCode,
		},
		"invalid issuer code more 6 digest": {
			issuerCode: "1234567",
			wantErr:    ErrInvalidIssuedCode,
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := IsPassportIssuerCodeValid(tt.issuerCode)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			require.NoError(t, err)

		})
	}
}

func Test_PassportIssueDate(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		issueDate time.Time
		birthday  time.Time
		checkDate time.Time
		wantErr   error
	}{
		"valid issue date 20 years old": {
			issueDate: time.Date(2017, 2, 20, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(1997, 2, 20, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
		"valid issue date 45 years old issue before 45": {
			issueDate: time.Date(2024, 2, 01, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(1979, 2, 16, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
		"valid issue date 45 years old corner case 90 days": {
			issueDate: time.Date(2024, 5, 16, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(1979, 2, 16, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 5, 27, 0, 0, 0, 0, time.UTC),
		},
		"valid issue date +2 month 20 years old": {
			issueDate: time.Date(2017, 4, 20, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(1997, 2, 20, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
		"valid issue date 18 years old": {
			issueDate: time.Date(2020, 2, 20, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(2006, 2, 20, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
		"valid issue date 45 years old": {
			issueDate: time.Date(2000, 2, 20, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(1955, 2, 20, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
		"valid issue date 55 years old": {
			issueDate: time.Date(2014, 2, 20, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(1969, 2, 20, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
		"invalid issue date less than 14 years old": {
			issueDate: time.Date(2017, 2, 20, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(2010, 2, 20, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidIssueDateBefore14Birthday,
		},
		"invalid issue date nil issueDate": {
			issueDate: time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(2010, 2, 20, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrEmptyIssueDate,
		},
		"invalid issue date nil birthday": {
			issueDate: time.Date(2010, 2, 20, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrEmptyBirthday,
		},
		"invalid issue date passport expire at 20 years old": {
			issueDate: time.Date(2019, 2, 20, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrIssueDatePassportExpiredAt20,
		},
		"invalid issue date passport expire at 45 years old": {
			issueDate: time.Date(1999, 2, 20, 0, 0, 0, 0, time.UTC),
			birthday:  time.Date(1955, 01, 01, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrIssueDatePassportExpiredAt45,
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := IsPassportIssueDateValid(tt.issueDate, tt.birthday, tt.checkDate)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			require.NoError(t, err)

		})
	}
}

func Test_PassportBirthday(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		birthday  time.Time
		wantErr   error
		checkDate time.Time
	}{
		"valid birthday": {
			birthday:  time.Date(1997, 2, 20, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
		"invalid birthday after time Now": {
			birthday:  time.Date(3010, 2, 20, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidBirthday,
		},
		"invalid birthday nil": {
			birthday:  time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidBirthday,
		},
		"17 year": {
			birthday:  time.Date(2006, 02, 28, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
			wantErr:   ErrInvalidBirthday,
		},
		"18 year": {
			birthday:  time.Date(2006, 02, 26, 0, 0, 0, 0, time.UTC),
			checkDate: time.Date(2024, 02, 27, 0, 0, 0, 0, time.UTC),
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := IsPassportBirthdayValid(tt.birthday, tt.checkDate)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			require.NoError(t, err)

		})
	}
}

func Test_PassportPlaceOfBirthNormalize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		placeOfBirth string
		want         string
	}{
		{
			placeOfBirth: "",
			want:         "",
		},
		{
			placeOfBirth: "  ",
			want:         " ",
		},
		{
			placeOfBirth: "--",
			want:         "-",
		},
		{
			placeOfBirth: "-  -",
			want:         "-",
		},
		{
			placeOfBirth: "U  S- SR",
			want:         "U S-SR",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.placeOfBirth, func(t *testing.T) {
			t.Parallel()

			assert.Equalf(t, tt.want, PassportPlaceOfBirthNormalize(tt.placeOfBirth), "PassportPlaceOfBirthNormalize(%v)", tt.placeOfBirth)
		})
	}
}

func Test_PassportIssuedByNormalize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		placeOfBirth string
		want         string
	}{
		{
			placeOfBirth: "",
			want:         "",
		},
		{
			placeOfBirth: "  ",
			want:         " ",
		},
		{
			placeOfBirth: "..",
			want:         ".",
		},
		{
			placeOfBirth: ",,",
			want:         ",",
		},
		{
			placeOfBirth: "ав ,, dsdd,  44 ;..",
			want:         "ав , dsdd, 44 ;.",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.placeOfBirth, func(t *testing.T) {
			t.Parallel()

			assert.Equalf(t, tt.want, PassportIssuedByNormalize(tt.placeOfBirth), "PassportIssuedByNormalize(%v)", tt.placeOfBirth)
		})
	}
}
