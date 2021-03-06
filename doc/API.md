Siad API
========

The siad API is currently under construction. Under semantic versioning, the
minor version will be incremented whenever API-breaking changes are introduced.
Once siad hits v1.0.0, the major version will be incremented instead.

All API calls return JSON objects. If there is an error, the error is returned
in plaintext with an appropriate HTTP error code. The standard response is {
"Success": true }. In this document, the API responses are defined as Go
structs. The structs will be encoded to JSON before being sent; they are used
here to provide type information.

There may be functional API calls which are not documented. These are not
guaranteed to be supported beyond the current release, and should not be used
in production.

Notes:
- Requests must set their User-Agent string to contain the substring "Sia-Agent".
- By default, siad listens on "localhost:9980". This can be changed using the '-a'
  flag when running siad.
- The `types.Currency` object is an arbitrary-precision unsigned integer. In JSON,
  it is represented as a base-10 string. You must use a "bignum" library to handle
  these values, or you risk losing precision.

Example GET curl call:  `curl -A "Sia-Agent" /wallet/transactions?startheight=1&endheight=250`

Example POST curl call: `curl -A "Sia-Agent" --data "amount=123&destination=abcd" /wallet/siacoins

Consensus
---------

Queries:

* /consensus       [GET]
* /consensus/block [GET]

#### /consensus [GET]

Function: Returns information about the consensus set, such as the current
block height.

Parameters: none

Response:
```
struct {
	height       types.BlockHeight (uint64)
	currentblock types.BlockID     (string)
	target       types.Target      (byte array)
}
```
'height' is the number of blocks in the blockchain.

'currentblock' is the hash of the current block.

'target' is the hash that needs to be met by a block for the block to be valid.
The target is inversely proportional to the difficulty.

#### /consensus/block [GET]

Function: Returns the block found at a given height.

Parameters:
```
height types.BlockHeight (uint64)
```
'height' is the height of the block that is being requested. The genesis block
is at height 0, it's child is at height 1, etc.

Response:
```
struct {
	block types.Block
}
```

Explorer
--------

Queries:

* /explorer         [GET]
* /explorer/block   [GET]
* /explorer/$(hash) [GET]

#### /explorer [GET]

Function: Returns the status of the blockchain and some
statistics. All Siacoin amounts are given in Hastings

Parameters: None

Response:
```
struct {
	height            types.BlockHeight (uint64)
	block             types.Block
	target            types.Target    (byte array)
	difficulty        types.Currency  (string)
	maturitytimestamp types.Timestamp (uint64)
	circulation       types.Currency  (string)

	transactioncount          uint64
	siacoininputcount         uint64
	siacoinoutputcount        uint64
	filecontractcount         uint64
	filecontractrevisioncount uint64
	storageproofcount         uint64
	siafundinputcount         uint64
	siafundoutputcount        uint64
	minerfeecount             uint64
	arbitrarydatacount        uint64
	transactionsignaturecount uint64

	activecontractcount uint64
	activecontractcost  types.Currency (string)
	activecontractsize  types.Currency (string)
	totalcontractcost   types.Currency (string)
	totalcontractsize   types.Currency (string)
}
```

#### /explorer/block [GET]

Function: Return a block at a given height.

Parameters:
```
height types.BlockHeight (uint64)
```

Response:
```
struct {
	block api.ExplorerBlock
}
```

#### /explorer/$(hash) [GET]

Function: Return information about an unknown hash.

Parameters: $(hash) is a url parameter specifying the hash that is being looked
up. The hash can be an unlock hash, a wallet address, a block id, a transaction
id, siacoin output id, file contract id, siafund output id, or any of the
derivatives of siacoin output ids (such as miner payout ids and file contract
payout ids).

Response:
```
struct {
	 hashtype     string
	 block        api.ExplorerBlock
	 blocks       []api.ExplorerBlock
	 transaction  api.ExplorerTransaction
	 transactions []api.ExplorerTransaction
}
```
`hashtype` indicates what type of hash it is. The options are 'blockid',
'transactionid', 'unlockhash', 'siacoinoutputid', 'filecontractid',
'siafundoutputid'. If the object is a block, only the 'block' field will be
filled out. If the object is a transaction, only the 'transaction' field will
be filled out. For all other types, the 'blocks' and 'transactions' fields will
be filled out, returning all of the blocks and transactions that feature the
provided hash.

Host
----

Queries:

* /host          [GET]
* /host          [POST]
* /host/announce [POST]

#### /host [GET]

Function: Fetches status information about the host.

Parameters: none

Response:
```
struct {
	collateral   types.Currency     (string)
	ipaddress    modules.NetAddress (string)
	maxduration  types.BlockHeight  (uint64)
	minduration  types.BlockHeight  (uint64)
	price        types.Currency     (string)
	totalstorage int64
	unlockhash   types.UnlockHash  (string)
	windowsize   types.BlockHeight (uint64)

