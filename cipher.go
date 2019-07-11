package goseine

import "crypto/cipher"

type Cipher struct {
	k []byte
}

func NewCipher(key []byte) cipher.Stream {
	return &Cipher{
		k: key,
	}
}

func NewDefaultCipher() cipher.Stream {
	return NewCipher([]byte("qmfaktnpgjs"))
}

func (c *Cipher) XORKeyStream(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		if src[i] == 0x00 {
			dst[i] = src[i]
		} else {
			dst[i] = src[i] ^ c.k[i%len(c.k)]
		}
	}
}
