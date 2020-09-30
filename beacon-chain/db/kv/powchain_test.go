package kv

import (
	"context"
	"testing"

	"github.com/prysmaticlabs/prysm/proto/beacon/db"
)

func TestStore_SavePowchainData(t *testing.T) {
	type args struct {
		data *db.ETH1ChainData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil data",
			args: args{
				data: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := setupDB(t)
			if err := store.SavePowchainData(context.Background(), tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SavePowchainData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}