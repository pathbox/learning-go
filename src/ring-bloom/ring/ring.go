Skip to content
Search or jump to…
Pull requests
Issues
Marketplace
Explore
 
@pathbox 
tannerryan
/
ring
Public
Code
Issues
Pull requests
1
Actions
Security
Insights
ring/ring.go /
@jameshartig
jameshartig Added Merge, MarshalBinary, UnmarshalBinary
Latest commit 2ef476e on 29 May 2019
 History
 1 contributor
139 lines (125 sloc)  3.71 KB
   
// Copyright (c) 2019 Tanner Ryan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ring

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sync"
)

var (
	errElements      = errors.New("error: elements must be greater than 0")
	errFalsePositive = errors.New("error: falsePositive must be greater than 0 and less than 1")
)

// Ring contains the information for a ring data store.
type Ring struct {
	size  uint64        // number of bits (bit array is size/8+1)
	bits  []uint8       // main bit array
	hash  uint64        // number of hash rounds
	mutex *sync.RWMutex // mutex for locking Add, Test, and Reset operations
}

// Init initializes and returns a new ring, or an error. Given a number of
// elements, it accurately states if data is not added. Within a falsePositive
// rate, it will indicate if the data has been added.
func Init(elements int, falsePositive float64) (*Ring, error) {
	if elements <= 0 {
		return nil, errElements
	}
	if falsePositive <= 0 || falsePositive >= 1 {
		return nil, errFalsePositive
	}

	r := Ring{}
	// number of bits
	m := (-1 * float64(elements) * math.Log(falsePositive)) / math.Pow(math.Log(2), 2)
	// number of hash operations
	k := (m / float64(elements)) * math.Log(2)

	r.mutex = &sync.RWMutex{}
	r.size = uint64(math.Ceil(m))
	r.hash = uint64(math.Ceil(k))
	r.bits = make([]uint8, r.size/8+1)
	return &r, nil
}

// Add adds the data to the ring.
func (r *Ring) Add(data []byte) {
	// generate hashes
	hash := generateMultiHash(data)
	r.mutex.Lock()
	for i := uint64(0); i < r.hash; i++ {
		index := getRound(hash, i) % r.size
		r.bits[index/8] |= (1 << (index % 8))
	}
	r.mutex.Unlock()
}

// Reset clears the ring.
func (r *Ring) Reset() {
	r.mutex.Lock()
	r.bits = make([]uint8, r.size/8+1)
	r.mutex.Unlock()
}

// Test returns a bool if the data is in the ring. True indicates that the data
// may be in the ring, while false indicates that the data is not in the ring.
func (r *Ring) Test(data []byte) bool {
	// generate hashes
	hash := generateMultiHash(data)
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for i := uint64(0); i < uint64(r.hash); i++ {
		index := getRound(hash, i) % r.size
		// check if index%8-th bit is not active
		if (r.bits[index/8] & (1 << (index % 8))) == 0 {
			return false
		}
	}
	return true
}

// Merges the sent Ring into itself.
func (r *Ring) Merge(m *Ring) error {
	if r.size != m.size || r.hash != m.hash {
		return errors.New("rings must have the same m/k parameters")
	}

	r.mutex.Lock()
	m.mutex.RLock()
	for i := 0; i < len(m.bits); i++ {
		r.bits[i] |= m.bits[i]
	}
	r.mutex.Unlock()
	m.mutex.RUnlock()
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (r *Ring) MarshalBinary() ([]byte, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	out := make([]byte, len(r.bits)+17)
	// store a version for future compatibility
	out[0] = 1
	binary.BigEndian.PutUint64(out[1:9], r.size)
	binary.BigEndian.PutUint64(out[9:17], r.hash)
	copy(out[17:], r.bits)
	return out, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (r *Ring) UnmarshalBinary(data []byte) error {
	// 17 bytes for version + size + hash and 1 byte at least for bits
	if len(data) < 17+1 {
		return fmt.Errorf("incorrect length: %d", len(data))
	}
	if data[0] != 1 {
		return fmt.Errorf("unexpected version: %d", data[0])
	}
	if r.mutex == nil {
		r.mutex = new(sync.RWMutex)
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.size = binary.BigEndian.Uint64(data[1:9])
	r.hash = binary.BigEndian.Uint64(data[9:17])
	// sanity check against the bits being the wrong size
	if len(r.bits) != int(r.size/8+1) {
		r.bits = make([]uint8, r.size/8+1)
	}
	copy(r.bits, data[17:])
	return nil
}
© 2022 GitHub, Inc.
Terms
Privacy
Security
Status
Docs
Contact GitHub
Pricing
API
Training
Blog
About
Loading complete