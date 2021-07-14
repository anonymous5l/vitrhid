package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"vitrhid/growcastle"

	"golang.org/x/sys/unix"
)

type Device struct {
	Addr      string
	Control   int
	Interrupt int
	Close     chan struct{}
	Disposed  bool
}

func (d *Device) Send(x, y, tip int8) {
	// reset
	buf := &bytes.Buffer{}
	buf.Write([]byte{0xA1, growcastle.MouseReportId})
	binary.Write(buf, binary.LittleEndian, tip) // Tip Switch
	binary.Write(buf, binary.LittleEndian, x)   // x
	binary.Write(buf, binary.LittleEndian, y)   // y

	b := buf.Bytes()

	if _, err := unix.Write(d.Interrupt, b); err != nil {
		d.Disposed = true
		d.Stop()
		return
	}

	time.Sleep(time.Millisecond * 80)
}

func (d *Device) internalRun(cnt, step int8) {
	// loop trigger
	for i := 0; i < 5; i++ {
		d.Send(0, 0, 0)       // first trigger
		d.Send(-127, -127, 0) // reset
		time.Sleep(time.Millisecond * 220)
	}

	// move to replay
	for i := int8(0); i < cnt; i++ {
		d.Send(0, step, 0)
	}

	d.Send(30, 0, 0)
	d.Send(0, 0, 1) // down
	d.Send(0, 0, 0) // up
	// move to first item
	d.Send(0, -40, 0)
	d.Send(0, 0, 1) // down
	d.Send(0, 0, 0) // up
}

func (d *Device) Start(sleep time.Duration) {
	if d.Close == nil {
		d.Close = make(chan struct{}, 1)
	}

	for {
		d.internalRun(8, 39)
		select {
		case <-time.After(time.Second * sleep):
			break
		case <-d.Close:
			return
		}
	}
}

func (d *Device) Stop() {
	if d.Close != nil {
		close(d.Close)
		d.Close = nil
	}
}

type Services struct {
	lock    sync.RWMutex
	devices map[string]*Device
	isStart byte
}

func NewServices() *Services {
	s := &Services{}
	s.devices = make(map[string]*Device)

	return s
}

func (s *Services) disconnect(addr string) {
	d, ok := s.devices[addr]
	if ok {
		unix.Close(d.Control)
		unix.Close(d.Interrupt)
		if d.Close != nil {
			close(d.Close)
		}
		delete(s.devices, addr)
	}
}

func (s *Services) AcceptControl() {
	for {
		fd, addr, err := unix.Accept(controlListenFd)
		if err != nil {
			log.Printf("accept: control %s", err)
			continue
		}
		l2addr := addr.(*unix.SockaddrL2)
		strAddr := hex.EncodeToString(l2addr.Addr[:])

		s.lock.Lock()
		d, ok := s.devices[strAddr]
		if ok {
			d.Control = fd
			d.Addr = strAddr
			d.Disposed = false
		} else {
			s.devices[strAddr] = &Device{Control: fd}
		}
		s.lock.Unlock()
	}
}

func (s *Services) AcceptInterrupt() {
	for {
		fd, addr, err := unix.Accept(interruptListenFd)
		if err != nil {
			log.Printf("accept: interrupt %s", err)
			continue
		}
		l2addr := addr.(*unix.SockaddrL2)
		strAddr := hex.EncodeToString(l2addr.Addr[:])

		s.lock.Lock()
		d, ok := s.devices[strAddr]
		if ok {
			d.Interrupt = fd
			d.Addr = strAddr
			d.Disposed = false
		} else {
			s.devices[strAddr] = &Device{Interrupt: fd}
		}
		s.lock.Unlock()
	}
}

func (s *Services) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/start" {
		if len(s.devices) == 0 {
			rw.Write([]byte("no devices"))
			return
		}
		if s.isStart == 0 {
			t := r.URL.Query().Get("t")
			d := r.URL.Query().Get("delay")
			it := time.Duration(43)
			if i, err := strconv.ParseInt(t, 10, 64); err != nil {
				rw.Write([]byte("invalid param"))
				return
			} else {
				it = time.Duration(i)
			}
			if d != "" {
				if i, err := strconv.ParseInt(d, 10, 64); err != nil {
					rw.Write([]byte("invalid delay param"))
					return
				} else {
					time.Sleep(time.Second * time.Duration(i))
				}
			}
			s.isStart = 1
			// pick first one
			for _, v := range s.devices {
				if !v.Disposed {
					log.Printf("Device %s Start", v.Addr)
					go v.Start(it)
				}
			}
		}
		rw.Write([]byte("success"))
	}

	if r.URL.Path == "/stop" {
		if s.isStart == 1 {
			for _, v := range s.devices {
				if !v.Disposed {
					log.Printf("Device %s Stop", v.Addr)
					v.Stop()
				}
			}
			s.isStart = 0
		}
		rw.Write([]byte("success"))
	}
}
