package modelconfig

type SSH struct {
	Keyfile *SSHKeyPath
	Port    *Port
	Verify  bool
}
