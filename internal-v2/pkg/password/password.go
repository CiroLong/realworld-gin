package password

func Hash(plain string) (string, error) {
	return bcryptHash(plain)
}

func Verify(hash string, plain string) bool {
	return bcryptVerify(hash, plain)
}
