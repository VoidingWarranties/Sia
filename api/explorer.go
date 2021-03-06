package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/NebulousLabs/Sia/build"
	"github.com/NebulousLabs/Sia/types"
)

type (
	// ExplorerBlock is a block with some extra information such as the id and
	// height. This information is provided for programs that may not be
	// complex enough to compute the ID on their own.
	ExplorerBlock struct {
		ID             types.BlockID           `json:"id"`
		Height         types.BlockHeight       `json:"height"`
		MinerPayoutIDs []types.SiacoinOutputID `json:"minerpayoutids"`
		Transactions   []ExplorerTransaction   `json:"transactions"`
		RawBlock       types.Block             `json:"rawblock"`

		// Transaction type counts.
		MinerPayoutCount          uint64
		TransactionCount          uint64
		SiacoinInputCount         uint64
		SiacoinOutputCount        uint64
		FileContractCount         uint64
		FileContractRevisionCount uint64
		StorageProofCount         uint64
		SiafundInputCount         uint64
		SiafundOutputCount        uint64
		MinerFeeCount             uint64
		ArbitraryDataCount        uint64
		TransactionSignatureCount uint64

		// Factoids about file contracts.
		ActiveContractCost  types.Currency
		ActiveContractCount uint64
		ActiveContractSize  types.Currency
		TotalContractCost   types.Currency
		TotalContractSize   types.Currency
		TotalRevisionVolume types.Currency
	}

	// ExplorerTransaction is a transcation with some extra information such as
	// the parent block. This information is provided for programs that may not
	// be complex enough to compute the extra information on their own.
	ExplorerTransaction struct {
		ID                                       types.TransactionID       `json:"id"`
		Height                                   types.BlockHeight         `json:"height"`
		Parent                                   types.BlockID             `json:"parent"`
		SiacoinOutputIDs                         []types.SiacoinOutputID   `json:"siacoinoutputids"`
		FileContractIDs                          []types.FileContractID    `json:"filecontractids"`
		FileContractValidProofOutputIDs          [][]types.SiacoinOutputID `json:"filecontractvalidproofoutputids"`          // outer array is per-contract
		FileContractMissedProofOutputIDs         [][]types.SiacoinOutputID `json:"filecontractmissedproofoutputids"`         // outer array is per-contract
		FileContractRevisionValidProofOutputIDs  [][]types.SiacoinOutputID `json:"filecontractrevisionvalidproofoutputids"`  // outer array is per-revision
		FileContractRevisionMissedProofOutputIDs [][]types.SiacoinOutputID `json:"filecontractrevisionmissedproofoutputids"` // outer array is per-revision
		SiafundOutputIDs                         []types.SiafundOutputID   `json:"siafundoutputids"`
		SiaClaimOutputIDs                        []types.SiacoinOutputID   `json:"siafundclaimoutputids"`
		RawTransaction                           types.Transaction         `json:"rawtransaction"`
	}

	// ExplorerGET is the object returned as a response to a GET request to
	// /explorer.
	ExplorerGET struct {
		// General consensus information.
		Height            types.BlockHeight `json:"height"`
		CurrentBlock      types.BlockID     `json:"currentblock"`
		Target            types.Target      `json:"target"`
		Difficulty        types.Currency    `json:"difficulty"`
		MaturityTimestamp types.Timestamp   `json:"maturitytimestamp"`
		TotalCoins        types.Currency    `json:"totalcoins"`

		// Information about transaction type usage.
		MinerPayoutCount          uint64 `json:"minerpayoutcount"`
		TransactionCount          uint64 `json:"transactioncount"`
		SiacoinInputCount         uint64 `json:"siacoininputcount"`
		SiacoinOutputCount        uint64 `json:"siacoinoutputcount"`
		FileContractCount         uint64 `json:"filecontractcount"`
		FileContractRevisionCount uint64 `json:"filecontractrevisioncount"`
		StorageProofCount         uint64 `json:"storageproofcount"`
		SiafundInputCount         uint64 `json:"siafundinputcount"`
		SiafundOutputCount        uint64 `json:"siafundoutputcount"`
		MinerFeeCount             uint64 `json:"minerfeecount"`
		ArbitraryDataCount        uint64 `json:"arbitrarydatacount"`
		TransactionSignatureCount uint64 `json:"transactionsignaturecount"`

		// Information about file contracts and file contract revisions.
		ActiveContractCount uint64         `json:"activecontractcount"`
		ActiveContractCost  types.Currency `json:"activecontractcost"`
		ActiveContractSize  types.Currency `json:"activecontractsize"`
		TotalContractCost   types.Currency `json:"totalcontractcost"`
		TotalContractSize   types.Currency `json:"totalcontractsize"`
	}

	// ExplorerBlockGET is the object returned by a GET request to
	// /explorer/block.
	ExplorerBlockGET struct {
		Block ExplorerBlock `json:"block"`
	}

	// ExplorerHashGET is the object returned as a response to a GET request to
	// /explorer/hash. The HashType will indicate whether the hash corresponds
	// to a block id, a transaction id, a siacoin output id, a file contract
	// id, or a siafund output id. In the case of a block id, 'Block' will be
	// filled out and all the rest of the fields will be blank. In the case of
	// a transaction id, 'Transaction' will be filled out and all the rest of
	// the fields will be blank. For everything else, 'Transactions' and
	// 'Blocks' will/may be filled out and everything else will be blank.
	ExplorerHashGET struct {
		HashType     string                `json:"hashtype"`
		Block        ExplorerBlock         `json:"block"`
		Blocks       []ExplorerBlock       `json:"blocks"`
		Transaction  ExplorerTransaction   `json:"transaction"`
		Transactions []ExplorerTransaction `json:"transactions"`
	}
)

