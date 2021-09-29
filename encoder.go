package ec11

import (
	"errors"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
	"strconv"
)


type stackT struct {
	s []uint8
}

const (
	clockwise        string = "032"
	counterclockwise string = "012"
)

type EncoderT struct {
	dtPin   string
	clkPin  string
	swPin   string
	p1      gpio.PinIO
	p2      gpio.PinIO
	p3      gpio.PinIO
	resChan chan int8
}

func New(dtPin, clkPin, swPin int) (EncoderT, error) {
	var err error
	e := EncoderT{
		dtPin:   strconv.Itoa(dtPin),
		clkPin:  strconv.Itoa(clkPin),
		swPin:   strconv.Itoa(swPin),
		resChan: make(chan int8, 20),
	}

	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		return e, err
	}

	// Lookup a pin by its number:
	e.p1 = gpioreg.ByName(e.dtPin) //dt
	if e.p1 == nil {
		return e, errors.New("Failed to find GPIO " + e.dtPin)
	}

	e.p2 = gpioreg.ByName(e.clkPin) //clk
	if e.p2 == nil {
		return e, errors.New("Failed to find GPIO " + e.clkPin)
	}

	e.p3 = gpioreg.ByName(e.swPin) //sw
	if e.p3 == nil {
		return e, errors.New("Failed to find GPIO " + e.swPin)
	}

	// setup pins mode
	err = e.p1.In(gpio.PullNoChange, gpio.BothEdges)
	if err != nil {
		return e, err
	}
	err = e.p2.In(gpio.PullNoChange, gpio.BothEdges)
	if err != nil {
		return e, err
	}
	err = e.p3.In(gpio.PullNoChange, gpio.BothEdges)
	if err != nil {
		return e, err
	}
	return e, nil
}

func (e EncoderT) Start() chan int8 {
	go e.encoder()
	return e.resChan
}

func (e EncoderT) encoder() {

	stack := stackT{
		make([]uint8, 3),
	}

	var current, previous uint8
	chW := make(chan bool, 10)

	// Wait for edges as detected by the hardware, and print the value read:
	go func() {
		for {
			e.p1.WaitForEdge(-1)
			if len(chW) > 0 {
				continue
			}
			chW <- true
		}
	}()

	go func() {
		for {
			e.p2.WaitForEdge(-1)
			if len(chW) > 0 {
				continue
			}
			chW <- true
		}
	}()

	go func() {
		for {
			e.p3.WaitForEdge(-1)
			if len(chW) > 0 {
				continue
			}
			chW <- true
		}
	}()

	var v3p gpio.Level
	for range chW {
		v1 := e.p1.Read()
		v2 := e.p2.Read()
		v3 := e.p3.Read()

		if v3 != v3p {
			if v3 == gpio.Low {
				e.resChan <- 0
				//fmt.Println("Button pressed")
			}
			v3p = v3
		}

		b1 := convertLevel(v1)
		b2 := convertLevel(v2)

		current = (b1 ^ b2) | b2<<1
		if current == previous {
			continue
		}

		previous = current
		stack.push(current)
		s := stack.popString()
		if s == clockwise {
			//fmt.Println("Rotate clockwise")
			e.resChan <- 1
		}
		if s == counterclockwise {
			//fmt.Println("Rotate counter clockwise")
			e.resChan <- -1
		}
	}
}

func convertLevel(l gpio.Level) uint8 {
	var b uint8
	switch l {
	case gpio.High:
		b = 1
	case gpio.Low:
		b = 0
	}
	return b
}

func (s *stackT) push(i uint8) {
	copy(s.s, s.s[1:3])
	s.s[2] = i
}

func (s stackT) popString() string {
	sum := ""
	for _, j := range s.s {
		sum += strconv.Itoa(int(j))
	}
	return sum
}
