package main

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	errDone    = errors.New("done generation")
	errBadLine = errors.New("bad line for trace format")
)

type generator func() (uint64, error)

func newZipf() generator {
	z := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), 1.0001, 10, 50000000)
	return func() (uint64, error) {
		return z.Uint64(), nil
	}
}

func newLIRS(pre string) generator {
	file, err := os.Open("./trace/" + pre + ".lirs.gz")
	if err != nil {
		panic(err)
	}
	trace, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	return newGenerator(parseLIRS, trace)
}

func newARC(pre string) generator {
	file, err := os.Open("./trace/" + pre + ".arc.gz")
	if err != nil {
		panic(err)
	}
	trace, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	return newGenerator(parseARC, trace)
}

type parser func(string, error) ([]uint64, error)

func newGenerator(parser parser, file io.Reader) generator {
	b := bufio.NewReader(file)
	s := make([]uint64, 0)
	i := -1
	var err error
	return func() (uint64, error) {
		if i++; i == len(s) {
			if s, err = parser(b.ReadString('\n')); err != nil {
				s = []uint64{0}
			}
			i = 0
		}
		return s[i], err
	}
}

func parseLIRS(line string, err error) ([]uint64, error) {
	if line = strings.TrimSpace(line); line != "" {
		key, err := strconv.ParseUint(line, 10, 64)
		return []uint64{key}, err
	}
	return nil, errDone
}

func parseARC(line string, err error) ([]uint64, error) {
	if line != "" {
		cols := strings.Fields(line)
		if len(cols) != 4 {
			return nil, errBadLine
		}
		start, err := strconv.ParseUint(cols[0], 10, 64)
		if err != nil {
			return nil, err
		}
		count, err := strconv.ParseUint(cols[1], 10, 64)
		if err != nil {
			return nil, err
		}
		seq := make([]uint64, count)
		for i := range seq {
			seq[i] = start + uint64(i)
		}
		return seq, nil
	}
	return nil, errDone
}
