package kv

import (
	"github.com/boltdb/bolt"
	"github.com/prysmaticlabs/prysm/beacon-chain/db/filters"
)

// createIndicesFromFilters takes in filter criteria and returns
// a list of of byte keys used to retrieve the values stored
// for the indices from the DB. Typically, these are list of hash tree roots
// or signing roots of objects.
func createIndicesFromFilters(f *filters.QueryFilter) [][]byte {
	keys := make([][]byte, 0)
	for k, v := range f.Filters() {
		switch k {
		case filters.Shard:
			idx := append(shardIdx, uint64ToBytes(v.(uint64))...)
			keys = append(keys, idx)
		case filters.ParentRoot:
			parentRoot := v.([]byte)
			idx := append(parentRootIdx, parentRoot...)
			keys = append(keys, idx)
		}
	}
	return keys
}

// lookupValuesForIndices takes in a list of indices and looks up
// their corresponding values in the DB, returning a list of
// roots which can then be used for batch lookups of their corresponding
// objects from the DB. For example, if we are fetching
// attestations and we have an index `[]byte("shard-5")`,
// we might find roots `0x23` and `0x45` stored under that index. We can then
// do a batch read for attestations corresponding to those roots.
func lookupValuesForIndices(indices [][]byte, bkt *bolt.Bucket) [][][]byte {
	values := make([][][]byte, 0)
	for _, k := range indices {
		roots := bkt.Get(k)
		splitRoots := make([][]byte, 0)
		for i := 0; i < len(roots); i += 32 {
			splitRoots = append(splitRoots, roots[i:i+32])
		}
		values = append(values, splitRoots)
	}
	return values
}

// updateIndices updates the value for each index by appending it to the previous
// values stored at said index. Typically, indices are roots of data that can then
// be used for reads or batch reads from the DB.
func updateIndices(indices [][]byte, root []byte, bkt *bolt.Bucket) error {
	for _, idx := range indices {
		valuesAtIndex := bkt.Get(idx)
		if valuesAtIndex == nil {
			if err := bkt.Put(idx, root); err != nil {
				return err
			}
		} else {
			if err := bkt.Put(idx, append(valuesAtIndex, root...)); err != nil {
				return err
			}
		}
	}
	return nil
}
