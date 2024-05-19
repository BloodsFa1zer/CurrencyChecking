package service

import (
	"CurrencyChecking/config"
	"CurrencyChecking/database"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/idna"
	"io"
	"net/http"
	"strings"
)

type UserService struct {
	DbUser   database.DbInterface
	validate *validator.Validate
}

func NewUserService(DbUser database.DbInterface, validate *validator.Validate) *UserService {
	return &UserService{DbUser: DbUser, validate: validate}
}

type Response struct {
	UAH struct {
		CurrencyData struct {
			Code  string  `json:"code"`
			Value float64 `json:"value"`
		} `json:"UAH"`
	} `json:"data"`
}

func (us *UserService) CreateUser(user database.User) (error, int) {
	isValid, err := isValidEmail(user.Email)
	if !isValid {
		return err, http.StatusBadRequest
	}

	err = us.DbUser.InsertUser(user)
	if err != nil {
		if err.Error() == "such Email already exists" {
			log.Warn().Msg("such Email already exists")
			return err, http.StatusConflict
		}
		log.Warn().Err(err).Msg("can`t insert user")
		return err, http.StatusBadRequest
	}

	return nil, http.StatusOK
}

func (us *UserService) GetRate() (string, error, int) {

	cfg := config.LoadENV(".env")

	// can be http.Get(URL as a const), but to make up scalable I decided to recreate URL
	url := cfg.URL + "?apikey=" + cfg.ApiKey + "&currencies=UAH"

	var currencyRate *Response
	response, err := http.Get(url)
	if err != nil {
		log.Warn().Err(err).Msg("can`t read a response")
		return "", err, http.StatusBadRequest
	}
	log.Info().Msg("successfully read and return response")

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Warn().Err(err).Msg("can`t read data")
		return "", err, http.StatusBadRequest
	}
	log.Info().Msg("successfully read data")
	fmt.Println(string(body))

	err = json.Unmarshal(body, &currencyRate)
	if err != nil {
		log.Warn().Err(err).Msg("can`t unmarshal data")
		return "", err, http.StatusBadRequest
	}
	log.Info().Msg("successfully unmarshal data and return it")

	var result []byte
	result = fmt.Appendf(result, "Currency rate is %.4f UAH to 1 USD", currencyRate.UAH.CurrencyData.Value)

	return string(result), nil, http.StatusOK
}

func isValidEmail(email string) (bool, error) {
	// Check the overall length of the email address
	if len(email) > 320 {
		return false, fmt.Errorf("email length exceeds 320 characters")
	}

	// Transform the local part to lowercase for case-insensitive unique storage
	parts := strings.Split(strings.ToLower(email), "@")
	if len(parts) != 2 {
		return false, fmt.Errorf("email must contain a single '@' character")
	}
	localPart, domainPart := parts[0], parts[1]

	// Check for empty local or domain parts
	if len(localPart) == 0 || len(domainPart) == 0 {
		return false, fmt.Errorf("local or domain part cannot be empty in the email")
	}

	// Check for consecutive special characters in the local part
	prevChar := rune(0)
	for _, char := range localPart {
		if strings.ContainsRune("!#$%&'*+-/=?^_`{|}~.", char) {
			if char == prevChar && char != '-' {
				return false, fmt.Errorf("consecutive special characters '%c' are not allowed in the local part", char)
			}
		}
		prevChar = char
	}

	// Check for spaces
	if strings.ContainsAny(email, " ") {
		return false, fmt.Errorf("spaces are not allowed in the email")
	}

	// Check the length of the local part and the domain part
	if len(localPart) > 64 || len(domainPart) > 255 {
		return false, fmt.Errorf("local part or domain part length exceeds the limit in the email")
	}

	// Convert international domain to ASCII (Punycode) only if needed
	asciiDomain, err := idna.ToASCII(domainPart)
	if err != nil {
		return false, fmt.Errorf("failed to convert domain to ASCII: %s", err)
	}

	// Convert international local part to ASCII (Punycode) only if needed
	_, err = idna.ToASCII(localPart)
	if err != nil {
		return false, fmt.Errorf("failed to convert local part to ASCII: %s", err)
	}

	// Check that the domain labels do not start or end with special characters and TLD is alphabetic
	domainLabels := strings.Split(asciiDomain, ".")
	for i, label := range domainLabels {
		// Check first and last character of each label
		if !strings.ContainsAny("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", string(label[0])) ||
			!strings.ContainsAny("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", string(label[len(label)-1])) {
			return false, fmt.Errorf("domain labels must not start or end with special characters in the email")
		}

		// Check label length
		if len(label) > 63 {
			return false, fmt.Errorf("domain label length exceeds the limit in the email")
		}

		// Validate that the TLD is alphabetic
		if i == len(domainLabels)-1 && !strings.HasPrefix(label, "xn--") {
			decodedTLD, err := idna.ToUnicode(label)
			if err != nil {
				return false, fmt.Errorf("failed to decode TLD: %s", err)
			}
			isAlpha := true
			for _, ch := range decodedTLD {
				if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') {
					isAlpha = false
					break
				}
			}
			if !isAlpha {
				return false, fmt.Errorf("TLD must be alphabetic in the email")
			}
		}
	}

	return true, nil
}

