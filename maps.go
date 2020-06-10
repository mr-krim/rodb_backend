package main

import "fmt"

type mobStruct struct {
	id, level, hp, def, mdef, hit, flee, atk, matk uint16
	name, race, element, size, location            string
}

func main() {
	simpleMob := make(map[uint8]mobStruct)
	simpleMob[1] = mobStruct{1, 15, 570, 16, 3, 151, 9, 16, 8, "Rocker", "Insect", "Earth", "Medium", "Prontera Field"}
	simpleMob[2] = mobStruct{2, 14, 338, 22, 0, 142, 9, 13, 7, "Savage Babe", "Brute", "Earth", "Small", "Geffen East Field, Prontera Field"}
	fmt.Printf("%+v\n", simpleMob[2])
}
