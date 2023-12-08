package mars

import (
	c_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"math/bits"
	m_rand "math/rand"
)

var K [40]uint32

var S = [512]uint32{
	0x09d0c479, 0x28c8ffe0, 0x84aa6c39, 0x9dad7287,
	0x7dff9be3, 0xd4268361, 0xc96da1d4, 0x7974cc93,
	0x85d0582e, 0x2a4b5705, 0x1ca16a62, 0xc3bd279d,
	0x0f1f25e5, 0x5160372f, 0xc695c1fb, 0x4d7ff1e4,
	0xae5f6bf4, 0x0d72ee46, 0xff23de8a, 0xb1cf8e83,
	0xf14902e2, 0x3e981e42, 0x8bf53eb6, 0x7f4bf8ac,
	0x83631f83, 0x25970205, 0x76afe784, 0x3a7931d4,
	0x4f846450, 0x5c64c3f6, 0x210a5f18, 0xc6986a26,
	0x28f4e826, 0x3a60a81c, 0xd340a664, 0x7ea820c4,
	0x526687c5, 0x7eddd12b, 0x32a11d1d, 0x9c9ef086,
	0x80f6e831, 0xab6f04ad, 0x56fb9b53, 0x8b2e095c,
	0xb68556ae, 0xd2250b0d, 0x294a7721, 0xe21fb253,
	0xae136749, 0xe82aae86, 0x93365104, 0x99404a66,
	0x78a784dc, 0xb69ba84b, 0x04046793, 0x23db5c1e,
	0x46cae1d6, 0x2fe28134, 0x5a223942, 0x1863cd5b,
	0xc190c6e3, 0x07dfb846, 0x6eb88816, 0x2d0dcc4a,
	0xa4ccae59, 0x3798670d, 0xcbfa9493, 0x4f481d45,
	0xeafc8ca8, 0xdb1129d6, 0xb0449e20, 0x0f5407fb,
	0x6167d9a8, 0xd1f45763, 0x4daa96c3, 0x3bec5958,
	0xababa014, 0xb6ccd201, 0x38d6279f, 0x02682215,
	0x8f376cd5, 0x092c237e, 0xbfc56593, 0x32889d2c,
	0x854b3e95, 0x05bb9b43, 0x7dcd5dcd, 0xa02e926c,
	0xfae527e5, 0x36a1c330, 0x3412e1ae, 0xf257f462,
	0x3c4f1d71, 0x30a2e809, 0x68e5f551, 0x9c61ba44,
	0x5ded0ab8, 0x75ce09c8, 0x9654f93e, 0x698c0cca,
	0x243cb3e4, 0x2b062b97, 0x0f3b8d9e, 0x00e050df,
	0xfc5d6166, 0xe35f9288, 0xc079550d, 0x0591aee8,
	0x8e531e74, 0x75fe3578, 0x2f6d829a, 0xf60b21ae,
	0x95e8eb8d, 0x6699486b, 0x901d7d9b, 0xfd6d6e31,
	0x1090acef, 0xe0670dd8, 0xdab2e692, 0xcd6d4365,
	0xe5393514, 0x3af345f0, 0x6241fc4d, 0x460da3a3,
	0x7bcf3729, 0x8bf1d1e0, 0x14aac070, 0x1587ed55,
	0x3afd7d3e, 0xd2f29e01, 0x29a9d1f6, 0xefb10c53,
	0xcf3b870f, 0xb414935c, 0x664465ed, 0x024acac7,
	0x59a744c1, 0x1d2936a7, 0xdc580aa6, 0xcf574ca8,
	0x040a7a10, 0x6cd81807, 0x8a98be4c, 0xaccea063,
	0xc33e92b5, 0xd1e0e03d, 0xb322517e, 0x2092bd13,
	0x386b2c4a, 0x52e8dd58, 0x58656dfb, 0x50820371,
	0x41811896, 0xe337ef7e, 0xd39fb119, 0xc97f0df6,
	0x68fea01b, 0xa150a6e5, 0x55258962, 0xeb6ff41b,
	0xd7c9cd7a, 0xa619cd9e, 0xbcf09576, 0x2672c073,
	0xf003fb3c, 0x4ab7a50b, 0x1484126a, 0x487ba9b1,
	0xa64fc9c6, 0xf6957d49, 0x38b06a75, 0xdd805fcd,
	0x63d094cf, 0xf51c999e, 0x1aa4d343, 0xb8495294,
	0xce9f8e99, 0xbffcd770, 0xc7c275cc, 0x378453a7,
	0x7b21be33, 0x397f41bd, 0x4e94d131, 0x92cc1f98,
	0x5915ea51, 0x99f861b7, 0xc9980a88, 0x1d74fd5f,
	0xb0a495f8, 0x614deed0, 0xb5778eea, 0x5941792d,
	0xfa90c1f8, 0x33f824b4, 0xc4965372, 0x3ff6d550,
	0x4ca5fec0, 0x8630e964, 0x5b3fbbd6, 0x7da26a48,
	0xb203231a, 0x04297514, 0x2d639306, 0x2eb13149,
	0x16a45272, 0x532459a0, 0x8e5f4872, 0xf966c7d9,
	0x07128dc0, 0x0d44db62, 0xafc8d52d, 0x06316131,
	0xd838e7ce, 0x1bc41d00, 0x3a2e8c0f, 0xea83837e,
	0xb984737d, 0x13ba4891, 0xc4f8b949, 0xa6d6acb3,
	0xa215cdce, 0x8359838b, 0x6bd1aa31, 0xf579dd52,
	0x21b93f93, 0xf5176781, 0x187dfdde, 0xe94aeb76,
	0x2b38fd54, 0x431de1da, 0xab394825, 0x9ad3048f,
	0xdfea32aa, 0x659473e3, 0x623f7863, 0xf3346c59,
	0xab3ab685, 0x3346a90b, 0x6b56443e, 0xc6de01f8,
	0x8d421fc0, 0x9b0ed10c, 0x88f1a1e9, 0x54c1f029,
	0x7dead57b, 0x8d7ba426, 0x4cf5178a, 0x551a7cca,
	0x1a9a5f08, 0xfcd651b9, 0x25605182, 0xe11fc6c3,
	0xb6fd9676, 0x337b3027, 0xb7c8eb14, 0x9e5fd030,

	0x6b57e354, 0xad913cf7, 0x7e16688d, 0x58872a69,
	0x2c2fc7df, 0xe389ccc6, 0x30738df1, 0x0824a734,
	0xe1797a8b, 0xa4a8d57b, 0x5b5d193b, 0xc8a8309b,
	0x73f9a978, 0x73398d32, 0x0f59573e, 0xe9df2b03,
	0xe8a5b6c8, 0x848d0704, 0x98df93c2, 0x720a1dc3,
	0x684f259a, 0x943ba848, 0xa6370152, 0x863b5ea3,
	0xd17b978b, 0x6d9b58ef, 0x0a700dd4, 0xa73d36bf,
	0x8e6a0829, 0x8695bc14, 0xe35b3447, 0x933ac568,
	0x8894b022, 0x2f511c27, 0xddfbcc3c, 0x006662b6,
	0x117c83fe, 0x4e12b414, 0xc2bca766, 0x3a2fec10,
	0xf4562420, 0x55792e2a, 0x46f5d857, 0xceda25ce,
	0xc3601d3b, 0x6c00ab46, 0xefac9c28, 0xb3c35047,
	0x611dfee3, 0x257c3207, 0xfdd58482, 0x3b14d84f,
	0x23becb64, 0xa075f3a3, 0x088f8ead, 0x07adf158,
	0x7796943c, 0xfacabf3d, 0xc09730cd, 0xf7679969,
	0xda44e9ed, 0x2c854c12, 0x35935fa3, 0x2f057d9f,
	0x690624f8, 0x1cb0bafd, 0x7b0dbdc6, 0x810f23bb,
	0xfa929a1a, 0x6d969a17, 0x6742979b, 0x74ac7d05,
	0x010e65c4, 0x86a3d963, 0xf907b5a0, 0xd0042bd3,
	0x158d7d03, 0x287a8255, 0xbba8366f, 0x096edc33,
	0x21916a7b, 0x77b56b86, 0x951622f9, 0xa6c5e650,
	0x8cea17d1, 0xcd8c62bc, 0xa3d63433, 0x358a68fd,
	0x0f9b9d3c, 0xd6aa295b, 0xfe33384a, 0xc000738e,
	0xcd67eb2f, 0xe2eb6dc2, 0x97338b02, 0x06c9f246,
	0x419cf1ad, 0x2b83c045, 0x3723f18a, 0xcb5b3089,
	0x160bead7, 0x5d494656, 0x35f8a74b, 0x1e4e6c9e,
	0x000399bd, 0x67466880, 0xb4174831, 0xacf423b2,
	0xca815ab3, 0x5a6395e7, 0x302a67c5, 0x8bdb446b,
	0x108f8fa4, 0x10223eda, 0x92b8b48b, 0x7f38d0ee,
	0xab2701d4, 0x0262d415, 0xaf224a30, 0xb3d88aba,
	0xf8b2c3af, 0xdaf7ef70, 0xcc97d3b7, 0xe9614b6c,
	0x2baebff4, 0x70f687cf, 0x386c9156, 0xce092ee5,
	0x01e87da6, 0x6ce91e6a, 0xbb7bcc84, 0xc7922c20,
	0x9d3b71fd, 0x060e41c6, 0xd7590f15, 0x4e03bb47,
	0x183c198e, 0x63eeb240, 0x2ddbf49a, 0x6d5cba54,
	0x923750af, 0xf9e14236, 0x7838162b, 0x59726c72,
	0x81b66760, 0xbb2926c1, 0x48a0ce0d, 0xa6c0496d,
	0xad43507b, 0x718d496a, 0x9df057af, 0x44b1bde6,
	0x054356dc, 0xde7ced35, 0xd51a138b, 0x62088cc9,
	0x35830311, 0xc96efca2, 0x686f86ec, 0x8e77cb68,
	0x63e1d6b8, 0xc80f9778, 0x79c491fd, 0x1b4c67f2,
	0x72698d7d, 0x5e368c31, 0xf7d95e2e, 0xa1d3493f,
	0xdcd9433e, 0x896f1552, 0x4bc4ca7a, 0xa6d1baf4,
	0xa5a96dcc, 0x0bef8b46, 0xa169fda7, 0x74df40b7,
	0x4e208804, 0x9a756607, 0x038e87c8, 0x20211e44,
	0x8b7ad4bf, 0xc6403f35, 0x1848e36d, 0x80bdb038,
	0x1e62891c, 0x643d2107, 0xbf04d6f8, 0x21092c8c,
	0xf644f389, 0x0778404e, 0x7b78adb8, 0xa2c52d53,
	0x42157abe, 0xa2253e2e, 0x7bf3f4ae, 0x80f594f9,
	0x953194e7, 0x77eb92ed, 0xb3816930, 0xda8d9336,
	0xbf447469, 0xf26d9483, 0xee6faed5, 0x71371235,
	0xde425f73, 0xb4e59f43, 0x7dbe2d4e, 0x2d37b185,
	0x49dc9a63, 0x98c39d98, 0x1301c9a2, 0x389b1bbf,
	0x0c18588d, 0xa421c1ba, 0x7aa3865c, 0x71e08558,
	0x3c5cfcaa, 0x7d239ca4, 0x0297d9dd, 0xd7dc2830,
	0x4b37802b, 0x7428ab54, 0xaeee0347, 0x4b3fbb85,
	0x692f2f08, 0x134e578e, 0x36d9e0bf, 0xae8b5fcf,
	0xedb93ecf, 0x2b27248e, 0x170eb1ef, 0x7dc57fd6,
	0x1e760f16, 0xb1136601, 0x864e1b9b, 0xd7ea7319,
	0x3ab871bd, 0xcfa4d76f, 0xe31bd782, 0x0dbeb469,
	0xabb96061, 0x5370f85d, 0xffb07e37, 0xda30d0fb,
	0xebc977b6, 0x0b98b40f, 0x3a4d0fe6, 0xdf4fc26b,
	0x159cf22a, 0xc298d6e2, 0x2b78ef6a, 0x61a94ac0,
	0xab561187, 0x14eea0f0, 0xdf0d4164, 0x19af70ee}

