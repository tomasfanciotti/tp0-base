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
	OP_CODE_AGENCY_READY   = 5
	OP_CODE_ASK_WINNER     = 6
	OP_CODE_WINNERS        = 7
	OP_CODE_SERVER_BUSSY   = 8
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

func (l *Lottery) ready(client_id string) (bool, error) {

	arguments := []string{client_id}
	packet := NewPacket(OP_CODE_AGENCY_READY, arguments)
	Send(*l.conn, packet)
	response, err := Receive(*l.conn)
	if err != nil {
		return false, err
	}

	switch response.opcode {
	case OP_CODE_ZERO:
		return false, errors.New("Conexión cerrada")

	case OP_CODE_ERROR:
		return false, errors.New(string(response.data))

	case OP_CODE_REGISTER_ACK:
		return true, nil

	default:
		return false, errors.New("unexpected response code: " + string(response.opcode))
	}
}

func (l *Lottery) winner(agency_id string) ([]string, error) {

	arguments := []string{agency_id}
	ask_winner := NewPacket(OP_CODE_ASK_WINNER, arguments)
	Send(*l.conn, ask_winner)
	response, err := Receive(*l.conn)
	if err != nil {
		return nil, err
	}

	// Conexion cerrada
	if response.opcode == OP_CODE_ZERO {
		return nil, errors.New("Conexión cerrada")
	}

	// Error Server
	if response.opcode == OP_CODE_ERROR {
		return nil, errors.New(string(response.data))
	}

	// Servidor ocupado
	if response.opcode == OP_CODE_SERVER_BUSSY {
		return nil, nil
	}

	// Otro OpCode
	if response.opcode != OP_CODE_WINNERS {
		return nil, errors.New("unexpected opcode: " + string(response.opcode))
	}

	// Winners
	winners, ok := decode(response.data).([]string)
	if !ok {
		return nil, errors.New("unexpected data: " + string(response.data))
	}

	return winners, nil
}
