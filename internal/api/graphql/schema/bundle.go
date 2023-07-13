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

    # InputData is the data supplied to the target of the transaction.
    # Contains smart contract byte code if this is contract creation.
    # Contains encoded contract state mutating function call if recipient
    # is a contract address.
    inputData: Bytes!

    # TrxIndex is the index of this transaction in the block. This will
    # be null if the transaction is in a pending pool.
    trxIndex: Long

    # Status is the return status of the transaction. This will be 1 if the
    # transaction succeeded, or 0 if it failed (due to a revert, or due to
    # running out of gas). If the transaction has not yet been processed, this
    # field will be null.
    status: Long
}
`
