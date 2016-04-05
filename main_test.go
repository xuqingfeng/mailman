package main

import (
    "os"
    "testing"
    "net"
)

func TestMain(m *testing.M) {

    os.Exit(m.Run())
}

func TestIsTCPPortAvailable(t *testing.T) {

    conn, err := net.Listen("tcp", "127.0.0.1:8888")
    if err != nil {
        t.Errorf("can't listen tcp 127.0.0.1:8888 %v", err)
    }
    if isTCPPortAvailable(8888) {
        t.Errorf("port 8888 is available")
    }
    conn.Close()
}
