package altair

import (
	"github.com/pkg/errors"
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	iface "github.com/prysmaticlabs/prysm/beacon-chain/state/interface"
	"github.com/prysmaticlabs/prysm/shared/mathutil"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// BaseReward takes state and validator index and calculate
// individual validator's base reward quotient.
//
// Spec code:
//  def get_base_reward(state: BeaconState, index: ValidatorIndex) -> Gwei:
//    increments = state.validators[index].effective_balance // EFFECTIVE_BALANCE_INCREMENT
//    return Gwei(increments * get_base_reward_per_increment(state))
func BaseReward(state iface.ReadOnlyBeaconState, index types.ValidatorIndex) (uint64, error) {
	val, err := state.ValidatorAtIndexReadOnly(index)
	if err != nil {
		return 0, err
	}
	totalBalance, err := helpers.TotalActiveBalance(state)
	if err != nil {
		return 0, errors.Wrap(err, "could not calculate active balance")
	}

	increments := val.EffectiveBalance() / params.BeaconConfig().EffectiveBalanceIncrement
	return baseRewardPerIncrement(totalBalance) * increments, nil
}

// baseRewardPerIncrement of the beacon state
//
// Spec code:
// def get_base_reward_per_increment(state: BeaconState) -> Gwei:
//    return Gwei(EFFECTIVE_BALANCE_INCREMENT * BASE_REWARD_FACTOR // integer_squareroot(get_total_active_balance(state)))
// This returns the base reward per increment for the beacon state.
//
// def get_base_reward_per_increment(state: BeaconState) -> Gwei:
//    return Gwei(EFFECTIVE_BALANCE_INCREMENT * BASE_REWARD_FACTOR // integer_squareroot(get_total_active_balance(state))
func baseRewardPerIncrement(activeBalance uint64) uint64 {
	return params.BeaconConfig().EffectiveBalanceIncrement * params.BeaconConfig().BaseRewardFactor / mathutil.IntegerSquareRoot(activeBalance)
}