	numCcntracts      uint64
	revenue           types.Currency (string)
	storageremaining  int64
	upcomingrevenue   types.Currency (string)
}
```
`collateral` is the number of hastings per byte per block that are put up as
collateral when making file contracts.

`ipaddress` is the network address of the host.

`maxduration` is the maximum allowed duration of a file contract.

`minduration` is the minimum allowed duration of a file contract.

`price` is the number of hastings per byte per block that the host is charging
when making file contracts.

`totalstorage` is the total amount of storage that has been allocated to the
host.

`unlockhash` is the address that hosting revenues will be sent to.

`windowsize` is the minimum required window that must be given to the host to
prove storage of a file. Due to potential spam attacks, bloat, DDOS, and host
downtime, 40 blocks is recommended as an absolute minimum. The current network
default is 288 blocks. The current software will break entirely below 20
blocks, though in theory something as low as 6 blocks could be safe.


`numcontracts` is the number of active contracts that the host is engaged in.

`revenue` is the total number of Hastings earned from hosting.

`storageremaining` is `TotalStorage` minus the number of bytes currently being
stored.

`upcomingrevenue` is the value of the contracts that have been created but not
fulfilled.

#### /host [POST]

Function: Configures hosting parameters. All parameters are optional;
unspecified parameters will be left unchanged.

Parameters:
```
collateral   int
maxduration  int
minduration  int
price        int
totalstorage int
windowsize   int
```
`collateral` is the number of hastings per byte per block that are put up as
collateral when making file contracts.

`maxduration` is the maximum allowed duration of a file contract.

`minduration` is the minimum allowed duration of a file contract.

`price` is the number of hastings per byte per block that the host is charging
when making file contracts.

`totalstorage` is the total amount of storage that has been allocated to the
host.

`windowsize` is the minimum required window that must be given to the host to
prove storage of a file. Due to potential spam attacks, bloat, DDOS, and host
downtime, 40 blocks is recommended as an absolute minimum. The current network
default is 288 blocks. The current software will break entirely below 20
blocks, though in theory something as low as 6 blocks could be safe.
`totalstorage` is how much storage (in bytes) the host will rent to the
network.

Response: standard

#### /host/announce [POST]

Function: The host will announce itself to the network as a source of storage.
Generally only needs to be called once.

Parameters:
```
netaddress string
```
`netaddress` is an optional parameter that specifies the address to be
announced. Supplying this parameters will also override standard connectivity
checks.

Response: standard

Miner
-----

Queries:

* /miner        [GET]
* /miner/start  [POST]
* /miner/stop   [POST]
* /miner/header [GET]
* /miner/header [POST]

#### /miner [GET]

Function: Return the status of the miner.

Parameters: none

Response:
```
struct {
	blocksmined      int
	cpuhashrate      int
	cpumining        bool
	staleblocksmined int
}
```
`cpumining` indicates whether the cpu miner is active or not.

`cpuhashrate` indicates how fast the cpu is hashing, in hashes per second.

`blocksmined` indicates how many blocks have been mined, this value is remembered after restarting.

`staleblocksmined` indicates how many stale blocks have been mined, this value is remembered after restarting.

#### /miner/start [POST]

Function: Starts a single threaded cpu miner. Does nothing if the cpu miner is
already running.

Parameters: none

Response: standard

#### /miner/stop [POST]

Function: Stops the cpu miner. Does nothing if the cpu miner is not running.

Parameters: none

Response: standard

#### /miner/header [GET]

Function: Provide a block header that is ready to be grinded on for work.

Parameters: none

Response:
```
[]byte
```
The response is a byte array containing a target followed by a block header
followed by a block. The target is the first 32 bytes. The block header is the
following 80 bytes, and the nonce is bytes 32-39 (inclusive) of the header
(bytes 64-71 of the whole array).

Layout:

0-31: target

32-111: header

#### /miner/header [POST]

Function: Submit a header that has passed the POW.

Parameters:
```
[]byte
```
The input byte array should be 80 bytes that form the solved block header.

Wallet
------

Queries:

* /wallet                      [GET]
* /wallet/address              [GET]
* /wallet/addresses            [GET]
* /wallet/backup               [GET]
* /wallet/init                 [POST]
* /wallet/load/033x            [POST]
* /wallet/load/seed            [POST]
* /wallet/load/siag            [POST]
* /wallet/lock                 [POST]
* /wallet/seeds                [GET]
* /wallet/siacoins             [POST]
* /wallet/siafunds             [POST]
* /wallet/transaction/$(id)    [GET]
* /wallet/transactions         [GET]
* /wallet/transactions/$(addr) [GET]
* /wallet/unlock               [POST]

The first time that the wallet is ever created, the wallet will be unencrypted
and locked. The wallet must be initialized and encrypted using a call to 
/wallet/init. After encrypting the wallet, the wallet must be unlocked. From 
that point forward (including restarting siad), the wallet will be encrypted,
and only the call to /wallet/unlock will be needed.

#### /wallet [GET]

Function: Returns basic information about the wallet, such as whether the
wallet is locked or unlocked.

Parameters: none

Response:
```
struct {
	encrypted bool
	unlocked  bool

