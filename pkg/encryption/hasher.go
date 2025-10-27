package encryption

type Hasher interface {
	Generate(password []byte, cost int) ([]byte, error)
	Compare(hashedPassword []byte, password []byte) error
}
