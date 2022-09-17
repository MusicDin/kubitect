package modelconfig

import "os"

func PathExists(value interface{}) error {
	path, _ := value.(string)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err != nil {
			// Internal error.
			return err
		}
		// Path exists.
		return nil
	}
	return ErrPathDoesNotExist
}
