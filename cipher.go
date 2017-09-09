package goseine

type Cipher struct {
	key []byte
}

func NewCipher(key []byte) *Cipher {
	return &Cipher {
		key: key,
	}
}

func (c *Cipher) Encrypt(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		if src[i] == 0x00 {
			dst[i] = src[i]
		} else {
			dst[i] = src[i] ^ c.key[i % len(c.key)]
		}
	}
}

func (c *Cipher) Decrypt(dst, src []byte) {
	c.Encrypt(dst, src)
}
