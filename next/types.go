// Copyright 2009 The Ninep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package next

import (
	"bytes"
	"io"
)

// 9P2000 message types
const (
	Tversion MType = 100 + iota
	Rversion
	Tauth
	Rauth
	Tattach
	Rattach
	Terror
	Rerror
	Tflush
	Rflush
	Twalk
	Rwalk
	Topen
	Ropen
	Tcreate
	Rcreate
	Tread
	Rread
	Twrite
	Rwrite
	Tclunk
	Rclunk
	Tremove
	Rremove
	Tstat
	Rstat
	Twstat
	Rwstat
	Tlast
)

const (
	MSIZE   = 2*1048576 + IOHDRSZ // default message size (1048576+IOHdrSz)
	IOHDRSZ = 24                  // the non-data size of the Twrite messages
	PORT    = 564                 // default port for 9P file servers
	NumFID  = 1 << 16
	QIDLen  = 13
)

// QID types
const (
	QTDIR     = 0x80 // directories
	QTAPPEND  = 0x40 // append only files
	QTEXCL    = 0x20 // exclusive use files
	QTMOUNT   = 0x10 // mounted channel
	QTAUTH    = 0x08 // authentication file
	QTTMP     = 0x04 // non-backed-up file
	QTSYMLINK = 0x02 // symbolic link (Unix, 9P2000.u)
	QTLINK    = 0x01 // hard link (Unix, 9P2000.u)
	QTFILE    = 0x00
)

// Flags for the mode field in Topen and Tcreate messages
const (
	OREAD   = 0x0    // open read-only
	OWRITE  = 0x1    // open write-only
	ORDWR   = 0x2    // open read-write
	OEXEC   = 0x3    // execute (== read but check execute permission)
	OTRUNC  = 0x10   // or'ed in (except for exec), truncate file first
	OCEXEC  = 0x20   // or'ed in, close on exec
	ORCLOSE = 0x40   // or'ed in, remove on close
	OAPPEND = 0x80   // or'ed in, append only
	OEXCL   = 0x1000 // or'ed in, exclusive client use
)

// File modes
const (
	DMDIR       = 0x80000000 // mode bit for directories
	DMAPPEND    = 0x40000000 // mode bit for append only files
	DMEXCL      = 0x20000000 // mode bit for exclusive use files
	DMMOUNT     = 0x10000000 // mode bit for mounted channel
	DMAUTH      = 0x08000000 // mode bit for authentication file
	DMTMP       = 0x04000000 // mode bit for non-backed-up file
	DMSYMLINK   = 0x02000000 // mode bit for symbolic link (Unix, 9P2000.u)
	DMLINK      = 0x01000000 // mode bit for hard link (Unix, 9P2000.u)
	DMDEVICE    = 0x00800000 // mode bit for device file (Unix, 9P2000.u)
	DMNAMEDPIPE = 0x00200000 // mode bit for named pipe (Unix, 9P2000.u)
	DMSOCKET    = 0x00100000 // mode bit for socket (Unix, 9P2000.u)
	DMSETUID    = 0x00080000 // mode bit for setuid (Unix, 9P2000.u)
	DMSETGID    = 0x00040000 // mode bit for setgid (Unix, 9P2000.u)
	DMREAD      = 0x4        // mode bit for read permission
	DMWRITE     = 0x2        // mode bit for write permission
	DMEXEC      = 0x1        // mode bit for execute permission
)

const (
	NOTAG Tag = 0xFFFF     // no tag specified
	NOFID FID = 0xFFFFFFFF // no fid specified
	// We reserve tag NOTAG and tag 0. 0 is a troublesome value to pass
	// around, since it is also a default value and using it can hide errors
	// in the code.
	NumTags = 1<<16 - 2
)

// Error values
const (
	EPERM   = 1
	ENOENT  = 2
	EIO     = 5
	EACCES  = 13
	EEXIST  = 17
	ENOTDIR = 20
	EINVAL  = 22
)

// Types contained in 9p messages.
type (
	MType      uint8
	Mode       uint8
	NumEntries uint16
	Tag        uint16
	FID        uint32
	Count      int32
	Perm       int32
	Offset     uint64
	Data       []byte
)

type ClientOpt func(*Client) error
type ServerOpt func(*Server) error
type Tracer func(string, ...interface{})

// Error represents a 9P2000 (and 9P2000.u) error
type Error struct {
	Err      string // textual representation of the error
	Errornum uint32 // numeric representation of the error (9P2000.u)
}

// File identifier
type QID struct {
	Type    uint8  // type of the file (high 8 bits of the mode)
	Version uint32 // version number for the path
	Path    uint64 // server's unique identification of the file
}

// Dir describes a file
type Dir struct {
	Size    uint16 // size-2 of the Dir on the wire
	Type    uint16
	Dev     uint32
	QID            // file's QID
	Mode    uint32 // permissions and flags
	Atime   uint32 // last access time in seconds
	Mtime   uint32 // last modified time in seconds
	Length  uint64 // file length in bytes
	Name    string // file name
	User    string // owner name
	Group   string // group name
	ModUser string // name of the last user that modified the file
	FID     uint64
}

// N.B. In all packets, the wire order is assumed to be the order in which you
// put struct members.

type TversionPkt struct {
	Msize   uint32
	Version string
}

type RversionPkt struct {
	Msize   uint32
	Version string
}

type TattachPkt struct {
	FID   uint64
	AFID  uint64
	Uname string
	Aname string
}

type RattachPkt struct {
	QID QID
}

type TwalkPkt struct {
	FID uint64
	NewFID uint64
	Paths []string
}

type RwalkPkt struct {
	QIDs []QID
}

type RerrorPkt struct {
	error string
}

type RPCCall struct {
	b     []byte
	Reply chan []byte
}

type RPCReply struct {
	b []byte
}

// Client implements a 9p client. It has a chan containing all tags,
// a scalar FID which is incremented to provide new FIDS (all FIDS for a given
// client are unique), an array of MaxTag-2 RPC structs, a ReadWriteCloser
// for IO, and two channels for a server goroutine: one down which RPCalls are
// pushed and another from which RPCReplys return.
// Once a client is marked Dead all further requests to it will fail.
// The ToNet/FromNet are separate so we can use io.Pipe for testing.
type Client struct {
	Tags       chan Tag
	FID        uint64
	RPC        []*RPCCall
	ToNet      io.WriteCloser
	FromNet    io.ReadCloser
	FromClient chan *RPCCall
	FromServer chan *RPCReply
	Msize      uint32
	Dead       bool
	Trace      Tracer
}

// Server is a 9p server.
// For now it's extremely serial. But we will use a chan for replies to ensure that
// we can go to a more concurrent one later.
type Server struct {
	NS      NineServer
	FromNet io.ReadCloser
	ToNet   io.WriteCloser
	Replies chan RPCReply
	Trace   Tracer
	Dead    bool
}

type NineServer interface {
	Dispatch(*bytes.Buffer, MType) error
	Rversion(uint32, string) (uint32, string, error)
	Rattach(uint64, uint64, string, string) (QID, error)
}