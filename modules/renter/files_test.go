package renter

import (
	"testing"

	"github.com/NebulousLabs/Sia/types"
)

// TestFileAvailable probes the available method of the file type.
func TestFileAvailable(t *testing.T) {
	rsc, _ := NewRSCode(1, 10)
	f := &file{
		size:        1000,
		erasureCode: rsc,
		pieceSize:   100,
	}

	if f.available() {
		t.Error("file should not be available")
	}

	var fc fileContract
	for i := uint64(0); i < f.numChunks(); i++ {
		fc.Pieces = append(fc.Pieces, pieceData{Chunk: i, Piece: 0})
	}
	f.contracts = map[types.FileContractID]fileContract{types.FileContractID{}: fc}

	if !f.available() {
		t.Error("file should be available")
	}
}

// TestFileExpiration probes the expiration method of the file type.
func TestFileExpiration(t *testing.T) {
	f := &file{
		contracts: make(map[types.FileContractID]fileContract),
	}

	if f.expiration() != 0 {
		t.Error("file with no pieces should report as having no time remaining")
	}

	// Add a contract.
	fc := fileContract{}
	fc.WindowStart = 100
	f.contracts[types.FileContractID{0}] = fc
	if f.expiration() != 100 {
		t.Error("file did not report lowest WindowStart")
	}

	// Add a contract with a lower WindowStart.
	fc.WindowStart = 50
	f.contracts[types.FileContractID{1}] = fc
	if f.expiration() != 50 {
		t.Error("file did not report lowest WindowStart")
	}

	// Add a contract with a higher WindowStart.
	fc.WindowStart = 75
	f.contracts[types.FileContractID{2}] = fc
	if f.expiration() != 50 {
		t.Error("file did not report lowest WindowStart")
	}
}

// TestRenterDeleteFile probes the DeleteFile method of the renter type.
func TestRenterDeleteFile(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	rt, err := newRenterTester("TestRenterDeleteFile")
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	// Delete a file from an empty renter.
	err = rt.renter.DeleteFile("dne")
	if err != ErrUnknownNickname {
		t.Error("Expected ErrUnknownNickname:", err)
	}

	// Put a file in the renter.
	rt.renter.files["1"] = &file{
		name: "one",
	}
	// Delete a different file.
	err = rt.renter.DeleteFile("one")
	if err != ErrUnknownNickname {
		t.Error("Expected ErrUnknownNickname:", err)
	}
	// Delete the file.
	err = rt.renter.DeleteFile("1")
	if err != nil {
		t.Error(err)
	}
	if len(rt.renter.FileList()) != 0 {
		t.Error("file was deleted, but is still reported in FileList")
	}

	/*
		// Put a file in the renter, then rename it.
		rt.renter.files["1"] = &file{
			name: "one",
		}
		rt.renter.RenameFile("1", "one")
		// Call delete on the previous name.
		err = rt.renter.DeleteFile("1")
		if err != ErrUnknownNickname {
			t.Error("Expected ErrUnknownNickname:", err)
		}
		// Call delete on the new name.
		err = rt.renter.DeleteFile("one")
		if err != nil {
			t.Error(err)
		}
	*/
}

// TestRenterFileList probes the FileList method of the renter type.
func TestRenterFileList(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	rt, err := newRenterTester("TestRenterFileList")
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	// Get the file list of an empty renter.
	if len(rt.renter.FileList()) != 0 {
		t.Error("FileList has non-zero length for empty renter?")
	}

	// Put a file in the renter.
	rsc, _ := NewRSCode(1, 1)
	rt.renter.files["1"] = &file{
		name:        "one",
		erasureCode: rsc,
		pieceSize:   1,
	}
	if len(rt.renter.FileList()) != 1 {
		t.Error("FileList is not returning the only file in the renter")
	}
	if rt.renter.FileList()[0].Nickname != "one" {
		t.Error("FileList is not returning the correct filename for the only file")
	}

	// Put multiple files in the renter.
	rt.renter.files["2"] = &file{
		name:        "two",
		erasureCode: rsc,
		pieceSize:   1,
	}
	if len(rt.renter.FileList()) != 2 {
		t.Error("FileList is not returning both files in the renter")
	}
	files := rt.renter.FileList()
	if !((files[0].Nickname == "one" || files[0].Nickname == "two") &&
		(files[1].Nickname == "one" || files[1].Nickname == "two") &&
		(files[0].Nickname != files[1].Nickname)) {
		t.Error("FileList is returning wrong names for the files:", files[0].Nickname, files[1].Nickname)
	}
}

// TestRenterRenameFile probes the rename method of the renter.
func TestRenterRenameFile(t *testing.T) {
	t.Skip("Renaming disabled")

	rt, err := newRenterTester("TestRenterRenameFile")
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	// Rename a file that doesn't exist.
	err = rt.renter.RenameFile("1", "1a")
	if err != ErrUnknownNickname {
		t.Error("Expecting ErrUnknownNickname:", err)
	}

	// Rename a file that does exist.
	rt.renter.files["1"] = &file{
		name: "1",
	}
	files := rt.renter.FileList()
	err = rt.renter.RenameFile("1", "1a")
	if err != nil {
		t.Fatal(err)
	}
	if len(rt.renter.FileList()) != 1 {
		t.Fatal("FileList has unexpected number of files:", len(rt.renter.FileList()))
	}
	if files[0].Nickname != "1a" {
		t.Error("RenameFile failed, new file nickname is not what is expected.")
	}

	// Rename a file to an existing name.
	rt.renter.files["1"] = &file{
		name: "1",
	}
	err = rt.renter.RenameFile("1", "1a")
	if err != ErrNicknameOverload {
		t.Error("Expecting ErrNicknameOverload:", err)
	}
	if files[0].Nickname != "1a" {
		t.Error("Side effect occured during rename:", files[0].Nickname)
	}

	// Rename a file to the same name.
	err = rt.renter.RenameFile("1", "1")
	if err != ErrNicknameOverload {
		t.Error("Expecting ErrNicknameOverload:", err)
	}
	if files[0].Nickname != "1a" {
		t.Error("Side effect occured during rename:", files[0].Nickname)
	}
}