// rotl выполняет побитовый циклический сдвиг влево для значения x на s позиций
func rotl(x uint32, s int) uint32 {
	return bits.RotateLeft32(x, s)
}

// rotr выполняет побитовый циклический сдвиг вправо для значения x на s позиций
func rotr(x uint32, s int) uint32 {
	return bits.RotateLeft32(x, -s) // иначе выполняем сдвиг вправо и объединяем результат с сдвигом влево
}

// setKey принимает срез беззнаковых целых чисел в качестве ключа
func setKey(key []uint32) {
	var T [15]uint32
	copy(K[:], key)
	copy(T[:], key)
	T[len(key)] = uint32(len(key))

	for j := 0; j < 4; j++ {
		for i := 0; i < 15; i++ {
			T[i] = T[i] ^ rotl(T[(i+8)%15]^T[(i+13)%15], 3) ^ uint32(4*i+j)
		}
		for i := 0; i < 15; i++ {
			T[i] = rotl(T[i]+S[511&T[(i+14)%15]], 9)
		}
		for i := 0; i < 10; i++ {
			K[10*j+i] = T[(4*i)%15]
		}
	}

	for i := 5; i <= 35; i += 2 {
		j := K[i] & 3
		w := K[i] | 3

		M := makeMask(w)
		p := rotl(S[j+256], int(31&K[i-1])) //?
		K[i] = w ^ (p & M)
	}
}

