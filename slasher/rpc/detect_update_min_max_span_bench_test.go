package rpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/prysmaticlabs/prysm/slasher/db"
)

func BenchmarkMinSpan(b *testing.B) {
	diffs := []uint64{2, 10, 100, 1000, 10000, 53999}
	dbs := db.SetupSlasherDB(b)
	defer db.TeardownSlasherDB(b, dbs)

	ctx := context.Background()
	slasherServer := &Server{
		SlasherDB: dbs,
	}
	for _, diff := range diffs {
		b.Run(fmt.Sprintf("MinSpan_diff_%d", diff), func(ib *testing.B) {
			for i := uint64(0); i < uint64(ib.N); i++ {
				_, err := slasherServer.DetectAndUpdateMinSpan(ctx, i+1, i+1+diff, i)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkMaxSpan(b *testing.B) {
	diffs := []uint64{2, 10, 100, 1000, 10000, 53999}
	dbs := db.SetupSlasherDB(b)
	defer db.TeardownSlasherDB(b, dbs)

	ctx := context.Background()
	slasherServer := &Server{
		SlasherDB: dbs,
	}
	for _, diff := range diffs {
		b.Run(fmt.Sprintf("MaxSpan_diff_%d", diff), func(ib *testing.B) {
			for i := uint64(0); i < uint64(ib.N); i++ {
				_, err := slasherServer.DetectAndUpdateMaxSpan(ctx, i, i+diff, i)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDetectSpan(b *testing.B) {
	diffs := []uint64{2, 10, 100, 1000, 10000, 53999}
	dbs := db.SetupSlasherDB(b)
	defer db.TeardownSlasherDB(b, dbs)

	slasherServer := &Server{
		SlasherDB: dbs,
	}
	for _, diff := range diffs {
		b.Run(fmt.Sprintf("Detect_MaxSpan_diff_%d", diff), func(ib *testing.B) {
			for i := uint64(0); i < uint64(ib.N); i++ {
				_, _, _, err := slasherServer.detectSpan(i, i+diff, i, detectMax)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
	for _, diff := range diffs {
		b.Run(fmt.Sprintf("Detect_MinSpan_diff_%d", diff), func(ib *testing.B) {
			for i := uint64(0); i < uint64(ib.N); i++ {
				_, _, _, err := slasherServer.detectSpan(i, i+diff, i, detectMin)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
