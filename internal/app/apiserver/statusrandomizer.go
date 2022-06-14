package apiserver

import "math/rand"

func randomizer() string {
	i := rand.Intn(100)
	if i < 90 {
		return "new"
	} else {
		return "error"
	}
}