// e_func выполняет функцию E на входном блоке с двумя ключами
// и возвращает кортеж из трех беззнаковых целых чисел
func e_func(in, key1, key2 uint32) (uint32, uint32, uint32) {
	var L, M, R uint32

	M = in + key1
	R = rotl(in, 13) * key2 // rotl - это функция для циклического сдвига влево
	i := 511 & M
	L = S[i] // S - это массив беззнаковых целых чисел
	R = rotl(R, 5)
	r := 31 & R
	M = rotl(M, int(r))
	L = L ^ R // ^ - это оператор исключающего ИЛИ
	R = rotl(R, 5)
	L = L ^ R
	r = 31 & R
	L = rotl(L, int(r))

	return L, M, R
}

// f_mix объявляет функцию, которая принимает четыре беззнаковых целых числа по указателю и не возвращает ничего
func f_mix(a, b, c, d *uint32) {
	// здесь можно написать тело функции, используя *a, *b, *c, *d для доступа к значениям переменных
	// rotr - функция, которая циклически сдвигает число вправо на заданное количество битов
	r := rotr(*a, 8)
	*b ^= S[*a&255]
	*b += S[(r&255)+256]
	r = rotr(*a, 16)
	*a = rotr(*a, 24)
	*c += S[r&255]
	*d ^= S[(*a&255)+256]
}

