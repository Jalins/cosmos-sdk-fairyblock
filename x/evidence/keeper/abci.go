package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/core/comet"
	"cosmossdk.io/x/evidence/types"
)

// BeginBlocker iterates through and handles any newly discovered evidence of
// misbehavior submitted by CometBFT. Currently, only equivocation is handled.
func (k Keeper) BeginBlocker(ctx context.Context, cometService comet.Service) error {
	//defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	bi := cometService.CometInfo(ctx)

	evidences := bi.Evidence
	for _, evidence := range evidences {
		switch evidence.Type {
		// It's still ongoing discussion how should we treat and slash attacks with
		// premeditation. So for now we agree to treat them in the same way.
		case comet.LightClientAttack, comet.DuplicateVote:
			evidence := types.FromABCIEvidence(evidence, k.stakingKeeper.ConsensusAddressCodec())
			err := k.handleEquivocationEvidence(ctx, evidence)
			if err != nil {
				return err
			}
		default:
			k.Logger.Error(fmt.Sprintf("ignored unknown evidence type: %x", evidence.Type))
		}
	}
	return nil
}