// buildExplorerTransaction takes a transaction and the height + id of the
// block it appears in an uses that to build an explorer transaction.
func buildExplorerTransaction(height types.BlockHeight, parent types.BlockID, txn types.Transaction) (et ExplorerTransaction) {
	et.ID = txn.ID()
	et.Height = height
	et.Parent = parent
	et.RawTransaction = txn

	for i := range txn.SiacoinOutputs {
		et.SiacoinOutputIDs = append(et.SiacoinOutputIDs, txn.SiacoinOutputID(uint64(i)))
	}
	for i, fc := range txn.FileContracts {
		fcid := txn.FileContractID(uint64(i))
		var fcvpoids []types.SiacoinOutputID
		var fcmpoids []types.SiacoinOutputID
		for j := range fc.ValidProofOutputs {
			fcvpoids = append(fcvpoids, fcid.StorageProofOutputID(types.ProofValid, uint64(j)))
		}
		for j := range fc.MissedProofOutputs {
			fcmpoids = append(fcmpoids, fcid.StorageProofOutputID(types.ProofMissed, uint64(j)))
		}
		et.FileContractIDs = append(et.FileContractIDs, fcid)
		et.FileContractValidProofOutputIDs = append(et.FileContractValidProofOutputIDs, fcvpoids)
		et.FileContractMissedProofOutputIDs = append(et.FileContractMissedProofOutputIDs, fcmpoids)
	}
	for _, fcr := range txn.FileContractRevisions {
		var fcrvpoids []types.SiacoinOutputID
		var fcrmpoids []types.SiacoinOutputID
		for j := range fcr.NewValidProofOutputs {
			fcrvpoids = append(fcrvpoids, fcr.ParentID.StorageProofOutputID(types.ProofValid, uint64(j)))
		}
		for j := range fcr.NewMissedProofOutputs {
			fcrmpoids = append(fcrmpoids, fcr.ParentID.StorageProofOutputID(types.ProofMissed, uint64(j)))
		}
		et.FileContractValidProofOutputIDs = append(et.FileContractValidProofOutputIDs, fcrvpoids)
		et.FileContractMissedProofOutputIDs = append(et.FileContractMissedProofOutputIDs, fcrmpoids)
	}
	for i := range txn.SiafundOutputs {
		et.SiafundOutputIDs = append(et.SiafundOutputIDs, txn.SiafundOutputID(uint64(i)))
	}
	for _, sfi := range txn.SiafundInputs {
		et.SiaClaimOutputIDs = append(et.SiaClaimOutputIDs, sfi.ParentID.SiaClaimOutputID())
	}
	return et
}