// b_mix объявляет функцию, которая принимает четыре беззнаковых целых числа по указателю и не возвращает ничего
func b_mix(a, b, c, d *uint32) {
	// здесь можно написать тело функции, используя *a, *b, *c, *d для доступа к значениям переменных
	r := rotl(*a, 8)      // r - локальная переменная, которая хранит результат циклического сдвига a на 8 бит влево
	*b ^= S[(*a&255)+256] // b изменяется побитовым исключающим ИЛИ с элементом массива S, индекс которого равен a, умноженному на 255 и прибавленному к 256
	*c -= S[r&255]        // c уменьшается на элемент массива S, индекс которого равен r, умноженному на 255
	r = rotl(*a, 16)      // r обновляется циклическим сдвигом a на 16 бит влево
	*a = rotl(*a, 24)     // a обновляется циклическим сдвигом a на 24 бита влево
	*d -= S[(r&255)+256]  // d уменьшается на элемент массива S, индекс которого равен r, умноженному на 255 и прибавленному к 256
	*d ^= S[*a&255]       // d изменяется побитовым исключающим ИЛИ с элементом массива S, индекс которого равен a, умноженному на 255
}

// f_trans - функция, которая принимает четыре беззнаковых целых числа по ссылке и одно беззнаковое целое число по значению
func f_trans(a, b, c, d *uint32, i uint32) {
	// вызываем функцию e_func с аргументами a, K[2*i+4] и K[2*i+5] и получаем три возвращаемых значения out1, out2 и out3
	out1, out2, out3 := e_func(*a, K[2*i+4], K[2*i+5])
	// выполняем циклический сдвиг влево на 13 бит для a и присваиваем результат обратно a
	*a = rotl(*a, 13)
	// прибавляем out2 к c и присваиваем результат обратно c
	*c += out2
	// прибавляем out1 к b и присваиваем результат обратно b
	*b += out1
	// выполняем побитовое исключающее или между d и out3 и присваиваем результат обратно d
	*d = *d ^ out3
}

