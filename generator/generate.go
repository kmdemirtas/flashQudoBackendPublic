package generator

import (
	"log"
	"strings"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/ironstone95/FlashQudoV2/model"
)

func GetRandomUser() *model.User {
	u := model.User{}
	err := faker.FakeData(&u)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &u
}

func GetRandomGroup() *model.Group {
	g := model.Group{}
	err := faker.FakeData(&g)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &g
}

func GetRandomBundle(groupID string) *model.Bundle {
	b := model.Bundle{}
	err := faker.FakeData(&b)
	if err != nil {
		log.Println(err)
		return nil
	}
	b.GroupID = groupID
	return &b
}

func GetRandomCard(bundleID string) *model.Card {
	c := model.Card{}
	err := faker.FakeData(&c)
	if err != nil {
		log.Println(err)
		return nil
	}
	c.BundleID = bundleID
	return &c
}

func CreateID() string {
	id := uuid.NewString()
	id = strings.ReplaceAll(id, "-", "")
	return id
}