// buildExplorerBlock takes a block and its height and uses it to construct an
// explorer block.
func (srv *Server) buildExplorerBlock(height types.BlockHeight, block types.Block) ExplorerBlock {
	var mpoids []types.SiacoinOutputID
	var etxns []ExplorerTransaction
	for i := range block.MinerPayouts {
		mpoids = append(mpoids, block.MinerPayoutID(uint64(i)))
	}
	for _, txn := range block.Transactions {
		etxns = append(etxns, buildExplorerTransaction(height, block.ID(), txn))
	}
	facts, exists := srv.explorer.BlockFacts(height)
	if build.DEBUG && !exists {
		panic("incorrect request to buildExplorerBlock - block does not exist")
	}
	return ExplorerBlock{
		ID:             block.ID(),
		Height:         height,
		Transactions:   etxns,
		MinerPayoutIDs: mpoids,
		RawBlock:       block,

		// Transaction type counts.
		MinerPayoutCount:          facts.MinerPayoutCount,
		TransactionCount:          facts.TransactionCount,
		SiacoinInputCount:         facts.SiacoinInputCount,
		SiacoinOutputCount:        facts.SiacoinOutputCount,
		FileContractCount:         facts.FileContractCount,
		FileContractRevisionCount: facts.FileContractRevisionCount,
		StorageProofCount:         facts.StorageProofCount,
		SiafundInputCount:         facts.SiafundInputCount,
		SiafundOutputCount:        facts.SiafundOutputCount,
		MinerFeeCount:             facts.MinerFeeCount,
		ArbitraryDataCount:        facts.ArbitraryDataCount,
		TransactionSignatureCount: facts.TransactionSignatureCount,

		// Factoids about file contracts.
		ActiveContractCost:  facts.ActiveContractCost,
		ActiveContractCount: facts.ActiveContractCount,
		ActiveContractSize:  facts.ActiveContractSize,
		TotalContractCost:   facts.TotalContractCost,
		TotalContractSize:   facts.TotalContractSize,
		TotalRevisionVolume: facts.TotalRevisionVolume,
	}
}

// explorerBlockHandlerGET handles GET requests to /explorer/block.
func (srv *Server) explorerBlockHandlerGET(w http.ResponseWriter, req *http.Request) {
	// Parse the height that's being requested.
	var height types.BlockHeight
	_, err := fmt.Sscan(req.FormValue("height"), &height)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch and return the explorer block.
	block, exists := srv.cs.BlockAtHeight(height)
	if !exists {
		writeError(w, "no block found at input height in call to /explorer/block", http.StatusBadRequest)
		return
	}
	writeJSON(w, ExplorerBlockGET{
		Block: srv.buildExplorerBlock(height, block),
	})
}

// explorerHandler handles API calls to /explorer/block.
func (srv *Server) explorerBlockHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "" || req.Method == "GET" {
		srv.explorerBlockHandlerGET(w, req)
	} else {
		writeError(w, "unrecognized method when calling /explorer/block", http.StatusBadRequest)
	}
}

// buildTransactionSet returns the blocks and transactions that are associated
// with a set of transaction ids.
func (srv *Server) buildTransactionSet(txids []types.TransactionID) (txns []ExplorerTransaction, blocks []ExplorerBlock) {
	for _, txid := range txids {
		block, height, exists := srv.explorer.Transaction(txid)
		if !exists && build.DEBUG {
			panic("explorer pointing to nonexistant txn")
		}
		if types.TransactionID(block.ID()) == txid {
			blocks = append(blocks, srv.buildExplorerBlock(height, block))
		} else {
			var txn types.Transaction
			for _, t := range block.Transactions {
				if t.ID() == txid {
					txn = t
					break
				}
			}
			txns = append(txns, buildExplorerTransaction(height, block.ID(), txn))
		}
	}
	return txns, blocks
}

// explorerHandlerGET handles GET requests to /explorer.
func (srv *Server) explorerHandlerGET(w http.ResponseWriter, req *http.Request) {
	stats := srv.explorer.Statistics()
	writeJSON(w, ExplorerGET{
		Height:            stats.Height,
		CurrentBlock:      stats.CurrentBlock,
		Target:            stats.Target,
		Difficulty:        stats.Difficulty,
		MaturityTimestamp: stats.MaturityTimestamp,
		TotalCoins:        stats.TotalCoins,

		MinerPayoutCount:          stats.MinerPayoutCount,
		TransactionCount:          stats.TransactionCount,
		SiacoinInputCount:         stats.SiacoinInputCount,
		SiacoinOutputCount:        stats.SiacoinOutputCount,
		FileContractCount:         stats.FileContractCount,
		FileContractRevisionCount: stats.FileContractRevisionCount,
		StorageProofCount:         stats.StorageProofCount,
		SiafundInputCount:         stats.SiafundInputCount,
		SiafundOutputCount:        stats.SiafundOutputCount,
		MinerFeeCount:             stats.MinerFeeCount,
		ArbitraryDataCount:        stats.ArbitraryDataCount,
		TransactionSignatureCount: stats.TransactionSignatureCount,

		ActiveContractCount: stats.ActiveContractCount,
		ActiveContractCost:  stats.ActiveContractCost,
		ActiveContractSize:  stats.ActiveContractSize,
		TotalContractCost:   stats.TotalContractCost,
		TotalContractSize:   stats.TotalContractSize,
	})
}

