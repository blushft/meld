package utility

import (
	"crypto/rand"

	randomdata "github.com/Pallinder/go-randomdata"
)

type Mock struct{}

func (g *Mock) Token(len int) string {
	t := make([]byte, len)
	rand.Read(t)
	return string(t)
}

func (g *Mock) FirstName() string {
	return randomdata.FirstName(randomdata.RandomGender)
}

func (g *Mock) FemaleFirstName() string {
	return randomdata.FirstName(randomdata.Female)
}

func (g *Mock) MaleFirstName() string {
	return randomdata.FirstName(randomdata.Male)
}

func (g *Mock) LastName() string {
	return randomdata.LastName()
}

func (g *Mock) Email() string {
	return randomdata.Email()
}
