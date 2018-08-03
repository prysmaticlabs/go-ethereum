package utils

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/beacon-chain/params"
)

func TestFaultyShuffleIndices(t *testing.T) {
	var list []int

	for i := 0; i < params.MaxValidators+1; i++ {
		list = append(list, i)
	}

	if _, err := ShuffleIndices(common.Hash{'a'}, list); err == nil {
		t.Error("Shuffle should have failed when validator count exceeds MaxValidators")
	}
}

func TestShuffleIndices(t *testing.T) {
	hash1 := common.BytesToHash([]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'a', 'b', 'c', 'd', 'e', 'f', 'g'})
	hash2 := common.BytesToHash([]byte{'1', '2', '3', '4', '5', '6', '7', '1', '2', '3', '4', '5', '6', '7', '1', '2', '3', '4', '5', '6', '7', '1', '2', '3', '4', '5', '6', '7', '1', '2', '3', '4', '5', '6', '7'})
	var list1 []int

	for i := 0; i < 100; i++ {
		list1 = append(list1, i)
	}

	list2 := make([]int, len(list1))
	copy(list2, list1)

	list1, err := ShuffleIndices(hash1, list1)
	if err != nil {
		t.Errorf("Shuffle failed with: %v", err)
	}

	list2, err = ShuffleIndices(hash2, list2)
	if err != nil {
		t.Errorf("Shuffle failed with: %v", err)
	}

	if reflect.DeepEqual(list1, list2) {
		t.Errorf("2 shuffled lists shouldn't be equal")
	}
}

func TestSplitIndices(t *testing.T) {
	var l []int
	validators := 64000
	for i := 0; i < validators; i++ {
		l = append(l, i)
	}
	split := SplitIndices(l, params.EpochLength)
	if len(split) != params.EpochLength {
		t.Errorf("Split list failed due to incorrect length, wanted:%v, got:%v", params.EpochLength, len(split))
	}

	for _, s := range split {
		if len(s) != validators/params.EpochLength {
			t.Errorf("Split list failed due to incorrect length, wanted:%v, got:%v", validators/params.EpochLength, len(s))
		}
	}
}
