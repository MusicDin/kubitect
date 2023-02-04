package keygen

import "fmt"

type KeyFileError struct {
	keyType string
	fName   string
	err     error
}

func NewKeyFileError(keyType, fileName string, err error) KeyFileError {
	return KeyFileError{
		keyType: keyType,
		fName:   fileName,
		err:     err,
	}
}

func (e KeyFileError) Error() string {
	return fmt.Sprintf("keygen: failed to read %s key file %s: %v", e.keyType, e.fName, e.err)
}
