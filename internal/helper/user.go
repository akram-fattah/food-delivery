package helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/akram-fattah/food-delivery/internal/models"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)


func isStrongPassword(p string) bool {
	if len(p) < 6 || len(p) > 20 {
		return false
	}

	var hasLower, hasUpper, hasNumber, hasSpecial bool

	for _, c := range p {
		switch {
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= '0' && c <= '9':
			hasNumber = true
		case strings.ContainsRune("@$!%*?&", c):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasNumber && hasSpecial
}

func ParseUser(r *http.Request) (models.User, error) {
	var u models.User

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&u)
	return u, err
}

func ValidateUser(u models.User) error {

	if strings.TrimSpace(u.Name) == "" {
		return errors.New("الاسم مطلوب")
	}

	if strings.TrimSpace(u.Username) == "" {
		return errors.New("اسم المستخدم مطلوب")
	}


	if strings.TrimSpace(u.Email) == "" {
		return errors.New("البريد الإلكتروني مطلوب")
	}

	if !emailRegex.MatchString(u.Email) {
		return errors.New("البريد الإلكتروني غير صالح")
	}
	 if !usernameRegex.MatchString(u.Username) {
		return errors.New("اسم المستخدم يجب أن يحتوي فقط على أحرف وأرقام وشرطات")
	}

	if len(u.Password) < 6 || len(u.Password) > 20{
		return errors.New("كلمة المرور يجب أن تكون 6 أحرف على الأقل و20 حرفًا كحد أقصى")
	}
	if !isStrongPassword(u.Password) {
		return errors.New("كلمة المرور يجب أن تحتوي على حرف كبير، حرف صغير، رقم، ورمز خاص")
	}

	return nil
}