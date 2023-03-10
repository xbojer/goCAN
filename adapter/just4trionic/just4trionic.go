package just4trionic

import (
	"context"

	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/roffe/gocan"
	"go.bug.st/serial"
)

var debug bool

func init() {
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		debug = true
	}
}

type Adapter struct {
	cfg        *gocan.AdapterConfig
	port       serial.Port
	send, recv chan gocan.CANFrame
	close      chan struct{}

	canRate, filter string
	closed          bool
}

func New(cfg *gocan.AdapterConfig) (gocan.Adapter, error) {
	adapter := &Adapter{
		cfg:   cfg,
		send:  make(chan gocan.CANFrame, 10),
		recv:  make(chan gocan.CANFrame, 10),
		close: make(chan struct{}, 1),
	}

	for _, f := range cfg.CANFilter {
		if f == 0x05 {
			adapter.filter = "t5"
			break
		}
		if f == 0x220 {
			adapter.filter = "f7"
			break
		}
		if f == 0x7E0 {
			adapter.filter = "f8"
			break
		}
	}

	if err := adapter.setCANrate(cfg.CANRate); err != nil {
		return nil, err
	}

	return adapter, nil
}

func (a *Adapter) Name() string {
	return "Just4Trionic"
}

func (a *Adapter) Init(ctx context.Context) error {
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	p, err := serial.Open(a.cfg.Port, mode)
	if err != nil {
		return fmt.Errorf("failed to open com port %q : %v", a.cfg.Port, err)
	}
	p.SetReadTimeout(1 * time.Millisecond)
	a.port = p

	p.ResetOutputBuffer()

	var cmds = []string{
		"\x1B", // Empty buffer
		"O",    // enter canbus mode
		a.filter,
		a.canRate, // Setup CAN bit-rates
		//a.mask,
	}

	delay := time.Duration(5 * time.Millisecond)

	for n, c := range cmds {
		if n == 3 {
			p.ResetInputBuffer()
		}
		if debug {
			log.Printf("sending: %s", c)
		}
		_, err := p.Write([]byte(c + "\r"))
		if err != nil {
			p.Close()
			return err
		}
		time.Sleep(delay)
	}

	go a.recvManager(ctx)
	go a.sendManager(ctx)

	return nil
}

func (a *Adapter) Recv() <-chan gocan.CANFrame {
	return a.recv
}

func (a *Adapter) Send() chan<- gocan.CANFrame {
	return a.send
}

func (a *Adapter) Close() error {
	a.closed = true
	a.close <- struct{}{}
	time.Sleep(50 * time.Millisecond)
	a.port.Write([]byte("\x1B"))
	time.Sleep(10 * time.Millisecond)
	return a.port.Close()
}

func (a *Adapter) setCANrate(rate float64) error {
	switch rate {
	case 10:
		a.canRate = "S0"
	case 20:
		a.canRate = "S1"
	case 50:
		a.canRate = "S2"
	case 100:
		a.canRate = "S3"
	case 125:
		a.canRate = "S4"
	case 250:
		a.canRate = "S5"
	case 500:
		a.canRate = "S6"
	case 615.384:
		a.canRate = "s2"
	case 800:
		a.canRate = "S7"
	case 1000:
		a.canRate = "S8"
	default:
		return fmt.Errorf("unknown rate: %f", rate)

	}
	return nil
}

func (a *Adapter) recvManager(ctx context.Context) {
	buff := bytes.NewBuffer(nil)
	readBuffer := make([]byte, 8)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		n, err := a.port.Read(readBuffer)
		if err != nil {
			if !a.closed {
				log.Printf("failed to read com port: %v", err)
			}
			return
		}
		if n == 0 {
			continue
		}
		a.parse(ctx, readBuffer[:n], buff)
	}
}

func (a *Adapter) parse(ctx context.Context, readBuffer []byte, buff *bytes.Buffer) {
	for _, b := range readBuffer {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if b == 0x0D || b == 0x0A {
			if buff.Len() == 0 {
				continue
			}
		}
		if b == 0x0A {
			by := buff.Bytes()
			switch by[0] {
			case 'w':
				f, err := a.decodeFrame(by[1 : buff.Len()-1])
				if err != nil {
					log.Printf("failed to decode frame: %v %s", err, by)
					continue
				}
				select {
				case a.recv <- f:
				default:
					log.Println("dropped frame")
				}
				buff.Reset()
			default:
				//log.Printf("COM>> %q\n", buff.String())
			}
			buff.Reset()
			continue
		}
		buff.WriteByte(b)
	}
}

func (*Adapter) decodeFrame(buff []byte) (gocan.CANFrame, error) {
	id, err := strconv.ParseUint(string(buff[0:3]), 16, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to decode identifier: %v", err)
	}
	data := make([]byte, hex.DecodedLen(int(buff[3]-0x30)*2))
	if _, err := hex.Decode(data, buff[4:]); err != nil {
		return nil, fmt.Errorf("failed to decode frame body: %v", err)
	}
	return gocan.NewFrame(
		uint32(id),
		data,
		gocan.Incoming,
	), nil
}

func (a *Adapter) sendManager(ctx context.Context) {
	var f string
	for {
		select {
		case v := <-a.send:
			f = "t" + strconv.FormatUint(uint64(v.Identifier()), 16) +
				strconv.Itoa(v.Length()) +
				hex.EncodeToString(v.Data())

			for i := v.Length(); i < 8; i++ {
				f += "00"
			}
			f += "\r"
			_, err := a.port.Write([]byte(f))
			if err != nil {
				log.Printf("failed to write to com port: %q, %v", f, err)
			}
			if debug {
				log.Printf("%q\n", f)
			}
			f = ""
		case <-ctx.Done():
			return
		case <-a.close:
			return
		}
	}
}
