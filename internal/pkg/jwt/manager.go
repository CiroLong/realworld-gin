package jwt

type Manager interface {
	Generate(userID int64) (string, error)
	Parse(token string) (int64, error)
}
