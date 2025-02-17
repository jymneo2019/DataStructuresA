package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// generateKeyPair generates an RSA key pair with primes of the specified bit size.
// It returns the modulus n, public exponent e, private exponent d, and phi (Euler's totient).
func generateKeyPair(bits int) (*big.Int, *big.Int, *big.Int) {
	// Generate two distinct random primes p and q.
	p, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	q, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		panic(err)
	}

	// Compute n = p * q.
	n := new(big.Int).Mul(p, q)

	// Compute phi(n) = (p - 1) * (q - 1).
	one := big.NewInt(1)
	pMinus1 := new(big.Int).Sub(p, one)
	qMinus1 := new(big.Int).Sub(q, one)
	phi := new(big.Int).Mul(pMinus1, qMinus1)

	// Choose a public exponent e (commonly 65537).
	e := big.NewInt(65537)
	// Ensure that gcd(e, phi) = 1.
	gcd := new(big.Int).GCD(nil, nil, e, phi)
	if gcd.Cmp(one) != 0 {
		// If not, choose a different odd number for e.
		e = big.NewInt(3)
		for {
			gcd = new(big.Int).GCD(nil, nil, e, phi)
			if gcd.Cmp(one) == 0 {
				break
			}
			e.Add(e, big.NewInt(2))
		}
	}

	// Compute d, the modular inverse of e modulo phi.
	d := new(big.Int).ModInverse(e, phi)
	if d == nil {
		panic("Failed to compute modular inverse")
	}

	return n, e, d
}

// encrypt encrypts a message (as a big.Int) using the public key (e, n).
// It returns the ciphertext c = message^e mod n.
func encrypt(message, e, n *big.Int) *big.Int {
	ciphertext := new(big.Int).Exp(message, e, n)
	return ciphertext
}

// decrypt decrypts a ciphertext (as a big.Int) using the private key (d, n).
// It returns the original message m = ciphertext^d mod n.
func decrypt(ciphertext, d, n *big.Int) *big.Int {
	plaintext := new(big.Int).Exp(ciphertext, d, n)
	return plaintext
}

func main() {
	// Choose a bit size for the primes (for demonstration purposes, 512 bits is used).
	// In real applications, use at least 2048 bits.
	bits := 512

	// Generate the RSA key pair.
	n, e, d := generateKeyPair(bits)
	fmt.Println("Public Key (n, e):")
	fmt.Printf("n = %s\n", n.String())
	fmt.Printf("e = %s\n\n", e.String())

	fmt.Println("Private Key (d):")
	fmt.Printf("d = %s\n\n", d.String())

	// For demonstration, use a small integer as the message.
	message := big.NewInt(12345)
	fmt.Println("Original Message:")
	fmt.Println(message)

	// Encrypt the message.
	ciphertext := encrypt(message, e, n)
	fmt.Println("\nEncrypted Ciphertext:")
	fmt.Println(ciphertext)

	// Decrypt the ciphertext.
	decryptedMessage := decrypt(ciphertext, d, n)
	fmt.Println("\nDecrypted Message:")
	fmt.Println(decryptedMessage)
}
