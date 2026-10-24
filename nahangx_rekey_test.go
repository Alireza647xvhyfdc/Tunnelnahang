package main

import (
	"bytes"
	"errors"
	"net"
	"testing"
)

// Crosses the rekey boundary (xRekey frames) in both directions.
func TestNahangxRekeyBothDirections(t *testing.T) {
	master := hsum([]byte("nahanx2|"), []byte("test-token"))
	a, b := net.Pipe()
	ch := make(chan *sess, 1)
	errc := make(chan error, 1)
	go func() {
		s, err := hsServerX(b, master)
		if err != nil {
			errc <- err
			return
		}
		ch <- s
	}()
	cs, err := hsClientX(a, master)
	if err != nil {
		t.Fatal("client handshake:", err)
	}
	var ss *sess
	select {
	case ss = <-ch:
	case err := <-errc:
		t.Fatal("server handshake:", err)
	}
	n := xRekey*2 + 7
	payload := []byte("x")
	done := make(chan error, 2)
	go func() {
		for i := 0; i < n; i++ {
			if err := cs.send(0, payload); err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()
	go func() {
		for i := 0; i < n; i++ {
			if err := ss.send(0, payload); err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()
	go func() {
		for i := 0; i < n; i++ {
			typ, d, err := ss.recv()
			if err == nil && (typ != 0 || !bytes.Equal(d, payload)) {
				err = errors.New("payload or type mismatch")
			}
			if err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()
	go func() {
		for i := 0; i < n; i++ {
			typ, d, err := cs.recv()
			if err == nil && (typ != 0 || !bytes.Equal(d, payload)) {
				err = errors.New("payload or type mismatch")
			}
			if err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()
	for i := 0; i < 4; i++ {
		if err := <-done; err != nil {
			t.Fatal("rekey transfer failed:", err)
		}
	}
}