//func main() {
//	// Standard ASCII email address (Valid)
//	fmt.Println(isValid("test@example.com")) // Should return "test@example.com"
//
//	// Unicode in local part (Valid)
//	fmt.Println(isValid("tést@example.com")) // Should return Punycode-converted email
//
//	// Punycode-encoded domain (IDNA) (Valid)
//	fmt.Println(isValid("test@xn--fsqu00a.xn--0zwm56d")) // Should return "test@xn--fsqu00a.xn--0zwm56d"
//
//	// Unicode in domain part (Valid)
//	fmt.Println(isValid("test@例子.测试")) // Should return Punycode-converted domain
//
//	// Unicode in both local and domain parts (Valid)
//	fmt.Println(isValid("tést@例子.测试")) // Should return fully Punycode-converted email
//
//	// Cyrillic script in local part (Valid)
//	fmt.Println(isValid("тест@пример.ру")) // Should return fully Punycode-converted email
//
//	// Arabic script in local part (Valid)
//	fmt.Println(isValid("اختبار@مثال.اختبار")) // Should return fully Punycode-converted email
//
//	// Hebrew script in local part (Valid)
//	fmt.Println(isValid("בדיקה@דוגמה.בדיקה")) // Should return fully Punycode-converted email
//
//	// Punycode-encoded local part and domain (Valid)
//	fmt.Println(isValid("xn--e1aybc@xn--80akhbyknj4f.xn--p1ai")) // Should return "xn--e1aybc@xn--80akhbyknj4f.xn--p1ai"
//
//	// Various other languages (All should be Valid)
//	fmt.Println(isValid("測試@例子.測試"))   // Should return fully Punycode-converted email (Chinese)
//	fmt.Println(isValid("テスト@例.テスト"))  // Should return fully Punycode-converted email (Japanese)
//	fmt.Println(isValid("테스트@예시.테스트")) // Should return fully Punycode-converted email (Korean)
//
//	// Email with consecutive special characters in local part (Invalid)
//	fmt.Println(isValid("test..test@example.com")) // Should return an error
//
//	// Email with special characters at the beginning or end of the local part (Invalid)
//	fmt.Println(isValid(".test@example.com")) // Should return an error
//	fmt.Println(isValid("test.@example.com")) // Should return an error
//
//	// Email with space (Invalid)
//	fmt.Println(isValid("test test@example.com")) // Should return an error
//
//	// Email with no local part (Invalid)
//	fmt.Println(isValid("@example.com")) // Should return an error
//
//	// Email with no domain part (Invalid)
//	fmt.Println(isValid("test@")) // Should return an error
//
//	// Email with empty string (Invalid)
//	fmt.Println(isValid("")) // Should return an error
//
//	// Email exceeding 320 character limit (Invalid)
//	fmt.Println(isValid(strings.Repeat("a", 320) + "@example.com")) // Should return an error
//}
