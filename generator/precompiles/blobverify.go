// Copyright 2020 Marius van der Wijden
// This file is part of the fuzzy-vm library.
//
// The fuzzy-vm library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The fuzzy-vm library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the fuzzy-vm library. If not, see <http://www.gnu.org/licenses/>.

package precompiles

import (
	"github.com/MariusVanDerWijden/FuzzyVM/filler"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/bls12381"
	"github.com/holiman/goevmlab/program"
)

var blobverifyAddr = common.HexToAddress("0x14")

type blobverifyCaller struct{}

func (*blobverifyCaller) call(p *program.Program, f *filler.Filler) error {

	var input []byte
	if f.Bool() {
		input = randomBlob(f)
	} else {
		input = correctBlob()
	}

	c := CallObj{
		Gas:       f.GasInt(),
		Address:   blobverifyAddr,
		InOffset:  0,
		InSize:    uint32(len(input)),
		OutOffset: 0,
		OutSize:   0,
		Value:     f.BigInt32(),
	}
	p.Mstore(input, 0)
	CallRandomizer(p, f, c)
	return nil
}

func randomBlob(f *filler.Filler) []byte {
	version := f.ByteSlice256()
	evaluationPoint := f.ByteSlice256()
	expectedOutput := f.ByteSlice256()
	dataKZG := f.ByteSlice(48)
	quotientKZG := f.ByteSlice(48)

	input := append(version, evaluationPoint...)
	input = append(input, expectedOutput...)
	input = append(input, dataKZG...)
	input = append(input, quotientKZG...)
	return input
}

func correctBlob() []byte {
	version := make([]byte, 32)
	evaluationPoint := bls12381.NewG1().Q()
	expectedOutput := bls12381.NewG1().Q()
	dataKZG := bls12381.NewG1().Q()
	quotientKZG := bls12381.NewG1().Q()

	input := append(version, evaluationPoint.Bytes()...)
	input = append(input, expectedOutput.Bytes()...)
	input = append(input, dataKZG.Bytes()...)
	input = append(input, quotientKZG.Bytes()...)
	return input
}
