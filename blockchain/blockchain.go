package blockchain

import (
	"log"

	"github.com/boltdb/bolt"
)

const blocksBucket = "blocks"

// Blockchain ...
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// Iterator ...
type Iterator struct {
	currentHash []byte
	db          *bolt.DB
}

// AddBlock ...
func (bc *Blockchain) AddBlock(data Chunks) {
	var lastHash []byte
	var err error

	// Check genesis block case
	if bc.tip != nil {
		err = bc.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			lastHash = b.Get([]byte("l"))

			return nil
		})

		if err != nil {
			log.Panic(err)
		}
	} else {

	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		slice := newBlock.Hash[:]

		b := tx.Bucket([]byte(blocksBucket))
		s := newBlock.Serialize()

		err := b.Put(slice, s)
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), slice)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = slice

		return nil
	})
}

// Iterator ...
func (bc *Blockchain) Iterator() *Iterator {
	bci := &Iterator{bc.tip, bc.db}

	return bci
}

// Next ...
func (i *Iterator) Next() *Block {
	var block *Block

	if i.currentHash == nil {
		return nil
	}

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash[:]

	return block
}

// NewBlockchain ...
func NewBlockchain(dbFile string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			// fmt.Println("No existing blockchain found. Creating a new one...")
			// genesis := NewGenesisBlock(data)

			_, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			// slice := genesis.Hash[:]

			// err = b.Put(slice, genesis.Serialize())
			// if err != nil {
			// 	log.Panic(err)
			// }

			// err = b.Put([]byte("l"), slice)
			// if err != nil {
			// 	log.Panic(err)
			// }
			// tip = slice
			tip = nil
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}