// r_trans - функция, которая принимает четыре беззнаковых целых числа по ссылке и одно беззнаковое целое число по значению
func r_trans(a, b, c, d *uint32, i uint32) {
	// rotr - функция, которая циклически сдвигает число a на 13 бит вправо
	*a = rotr(*a, 13)
	// e_func - функция, которая принимает число a и два элемента массива K с индексами 2*i+4 и 2*i+5 и возвращает три числа out1, out2, out3
	out1, out2, out3 := e_func(*a, K[2*i+4], K[2*i+5])
	// c уменьшается на out2, b уменьшается на out1, d выполняет побитовое исключающее или с out3
	*c -= out2
	*b -= out1
	*d = *d ^ out3
}

// Функция makeMask принимает беззнаковое целое число w в качестве параметра
// и возвращает беззнаковое целое число, которое является битовой маской
// с w единичных битов в младших разрядах
func makeMask(w uint32) uint32 {
	// m - это результат побитового XOR между дополнением w и w, сдвинутым на 1 бит вправо, с маской 0x7fffffff
	m := (^w ^ (w >> 1)) & 0x7fffffff

	// m - это результат побитового AND между m и его сдвигами на 1, 2, 3 и 6 бит вправо
	m &= (m >> 1) & (m >> 2)
	m &= (m >> 3) & (m >> 6)

	m <<= 1
	m |= m << 1
	m |= m << 2
	m |= m << 4
	m |= m << 1

	m |= m >> 1

	// m - это результат побитового AND между m и маской 0x3FFFFFFE
	m &= 0x3FFFFFFE

	m &= ^((^(w << 1) ^ w) ^ (^(w >> 1) ^ w))

	// возвращаем m
	return m
}

// Encrypt принимает входной блок из 4 беззнаковых целых чисел и ключ из вектора беззнаковых целых чисел
// и возвращает выходной блок из 4 беззнаковых целых чисел, зашифрованный с помощью ключа
func encryptBlock(inBlock [4]uint32, key []uint32) [4]uint32 {
	setKey(key)

	var D [4]uint32
	for i := 0; i < 4; i++ {
		D[i] = inBlock[i] + K[i]
	}

	for i := 0; i < 8; i++ {
		f_mix(&D[0], &D[1], &D[2], &D[3])
		if i == 0 || i == 4 {
			D[0] += D[3]
		}
		if i == 1 || i == 5 {
			D[0] += D[1]
		}
		D = rotateLeft(D)
	}

	for i := 0; i < 16; i++ {
		f_trans(&D[0], &D[1], &D[2], &D[3], uint32(i))
		if i < 8 {
			D = rotateLeft(D)
		} else {
			D = rotateRight(D)
		}

		if i == 7 || i == 15 {
			D[1], D[3] = D[3], D[1]
		}
	}

	for i := 0; i < 8; i++ {
		b_mix(&D[0], &D[1], &D[2], &D[3])
		if i == 1 || i == 5 {
			D[1] -= D[0]
		}
		if i == 2 || i == 6 {
			D[1] -= D[2]
		}

		D = rotateLeft(D)
	}

	for i := 0; i < 4; i++ {
		D[i] -= K[i+36]
	}

	return D
}

