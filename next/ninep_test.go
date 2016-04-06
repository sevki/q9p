// Copyright 2009 The Ninep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package next

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	//	"math"
	"reflect"
	"testing"
)

var debug = flag.Int("debug", 0, "print debug messages")

// Two files, dotu was true.
var testunpackbytes = []byte{
	79, 0, 0, 0, 0, 0, 0, 0, 0, 228, 193, 233, 248, 44, 145, 3, 0, 0, 0, 0, 0, 164, 1, 0, 0, 0, 0, 0, 0, 47, 117, 180, 83, 102, 3, 0, 0, 0, 0, 0, 0, 6, 0, 112, 97, 115, 115, 119, 100, 4, 0, 110, 111, 110, 101, 4, 0, 110, 111, 110, 101, 4, 0, 110, 111, 110, 101, 0, 0, 232, 3, 0, 0, 232, 3, 0, 0, 255, 255, 255, 255, 78, 0, 0, 0, 0, 0, 0, 0, 0, 123, 171, 233, 248, 42, 145, 3, 0, 0, 0, 0, 0, 164, 1, 0, 0, 0, 0, 0, 0, 41, 117, 180, 83, 195, 0, 0, 0, 0, 0, 0, 0, 5, 0, 104, 111, 115, 116, 115, 4, 0, 110, 111, 110, 101, 4, 0, 110, 111, 110, 101, 4, 0, 110, 111, 110, 101, 0, 0, 232, 3, 0, 0, 232, 3, 0, 0, 255, 255, 255, 255,
}

