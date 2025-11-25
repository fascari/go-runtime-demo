package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type (
	Block struct {
		Index        int       `json:"index"`
		Timestamp    time.Time `json:"timestamp"`
		Data         string    `json:"data"`
		PreviousHash string    `json:"previous_hash"`
		Hash         string    `json:"hash"`
		Nonce        int       `json:"nonce"`
	}

	Blockchain struct {
		chain      []Block
		difficulty int
		mu         sync.RWMutex
	}
)

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

func (bc *Blockchain) Chain() []Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	chainCopy := make([]Block, len(bc.chain))
	copy(chainCopy, bc.chain)
	return chainCopy
}

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

	bc.mineBlock(&newBlock)
	bc.chain = append(bc.chain, newBlock)

	return newBlock
}

// MineParallel demonstrates work-stealing and goroutine distribution across Ps
func (bc *Blockchain) MineParallel(data string, numGoroutines int) ([]Block, time.Duration) {
	start := time.Now()
	var wg sync.WaitGroup
	blocks := make([]Block, 0, numGoroutines)
	blocksChan := make(chan Block, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			blockData := data + "-worker-" + strconv.Itoa(id)
			block := bc.AddBlock(blockData)
			blocksChan <- block
		}(i)
	}

	go func() {
		wg.Wait()
		close(blocksChan)
	}()

	for block := range blocksChan {
		blocks = append(blocks, block)
	}

	duration := time.Since(start)
	return blocks, duration
}

func (bc *Blockchain) mineBlock(block *Block) {
	target := ""
	for i := 0; i < bc.difficulty; i++ {
		target += "0"
	}

	for {
		block.Hash = calculateHash(*block)

		if block.Hash[:bc.difficulty] == target {
			break
		}

		block.Nonce++

		// Yield to scheduler every 100k iterations to allow other goroutines to execute
		if block.Nonce%100000 == 0 {
			runtime.Gosched()
		}
	}
}

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

func (bc *Blockchain) Difficulty() int {
	return bc.difficulty
}

func (bc *Blockchain) Length() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return len(bc.chain)
}
