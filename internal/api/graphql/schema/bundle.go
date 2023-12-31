package schema

// Auto generated GraphQL schema bundle
const schema = `
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

    # fullTransactions is the list of transactions assigned to the block.
    # this method will fetch all transaction details from the rpc.
    fullTransactions: [Transaction!]!

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

    # Block is the block this transaction was assigned to.
    block: Block

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

    # Type is the type of the transaction.
    type: String!
}
type Tick {
    # The timestamp of the tick
    timestamp: Int!

    # The value of the tick
    value: Long!
}

type TtfTick {
    # The timestamp of the tick
    timestamp: Int!

    # The value of the tick
    value: Float!
}
# AggSubject is the subject of the aggregation
enum AggSubject {
    TXS_COUNT,
    GAS_USED
}
# MazePathDirection is an enum that represents the four directions that a MazePath can go in.
enum MazePathDirection {
    NORTH,
    EAST,
    SOUTH,
    WEST
}

# MazePath is a type that represents a path in a maze.
type MazePath {
    # The direction that the path goes in.
    direction: MazePathDirection!,

    # The length of the path, which represents the number of steps before encountering an obstacle or a turn.
    length: Int!

    # The id of the next tile in the path.
    id: Int!

    # X is the x coordinate of the path.
    x: Int!,

    # Y is the y coordinate of the path.
    y: Int!,
}

# MazePosition is a type that represents a position in a maze.
type MazePosition {
    # Id is the id of the tile.
    id: Int!

    # X is the x coordinate of the tile.
    x: Int!,

    # Y is the y coordinate of the tile.
    y: Int!,

    # Paths is a list of paths that can be taken from this tile.
    paths: [MazePath!]!
}

# Maze is a type that represents a maze.
type Maze {
    # Width is the width of the maze.
    width: Int!,

    # Height is the height of the maze.
    height: Int!,

    # VisibilityRange is the range that the player can see.
    visibilityRange: Int!,

    # Address is the address of the maze.
    address: Address!

    # Name is the name of the maze.
    name: String!

    # StartX is the x coordinate of the starting tile.
    startX: Int!

    # StartY is the y coordinate of the starting tile.
    startY: Int!

    # EndX is the x coordinate of the ending tile.
    endX: Int!

    # EndY is the y coordinate of the ending tile.
    endY: Int!
}
type CurrentState {
    # Get current block height
    currentBlockHeight:Long

    # Get total number of accounts
    numberOfAccounts:Int!

    # Get total number of transactions
    numberOfTransactions:Long!

    # Get total number of validators
    numberOfValidators:Int!

    # Get disk size per 100M transactions in bytes
    diskSizePer100MTxs:Long!

    # Get disk size pruned per 100M transactions in bytes
    diskSizePrunedPer100MTxs:Long!

    # Get time to finality in seconds (rounded to 2 decimal places)
    timeToFinality:Float!

    # Get time to block in seconds (rounded to 2 decimal places)
    timeToBlock:Float!

    # Get idle state of the blockchain.
    isIdle: Boolean!
}
# Account defines block-chain account information container
type Account {
    # Address is the address of the account.
    address: Address!

    # Balance is the current balance of the Account in WEI.
    balance: BigInt!

    # transactions is the list of transactions that are linked to this account.
    transactions: [Transaction!]!
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
# Root schema definition
schema {
    query: Query
    mutation: Mutation
}

# Entry points for querying the API
type Query {
    # State represents the current state of the blockchain and network.
    state: CurrentState!

    # Get transaction information for given transaction hash.
    transaction(hash:Bytes32!):Transaction

    # Get block information by number.
    block(number:Long!):Block

    # Get recent observed blocks
    recentBlocks(limit:Int!):[Block!]!

    # Get current block height
    currentBlockHeight:Long

    # Get total number of accounts
    numberOfAccounts:Int!

    # Get total number of transactions
    numberOfTransactions:Long!

    # Get total number of validators
    numberOfValidators:Int!

    # Get disk size per 100M transactions in bytes
    diskSizePer100MTxs:Long!

    # Get disk size pruned per 100M transactions in bytes
    diskSizePrunedPer100MTxs:Long!

    # Get time to finality in seconds (rounded to 2 decimal places)
    timeToFinality:Float!

    # Get time to block in seconds (rounded to 2 decimal places)
    timeToBlock:Float!

    # Get block aggregated data by timestamp. It returns last 60 ticks aggregated
    # by 10 seconds.
    # parameters:
    #   subject: the subject of the aggregation - value of AggSubject enum
    blockTimestampAggregations(subject: AggSubject!):[Tick!]!

    # Get ttf aggregated data by timestamp. It returns last 60 ticks aggregated
    # by 10 seconds.
    ttfTimestampAggregations:[TtfTick!]!

    # Get an Account information by hash address.
    account(address:Address!):Account!

    # Get idle state of the blockchain.
    isIdle: Boolean!

    # Get list of maze games.
    mazeList: [Maze!]!

    # Get maze metadata.
    maze(address:Address!): Maze!
}

type Mutation {
    # Send request to obtain tokens from faucet. Returns phrase that should be signed by the user.
    requestTokens(symbol: String): String!

    # Send signed phrase to faucet to obtain tokens.
    claimTokens(address: Address!, challenge: String!, signature: String!, erc20Address: Address): Boolean!

    # Generate challenge that should be signed by the user to obtain position.
    mazeGameSession: String!

    # Send signed challenge to obtain position.
    mazeMyPosition(address: Address!, challenge: String!, signature: String!, mazeAddress: Address!): MazePosition
}

`