	confirmedsiacoinbalance     types.Currency (string)
	unconfirmedoutgoingsiacoins types.Currency (string)
	unconfirmedincomingsiacoins types.Currency (string)

	siafundbalance      types.Currency (string)
	siacoinclaimBalance types.Currency (string)
}
```
'encrypted' indicates whether the wallet has been encrypted or not. If the
wallet has not been encrypted, then no data has been generated at all, and the
first time the wallet is unlocked, the password given will be used as the
password for encrypting all of the data. 'encrypted' will only be set to false
if the wallet has never been unlocked before (the unlocked wallet is still
encryped - but the encryption key is in memory).

'unlocked' indicates whether the wallet is currently locked or unlocked. Some
calls become unavailable when the wallet is locked.

'confirmedsiacoinbalance' is the number of siacoins available to the wallet as
of the most recent block in the blockchain.

'unconfirmedoutgoingsiacoins' is the number of siacoins that are leaving the
wallet according to the set of unconfirmed transactions. Often this number
appears inflated, because outputs are frequently larger than the number of
coins being sent, and there is a refund. These coins are counted as outgoing,
and the refund is counted as incoming. The difference in balance can be
calculated using 'unconfirmedincomingsiacoins' - 'unconfirmedoutgoingsiacoins'

'unconfirmedincomingsiacoins' is the number of siacoins are entering the wallet
according to the set of unconfirmed transactions. This number is often inflated
by outgoing siacoins, because outputs are frequently larger than the amount
being sent. The refund will be included in the unconfirmed incoming siacoins
balance.

'siafundbalance' is the number of siafunds available to the wallet as
of the most recent block in the blockchain.

'siacoinclaimbalance' is the number of siacoins that can be claimed from the
siafunds as of the most recent block. Because the claim balance increases every
time a file contract is created, it is possible that the balance will increase
before any claim transaction is confirmed.

#### /wallet/address [GET]

Function: Get a new address from the wallet generated by the primary seed. An
error will be returned if the wallet is locked.

Parameters: none

Response:
```
struct {
	address types.UnlockHash (string)
}
```
'address' is a Sia address that can receive siacoins or siafunds.

#### /wallet/addresses [GET]

Function: Fetch the list of addresses from the wallet.

Parameters: none

Response:
```
struct {
	addresses []WalletAddress
}
```
'addresses' is an array of wallet addresses. Wallet addresses take the
following form:
```
struct {
	address types.UnlockHash (string)
}
```
A struct is used to allow future fields such as a list of transactions or a set
of balances to be provided without breaking JSON compatibility.

#### /wallet/backup [GET]

Function: Create a backup of the wallet settings file. Though this can easily
be done manually, the settings file is often in an unknown or difficult to find
location. The /wallet/backup call can spare users the trouble of needing to
find their wallet file.

Parameters:
```
filepath string
```
'filepath' is the on-disk location that the file is going to be saved.

Response: standard

#### /wallet/init [POST]

Function: Initialize the wallet. After the wallet has been initialized once, it
does not need to be initialized again, and future calls to /wallet/init will
return an error. The encryption password is provided by the api call. If the
password is blank, then the password will be set to the same as the seed.

Parameters:
```
encryptionpassword string
dictionary string
```
'encryptionpassword' is the password that will be used to encrypt the wallet.
All subsequent calls should use this password. If left blank, the seed that
gets returned will also be the encryption password.

'dictionary' is the name of the dictionary that should be used when encoding
the seed. 'english' is the most common choice when picking a dictionary.

Response:
```
struct {
	primaryseed string
}
```
'primaryseed' is the dictionary encoded seed that is used to generate addresses
that the wallet is able to spend.

#### /wallet/load/033x [POST]

Function: Load a v0.3.3.x wallet into the current wallet, harvesting all of the
secret keys. All spendable addresses in the loaded wallet will become spendable
from the current wallet.

Parameters:
```
filepath           string
encryptionpassword string
```
'filepath' is the filepath of the v0.3.3.x wallet that is being loaded into the
current wallet.

'encryptionpassword' is the encryption key of the wallet. An error will be
returned if the wrong key is provided.

Response: standard.

#### /wallet/load/seed [POST]

Function: Give the wallet a seed to track when looking for incoming
transactions. The wallet will be able to spend outputs related to addresses
created by the seed. The seed is added as an auxiliary seed, and does not
replace the primary seed. Only the primary seed will be used for generating new
addresses.

Parameters:
```
encryptionpassword string
dictionary         string
seed               string
```
'encryptionpassword' is the key that is used to encrypt the new seed when it is
saved to disk.

'dictionary' is the name of the dictionary that should be used when encoding
the seed. 'english' is the most common choice when picking a dictionary.

'seed' is the dictionary-encoded phrase that corresponds to the seed being
added to the wallet.

Response: standard

#### /wallet/load/siag [POST]

Function: Load a key into the wallet that was generated by siag. Most siafunds
are currently in addresses created by siag.

Parameters:
```
encryptionpassword string
keyfiles           string
```
'encryptionpassword' is the key that is used to encrypt the siag key when it is
imported to the wallet.

'keyfiles' is a list of filepaths that point to the keyfiles that make up the
siag key. There should be at least one keyfile per required signature. The
filenames need to be commna separated (no spaces), which means filepaths that
contain a comma are not allowed.

#### /wallet/lock [POST]

Function: Locks the wallet, wiping all secret keys. After being locked, the
keys are encrypted. Queries for the seed, to send siafunds, and related queries
become unavailable. Queries concerning transaction history and balance are
still available.

Parameters: none

Response: standard.

#### /wallet/seeds [GET]

Function: Return a list of seeds in use by the wallet. The primary seed is the
only seed that gets used to generate new addresses. This call is unavailable
when the wallet is locked.

Parameters:
```
dictionary string
```
'dictionary' is the name of the dictionary that should be used when encoding
the seed. 'english' is the most common choice when picking a dictionary.

Response:
```
struct {
	primaryseed        mnemonics.Phrase   (string)
	addressesremaining int
	allseeds           []mnemonics.Phrase (array of strings)
}
```
'primaryseed' is the seed that is actively being used to generate new addresses
for the wallet.

'addressesremaining' is the number of addresses that remain in the primary seed
until exhaustion has been reached. Once exhaustion has been reached, new
addresses will continue to be generated but they will be more difficult to
recover in the event of a lost wallet file or encryption password.

'allseeds' is a list of all seeds that the wallet references when scanning the
blockchain for outputs. The wallet is able to spend any output generated by any
of the seeds, however only the primary seed is being used to generate new
addresses.

A seed is an encoded version of a 128 bit random seed. The output is 15 words
chosen from a small dictionary as indicated by the input. The most common
choice for the dictionary is going to be 'english'. The underlying seed is the
same no matter what dictionary is used for the encoding. The encoding also
contains a small checksum of the seed, to help catch simple mistakes when
copying. The library
[entropy-mnemonics](https://github.com/NebulousLabs/entropy-mnemonics) is used
when encoding.

#### /wallet/siacoins [POST]

Function: Send siacoins to an address. The outputs are arbitrarily selected
from addresses in the wallet.

Parameters:
```
amount      int
destination string
```
'amount' is the number of hastings being sent. A hasting is the smallest unit
in Sia. There are 10^24 hastings in a siacoin.

'destination' is the address that is receiving the coins.

Response:
```
struct {
	transactionids []types.TransactionID (string)
}
```
'transactionids' are the ids of the transactions that were created when sending
the coins. The last transaction contains the output headed to the
'destination'.

#### /wallet/siafunds [POST]

Function: Send siafunds to an address. The outputs are arbitrarily selected
from addresses in the wallet. Any siacoins available in the siafunds being sent
(as well as the siacoins available in any siafunds that end up in a refund
address) will become available to the wallet as siacoins after 144
confirmations. To access all of the siacoins in the siacoin claim balance, send
all of the siafunds to an address in your control (this will give you all the
siacoins, while still letting you control the siafunds).

Parameters:
```
amount      int
destination string
```
'amount' is the number of siafunds being sent.

'destination' is the address that is receiving the funds.

Response:
```
struct {
	transactionids []types.TransactionID (string)
}
```
'transactionids' are the ids of the transactions that were created when sending
the coins. The last transaction contains the output headed to the
'destination'.

#### /wallet/transaction/$(id) [GET]

Function: Get the transaction associated with a specific transaction id.

Parameters: $(id) is a url parameter specifying the id of the transaction that
should be returned.

Response:
```
struct {
	transaction modules.ProcessedTransaction
}
```

Processed transactions are transactions that have been processed by the wallet
and given more information, such as a confirmation height and a timestamp.
Processed transactions will always be returned in chronological order.

A processed transaction takes the following form:
```
struct modules.ProcessedTransaction {
	transaction           types.Transaction
	transactionid         types.TransactionID (string)
	confirmationheight    types.BlockHeight   (int)
	confirmationtimestamp types.Timestamp     (uint64)

