package datastruct

var (
	aeadNames []string
)

func init() {
	aeadNames = []string{
		"",
		"AEAD_AES_128_GCM",           // 1
		"AEAD_AES_256_GCM",           // 2
		"AEAD_AES_128_CCM",           // 3
		"AEAD_AES_256_CCM",           // 4
		"AEAD_AES_128_GCM_8",         // 5
		"AEAD_AES_256_GCM_8",         // 6
		"AEAD_AES_128_GCM_12",        // 7
		"AEAD_AES_256_GCM_12",        // 8
		"AEAD_AES_128_CCM_SHORT",     // 9
		"AEAD_AES_256_CCM_SHORT",     // 10
		"AEAD_AES_128_CCM_SHORT_8",   // 11
		"AEAD_AES_256_CCM_SHORT_8",   // 12
		"AEAD_AES_128_CCM_SHORT_12",  // 13
		"AEAD_AES_256_CCM_SHORT_12",  // 14
		"AEAD_AES_SIV_CMAC_256",      // 15
		"AEAD_AES_SIV_CMAC_384",      // 16
		"AEAD_AES_SIV_CMAC_512",      // 17
		"AEAD_AES_128_CCM_8",         // 18
		"AEAD_AES_256_CCM_8",         // 19
		"AEAD_AES_128_OCB_TAGLEN128", // 20
		"AEAD_AES_128_OCB_TAGLEN96",  // 21
		"AEAD_AES_128_OCB_TAGLEN64",  // 22
		"AEAD_AES_192_OCB_TAGLEN128", // 23
		"AEAD_AES_192_OCB_TAGLEN96",  // 24
		"AEAD_AES_192_OCB_TAGLEN64",  // 25
		"AEAD_AES_256_OCB_TAGLEN128", // 26
		"AEAD_AES_256_OCB_TAGLEN96",  // 27
		"AEAD_AES_256_OCB_TAGLEN64",  // 28
		"AEAD_CHACHA20_POLY1305",     // 29
		"AEAD_AES_128_GCM_SIV",       // 30
		"AEAD_AES_256_GCM_SIV",       // 31
		"AEAD_AEGIS128L",             // 32
		"AEAD_AEGIS256",              // 33
	}
}

func GetAEADName(id byte) string {
	return aeadNames[id]
}
