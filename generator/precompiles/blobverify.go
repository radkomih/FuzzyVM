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

/*
	Precompile:
		[32]byte: VersionedHash
		[32]byte: EvaluationPoint (LE)
		[32]byte: ExpectedOutput (LE)
		[48]byte: DataKZG
		[48]byte: QuotientKZG

		Assertions:
		EP < BLS_MODULUS
		EO < BLS_MODULUS
		KZGToVersionedHash(DataKZG) == VersionedHash
		VerifyKZG(DataKZG, EP, EO, QuotientKZG)
*/

// VersionedHash = kzg.KZGToVersionedHash(DataKZG)
// BLS_MODULUS = gokzg4844.BlsModulus[:]

/*
blob := GetRandBlob(123)
	commitment, err := ctx.BlobToKZGCommitment(blob, NumGoRoutines)
	require.NoError(t, err)
	proof, err := ctx.ComputeBlobKZGProof(blob, commitment, NumGoRoutines)
	require.NoError(t, err)
	err = ctx.VerifyBlobKZGProof(blob, commitment, proof)
	require.NoError(t, err)
*/

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

func point(f *filler.Filler) []byte {
	g1 := bls12381.NewG1()
	switch f.Byte() % 3 {
	case 0:
		return g1.Q().Bytes()
	case 1:
		return g1.ToBytes(g1.Zero())
	case 2:
		return g1.ToBytes(g1.One())
	case 3:
		rnd := f.BigInt32()
		res := g1.One()
		g1.MulScalar(res, res, rnd)
		return g1.ToBytes(res)
	}
	return nil
}
