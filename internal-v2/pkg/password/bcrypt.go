package password

import "golang.org/x/crypto/bcrypt"

const defaultCost = bcrypt.DefaultCost

func bcryptHash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), defaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func bcryptVerify(hash string, plain string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(plain),
	)
	return err == nil
}
