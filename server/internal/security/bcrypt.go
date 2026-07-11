package security

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func NewBcryptHaser() *BcryptHasher {
	return &BcryptHasher{}
}

func (b *BcryptHasher) HashPassword(password string) (string, error) {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedpassword), err
}

func (b *BcryptHasher) HashedToPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
