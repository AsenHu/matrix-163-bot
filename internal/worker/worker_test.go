package worker

import (
	"testing"

	"maunium.net/go/mautrix/id"
)

func TestWorker_music(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input   string
		roomId  id.RoomID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var w Worker
			gotErr := w.music(tt.input, tt.roomId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("music() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("music() succeeded unexpectedly")
			}
		})
	}
}

func TestWorker_music(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input   string
		roomId  id.RoomID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var w Worker
			gotErr := w.music(tt.input, tt.roomId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("music() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("music() succeeded unexpectedly")
			}
		})
	}
}
