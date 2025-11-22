package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Index        int       `json:"index"`
	Timestamp    time.Time `json:"timestamp"`
	Data         string    `json:"data"`
	PreviousHash string    `json:"previous_hash"`
	Hash         string    `json:"hash"`
	Nonce        int       `json:"nonce"`
}

// Blockchain manages the chain of blocks
type Blockchain struct {
	chain      []Block
	difficulty int
	mu         sync.RWMutex
}

// NewBlockchain creates a new blockchain with genesis block
func NewBlockchain(difficulty int) *Blockchain {
	bc := &Blockchain{
		chain:      make([]Block, 0),
		difficulty: difficulty,
	}

	genesis := Block{
		Index:        0,
		Timestamp:    time.Now(),
		Data:         "Genesis Block",
		PreviousHash: "0",
		Nonce:        0,
	}
	genesis.Hash = calculateHash(genesis)
	bc.chain = append(bc.chain, genesis)

	return bc
}

// Chain returns a copy of the blockchain
func (bc *Blockchain) Chain() []Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	chainCopy := make([]Block, len(bc.chain))
	copy(chainCopy, bc.chain)
	return chainCopy
}

// AddBlock adds a new block to the blockchain
// This demonstrates goroutine execution and scheduler behavior
func (bc *Blockchain) AddBlock(data string) Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	previousBlock := bc.chain[len(bc.chain)-1]

	newBlock := Block{
		Index:        previousBlock.Index + 1,
		Timestamp:    time.Now(),
		Data:         data,
		PreviousHash: previousBlock.Hash,
		Nonce:        0,
	}

	// CPU-intensive mining - demonstrates scheduler behavior
	bc.mineBlock(&newBlock)
	bc.chain = append(bc.chain, newBlock)

	return newBlock
}

// MineParallel mines multiple blocks in parallel
// Demonstrates work-stealing and goroutine distribution across Ps
func (bc *Blockchain) MineParallel(data string, numGoroutines int) ([]Block, time.Duration) {
	start := time.Now()
	var wg sync.WaitGroup
	blocks := make([]Block, 0, numGoroutines)
	blocksChan := make(chan Block, numGoroutines)

	// Create N goroutines - scheduler distributes across Ps
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Each goroutine enters _Grunnable state
			// Scheduler assigns to available P
			// Transitions to _Grunning when executing
			blockData := data + "-worker-" + strconv.Itoa(id)
			block := bc.AddBlock(blockData)
			blocksChan <- block
		}(i)
	}

	// Wait in separate goroutine to close channel
	go func() {
		wg.Wait()
		close(blocksChan)
	}()

	// Collect results
	for block := range blocksChan {
		blocks = append(blocks, block)
	}

	duration := time.Since(start)
	return blocks, duration
}

// mineBlock performs Proof of Work
// CPU-intensive operation perfect for studying scheduler behavior
func (bc *Blockchain) mineBlock(block *Block) {
	target := ""
	for i := 0; i < bc.difficulty; i++ {
		target += "0"
	}

	// This loop demonstrates:
	// 1. CPU-bound goroutine behavior
	// 2. Cooperative scheduling with Gosched
	// 3. Preemption points (Go 1.14+)
	for {
		block.Hash = calculateHash(*block)

		if block.Hash[:bc.difficulty] == target {
			break
		}

		block.Nonce++

		// Cooperative scheduling: yield to scheduler every 100k iterations
		// Without this, goroutine could monopolize P for long time
		// Transitions: _Grunning -> _Grunnable
		// Allows other goroutines on same P to execute
		if block.Nonce%100000 == 0 {
			runtime.Gosched()
		}
	}
}

// calculateHash computes SHA-256 hash of block
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) +
		block.Timestamp.String() +
		block.Data +
		block.PreviousHash +
		strconv.Itoa(block.Nonce)

	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

// Difficulty returns the current mining difficulty
func (bc *Blockchain) Difficulty() int {
	return bc.difficulty
}

// Length returns the number of blocks in the chain
func (bc *Blockchain) Length() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return len(bc.chain)
}
