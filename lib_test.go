package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const MGASPERSECOND = 35000000

type PureBenchFunction = func(input []byte) ([]byte, error)

func makeBench(runner PureBenchFunction, input []byte) func(b *testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := runner(input)
			if err != nil {
				log.Panicln("precompile returned error", err)
			}
		}
	}
}

func TestAndBenchSha256(t *testing.T) {
	for i := 0; i < 1024; i = i + 8 {
		log.Printf("Benchmarking SHA256 on %d bytes\n", i)
		input := make([]byte, i)
		rand.Read(input)
		funcToRun := vm.PrecompiledContractsIstanbul[common.BytesToAddress([]byte{0x02})].Run
		runnable := makeBench(funcToRun, input)
		result := testing.Benchmark(runnable)
		runningNs := result.NsPerOp()
		gas := MGASPERSECOND * runningNs / 1000000000
		t.Log("Gas = ", gas)
	}
}

func TestAndBenchRipemd(t *testing.T) {
	for i := 0; i < 1024; i = i + 8 {
		log.Printf("Benchmarking RIPEMD160 on %d bytes\n", i)
		input := make([]byte, i)
		rand.Read(input)
		funcToRun := vm.PrecompiledContractsIstanbul[common.BytesToAddress([]byte{0x03})].Run
		runnable := makeBench(funcToRun, input)
		result := testing.Benchmark(runnable)
		runningNs := result.NsPerOp()
		gas := MGASPERSECOND * runningNs / 1000000000
		t.Log("Gas = ", gas)
	}
}

func TestAndBenchBlake2f(t *testing.T) {
	lengths := []int{0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192}
	for _, l := range lengths {
		log.Printf("Benchmarking Blake2f on %d iterations\n", l)
		input := make([]byte, 213)
		rand.Read(input)
		input[212] &= 1
		b := bytes.NewBuffer(make([]byte, 0))
		binary.Write(b, binary.BigEndian, uint32(l))
		inp := b.Bytes()
		for k := 0; k < 4; k++ {
			input[k] = inp[k]
		}
		funcToRun := vm.PrecompiledContractsIstanbul[common.BytesToAddress([]byte{0x09})].Run
		runnable := makeBench(funcToRun, input)
		result := testing.Benchmark(runnable)
		runningNs := result.NsPerOp()
		gas := MGASPERSECOND * runningNs / 1000000000
		t.Log("Gas = ", gas)
	}
	for i := 0; i < 1024; i = i + 8 {

	}
}
