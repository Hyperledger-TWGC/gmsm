// +build !amd64,!386

package sm3

var (
	FF0 = ff0Generic
	FF1 = ff1Generic
	GG0 = gg0Generic
	GG1 = gg1Generic
	P0  = p0Generic
	P1  = p1Generic
)

func ff0Generic(x, y, z uint32) uint32 { return x ^ y ^ z }

func ff1Generic(x, y, z uint32) uint32 { return (x & y) | (x & z) | (y & z) }

func gg0Generic(x, y, z uint32) uint32 { return x ^ y ^ z }

func gg1Generic(x, y, z uint32) uint32 { return (x & y) | (^x & z) }

func leftRotate(x uint32, i uint32) uint32 { return (x<<(i%32) | x>>(32-i%32)) }

func p0Generic(x uint32) uint32 {
	return x ^ leftRotate(x, 9) ^ leftRotate(x, 17)
}

func p1Generic(x uint32) uint32 {
	return x ^ leftRotate(x, 15) ^ leftRotate(x, 23)
}
