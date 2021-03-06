package renter

import (
	"bytes"
	"crypto/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/NebulousLabs/Sia/crypto"
	"github.com/NebulousLabs/Sia/modules"
	"github.com/NebulousLabs/Sia/modules/renter/hostdb"
	"github.com/NebulousLabs/Sia/types"
)

// a testHost simulates a host. It implements the hostdb.Uploader interface.
type testHost struct {
	ip   modules.NetAddress
	data []byte

	// used to simulate real-world conditions
	delay    time.Duration // transfers will take this long
	failRate int           // transfers will randomly fail with probability 1/failRate

	sync.Mutex
}

func (h *testHost) Address() modules.NetAddress  { return h.ip }
func (h *testHost) EndHeight() types.BlockHeight { return 0 }
func (h *testHost) Close() error                 { return nil }

func (h *testHost) ContractID() types.FileContractID {
	var fcid types.FileContractID
	copy(fcid[:], h.ip)
	return fcid
}

// Upload adds a piece to the testHost. It randomly fails according to the
// testHost's parameters.
func (h *testHost) Upload(data []byte) (offset uint64, err error) {
	// simulate I/O delay
	time.Sleep(h.delay)

	h.Lock()
	defer h.Unlock()

	// randomly fail
	if n, _ := crypto.RandIntn(h.failRate); n == 0 {
		return 0, crypto.ErrNilInput
	}

	h.data = append(h.data, data...)
	return uint64(len(h.data) - len(data)), nil
}

// TestErasureUpload tests parallel uploading of erasure-coded data.
func TestErasureUpload(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	// generate data
	const dataSize = 777
	data := make([]byte, dataSize)
	rand.Read(data)

	// create Reed-Solomon encoder
	rsc, err := NewRSCode(2, 10)
	if err != nil {
		t.Fatal(err)
	}

	// create hosts
	const pieceSize = 10
	hosts := make([]hostdb.Uploader, rsc.NumPieces())
	for i := range hosts {
		hosts[i] = &testHost{
			ip:       modules.NetAddress(strconv.Itoa(i)),
			delay:    time.Duration(i) * time.Millisecond,
			failRate: 5, // 20% failure rate
		}
	}
	// make one host really slow
	hosts[0].(*testHost).delay = 100 * time.Millisecond
	// make one host always fail
	hosts[1].(*testHost).failRate = 1

	// upload data to hosts
	f := newFile("foo", rsc, pieceSize, dataSize)
	r := bytes.NewReader(data)
	for chunk, pieces := range f.incompleteChunks() {
		err = f.repair(chunk, pieces, r, hosts)
		if err != nil {
			t.Fatal(err)
		}
	}

	// download data
	chunks := make([][][]byte, f.numChunks())
	for i := uint64(0); i < f.numChunks(); i++ {
		chunks[i] = make([][]byte, rsc.NumPieces())
	}
	for _, h := range hosts {
		contract, exists := f.contracts[h.ContractID()]
		if !exists {
			continue
		}
		for _, p := range contract.Pieces {
			encPiece := h.(*testHost).data[p.Offset : p.Offset+pieceSize+crypto.TwofishOverhead]
			piece, err := deriveKey(f.masterKey, p.Chunk, p.Piece).DecryptBytes(encPiece)
			if err != nil {
				t.Fatal(err)
			}
			chunks[p.Chunk][p.Piece] = piece
		}
	}
	buf := new(bytes.Buffer)
	for _, chunk := range chunks {
		err = rsc.Recover(chunk, f.chunkSize(), buf)
		if err != nil {
			t.Fatal(err)
		}
	}
	buf.Truncate(dataSize)

	if !bytes.Equal(buf.Bytes(), data) {
		t.Fatal("recovered data does not match original")
	}

	/*
		for i, h := range hosts {
			host := h.(*testHost)
			pieces := 0
			for _, p := range host.pieceMap {
				pieces += len(p)
			}
			t.Logf("Host #: %d\tDelay: %v\t# Pieces: %v\t# Chunks: %d", i, host.delay, pieces, len(host.pieceMap))
		}
	*/
}
