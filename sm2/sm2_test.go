/*
Copyright Suzhou Tongji Fintech Research Institute 2017 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sm2

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestSm2(t *testing.T) {
	priv, err := GenerateKey() // 生成密钥对
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", priv.Curve.IsOnCurve(priv.X, priv.Y)) // 验证是否为sm2的曲线
	pub := &priv.PublicKey
	msg := []byte("123456")
	d0, err := pub.Encrypt(msg)
	if err != nil {
		fmt.Printf("Error: failed to encrypt %s: %v\n", msg, err)
		return
	}
	// fmt.Printf("Cipher text = %v\n", d0)
	d1, err := priv.Decrypt(d0)
	if err != nil {
		fmt.Printf("Error: failed to decrypt: %v\n", err)
	}
	fmt.Printf("clear text = %s\n", d1)

	msg, _ = ioutil.ReadFile("ifile")             // 从文件读取数据
	sign, err := priv.Sign(rand.Reader, msg, nil) // 签名
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("ofile", sign, os.FileMode(0644))
	if err != nil {
		log.Fatal(err)
	}
	signdata, _ := ioutil.ReadFile("ofile")
	ok := priv.Verify(msg, signdata) // 密钥验证
	if ok != true {
		fmt.Printf("Verify error\n")
	} else {
		fmt.Printf("Verify ok\n")
	}
	pubKey := priv.PublicKey
	ok = pubKey.Verify(msg, signdata) // 公钥验证
	if ok != true {
		fmt.Printf("Verify error\n")
	} else {
		fmt.Printf("Verify ok\n")
	}

}

func BenchmarkSM2(t *testing.B) {
	t.ReportAllocs()
	msg := []byte("test")
	priv, err := GenerateKey() // 生成密钥对
	if err != nil {
		log.Fatal(err)
	}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		sign, err := priv.Sign(rand.Reader, msg, nil) // 签名
		if err != nil {
			log.Fatal(err)
		}
		priv.Verify(msg, sign) // 密钥验证
		// if ok != true {
		// 	fmt.Printf("Verify error\n")
		// } else {
		// 	fmt.Printf("Verify ok\n")
		// }
	}
}

func TestKEB2(t *testing.T) {
	ida := []byte{'1', '2', '3', '4', '5', '6', '7', '8',
		'1', '2', '3', '4', '5', '6', '7', '8'}
	idb := []byte{'1', '2', '3', '4', '5', '6', '7', '8',
		'1', '2', '3', '4', '5', '6', '7', '8'}
	daBuf := []byte{0x81, 0xEB, 0x26, 0xE9, 0x41, 0xBB, 0x5A, 0xF1,
		0x6D, 0xF1, 0x16, 0x49, 0x5F, 0x90, 0x69, 0x52,
		0x72, 0xAE, 0x2C, 0xD6, 0x3D, 0x6C, 0x4A, 0xE1,
		0x67, 0x84, 0x18, 0xBE, 0x48, 0x23, 0x00, 0x29}
	dbBuf := []byte{0x78, 0x51, 0x29, 0x91, 0x7D, 0x45, 0xA9, 0xEA,
		0x54, 0x37, 0xA5, 0x93, 0x56, 0xB8, 0x23, 0x38,
		0xEA, 0xAD, 0xDA, 0x6C, 0xEB, 0x19, 0x90, 0x88,
		0xF1, 0x4A, 0xE1, 0x0D, 0xEF, 0xA2, 0x29, 0xB5}
	raBuf := []byte{0XD4, 0XDE, 0X15, 0X47, 0X4D, 0XB7, 0X4D, 0X06,
		0X49, 0X1C, 0X44, 0X0D, 0X30, 0X5E, 0X01, 0X24,
		0X00, 0X99, 0X0F, 0X3E, 0X39, 0X0C, 0X7E, 0X87,
		0X15, 0X3C, 0X12, 0XDB, 0X2E, 0XA6, 0X0B, 0XB3}

	rbBuf := []byte{0X7E, 0x07, 0x12, 0x48, 0x14, 0xB3, 0x09, 0x48,
		0x91, 0x25, 0xEA, 0xED, 0x10, 0x11, 0x13, 0x16,
		0x4E, 0xBF, 0x0F, 0x34, 0x58, 0xC5, 0xBD, 0x88,
		0x33, 0x5C, 0x1F, 0x9D, 0x59, 0x62, 0x43, 0xD6}

	expk := []byte{0x6C, 0x89, 0x34, 0x73, 0x54, 0xDE, 0x24, 0x84,
		0xC6, 0x0B, 0x4A, 0xB1, 0xFD, 0xE4, 0xC6, 0xE5}

	curve := P256Sm2()
	curve.ScalarBaseMult(daBuf)
	da := new(PrivateKey)
	da.PublicKey.Curve = curve
	da.D = new(big.Int).SetBytes(daBuf)
	da.PublicKey.X, da.PublicKey.Y = curve.ScalarBaseMult(daBuf)

	db := new(PrivateKey)
	db.PublicKey.Curve = curve
	db.D = new(big.Int).SetBytes(dbBuf)
	db.PublicKey.X, db.PublicKey.Y = curve.ScalarBaseMult(dbBuf)

	ra := new(PrivateKey)
	ra.PublicKey.Curve = curve
	ra.D = new(big.Int).SetBytes(raBuf)
	ra.PublicKey.X, ra.PublicKey.Y = curve.ScalarBaseMult(raBuf)

	rb := new(PrivateKey)
	rb.PublicKey.Curve = curve
	rb.D = new(big.Int).SetBytes(rbBuf)
	rb.PublicKey.X, rb.PublicKey.Y = curve.ScalarBaseMult(rbBuf)

	k1, Sb, S2, err := KeyExchangeB(16, ida, idb, db, &da.PublicKey, rb, &ra.PublicKey)
	if err != nil {
		t.Error(err)
	}
	k2, S1, Sa, err := KeyExchangeA(16, ida, idb, da, &db.PublicKey, ra, &rb.PublicKey)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(k1, k2) != 0 {
		t.Error("key exchange differ")
	}
	if bytes.Compare(k1, expk) != 0 {
		t.Errorf("expected %x, found %x", expk, k1)
	}
	if bytes.Compare(S1, Sb) != 0 {
		t.Error("hash verfication failed")
	}
	if bytes.Compare(Sa, S2) != 0 {
		t.Error("hash verfication failed")
	}
}

func TestReadCertificateFromPem(t *testing.T) {
        // gm版本 cryptogen 工具 生成 的自签名 CA
	var certPem = `-----BEGIN CERTIFICATE-----
MIICUjCCAfegAwIBAgIQYYTprq/7P3K7xn2w6qhTkDAKBggqgRzPVQGDdTBzMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eu
b3JnMS5leGFtcGxlLmNvbTAeFw0yMDA0MTMwMzMzNTFaFw0zMDA0MTEwMzMzNTFa
MHMxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1T
YW4gRnJhbmNpc2NvMRkwFwYDVQQKExBvcmcxLmV4YW1wbGUuY29tMRwwGgYDVQQD
ExNjYS5vcmcxLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAE
eq64uveK/KX8jnZZ/5IOoIPpdSIV+gGmD2N8abr2/EKz5KE2zbxNRXcCvTUnO1pN
360Bk2YEk+T/BW4FFDjX+aNtMGswDgYDVR0PAQH/BAQDAgGmMB0GA1UdJQQWMBQG
CCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdDgQiBCBx
R+y7n5LEXeBKUxxXE/uhAEsC+ZWxpUdTkkVK0VAU1TAKBggqgRzPVQGDdQNJADBG
AiEAzIefgCfH8xEEOCVzMFwn3sBHlxT62qiVBAQa/RmfjuACIQDrDFELRK4UnDfp
Y2IzEADy1jvAgdSAJiU2EPCba5VpNg==
-----END CERTIFICATE-----`

        // TestSm2 生成的 CA 证书
	var selfTestGenerateCertPem = `-----BEGIN CERTIFICATE-----
MIIDMzCCAtqgAwIBAgIB/zAKBggqgRzPVQGDdTBIMQ0wCwYDVQQKEwRURVNUMRkw
FwYDVQQDExB0ZXN0LmV4YW1wbGUuY29tMQ8wDQYDVQQqEwZHb3BoZXIxCzAJBgNV
BAYTAk5MMB4XDTcwMDEwMTAwMTY0MFoXDTcwMDEwMjAzNDY0MFowSDENMAsGA1UE
ChMEVEVTVDEZMBcGA1UEAxMQdGVzdC5leGFtcGxlLmNvbTEPMA0GA1UEKhMGR29w
aGVyMQswCQYDVQQGEwJOTDBZMBMGByqGSM49AgEGCCqBHM9VAYItA0IABNK3zTaa
4T9a8w7LmIImcDLWl4Fqy7bk0LUlFlhU7ZAnP8CCIAP/ijlc5jvdRHFXkVY+GO6M
UdGasbNs/LMhnDijggGzMIIBrzAOBgNVHQ8BAf8EBAMCAgQwJgYDVR0lBB8wHQYI
KwYBBQUHAwIGCCsGAQUFBwMBBgIqAwYDgQsBMA8GA1UdEwEB/wQFMAMBAf8wXwYI
KwYBBQUHAQEEUzBRMCMGCCsGAQUFBzABhhdodHRwOi8vb2NzcC5leGFtcGxlLmNv
bTAqBggrBgEFBQcwAoYeaHR0cDovL2NydC5leGFtcGxlLmNvbS9jYTEuY3J0MEYG
A1UdEQQ/MD2CEHRlc3QuZXhhbXBsZS5jb22BEWdvcGhlckBnb2xhbmcub3JnhwR/
AAABhxAgAUhgAAAgAQAAAAAAAABoMA8GA1UdIAQIMAYwBAYCKgMwKgYDVR0eBCMw
IaAfMA6CDC5leGFtcGxlLmNvbTANggtleGFtcGxlLmNvbTBXBgNVHR8EUDBOMCWg
I6Ahhh9odHRwOi8vY3JsMS5leGFtcGxlLmNvbS9jYTEuY3JsMCWgI6Ahhh9odHRw
Oi8vY3JsMi5leGFtcGxlLmNvbS9jYTEuY3JsMBYGAyoDBAQPZXh0cmEgZXh0ZW5z
aW9uMA0GA1UdDgQGBAQEAwIBMAoGCCqBHM9VAYN1A0cAMEQCIAFl6HuA0qntdsGh
9SBSf6/JCtZmeSGuNbr1PgNRqDupAiA4UXAzrPBgAbIN3CWjQV28QCorLCvQ3Xct
fyXNzlRCtA==
-----END CERTIFICATE-----`
	cert, err := ReadCertificateFromMem([]byte(certPem))
	assert.NoError(t,err,"ReadCertificateFromMem Failed.")

	err = cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)
	assert.NoError(t,err,"CheckSignature Failed")


	cert2,err := ReadCertificateFromMem([]byte(selfTestGenerateCertPem))
	assert.NoError(t,err,"ReadCertificateFromMem Failed.")

	err = cert.CheckSignature(cert2.SignatureAlgorithm, cert2.RawTBSCertificate, cert2.Signature)
	assert.NoError(t,err,"CheckSignature Failed")
}