/*
func testUnpackDir(t *testing.T) {
	b := testunpackbytes
	for len(b) > 0 {
		var err error
		if _, b, _, err = UnpackDir(b, true); err != nil {
			t.Fatalf("Unpackdir: %v", err)
		}
	}
}
*/
func TestEncode(t *testing.T) {
	// The traces used in this array came from running 9p servers and clients.
	// Except for flush, which we made up.
	// TODO: put the replies in, then the decode testing when we get decode done.
	var tests = []struct {
		n string
		b []byte
		f func(b *bytes.Buffer)
	}{
		{
			"TVersion test with 8192 byte msize and 9P2000",
			[]byte{19, 0, 0, 0, 100, 0x55, 0xaa, 0, 32, 0, 0, 6, 0, 57, 80, 50, 48, 48, 48},
			func(b *bytes.Buffer) { MarshalTversionPkt(b, Tag(0xaa55), 8192, "9P2000") },
		},
		{
			"RVersion test with 8192 byte msize and 9P2000",
			[]byte{19, 0, 0, 0, 101, 0xaa, 0x55, 0, 32, 0, 0, 6, 0, 57, 80, 50, 48, 48, 48},
			func(b *bytes.Buffer) { MarshalRversionPkt(b, Tag(0x55aa), 8192, "9P2000") },
		},
		/*
			{
				"Flush test with tag 1 and oldtag 2",
				[]byte{9, 0, 0, 0, 108, 1, 0, 2, 0},
				[]interface{}{Tflush, Tag(1), Tag(2)},
			},
			{
				"Auth test with tag 0, fid 0,uname rminnich",
				[]byte{21, 0, 0, 0, 102, 0, 0, 0, 0, 0, 0, 8, 0, 114, 109, 105, 110, 110, 105, 99, 104},
				[]interface{}{Tauth, Tag(0), FID(0), "rminnich"},
			},
			{
				"Attach test with tag 0, fid 0, afid -1, uname rminnich",
				[]byte{28, 0, 0, 0, 104, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 8, 0, 114, 109, 105, 110, 110, 105, 99, 104, 1, 0, 47},
				[]interface{}{Tattach, Tag(0), FID(0), NOFID, "rminnich", "/"},
			},
			{
				"Tauth with an rerror of no user required",
				//Tauth tag 1 afid 45 uname 'rminnich' nuname 4294967295 aname ''
				[]byte{23,0,0,0,102,1,0,45,0,0,0,8,0,114,109,105,110,110,105,99,104,0,0},
				[]interface{}{Tauth, Tag(1), FID(45), "rminnich", ""},
				// [39 0 0 0 107 1 0 30 0 110 111 32 97 117 116 104 101 110 116 105 99 97 116 105 111 110 32 114 101 113 117 105 114 101 100 58 32 50 50]
				//Rerror tag 1 ename 'no authentication required: 22' ecode 0
			},
			{
				"Tattach from Harvey to ninep: Tattach tag 1 fid 48 afid 4294967295 uname 'rminnich' nuname 4294967295 aname ''",
				[]byte{27,0,0,0,104,1,0,48,0,0,0,255,255,255,255,8,0,114,109,105,110,110,105,99,104,0,0},
				[]interface{}{Tattach, Tag(1), FID(48), NOFID, "rminnich", ""},
				// 20 0 0 0 105 1 0 128 99 207 44 145 115 221 96 0 0 0 0 0]
				// Rattach tag 1 aqid (60dd73 912ccf63 'd')
			},
			{
				"Twalk tag 0 fid 0 newfid 1 to null",
				[]byte{23, 0, 0, 0, 110, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 4, 0, 110, 117, 108, 108},
				[]interface{}{Twalk, Tag(0), FID(0), FID(1), NumEntries(1), "null"},
			},
			{
				"Topen tag 0 fid 1 mode 2",
				[]byte{12, 0, 0, 0, 112, 0, 0, 1, 0, 0, 0, 2},
				[]interface{}{Topen, Tag(0), FID(1), Mode(2)},
			},
			{
				"Tread tag 0 fid 1 offset 0 count 8192",
				[]byte{23, 0, 0, 0, 116, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0},
				[]interface{}{Tread, Tag(0), FID(1), Offset(0), Count(8192)},
			},
			{
				"Tstat tag 1 fid 49",
				[]byte{11, 0, 0, 0, 124, 1, 0, 49, 0, 0, 0},
				// Rstat
				//
				//[84,0,0,0,125,1,0,75,0,73,0,0,0,0,0,0,0,128,99,207,44,145,115,221,96,0,0,0,0,0,253,1,0,128,109,185,47,86,196,66,41,86,0,16,0,0,0,0,0,0,6,0,104,97,114,118,101,121,8,0,114,109,105,110,110,105,99,104,8,0,114,109,105,110,110,105,99,104,4,0,110,111,110,101]

				//Rstat tag 1 st ('harvey' 'rminnich' 'rminnich' 'none' q (60dd73 912ccf63 'd') m d775 at 1445968237 mt 1445544644 l 4096 t 0 d 0 ext )
				[]interface{}{Tstat, Tag(1), FID(49)},
			},
			{
				"Twrite tag 3 fid 139 offset 0 count 3",
				[]byte{26, 0, 0, 0, 118, 3, 0, 139, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 104, 105, 10},
				// rwrite []byte{11,0,0,0,119,3,0,3,0,0,0}
				[]interface{}{Twrite, Tag(3), FID(139), Offset(0), Count(3), []byte("hi\n")},
			},
			{
				"Tclunk tag 1 fid 49",
				[]byte{11, 0, 0, 0, 120, 1, 0, 49, 0, 0, 0},
				// rclunk 7 0 0 0 121 1 0]
				[]interface{}{Tclunk, Tag(1), FID(49)},
			},
			{
				"Tremove tag 1 fid 49",
				[]byte{11, 0, 0, 0, 122, 1, 0, 49, 0, 0, 0},
				// rclunk 7 0 0 0 121 1 0]
				[]interface{}{Tremove, Tag(1), FID(49)},
			},
			{
				"Twstat tag 3 fid 49 ",
				//Twstat tag 3 fid 49 st ('' '' '' '' q (ffffffffffffffff ffffffff 'daAltL') m daAltDSPL777 at 4294967295 mt 1445968327 l 18446744073709551615 t 65535 d 4294967295 ext )
				[]byte{62, 0, 0, 0, 126, 3, 0, 49, 0, 0, 0, 49, 0, 47, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 199, 185, 47, 86, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0},
				// Rwstat [11 0 0 0 120 3 0 49 0 0 0]
				[]interface{}{Twstat, Tag(3), FID(49), &Dir{ /* TODO: remove this size
					Size:   47,
					Type:   math.MaxUint16,
					Dev:    math.MaxUint32,
					Qid:    Qid{Type: math.MaxUint8, Version: math.MaxUint32, Path: math.MaxUint64},
					Mode:   math.MaxUint32,
					Atime:  4294967295,
					Mtime:  1445968327,
					Length: 18446744073709551615,
					Name:   "",
					Uid:    "",
					Gid:    "",
					Muid:   "",
				},
				},
			},
			{
				"Tcreate tag 3 fid 74 name 'y' perm 666 mode 0",
				[]byte{19,0,0,0,114,3,0,74,0,0,0,1,0,121,182,1,0,0,0},
				[]interface{}{Tcreate, Tag(3), FID(74), "y", Perm(0666), Mode(0)},
				/// rcreate [24 0 0 0 115 3 0 0 226 200 71 172 45 166 98 0 0 0 0 0 0 0 0 0]
				// Rcreate tag 3 qid (62a62d ac47c8e2 '') iounit 0
			},*/
	}

	for _, v := range tests {
		var b bytes.Buffer
		v.f(&b)
		if !reflect.DeepEqual(v.b, b.Bytes()) {
			t.Errorf("Mismatch on %v: Got\n%v[%v], want\n%v[%v]", v.n, b.Bytes(), len(b.Bytes()), v.b, len(v.b))
		}
	}

}

