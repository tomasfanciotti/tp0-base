package common

import (
	"errors"
	"net"
	"strconv"
)

// OpCodes
const (
    OP_CODE_ZERO           = 0
	OP_CODE_REGISTER       = 1
	OP_CODE_REGISTER_ACK   = 2
	OP_CODE_REGISTER_BATCH = 3
	OP_CODE_ERROR          = 4
)

// Server interface

type Lottery struct {
	conn *net.Conn
}

func NewLottery(conn *net.Conn) *Lottery {
	lottery := &Lottery{
		conn: conn,
	}
	return lottery
}

func (l *Lottery) getArgumentsFromBet(bet *Bet) []string {

	arguments := []string{
		bet.Agencia,
		bet.Nombre,
		bet.Apellido,
		bet.Dni,
		bet.Nacimiento,
		bet.Numero,
	}

	return arguments
}

func (l *Lottery) almacenar_apuesta(bet *Bet, client_id string) (bool, error) {

	arguments := l.getArgumentsFromBet(bet)
	new_packet := NewPacket(OP_CODE_REGISTER, arguments)
	Send(*l.conn, new_packet)

	packet_response, err := Receive(*l.conn)

	if err != nil {
		return false, err
	}

	if packet_response.opcode != OP_CODE_REGISTER_ACK {
		return false, errors.New("Servidor NO devolvió ACK")
	}

	return true, nil
}

func (l *Lottery) almacenar_bacth(batch []Bet) (bool, error) {

	arguments := []string{strconv.Itoa(len(batch))}
	for i, v := range batch {
		arguments = append(arguments, "@"+strconv.Itoa(i))
		arguments = append(arguments, l.getArgumentsFromBet(&v)...)
	}

	new_packet := NewPacket(OP_CODE_REGISTER_BATCH, arguments)
	Send(*l.conn, new_packet)

	packet_response, err := Receive(*l.conn)

	if err != nil {
		return false, err
	}

	if packet_response.opcode != OP_CODE_REGISTER_ACK {
		return false, errors.New("Servidor NO devolvió ACK")
	}

	return true, nil
}
