package trieutil

import (
	"math"

	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// NextPowerOf2 returns the next power of 2 >= the input
//
// Spec pseudocode definition:
//   def get_next_power_of_two(x: int) -> int:
//    """
//    Get next power of 2 >= the input.
//    """
//    if x <= 2:
//        return x
//    else:
//        return 2 * get_next_power_of_two((x + 1) // 2)
func NextPowerOf2(n int) int {
	if n <= 2 {
		return n
	}
	return 2 * NextPowerOf2((n+1)/2)
}

// PrevPowerOf2 returns the previous power of 2 >= the input
//
// Spec pseudocode definition:
//   def get_previous_power_of_two(x: int) -> int:
//    """
//    Get the previous power of 2 >= the input.
//    """
//    if x <= 2:
//        return x
//    else:
//        return 2 * get_previous_power_of_two(x // 2)
func PrevPowerOf2(n int) int {
	if n <= 2 {
		return n
	}
	return 2 * PrevPowerOf2(n/2)
}

// MerkleTree returns all the nodes in a merkle tree from inputting merkle leaves.
//
// Spec pseudocode definition:
//   def merkle_tree(leaves: Sequence[Hash]) -> Sequence[Hash]:
//    padded_length = get_next_power_of_two(len(leaves))
//    o = [Hash()] * padded_length + list(leaves) + [Hash()] * (padded_length - len(leaves))
//    for i in range(padded_length - 1, 0, -1):
//        o[i] = hash(o[i * 2] + o[i * 2 + 1])
//    return o
func MerkleTree(leaves [][]byte) [][]byte {
	paddedLength := NextPowerOf2(len(leaves))
	parents := make([][]byte, paddedLength)
	paddedLeaves := make([][]byte, paddedLength-len(leaves))

	for i := 0; i < len(parents); i++ {
		parents[i] = params.BeaconConfig().ZeroHash[:]
	}
	for i := 0; i < len(paddedLeaves); i++ {
		paddedLeaves[i] = params.BeaconConfig().ZeroHash[:]
	}

	merkleTree := make([][]byte, 0, len(parents)+len(leaves)+len(paddedLeaves))
	merkleTree = append(merkleTree, parents...)
	merkleTree = append(merkleTree, leaves...)
	merkleTree = append(merkleTree, paddedLeaves...)

	for i := len(paddedLeaves) - 1; i > 0; i-- {
		a := append(merkleTree[2*i], merkleTree[2*i+1]...)
		b := hashutil.Hash(a)
		merkleTree[i] = b[:]
	}

	return merkleTree
}

// ConcatGeneralizedIndices concats the generalized indices together.
//
// Spec pseudocode definition:
//   def concat_generalized_indices(*indices: GeneralizedIndex) -> GeneralizedIndex:
//    """
//    Given generalized indices i1 for A -> B, i2 for B -> C .... i_n for Y -> Z, returns
//    the generalized index for A -> Z.
//    """
//    o = GeneralizedIndex(1)
//    for i in indices:
//        o = GeneralizedIndex(o * get_previous_power_of_two(i) + (i - get_previous_power_of_two(i)))
//    return o
func ConcatGeneralizedIndices(indices []int) int {
	index := 1
	for _, i := range indices {
		index = index * PrevPowerOf2(i) + (i - PrevPowerOf2(i))
	}
	return index
}

// GeneralizedIndexLength returns the generalized index length from a given index.
//
// Spec pseudocode definition:
//   def get_generalized_index_length(index: GeneralizedIndex) -> int:
//    """
//    Return the length of a path represented by a generalized index.
//    """
//    return int(log2(index))
func GeneralizedIndexLength(index int) int {
	return int(math.Log2(float64(index)))
}