/*
func testDecode(t *testing.T) {
	var tests = []struct {
		n string
		b []byte
		f func(b *bytes.Buffer) error
	}{
		{
			"TVersion test with 8192 byte msize and 9P2000",
			[]byte{19, 0, 0, 0, 100, 255, 255, 0, 32, 0, 0, 6, 0, 57, 80, 50, 48, 48, 48},
			func (b *bytes.Buffer){ MarshalTversionPkt(b, NOTAG, 8192, "9P2000")},
		},

	for _, v := range tests {
		var b bytes.Buffer
		v.f(&b)
		if !reflect.DeepEqual(v.b, b.Bytes()) {
			t.Errorf("Mismatch on %v: Got %v[%v], want %v[%v]", v.n, b.Bytes(), len(b.Bytes()), v.b, len(v.b))
		}
	}

}
*/

func TestTags(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fatalf("%v", err)
	}
	_ = c.GetTag()
	if len(c.Tags) != NumTags-1 {
		t.Errorf("Got one tag, len(tags) is %d, want %d", len(c.Tags), NumTags-1)
	}
}

type echo struct {
	*Server
	Versioned bool
}

// Dispatch dispatches request to different functions.
// We could do this with interface assertions and such a al rsc/fuse
// but most people I talked do disliked that. So we don't. If you want
// to make things optional, just define the ones you want to implement in this case.
func (e *echo) Dispatch(b *bytes.Buffer, t MType) error {
	switch t {
	case Tversion:
		return e.SrvRversion(b)
	case Tattach:
		return e.SrvRattach(b)
	}
	// This has been tested by removed Attach from the switch.
	ServerError(b, fmt.Sprintf("Dispatch: %v not supported", RPCNames[t]))
	return nil
}

func (e *echo) Rversion(msize uint32, version string) (uint32, string, error) {
	if version != "9P2000" {
		return 0, "", fmt.Errorf("%v not supported; only 9P2000", version)
	}
	e.Versioned = true
	return msize, version, nil
}

func (e *echo) Rattach(uint64, uint64, string, string) (QID, error) {
	if !e.Versioned {
		return QID{}, fmt.Errorf("Version must be one first")
	}
	return QID{}, nil
}

func TestTVersion(t *testing.T) {
	sr, cw := io.Pipe()
	cr, sw := io.Pipe()
	c, err := NewClient(func(c *Client) error {
		c.FromNet, c.ToNet = cr, cw
		return nil
	},
		func(c *Client) error {
			c.Msize = 8192
			c.Trace = t.Logf
			return nil
		})
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("Client is %v", c.String())

	e := &echo{}
	e.Server, err = NewServer(func(s *Server) error {
		s.FromNet, s.ToNet = sr, sw
		s.Trace = t.Logf
		s.NS = e
		return nil
	})

	if err != nil {
		t.Fatalf("NewServer: want nil, got %v", err)
	}

	e.Start()
	t.Logf("Server is %v", e.String())
	if _, err = c.CallTattach(0, 0, "", ""); err == nil {
		t.Fatalf("CallTattach: want err, got nil")
	}
	t.Logf("CallTattach: wanted an error and got %v", err)

	m, v, err := c.CallTversion(8000, "9p3000")
	if err == nil {
		t.Fatalf("CallTversion: want err, got nil")
	}
	t.Logf("CallTversion: wanted an error and got %v", err)

	m, v, err = c.CallTversion(8000, "9P2000")
	if err != nil {
		t.Fatalf("CallTversion: want nil, got %v", err)
	}
	t.Logf("CallTversion: msize %v version %v", m, v)

	t.Logf("Server is %v", e.String())
	a, err := c.CallTattach(0, 0, "", "")
	if err != nil {
		t.Fatalf("CallTattach: want nil, got %v", err)
	}
	t.Logf("Attach is %v", a)

}