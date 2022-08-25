package main

type Helper interface {
	GetBlockByNumber(n int) Block
	GetLatestBlock() Block
	GetBlockByHash(hex_address string) Block
}

type Block struct {
	PreviousBlockHash string
	CurrentBlockHash  string
	Number            int
}

type Service struct {
	helper Helper

	blockByHash map[string]*Block
}

func (s Service) getLatestBlocks(previousBlocks []*Block) []*Block {
	s.storeBlocks(previousBlocks)
	s.downloadNewerBlocks()
	h := s.getLongestChainsLatestHash()
	return s.getChainByLatestHash(h)
}

func (s Service) storeBlocks(blocks []*Block) {
	for _, b := range blocks {
		s.blockByHash[b.CurrentBlockHash] = b
	}
}

func (s Service) downloadNewerBlocks() {
	currentBlock := s.helper.GetLatestBlock()
	s.blockByHash[currentBlock.CurrentBlockHash] = &currentBlock
	// keep downloading until reach one that we have seen already
	for _, ok := s.blockByHash[currentBlock.CurrentBlockHash]; !ok; {
		currentBlock = s.helper.GetBlockByHash(currentBlock.PreviousBlockHash)
		s.blockByHash[currentBlock.CurrentBlockHash] = &currentBlock
	}
}

func (s Service) getLongestChainsLatestHash() string {
	blockNumber := 0
	longestBlocksHash := ""
	for _, b := range s.blockByHash {
		if b.Number > blockNumber {
			longestBlocksHash = b.CurrentBlockHash
			blockNumber = b.Number
		}
	}
	return longestBlocksHash
}

func (s Service) getChainByLatestHash(h string) []*Block {
	var chain []*Block
	for b, exists := s.blockByHash[h]; exists; {
		chain = append(chain, b)
	}
	reverse(chain)
	return chain
}

func reverse(arr []*Block) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}