// rotateLeft shifts the elements of the array to the left by one position
func rotateLeft(arr [4]uint32) [4]uint32 {
	var res [4]uint32
	for i := 0; i < 4; i++ {
		res[i] = arr[(i+1)%4]
	}
	return res
}

// rotateRight shifts the elements of the array to the right by one position
func rotateRight(arr [4]uint32) [4]uint32 {
	var res [4]uint32
	for i := 0; i < 4; i++ {
		res[i] = arr[(i+3)%4]
	}
	return res
}

func Encrypt(in, arr_key []byte) []byte {
	arr := convertByteToUint32(in)
	key := convertByteToUint32(arr_key)
	// Создаем срез для хранения зашифрованных данных
	encrypt := make([]uint32, 0)
	// Проходим по срезу arr с шагом 4
	log.Println(fmt.Sprintf("Len arr : %d", len(arr)))
	for i := 0; i < len(arr); i += 4 {
		// Вызываем функцию mars.encrypt с четырьмя элементами arr и ключом key
		// Получаем срез ciphertext
		ciphertext := encryptBlock([4]uint32{arr[i], arr[i+1], arr[i+2], arr[i+3]}, key)
		// Добавляем ciphertext в конец среза encrypt
		encrypt = append(encrypt, ciphertext[:]...)
	}
	encrypted := convertUint32ToByte(encrypt)

	return encrypted
}

func convertByteToUint32(arr []byte) []uint32 {
	// Создаем срез для хранения uint32 значений
	res := make([]uint32, 0)
	// Читаем байты в uint32 с помощью пакета encoding/binary
	for i := 0; i < len(arr); i += 4 {
		// Проверяем, достаточно ли байтов для чтения uint32
		if i+4 >= len(arr) {
			// Дополняем байтовый срез нулями, чтобы получить 4 байта
			arr = append(arr, make([]byte, 4-(len(arr)-i))...)
		}
		// Читаем 4 байта в uint32
		res = append(res, binary.LittleEndian.Uint32(arr[i:i+4]))
	}
	if len(res)%4 != 0 {
		res = append(res, make([]uint32, 4-len(res)%4)...)
	}
	return res
}

func convertUint32ToByte(arr []uint32) []byte {
	// Преобразование среза uint32 в срез байтов
	dest := make([]byte, 4*len(arr)) // создаем срез байтов нужной длины
	for i, v := range arr {          // перебираем элементы среза source
		binary.LittleEndian.PutUint32(dest[i*4:(i+1)*4], v) // записываем каждый элемент в срез dest
	}
	for i := len(dest) - 1; i > 0; i-- {
		if dest[i] == 0 {
			dest = dest[:len(dest)-1]
		}
	}
	return dest
}

// Decrypt принимает входной блок из 4 беззнаковых целых чисел и ключ из вектора беззнаковых целых чисел
// и возвращает выходной блок из 4 беззнаковых целых чисел, расшифрованный с помощью ключа
func decryptBlock(inBlock [4]uint32, key []uint32) [4]uint32 {
	setKey(key)

	var D [4]uint32
	for i := 0; i < 4; i++ {
		D[3-i] = inBlock[i] + K[i+36]
	}

	for i := 0; i < 8; i++ {
		f_mix(&D[0], &D[1], &D[2], &D[3])
		if i == 0 || i == 4 {
			D[0] += D[3]
		}
		if i == 1 || i == 5 {
			D[0] += D[1]
		}
		D = rotateLeft(D)
	}

	for i := 0; i < 16; i++ {
		r_trans(&D[0], &D[1], &D[2], &D[3], uint32(15-i))
		if i < 8 {
			D = rotateLeft(D)
		} else {
			D = rotateRight(D)
		}

		if i == 7 || i == 15 {
			D[1], D[3] = D[3], D[1]
		}
	}

	for i := 0; i < 8; i++ {
		b_mix(&D[0], &D[1], &D[2], &D[3])
		if i == 1 || i == 5 {
			D[1] -= D[0]
		}
		if i == 2 || i == 6 {
			D[1] -= D[2]
		}

		D = rotateLeft(D)
	}

	for i := 0; i < 4; i++ {
		D[i] -= K[3-i]
	}
	D = reverse(D)

	return D
}

