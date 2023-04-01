package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	WAITING_PERIOD = 4
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config  ClientConfig
	conn    net.Conn
	running bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config:  config,
		running: false,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

func (c *Client) SetupGracefulShutdown() {

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		c.conn.Close()
		log.Errorf("action: receive_message | result: fail | client_id: %v ", c.config.ID)
	}()
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// autoincremental msgID to identify every message sent
	msgID := 1
	c.running = true

loop:
	// Send messages if the loopLapse threshold has not been surpassed
	for timeout := time.After(c.config.LoopLapse); ; {
		select {
		case <-timeout:
			log.Infof("action: timeout_detected | result: success | client_id: %v",
				c.config.ID,
			)
			break loop
		default:
			if !c.running {
				break
			}
		}

		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		// TODO: Modify the send to avoid short-write
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		msgID++
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) LoadSingleBet(bet *Bet) {

	c.createClientSocket()

	lottery := NewLottery(&c.conn)
	res, err := lottery.almacenar_apuesta(bet, c.config.ID)

	if err != nil || !res {
		log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	c.conn.Close()

	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", bet.Dni, bet.Numero)

}

func (c *Client) LoadBatchBets(chunkFile string, batchSize int) {

	c.running = true
	file, err := os.Open(chunkFile)
	if err != nil {
		log.Errorf("action: open_chunk_file | result: fail | err: %s", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	chunk := []Bet{}
	total_bets := 0

	lottery := NewLottery(&c.conn)

	for scanner.Scan() {
		campos := strings.Split(strings.TrimRight(scanner.Text(), "\n"), ",")

		if len(campos) != 5 {
			log.Info("action: scan_chunk_file | result: warning | msg: line fields does not match with a bet register. ignoring")
			continue
		}

		bet := newBet(c.config.ID, campos[0], campos[1], campos[2], campos[3], campos[4])
		chunk = append(chunk, bet)

		if len(chunk) >= batchSize {

			c.createClientSocket()
			if _, err := lottery.almacenar_bacth(chunk); err != nil {
				log.Errorf("action: send_chunk | result: fail | err: %s", err)
			} else {
				total_bets += 1
			}
			c.conn.Close()
			chunk = []Bet{}
		}
	}

	c.createClientSocket()
	if _, err := lottery.almacenar_bacth(chunk); err != nil {
		log.Errorf("action: send_chunk | result: fail | err: %s", err)
	}
	c.conn.Close()

	// Verificar si hubo algún error durante la lectura del archivo.
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	log.Infof("action: send_chunk | result: success | total: %d", total_bets)

	// Termine
	c.createClientSocket()
	if _, err := lottery.ready(c.config.ID); err != nil {
		log.Errorf("action: ready | result: fail | err: %s", err)
	}
	c.conn.Close()

	// Ask for winner
	for c.running {

		c.createClientSocket()
		result, err := lottery.winner(c.config.ID)
		if err != nil {
			log.Errorf("action: consulta_ganadores | result: fail | err: %s", err)
			break
		}
		c.conn.Close()

		if result != nil {
			log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d | ganadores: %v", len(result), result)
			break
		}
		log.Errorf("action: consulta_ganadores | result: fail | msgs: servidor_ocupado")
		time.Sleep(WAITING_PERIOD * time.Second)
	}
}
