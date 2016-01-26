package main

import (
  "fmt"
)

type NameDescription interface {
  GetName() string
  GetDescription() string
}

type NameDescriptionTable struct {
  RowId int
  Name string
  Description string
}


type Platform NameDescriptionTable

func (p Platform) GetName() string {
  return p.Name
}

func (p Platform) GetDescription() string {
  return p.Description
}

func (p Platform) String() string {
  return fmt.Sprintf("%s (%s)", p.GetName(), p.GetDescription())
}



type Genre NameDescriptionTable

func (g Genre) GetName() string {
  return g.Name
}

func (g Genre) GetDescription() string {
  return g.Description
}

func (g Genre) String() string {
  return fmt.Sprintf("%s (%s)", g.GetName(), g.GetDescription())
}


type HardwareType NameDescriptionTable

func (ht HardwareType) GetName() string {
  return ht.Name
}

func (ht HardwareType) GetDescription() string {
  return ht.Description
}

func (ht HardwareType) String() string {
  return fmt.Sprintf("%s (%s)", ht.GetName(), ht.GetDescription())
}


type GameList struct {
  Games []Game
}

type Game struct {
  RowId int
  Title string
  Genre string
  Platform string
  NumberOwned int
  NumberBoxed int
  NumberOfManuals int
  DatePurchased string
  ApproximatePurchaseDate bool
  Notes string
}

func (g Game) String() string {
  return fmt.Sprintf("%s (%s)", g.Title, g.Platform)
}
