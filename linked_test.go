package main

import "testing"

// خروجی ss -Htn state established: State Recv-Q Send-Q Local Peer
func TestLinkedCountsServerAndClientConnections(t *testing.T) {
	lines := []string{
		"ESTAB 0 0 10.0.0.5:3080 203.0.113.9:51234",
		"ESTAB 0 0 10.0.0.5:3080 203.0.113.9:51235",
		"ESTAB 0 0 10.0.0.5:443 198.51.100.7:40000",
		"ESTAB 0 0 192.168.1.2:45678 10.0.0.5:3080",
		"ESTAB 0 0 192.168.1.2:45679 10.0.0.9:3080",
		"",
	}
	server := Tunnel{Role: "server", Addr: "0.0.0.0:3080"}
	if n := linked(server, lines); n != 2 {
		t.Fatalf("server: want 2 connections on :3080, got %d", n)
	}
	client := Tunnel{Role: "client", Addr: "10.0.0.5:3080"}
	if n := linked(client, lines); n != 1 {
		t.Fatalf("client: want 1 connection to 10.0.0.5:3080, got %d", n)
	}
	if n := linked(Tunnel{Role: "client", Addr: "example.com:3080"}, lines); n != -1 {
		t.Fatalf("client with hostname should be unknown (-1), got %d", n)
	}
}
