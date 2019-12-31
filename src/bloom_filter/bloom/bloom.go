package bloom

import (
	"bytes"
	"encoding/json"
	"hash/fnv"
	"log"
	"math"
)

// bloom
// interface for bloom32 / bloom64 objects
type bloom interface {
	Add([]byte)
	Has([]byte) bool
	JSONMarshal() []byte
}

// getSize
// helperfunction to calc ceil to the next base 2 power exponent for the bloom boolset-length
// returns 2**exponent and exponent
func getSize(ui64 uint64) (size uint64, exponent uint64) {
	size = uint64(1)
	for size < ui64 {
		size <<= 1
		exponent++
	}
	return size, exponent
}

func calcSizeByWrongPositives(numEntries, wrongs float64) (uint64, uint64) {
	size := -1 * numEntries * math.Log(wrongs) / math.Pow(float64(0.69314718056), 2)
	locs := math.Ceil(float64(0.69314718056) * size / numEntries)
	return uint64(size), uint64(locs)
}

// New
// returns a new bloom32/bloom64 bloomfilter
func New(params ...float64) (bloomfilter bloom) {
	var entries, locs uint64
	if len(params) == 2 {
		if params[1] < 1 {
			entries, locs = calcSizeByWrongPositives(params[0], params[1])
		} else {
			entries, locs = uint64(params[0]), uint64(params[1])
		}
	} else {
		log.Fatal("usage: New(float64(number_of_entries), float64(number_of_hashlocations)) i.e. New(float64(1000), float64(3)) or New(float64(number_of_entries), float64(number_of_hashlocations)) i.e. New(float64(1000), float64(0.03))")
	}
	size, exponent := getSize(uint64(entries))
	switch exponent > 32 {
	case true:
		bloomfilter = bloom64{
			boolSet: make([]bool, size),
			sizeExp: exponent,
			size:    size - 1,
			setLocs: uint64(locs),
		}
	case false:
		bloomfilter = bloom32{
			boolSet: make([]bool, size),
			sizeExp: uint32(exponent),
			size:    uint32(size - 1),
			setLocs: uint32(locs),
			shift:   32 - uint32(exponent),
		}
	}
	return bloomfilter
}

// New32
// returns a new bloom32/bloom64 bloomfilter
func NewWithBoolset(bs *[]bool, locs int) (bloomfilter bloom) {
	size, exponent := getSize(uint64(len(*bs)))
	switch exponent > 32 {
	case true:
		bloomfilter = bloom64{
			boolSet: *bs,
			sizeExp: exponent,
			size:    size - 1,
			setLocs: uint64(locs),
		}
	case false:
		bloomfilter = bloom32{
			boolSet: *bs,
			sizeExp: uint32(exponent),
			size:    uint32(size - 1),
			setLocs: uint32(locs),
			shift:   32 - uint32(exponent),
		}
	}
	return bloomfilter
}

// bloomJSONImExport
// Im/Export structure used by JSONMarshal / JSONUnmarshal
type bloomJSONImExport struct {
	FilterSet []byte
	SetLocs   uint64
}

// JSONUnmarshal
// takes JSON-Object (type bloomJSONImExport) as []bytes
// returns bloom32 / bloom64 object
func JSONUnmarshal(dbData []byte) bloom {
	bloomImEx := bloomJSONImExport{}
	json.Unmarshal(dbData, &bloomImEx)
	buf := bytes.NewBuffer(bloomImEx.FilterSet)
	bs := ToBool(buf)
	bf := NewWithBoolset(&bs, int(bloomImEx.SetLocs))
	return bf
}

// ToBool
// takes *bytes.Buffer with bloom.boolset ([]bool) converted to []uint8 by ToBytes()
// returns []bool
func ToBool(buf *bytes.Buffer) []bool {
	var calc = []uint8{128, 64, 32, 16, 8, 4, 2, 1}
	hlpBool := make([]bool, (8 * (*buf).Len()))
	index := 0
	for {
		if r, err := buf.ReadByte(); err != nil {
			break
		} else {
			for _, c := range calc {
				if r/c > 0 {
					hlpBool[index] = true
					r -= c
				}
				index++
			}

		}
	}
	return hlpBool[:index]
}

