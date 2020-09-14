package main

import (
    "log"
    "math"
    "strings"
    "strconv"
    pkgerrors "github.com/pkg/errors"
)

type Range struct {
    ip string
    min int
    max int
    masks [32]byte
}

func NewRange(min int, max int) *Range {
    r := Range{ip: "192.168.0.0", min: min, max: max, masks: [32]byte{}}
    for i:=0; i<32; i++ {
        r.masks[i] = 0
    }

    return &r
}

func(r *Range) base() string {
    index := strings.LastIndex(r.ip, ".")
    if index == -1 {
        return r.ip
    } else {
        return r.ip[0:index+1]
    }
}

func(r *Range) Print() {
    log.Println(r)
}

func(r *Range) Allocate() (string, error) {
    i := r.min
    index := (r.min-1)/8
    base := byte(math.Exp2(float64(7-(r.min-1)%8)))
    for i <= r.max {
        if r.masks[index] & base == 0 {
            r.masks[index] |= base
            return r.base() + strconv.Itoa(i), nil
        }
        if (i % 8 == 0) {
            base = 0x80
            index += 1
            for r.masks[index] == 0xff {
                log.Println("by pass", index)
                i += 8
                index += 1
            }
        } else {
            base = base / 2
        }
        i = i + 1
    }
    
    return "", pkgerrors.New("No available IP")
}

func(r *Range) Free(sip string) error {
    ip := 0
    i := strings.LastIndex(sip, ".")
    if i == -1 {
        return pkgerrors.New("invalid ip")
    } else {
        base_ip := sip[0:i+1]
        if r.base() != base_ip {
            return pkgerrors.New("ip is not in range")
        }

        ip, _ = strconv.Atoi(sip[i+1:len(sip)])
    }
    
    if ip < r.min || ip > r.max {
        return pkgerrors.New("ip is not in range")
    }

    index := (ip-1)/8
    base := byte(math.Exp2(float64(7-(ip-1)%8)))
    if r.masks[index] & base == 0 {
        return pkgerrors.New("ip is not allocated")
    }
    
    r.masks[index] &= (^base)
    return nil
}

func allocate(r *Range, pr bool) bool {
    ip, err := r.Allocate()
    if err != nil {
        log.Println(err)
        return false
    } else {
        log.Println("Allocate: " + ip)
        if pr {
            r.Print()
        }
        return true
    }
}

func main() {
    r := NewRange(10, 50)
    r.Print()
    
    for allocate(r, false) {
    }
    
    r.Print()    
    r.Free("192.168.0.17")
    r.Print()
    allocate(r, false)
    r.Print()
    
    fa := []int{50, 23, 40, 10}
    for _, i := range fa {
        r.Free("192.168.0." + strconv.Itoa(i))
    }
    r.Print()
    for allocate(r, true) {
    }
}
