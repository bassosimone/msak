package measurer_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/m-lab/go/rtx"
	"github.com/m-lab/msak/internal/measurer"
	"github.com/m-lab/msak/internal/netx"
)

func TestNdt8MeasurerUnix_Start(t *testing.T) {
	// Create a TCP socket to test. net.FileConn is only supported on unix-like
	// operating systems.
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	rtx.Must(err, "cannot create test socket")
	fp := os.NewFile(uintptr(fd), "test-socket")
	conn, err := net.FileConn(fp)
	rtx.Must(err, "cannot create net.Conn")

	netxConn := &netx.Conn{
		Conn: conn,
	}
	mc := &mockWSConn{
		underlyingConn: netxConn,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mchan := measurer.Start(ctx, mc)
	select {
	case <-mchan:
		fmt.Println("received measurement")
	case <-time.After(1 * time.Second):
		t.Fatalf("did not receive any measurement")
	}
}
