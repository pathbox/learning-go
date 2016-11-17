package main

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Primitive) DecodeMsg(dc *msgp.Reader) (err error) {
	var xvk uint32
	xvk, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if xvk != 3 {
		err = msgp.ArrayError{Wanted: 3, Got: xvk}
		return
	}
	z.One, err = dc.ReadInt()
	if err != nil {
		return
	}
	z.Two, err = dc.ReadUint()
	if err != nil {
		return
	}
	z.Three, err = dc.ReadFloat64()
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Primitive) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 3
	err = en.Append(0x93)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.One)
	if err != nil {
		return
	}
	err = en.WriteUint(z.Two)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.Three)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Primitive) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 3
	o = append(o, 0x93)
	o = msgp.AppendInt(o, z.One)
	o = msgp.AppendUint(o, z.Two)
	o = msgp.AppendFloat64(o, z.Three)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Primitive) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var bzg uint32
	bzg, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if bzg != 3 {
		err = msgp.ArrayError{Wanted: 3, Got: bzg}
		return
	}
	z.One, bts, err = msgp.ReadIntBytes(bts)
	if err != nil {
		return
	}
	z.Two, bts, err = msgp.ReadUintBytes(bts)
	if err != nil {
		return
	}
	z.Three, bts, err = msgp.ReadFloat64Bytes(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

func (z Primitive) Msgsize() (s int) {
	s = 1 + msgp.IntSize + msgp.UintSize + msgp.Float64Size
	return
}
