package common

import (
	"errors"
	"net"
)

// OpCodes
const (
	OP_CODE_REGISTER     = 1
	OP_CODE_REGISTER_ACK = 2
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

func (l *Lottery) almacenar_apuesta(bet *Bet, client_id string) (bool, error) {

	arguments := []string{
		client_id,
		bet.Nombre,
		bet.Apellido,
		bet.Dni,
		bet.Nacimiento,
		bet.Numero,
	}
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
