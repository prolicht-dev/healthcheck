package healthcheck

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Test Plain",
			args: args{},
			want: &Service{
				ctx:           context.Background(),
				hasContext:    false,
				listenAddress: ":11223",
				checkFunc:     nil, // we cannot compare the check function
			},
		},
		{
			name: "Test With Empty Options",
			args: args{
				opts: []Option{},
			},
			want: &Service{
				ctx:           context.Background(),
				hasContext:    false,
				listenAddress: ":11223",
				checkFunc:     nil, // we cannot compare the check function
			},
		},
		{
			name: "Test With Options",
			args: args{
				opts: []Option{ListenOn(":123456"), WithContext(context.Background())},
			},
			want: &Service{
				ctx:           context.Background(),
				hasContext:    true,
				listenAddress: ":123456",
				checkFunc:     nil, // we cannot compare the check function
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.opts...)
			got.checkFunc = nil
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Start(t *testing.T) {
	ctxC, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name    string
		srvc    *Service
		wantErr bool
	}{
		{
			name:    "Test Plain Start",
			srvc:    New(),
			wantErr: false,
		},
		{
			name:    "Test Duplicate Port",
			srvc:    New(),
			wantErr: true,
		},
		{
			name:    "Test Context Start",
			srvc:    New(ListenOn(":11224"), WithContext(ctxC)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.srvc.Start()

			time.Sleep(200 * time.Millisecond) // let web server goroutine start

			client := &http.Client{
				Timeout: 2 * time.Second,
			}

			resp, err := client.Get("http://localhost" + tt.srvc.listenAddress + "/health")
			if resp != nil {
				if resp.StatusCode != 200 && !tt.wantErr {
					t.Errorf("Start() failed, invalid response status %d", resp.StatusCode)
				}
			} else if !tt.wantErr {
				t.Errorf("Start() failed, web server not reachable: %v", err)
			}

			if tt.srvc.hasContext {
				cancel()
				time.Sleep(200 * time.Millisecond) // let web server goroutine stop

				resp, err = client.Get("http://localhost" + tt.srvc.listenAddress + "/health")
				if err == nil {
					t.Errorf("Start() failed, web-server is still alive: %d", resp.StatusCode)
				}
			}
		})
	}
}

func TestListenOn(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Test Port Only",
			args: args{
				addr: ":8080",
			},
			want: &Service{
				ctx:           New().ctx,
				hasContext:    New().hasContext,
				listenAddress: ":8080",
				checkFunc:     nil, // cannot deeply compare check func,
			},
		},
		{
			name: "Test Addr:Port Only",
			args: args{
				addr: "localhost:8080",
			},
			want: &Service{
				ctx:           New().ctx,
				hasContext:    New().hasContext,
				listenAddress: "localhost:8080",
				checkFunc:     nil, // cannot deeply compare check func,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(ListenOn(tt.args.addr))
			got.checkFunc = nil
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListenOn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithContext(t *testing.T) {
	ctxC, cancel := context.WithCancel(context.Background())
	defer cancel()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Test Normal Context",
			args: args{
				ctx: context.Background(),
			},
			want: &Service{
				ctx:           context.Background(),
				hasContext:    true,
				listenAddress: New().listenAddress,
				checkFunc:     nil, // cannot deeply compare check func,
			},
		},
		{
			name: "Test Cancel Context",
			args: args{
				ctx: ctxC,
			},
			want: &Service{
				ctx:           ctxC,
				hasContext:    true,
				listenAddress: New().listenAddress,
				checkFunc:     nil, // cannot deeply compare check func,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(WithContext(tt.args.ctx))
			got.checkFunc = nil
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCustomCheck(t *testing.T) {
	customFnc := func() int { return 123 }

	type args struct {
		fnc func() int
	}
	tests := []struct {
		name    string
		args    args
		want    *Service
		wantFnc func() int
	}{
		{
			name: "Test Custom Function",
			args: args{
				fnc: customFnc,
			},
			want: &Service{
				ctx:           New().ctx,
				hasContext:    New().hasContext,
				listenAddress: New().listenAddress,
				checkFunc:     nil, // cannot deeply compare check func,
			},
			wantFnc: customFnc,
		},
		{
			name: "Test Nil Function",
			args: args{
				fnc: nil,
			},
			want: &Service{
				ctx:           New().ctx,
				hasContext:    New().hasContext,
				listenAddress: New().listenAddress,
				checkFunc:     nil, // cannot deeply compare check func,
			},
			wantFnc: New().checkFunc,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(WithCustomCheck(tt.args.fnc))
			gotFnc := got.checkFunc
			got.checkFunc = nil
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithContext() = %v, want %v", got, tt.want)
			}

			if reflect.ValueOf(gotFnc).Pointer() != reflect.ValueOf(tt.wantFnc).Pointer() {
				t.Error("WithContext() function mismatch")
			}
		})
	}
}
