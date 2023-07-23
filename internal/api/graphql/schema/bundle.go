package schema

// Auto generated GraphQL schema bundle
const schema = `
# Root schema definition
schema {
    query: Query
}

# Entry points for querying the API
type Query {
    # Get transaction information for given transaction hash.
    transaction(hash:Bytes32!):Transaction

    # Get block information by number.
    block(number:Long!):Block

    # Get recent observed blocks
    recentBlocks(limit:Int!):[Block!]!

    # Get current block height
    currentBlockHeight:Long

    # Get block aggregated data by timestamp
    # parameters:
    #   subject: the subject of the aggregation - value of AggSubject enum
    #   resolution: the resolution of the aggregation - value of AggResolution enum
    #   ticks: the number of ticks to return
    #   endTime: the end timestamp of the aggregation, if not specified, last block's timestamp is used
    blockTimestampAggregations(subject: AggSubject!, resolution: AggResolution!, ticks:Int!, endTime:Int):[Tick!]!
}
type Tick {
    # The timestamp of the tick
    timestamp: Int!

    # The value of the tick
    value: Long!
}
# Bytes32 is a 32 byte binary string, represented by 0x prefixed hexadecimal hash.
scalar Bytes32

# Address is a 20 byte Opera address, represented as 0x prefixed hexadecimal number.
scalar Address

# BigInt is a large integer value. Input is accepted as either a JSON number,
# or a hexadecimal string alternatively prefixed with 0x. Output is 0x prefixed hexadecimal.
scalar BigInt

# Long is a 64 bit unsigned integer value.
scalar Long

# Bytes is an arbitrary length binary string, represented as 0x-prefixed hexadecimal.
# An empty byte string is represented as '0x'.
scalar Bytes

# Cursor is a string representing position in a sequential list of edges.
scalar Cursor

# Time represents date and time including time zone information in RFC3339 format.
scalar Time
# AggResolution is the resolution of the aggregation
enum AggResolution {
    SECONDS,
    MINUTE,
    HOUR,
    DAY
}

# AggSubject is the subject of the aggregation
enum AggSubject {
    TXS_COUNT,
    GAS_USED
}
type Block {
    # Number is the number of this block.
    number: Long!

    # Epoch is the number of this block's epoch.
    epoch: Long!

    # Hash is the unique block hash of this block.
    hash: Bytes32!

    # ParentHash is the unique block hash of the parent block of this block.
    parentHash: Bytes32!

    # Timestamp is the unix timestamp at which this block was mined.
    timestamp: Long!

    # GasLimit represents the maximum gas allowed in this block.
    gasLimit: Long!

    # GasUsed represents the actual total used gas by all transactions in this block.
    gasUsed: Long!

    # transactions is the list of unique hash values of transaction
    # assigned to the block.
    transactions: [Bytes32!]!

    # TransactionCount is the number of transactions in this block.
    transactionsCount: Int!
}
type Transaction {
    # Hash of the transaction
    hash: Bytes32!

    # BlockHash is the hash of the block this transaction was assigned to.
    # Null if the transaction is pending.
    blockHash: Bytes32

    # BlockHash is the hash of the block this transaction was assigned to.
    # Null if the transaction is pending.
    blockNumber: Long

    # From is the address of the account that sent this transaction
    from: Address!

    # To is the account the transaction was sent to.
    # This is null for contract creating transactions.
    to: Address

    # ContractAddress represents the address of smart contract
    # deployed by this transaction;
    # null if the transaction is not contract creation
    contractAddress: Address

    # Nonce is the number of transactions sent by the account prior to this transaction.
    nonce: Long!

    # Gas represents gas provided by the sender.
    gas: Long!

    # GasUsed is the amount of gas that was used on processing this transaction.
    # If the transaction is pending, this field will be null.
    gasUsed: Long

    # GasUsed is the amount of gas used when this transaction was executed in the block.
    # If the transaction is pending, this field will be null.
    cumulativeGasUsed: Long

    # GasPrice is the price of gas per unit in WEI.
    gasPrice: BigInt!

    # Value is the value sent along with this transaction in WEI.
    value: BigInt!

    # Input is the data supplied to the target of the transaction.
    # Contains smart contract byte code if this is contract creation.
    # Contains encoded contract state mutating function call if recipient
    # is a contract address.
    input: Bytes!

    # TransactionIndex is the index of this transaction in the block. This will
    # be null if the transaction is in a pending pool.
    transactionIndex: Long

    # Status is the return status of the transaction. This will be 1 if the
    # transaction succeeded, or 0 if it failed (due to a revert, or due to
    # running out of gas). If the transaction has not yet been processed, this
    # field will be null.
    status: Long
}
`
