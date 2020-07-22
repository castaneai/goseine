package packets

var key = []byte("qmfaktnpgjs")

func Encrypt(b []byte) {
	for i := 0; i < len(b); i++ {
		if b[i] != 0x00 {
			e := b[i] ^ key[i%len(key)]
			if e != 0x00 {
				b[i] = e
			}
		}
	}
}

func Decrypt(b []byte) {
	Encrypt(b)
}
