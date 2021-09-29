package ec11

import (
	"fmt"
	"testing"
)

var e EncoderT

func TestEncoder(t *testing.T) {
	var err error
	e, err = New(15, 16, 14)
	if err != nil {
		t.Error(err)
	}
}

func TestEncoder_t_Encoder(t *testing.T) {
	ch:=e.Start()
	for i := range ch {
		fmt.Println(i)
	}
}