// explorerHandlerGEThash handles GET requests to /explorer/$(hash).
func (srv *Server) explorerHandlerGEThash(w http.ResponseWriter, req *http.Request) {
	// The hash is scanned as an address, because an address can be typecast to
	// all other necessary types, and will correclty decode hashes whether or
	// not they have a checksum.
	encodedHash := strings.TrimPrefix(req.URL.Path, "/explorer/")
	hash, err := scanAddress(encodedHash)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Try the hash as a block id.
	block, height, exists := srv.explorer.Block(types.BlockID(hash))
	if exists {
		writeJSON(w, ExplorerHashGET{
			HashType: "blockid",
			Block:    srv.buildExplorerBlock(height, block),
		})
		return
	}

	// Try the hash as a transaction id.
	block, height, exists = srv.explorer.Transaction(types.TransactionID(hash))
	if exists {
		var txn types.Transaction
		for _, t := range block.Transactions {
			if t.ID() == types.TransactionID(hash) {
				txn = t
			}
		}
		writeJSON(w, ExplorerHashGET{
			HashType:    "transactionid",
			Transaction: buildExplorerTransaction(height, block.ID(), txn),
		})
		return
	}

	// Try the hash as a siacoin output id.
	txids := srv.explorer.SiacoinOutputID(types.SiacoinOutputID(hash))
	if len(txids) != 0 {
		txns, blocks := srv.buildTransactionSet(txids)
		writeJSON(w, ExplorerHashGET{
			HashType:     "siacoinoutputid",
			Blocks:       blocks,
			Transactions: txns,
		})
		return
	}

	// Try the hash as a file contract id.
	txids = srv.explorer.FileContractID(types.FileContractID(hash))
	if len(txids) != 0 {
		txns, blocks := srv.buildTransactionSet(txids)
		writeJSON(w, ExplorerHashGET{
			HashType:     "filecontractid",
			Blocks:       blocks,
			Transactions: txns,
		})
		return
	}

	// Try the hash as a siafund output id.
	txids = srv.explorer.SiafundOutputID(types.SiafundOutputID(hash))
	if len(txids) != 0 {
		txns, blocks := srv.buildTransactionSet(txids)
		writeJSON(w, ExplorerHashGET{
			HashType:     "siafundoutputid",
			Blocks:       blocks,
			Transactions: txns,
		})
		return
	}

	// Try the hash as an unlock hash. Unlock hash is checked last because
	// unlock hashes do not have collision-free guarantees. Someone can create
	// an unlock hash that collides with another object id. They will not be
	// able to use the unlock hash, but they can disrupt the explorer. This is
	// handled by checking the unlock hash last. Anyone intentionally creating
	// a colliding unlock hash (such a collision can only happen if done
	// intentionally) will be unable to find their unlock hash in the
	// blockchain through the explorer hash lookup.
	txids = srv.explorer.UnlockHash(types.UnlockHash(hash))
	if len(txids) != 0 {
		txns, blocks := srv.buildTransactionSet(txids)
		writeJSON(w, ExplorerHashGET{
			HashType:     "unlockhash",
			Blocks:       blocks,
			Transactions: txns,
		})
		return
	}

	// Hash not found, return an error.
	writeError(w, "unrecognized hash used as input to /explorer/hash", http.StatusBadRequest)
}

// explorerHandler handles API calls to /explorer and /explorer/
func (srv *Server) explorerHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/explorer" && (req.Method == "" || req.Method == "GET") {
		srv.explorerHandlerGET(w, req)
	} else if strings.HasPrefix(req.URL.Path, "/explorer/") && (req.Method == "" || req.Method == "GET") {
		srv.explorerHandlerGEThash(w, req)
	} else {
		writeError(w, "unrecognized call to /explorer", http.StatusBadRequest)
	}
}
