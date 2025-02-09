package matrix_test

import (
	"matrix-163-bot/internal/config"
	"matrix-163-bot/internal/matrix"
	"matrix-163-bot/internal/netease"
	"testing"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     *config.Account
		want    *mautrix.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := matrix.Login(tt.cfg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Login() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Login() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     *config.Account
		want    *mautrix.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := matrix.Login(tt.cfg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Login() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Login() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     *config.Account
		want    *mautrix.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := matrix.Login(tt.cfg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Login() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Login() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     *config.Account
		want    *mautrix.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := matrix.Login(tt.cfg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Login() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Login() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     *config.Account
		want    *mautrix.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := matrix.Login(tt.cfg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Login() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Login() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendMusic(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		client  *mautrix.Client
		song    netease.Song
		roomId  id.RoomID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := matrix.SendMusic(tt.client, tt.song, tt.roomId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SendMusic() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SendMusic() succeeded unexpectedly")
			}
		})
	}
}

func TestSendMusic(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		client  *mautrix.Client
		song    netease.Song
		roomId  id.RoomID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := matrix.SendMusic(tt.client, tt.song, tt.roomId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SendMusic() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SendMusic() succeeded unexpectedly")
			}
		})
	}
}

func TestSendMusic(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		client  *mautrix.Client
		song    netease.Song
		roomId  id.RoomID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := matrix.SendMusic(tt.client, tt.song, tt.roomId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SendMusic() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SendMusic() succeeded unexpectedly")
			}
		})
	}
}

func TestSendMusic(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		client  *mautrix.Client
		song    netease.Song
		roomId  id.RoomID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := matrix.SendMusic(tt.client, tt.song, tt.roomId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SendMusic() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SendMusic() succeeded unexpectedly")
			}
		})
	}
}