func reverse(a [4]uint32) [4]uint32 {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func Decrypt(encrypted, arr_key []byte) []byte {
	encrypt := convertByteToUint32(encrypted)
	key := convertByteToUint32(arr_key)
	// Создаем срез для хранения зашифрованных данных
	decrypt := make([]uint32, 0)
	for i := 0; i < len(encrypt); i += 4 {
		plaintext := decryptBlock([4]uint32{encrypt[i], encrypt[i+1], encrypt[i+2], encrypt[i+3]}, key)
		decrypt = append(decrypt, plaintext[:]...)
	}
	decrypted := convertUint32ToByte(decrypt)
	return decrypted
}

//func GetRandomKey() []byte {
//	m_rand.Seed(time.Now().UnixNano())
//	unfLen := m_rand.Intn(11) + 4 // random integer between 4 and 14
//	unfNum := m_rand.Uint32()     // random unsigned 32-bit integer
//
//	key := make([]uint32, unfLen)
//	for i := 0; i < unfLen; i++ {
//		key[i] = unfNum
//	}
//
//	return convertUint32ToByte(key)
//}

// generateNumber функция, которая генерирует 64-битное число в заданном диапазоне
func generateNumber(min, max *big.Int) *big.Int {
	// вычисляем разность между максимальным и минимальным значением
	diff := new(big.Int).Sub(max, min)
	n, _ := c_rand.Int(c_rand.Reader, diff)
	n.Add(n, min)
	// прибавляем к нему минимальное значение
	return n
}

// Функция для вычисления модуля при бинарном возведении в степень
func modulo(base, exponent, mod *big.Int) *big.Int {
	base_copy := new(big.Int).Set(base)
	exponent_copy := new(big.Int).Set(exponent)
	mod_copy := new(big.Int).Set(mod)

	x := big.NewInt(1)
	y := new(big.Int).Set(base_copy)

	for exponent_copy.Cmp(big.NewInt(0)) > 0 {
		if exponent_copy.Bit(0) == 1 {
			x.Mul(x, y)
			x.Mod(x, mod_copy)
		}

		y.Mul(y, y)
		y.Mod(y, mod_copy)
		exponent_copy.Rsh(exponent_copy, 1)
	}
	return x.Mod(x, mod_copy)
}

// Функция для вычисления символа Якоби данного числа
func calculateJacobian(a, n *big.Int) *big.Int {
	a_copy := new(big.Int).Set(a)
	n_copy := new(big.Int).Set(n)
	if n_copy.Cmp(big.NewInt(0)) <= 0 || n_copy.Bit(0) == 0 {
		return big.NewInt(0)
	}

	ans := big.NewInt(1)

	if a_copy.Cmp(big.NewInt(0)) < 0 {
		a_copy.Neg(a_copy) // (a/n) = (-a/n)*(-1/n)
		if new(big.Int).Mod(n_copy, big.NewInt(4)).Cmp(big.NewInt(3)) == 0 {
			ans.Neg(ans) // (-1/n) = -1 если n = 3 (mod 4)
		}
	}

	if a_copy.Cmp(big.NewInt(1)) == 0 {
		return ans // (1/n) = 1
	}

	for a_copy.Cmp(big.NewInt(0)) != 0 {
		if a_copy.Cmp(big.NewInt(0)) < 0 {
			a_copy.Neg(a_copy) // (a/n) = (-a/n)*(-1/n)
			if new(big.Int).Mod(n_copy, big.NewInt(4)).Cmp(big.NewInt(3)) == 0 {
				ans.Neg(ans) // (-1/n) = -1 если n = 3 (mod 4)
			}
		}

		for a_copy.Bit(0) == 0 {
			a_copy.Rsh(a_copy, 1)
			if new(big.Int).Mod(n_copy, big.NewInt(8)).Cmp(big.NewInt(3)) == 0 || new(big.Int).Mod(n_copy, big.NewInt(8)).Cmp(big.NewInt(5)) == 0 {
				ans.Neg(ans)
			}
		}

		temp := new(big.Int).Set(a_copy)
		a_copy.Set(n_copy)
		n_copy.Set(temp)

		if new(big.Int).Mod(a_copy, big.NewInt(4)).Cmp(big.NewInt(3)) == 0 && new(big.Int).Mod(n_copy, big.NewInt(4)).Cmp(big.NewInt(3)) == 0 {
			ans.Neg(ans)
		}

		a_copy.Mod(a_copy, n_copy)
		if a_copy.Cmp(new(big.Int).Div(n_copy, big.NewInt(2))) > 0 {
			a_copy.Sub(a_copy, n_copy)
		}
	}

	if n_copy.Cmp(big.NewInt(1)) == 0 {
		return ans
	}

	return big.NewInt(0)
}

// Функция для выполнения теста простоты Соловея-Штрассена
func solovoyStrassen(p *big.Int, iteration int) bool {
	p_copy := new(big.Int).Set(p)
	if p_copy.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	if p_copy.Cmp(big.NewInt(2)) != 0 && p_copy.Bit(0) == 0 {
		return false
	}

	// Создаем объект для генератора случайных чисел
	rand := m_rand.New(m_rand.NewSource(0))

	for i := 0; i < iteration; i++ {

		// Генерируем случайное число r
		r := big.NewInt(0)
		r.Rand(rand, big.NewInt(0).SetBit(r, 63, 1))
		a := big.NewInt(0)
		a.Mod(r, big.NewInt(0).Sub(p_copy, big.NewInt(1))) // a = r % (p - 1)
		a.Add(a, big.NewInt(1))                            // a = a + 1
		jacobian := new(big.Int).Mod(new(big.Int).Add(p_copy, calculateJacobian(a, p_copy)), p_copy)
		mod := modulo(a, new(big.Int).Div(new(big.Int).Sub(p_copy, big.NewInt(1)), big.NewInt(2)), p_copy)

		if jacobian.Cmp(big.NewInt(0)) == 0 || mod.Cmp(jacobian) != 0 {
			return false
		}
	}
	return true
}

func GetCommonPrimeNumbers() (arr_p, arr_g []byte) {
	pMin := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)
	pMax := new(big.Int).Exp(big.NewInt(2), big.NewInt(384), nil)
	var p, g *big.Int
	for {
		p = generateNumber(pMin, pMax)
		if solovoyStrassen(p, 50) {
			break
		}
	}
	gMin := big.NewInt(2)
	gMax := big.NewInt(100)
	for {
		g = generateNumber(gMin, gMax)
		if solovoyStrassen(g, 50) {
			break
		}
	}
	return p.Bytes(), g.Bytes()
}

func GenerateKeyPair(p, g []byte) (publicKey_arr, privateKey_arr []byte) {
	min := big.NewInt(1)
	max := new(big.Int).Sub(new(big.Int).SetBytes(p), big.NewInt(1))

	var privateKey, publicKey *big.Int
	privateKey = generateNumber(min, max)

	publicKey = modExp(new(big.Int).SetBytes(g), privateKey, new(big.Int).SetBytes(p))

	return publicKey.Bytes(), privateKey.Bytes()
}

// ModExp Функция, которая возводит число x в степень y по модулю m
func modExp(x, y, m *big.Int) *big.Int {
	z := new(big.Int)
	return z.Exp(x, y, m)
}

func GetCommonSecretKey(publicKeyReceiver, privateKey, p []byte) []byte {
	return modExp(new(big.Int).SetBytes(publicKeyReceiver),
		new(big.Int).SetBytes(privateKey),
		new(big.Int).SetBytes(p)).Bytes()
}