// ToBytes
// takes pointer to bloom.boolSet (*[]bool) converts to []uint8
// returns []uint8
func ToBytes(bf *[]bool) []byte {
	var (
		outBytes bytes.Buffer
		calc     = []uint8{128, 64, 32, 16, 8, 4, 2, 1}
		lc       = len(calc)
		l        = len(*bf)
		v        uint8
		bools    = make([]bool, 8)
	)

	if l%lc != 0 {
		log.Fatalf("Len([]bool) must be multiple of %v but is %v !", lc, l)
	}

	for i := 0; i < l; i += lc {
		bools = (*bf)[i : i+lc]
		v = uint8(0)
		for i, b := range bools {
			if b {
				v += calc[i]
			}
		}
		outBytes.WriteByte(v)
	}
	return outBytes.Bytes()
}

// bloom32
// bloom filter with max. 2**32 elements
type bloom32 struct {
	boolSet []bool
	sizeExp uint32
	size    uint32
	setLocs uint32
	shift   uint32
}

// <--- http://www.cse.yorku.ca/~oz/hash.html
// Berkeley DB Hash (32bit)
// modified to fit with boolset-length
// hash is casted to l, h = 16bit fragments
// returns l,h
func (bf bloom32) absdbm(b *[]byte) (l, h uint32) {
	hash := uint32(len(*b))
	for _, c := range *b {
		hash = uint32(c) + (hash << 6) + (hash << bf.sizeExp) - hash
	}
	h = hash >> bf.shift
	l = hash << bf.sizeExp >> bf.sizeExp
	return l, h
}

// Add
// Add entry to Bloom filter
func (bf bloom32) Add(entry []byte) {
	l, h := bf.absdbm(&entry)
	for i := uint32(0); i < bf.setLocs; i++ {
		bf.boolSet[(h+i*l)&bf.size] = true
	}
}

// Has
// Check if entry had been added to Bloom filter
func (bf bloom32) Has(entry []byte) bool {
	l, h := bf.absdbm(&entry)
	for i := uint32(0); i < bf.setLocs; i++ {
		switch bf.boolSet[(h+i*l)&bf.size] {
		case false:
			return false
		}
	}
	return true
}

// JSONMarshal
// returns JSON-object (type bloomJSONImExport) as []byte
func (bf bloom32) JSONMarshal() []byte {
	bloomImEx := bloomJSONImExport{}
	bloomImEx.SetLocs = uint64(bf.setLocs)
	bloomImEx.FilterSet = ToBytes(&(bf.boolSet))
	data, err := json.Marshal(bloomImEx)
	if err != nil {
		log.Fatal("json.Marshal failed: ", err)
	}
	return data
}

// bloom64
// bloom filter with max. 2**64 elements
type bloom64 struct {
	boolSet []bool
	sizeExp uint64
	size    uint64
	setLocs uint64
}

// fnv64a
// use FNV64a Hash to set start (l) and value to be added (h)
// returns l,h
func (bf bloom64) fnv64a(b []byte) (l, h uint64) {
	h64 := fnv.New64a()
	h64.Write(b)
	hash64 := h64.Sum64()
	l = hash64 << 32 >> 32
	h = hash64 >> 32
	return l, h
}

// Add
// Add entry to Bloom filter
func (bf bloom64) Add(entry []byte) {
	l, h := bf.fnv64a(entry)
	for i := uint64(0); i < bf.setLocs; i++ {
		bf.boolSet[(h+i*l)&bf.size] = true
	}
}

// Has
// Check if entry had been added to Bloom filter
func (bf bloom64) Has(entry []byte) bool {
	l, h := bf.fnv64a(entry)
	for i := uint64(0); i < bf.setLocs; i++ {
		switch bf.boolSet[(h+i*l)&bf.size] {
		case false:
			return false
		}
	}
	return true
}

// JSONMarshal
// returns JSON-object (type bloomJSONImExport) as []byte
func (bf bloom64) JSONMarshal() []byte {
	bloomImEx := bloomJSONImExport{}
	bloomImEx.SetLocs = bf.setLocs
	bloomImEx.FilterSet = ToBytes(&(bf.boolSet))
	data, err := json.Marshal(bloomImEx)
	if err != nil {
		log.Fatal("json.Marshal failed: ", err)
	}
	return data
}

// Alternative Hash Function
// http://cr.yp.to/cdb/cdb.txt
// The cdb hash function is ``h = ((h << 5) + h) ^ c'', with a starting
// hash of 5381.

// func (bf bloom) djb2(f []byte) uint64 {
// 	hash := 5381
// 	for _, c := range f {
// 		hash = (hash<<5 +hash) ^ int(c)
// 	}
// 	return uint64(hash)
// }