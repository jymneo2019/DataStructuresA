package main

import (
    "fmt"
)

// Rotor represents a single rotor in the Enigma machine.
type Rotor struct {
    wiring   string // The wiring of the rotor
    position int    // Current position of the rotor
}

// NewRotor creates a new rotor with the given wiring and initial position.
func NewRotor(wiring string, position int) *Rotor {
    return &Rotor{wiring: wiring, position: position}
}

// Rotate advances the rotor by one position.
func (r *Rotor) Rotate() {
    r.position = (r.position + 1) % 26 // Rotate position within the alphabet
}

// Encrypt encrypts a single character using the rotor.
func (r *Rotor) Encrypt(char byte) byte {
    offset := int(char-'A'+byte(r.position)) % 26
    return r.wiring[offset] // Return the encoded character
}

// Enigma represents the whole Enigma machine with multiple rotors.
type Enigma struct {
    rotors []*Rotor
}

// NewEnigma creates a new Enigma machine with the specified rotors.
func NewEnigma(rotors []*Rotor) *Enigma {
    return &Enigma{rotors: rotors}
}

// EncryptMessage encrypts a message using the Enigma machine.
func (e *Enigma) EncryptMessage(message string) string {
    encrypted := ""
    for i := 0; i < len(message); i++ {
        char := message[i]
        if char < 'A' || char > 'Z' {
            encrypted += string(char) // Non-alphabetic characters are unchanged
            continue
        }
        for _, rotor := range e.rotors {
            char = rotor.Encrypt(char)
            rotor.Rotate() // Rotate the rotor after each character
        }
        encrypted += string(char)
    }
    return encrypted
}

func main() {
    // Define rotors with simple wiring
    rotor1 := NewRotor("EKMFLGDQVZNTOWYHXUSPAIBRCJ", 0) // Rotor I wiring
    rotor2 := NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", 0) // Rotor II wiring
    rotor3 := NewRotor("BDFHJLCPRTXVZNYEIWGAKMUSQO", 0) // Rotor III wiring

    enigma := NewEnigma([]*Rotor{rotor1, rotor2, rotor3})

    message := "HELLO"
    encryptedMessage := enigma.EncryptMessage(message)
    fmt.Printf("Encrypted Message: %s\n", encryptedMessage)
}
