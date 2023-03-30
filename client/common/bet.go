package common

type Bet struct {
	Agencia    string
	Nombre     string
	Apellido   string
	Dni        string
	Nacimiento string
	Numero     string
}

func newBet(agencia string, nombre string, apellido string, dni string, nacimiento string, numero string) Bet {
	return Bet{
		Agencia:    agencia,
		Nombre:     nombre,
		Apellido:   apellido,
		Dni:        dni,
		Nacimiento: nacimiento,
		Numero:     numero,
	}
}
