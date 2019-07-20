package utility

import (
	"crypto/rand"

	randomdata "github.com/Pallinder/go-randomdata"
)

type Generate struct{}

func (g *Generate) Token(len int) string {
	t := make([]byte, len)
	rand.Read(t)
	return string(t)
}

func (g *Generate) FirstName() string {
	return randomdata.FirstName(randomdata.RandomGender)
}

func (g *Generate) FemaleFirstName() string {
	return randomdata.FirstName(randomdata.Female)
}

func (g *Generate) MaleFirstName() string {
	return randomdata.FirstName(randomdata.Male)
}

func (g *Generate) LastName() string {
	return randomdata.LastName()
}

func (g *Generate) Email() string {
	return randomdata.Email()
}
