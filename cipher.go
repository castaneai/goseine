package goseine

type Cipher interface {
	Encrypt(dst, src []byte)
	Decrypt(dst, src []byte)
}

type GoseineCipher struct {
	key []byte
}

func NewGoseineCipher(key []byte) Cipher {
	return &GoseineCipher{
		key: key,
	}
}

func NewCurrentGoseineCipher() Cipher {
	return NewGoseineCipher([]byte("qmfaktnpgjs"))
}

func (c *GoseineCipher) Encrypt(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		if src[i] == 0x00 {
			dst[i] = src[i]
		} else {
			dst[i] = src[i] ^ c.key[i%len(c.key)]
		}
	}
}

func (c *GoseineCipher) Decrypt(dst, src []byte) {
	c.Encrypt(dst, src)
}