	inputs  []modules.ProcessedInput
	outputs []modules.ProcessedOutput
}
```
'transaction' is a types.Transaction, and is defined in types.transaction.go

'transactionid' is the id of the transaction from which the wallet transaction
was derived.

'confirmationheight' is the height at which the transaction was confirmed. The
height will be set to 'uint64max' if the transaction has not been confirmed.

'confirmationtimestamp' is the time at which a transaction was confirmed. The
timestamp is an unsigned 64bit unix timestamp, and will be set to 'uint64max'
if the transaction is unconfirmed.

'inputs' is an array of processed inputs detailing the inputs to the
transaction. More information below.

'outputs' is an array of processed outputs detailing the outputs of
the transaction. Outputs related to file contracts are excluded.

A modules.ProcessedInput takes the following form:
```
struct modules.ProcessedInput {
	fundtype       types.Specifier  (string)
	walletaddress  bool
	relatedaddress types.UnlockHash (string)
	value          types.Currency   (string)
}
```

'fundtype' indicates what type of fund is represented by the input. Inputs can
be of type 'siacoin input', and 'siafund input'.

'walletaddress' indicates whether the address is owned by the wallet.
 
'relatedaddress' is the address that is affected. For inputs (outgoing money),
the related address is usually not important because the wallet arbitrarily
selects which addresses will fund a transaction. For outputs (incoming money),
the related address field can be used to determine who has sent money to the
wallet.

'value' indicates how much money has been moved in the input or output.

A modules.ProcessedOutput takes the following form:
```
struct modules.ProcessedOutput {
	fundtype       types.Specifier   (string)
	maturityheight types.BlockHeight (int)
	walletaddress  bool
	relatedaddress types.UnlockHash  (string)
	value          types.Currency    (string)
}
```

'fundtype' indicates what type of fund is represented by the output. Outputs
can be of type 'siacoin output', 'siafund output', and 'claim output'. Siacoin
outputs and claim outputs both relate to siacoins. Siafund outputs relate to
siafunds. Another output type, 'miner payout', points to siacoins that have been
spent on a miner payout. Because the destination of the miner payout is determined by
the block and not the transaction, the data 'maturityheight', 'walletaddress',
and 'relatedaddress' are left blank.

'maturityheight' indicates what height the output becomes available to be
spent. Siacoin outputs and siafund outputs mature immediately - their maturity
height will always be the confirmation height of the transaction. Claim outputs
cannot be spent until they have had 144 confirmations, thus the maturity height
of a claim output will always be 144 larger than the confirmation height of the
transaction.

'walletaddress' indicates whether the address is owned by the wallet.
 
'relatedaddress' is the address that is affected.

'value' indicates how much money has been moved in the input or output.

#### /wallet/transactions [GET]

Function: Return a list of transactions related to the wallet.

Parameters:
```
startheight types.BlockHeight (uint64)
endheight   types.BlockHeight (uint64)
```
'startheight' refers to the height of the block where transaction history
should begin.

'endheight' refers to the height of of the block where the transaction history
should end. If 'endheight' is greater than the current height, all transactions
up to and including the most recent block will be provided.

Response:
```
struct {
	confirmedtransactions   []modules.ProcessedTransaction
	unconfirmedtransactions []modules.ProcessedTransaction
}
```
'confirmedtransactions' lists all of the confirmed transactions appearing between
height 'startheight' and height 'endheight' (inclusive).

'unconfirmedtransactions' lists all of the unconfirmed transactions.

#### /wallet/transactions/$(addr) [GET]

Function: Return all of the transaction related to a specific address.

Parameters: $(addr) is a url parameter specifiy the unlock hash or wallet
address that is being looked up. All transactions related to that address will
be returned.

Response:
```
struct {
	transactions []modules.ProcessedTransaction.
}
```
'transactions' is a list of processed transactions that relate to the input
address.  See the documentation for '/wallet/transaction' for more information.

#### /wallet/unlock [POST]

Function: Unlock the wallet. The wallet is capable of knowing whether the
correct password was provided.

Parameters:
```
encryptionpassword string
```
'encryptionpassword' is the password that gets used to decrypt the file. Most
frequently, the encryption password is the same as the primary wallet seed.

Response: standard


*Routes below this line have not been made RESTful*

---

Daemon
------

Queries:

* /daemon/stop
* /daemon/version
* /daemon/updates/apply
* /daemon/updates/check

#### /daemon/stop

Function: Cleanly shuts down the daemon. May take a while.

Parameters: none

Response: standard

#### /daemon/version

Function: Returns the daemon's version

Parameters: none

Response:
```
struct {
	Version   string
}
```

#### /daemon/updates/apply:

Function: Applies the update specified by `version`.
                
Parameters:
```
version string
```

Response: standard

#### /daemon/updates/check:

Function: Checks for an update, returning a bool indicating whether
there is an update and a version indicating the version of the update.

Parameters: none

Response:
```
struct {
	Available bool
	Version   string
}
```

Gateway
-------

Queries:

* /gateway/status
* /gateway/peers/add
* /gateway/peers/remove

#### /gateway/status

Function: Returns information about the gateway, including the list of peers.

Parameters: none

Response:
```
struct {
	Address NetAddress
	Peers   []string
}
```

#### /gateway/peers/add

Function: Will add a peer to the gateway.

Parameters:
```
address string
```
`address` should be a reachable hostname + port number, typically of the form
"a.b.c.d:xxxx".

Response: standard

#### /gateway/peers/remove

Function: Will remove a peer from the gateway.

Parameters:
```
address string
```
`address` should be a reachable hostname + port number, typically of the form
"a.b.c.d:xxxx".

Response: standard

HostDB
------

Queries:

* /hostdb/hosts/active
* /hostdb/hosts/all

#### /hostdb/hosts/active

Function: Lists all of the active hosts in the hostdb.

Parameters: none

Response:
```
struct {
	Hosts []HostSettings
```

#### /hostdb/hosts/all

Function: Lists all of the hosts in the hostdb.

Parameters: none

Response:
```
struct {
	Hosts []HostSettings
}
```

Renter
------

Queries:

* /renter/downloadqueue
* /renter/files/delete
* /renter/files/download
* /renter/files/list
* /renter/files/load
* /renter/files/loadascii
* /renter/files/rename
* /renter/files/share
* /renter/files/shareascii
* /renter/files/upload

#### /renter/downloadqueue

Function: Lists all files in the download queue.

Parameters: none

Response:
```
[]struct{
	Filesize    uint64
	Received    uint64
	Destination string
	Nickname    string
}
```
Each file in the queue is represented by the above struct.

`Filesize` is the size of the file being download.

`Received` is the number of bytes downloaded thus far.

`Destination` is the path that the file was downloaded to.

`Nickname` is the nickname given to the file when it was uploaded.

#### /renter/files/delete

Function: Deletes a renter file entry. Does not delete any downloads or
original files, only the entry in the renter.

Parameters:
```
nickname string
```
`nickname` is the nickname of the file that has been uploaded to the network.

Response: standard

#### /renter/files/download

Function: Starts a file download.

Parameters:
```
nickname    string
destination string
```
`nickname` is the nickname of the file that has been uploaded to the network.

`destination` is the path that the file will be downloaded to.

Response: standard

#### /renter/files/list

Function: Lists the status of all files.

Parameters: none

Response:
```
[]struct {
	Available      bool
	UploadProgress float32
	Nickname       string
	Filesize       uint64
	TimeRemaining  types.BlockHeight (uint64)
}
```
Each uploaded file is represented by the above struct.

`Available` indicates whether or not the file can be downloaded immediately.

`UploadProgress` is the current upload percentage of the file, including
redundancy. In general, files will be available for download before
UploadProgress == 100.

`Nickname` is the nickname given to the file when it was uploaded.

`Filesize` is the size of the file in bytes.

`TimeRemaining` indicates how many blocks the file will be available for.

#### /renter/files/load

Function: Load a '.sia' into the renter.

Parameters:
```
filename string
```
`filename` is the filepath of the '.sia' that is being loaded.

Response:
```
struct {
	FilesAdded []string
}
```

#### /renter/files/loadascii

Function: Load a '.sia' into the renter.

Parameters:
```
file string
```
`file` is the ASCII representation of the '.sia' file being loaded into the
renter.

Response:
```
struct {
	FilesAdded []string
}
```

#### /renter/files/rename

Function: Rename a file. Does not rename any downloads or source files, only
renames the entry in the renter.

Parameters:
```
nickname string
newname  string
```
`nickname` is the current name of the file entry.

`newname` is the new name for the file entry.

Response: standard.

#### /renter/files/share

Function: Create a '.sia' that can be shared with other people.

Parameters:
```
nickname string
filepath string
```
`nickname` is the nickname of the file that will be shared.

`filepath` is the filepath of the '.sia' that will be created to share the
file. `filepath` must have the suffix '.sia'.

Response: standard.

#### /renter/files/shareascii

Function: Create a '.sia' that can be shared with other people.

Parameters:
```
nickname string
```
`nickname` is the nickname of the file that will be shared.

Response:
```
struct {
	File string
}
```
`file` is the ASCII representation of the '.sia' that would have been created.

#### /renter/files/upload

Function: Upload a file.

Parameters:
```
source   string
nickname string
```
`source` is the path to the file to be uploaded.

`nickname` is the name that will be used to reference the file.

Response: standard.

Transaction Pool
----------------

Queries:

* /transactionpool/transactions

#### /transactionpool/transactions

Function: Returns all of the transactions in the transaction pool.

Parameters: none

Response:
```
struct {
	Transactions []consensus.Transaction
}
```
Please see consensus/types/transactions.go for a more detailed explanation on
what a transaction looks like. There are many fields.
