//
// miner.go
// Staring template for PointCoint miner.
//
// cs4501: Cryptocurrency Cafe
// University of Virginia, Spring 2015
// Project 2
//

package main

import (
	//"fmt"
	"log"
	//"math/big"
	//"math/rand"
	//"time"

	"github.com/PointCoin/btcjson"
	"github.com/PointCoin/btcutil"
	//"github.com/PointCoin/pointcoind/blockchain"
)

const (
	// This should match your settings in pointcoind.conf
	rpcuser = "jakerose27"
	rpcpass = "taylor1"
	// This file should exist if pointcoind was setup correctly
	cert    = "/home/ubuntu/.pointcoind/rpc.cert"
)

func main() {
	// Setup the client using application constants, fail horribly if there's a problem
	client := setupRpcClient(cert, rpcuser, rpcpass)
	var prevHeight int64
	prevHeight = 0
	for { // Loop forever (you may want to do something smarter!)
		// Get a new block template from pointcoind.
		log.Printf("Requesting a block template\n")
		template, err := client.GetBlockTemplate(&btcjson.TemplateRequest{})
		if err != nil {
			log.Fatal(err)
		}

		// The template returned by GetBlockTemplate provides these fields that 
		// you will need to use to create a new block:

		// The hash of the previous block
		prevHash := template.PreviousHash

		// The difficulty target
		difficulty := formatDiff(template.Bits)

		// The height of the next block (number of blocks between genesis block and next block)
		height := template.Height
		if prevHeight == 0 {
			prevHeight = height
		}
		if height != prevHeight {
			log.Printf("Block height: %d", height)
			prevHeight = height
		}

		// The transactions from the network	
		txs := formatTransactions(template.Transactions) 

		// These are configurable parameters to the coinbase transaction
		msg := "jar8fe" // replace with your UVa Computing ID (e.g., "dee2b")
		a := "Pns4GnMUE3FJKPuuTDWYC17r4pDrUGEukm" // replace with the address you want mining fees to go to (or leave it like this and Nick gets them)

		coinbaseTx := CreateCoinbaseTx(height, a, msg)

		txs = prepend(coinbaseTx.MsgTx(), txs)
		merkleRoot := createMerkleRoot(txs)

		// Finish the miner!
		
		block := CreateBlock(prevHash, merkleRoot, difficulty, 0, txs)		
		//log.Printf("Source: %d", s)
		var i uint32
		for i = 0; i<10000000 ; i++ {
			//log.Printf("Nonce number %d", i)
			block.Header.Nonce = i
			hash, _err := block.Header.BlockSha()
			if _err != nil{
				continue
			}
			if lessThanDiff(hash, difficulty) {
				//log.Printf("Submitting block...")
				err := client.SubmitBlock(btcutil.NewBlock(block), nil)
				if err != nil {
					log.Print(err)
					break
				} else {
					log.Printf("WIN")
					break
				}
			}
		}
	}
}
