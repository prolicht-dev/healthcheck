package main

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func Test_checkWebEndpointFromArgs(t *testing.T) {
	// We manipuate the Args to set them up for the testcases
	// after this test we restore the initial args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name string
		args []string
		want int
	}{
		{
			name: "Test Invalid Arguments",
			args: nil,
			want: 1,
		},
		{
			name: "Test Valid Arguments",
			args: []string{"https://httpstat.us/200"},
			want: 0,
		},
		{
			name: "Test Valid Arguments Wrong Url",
			args: []string{"https://httpstat.us/500"},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// we need a value to set Args[0] to the binary name, cause normal program arguments begin at Args[1]
			os.Args = append([]string{"hc"}, tt.args...)

			if got := checkWebEndpointFromArgs(); got != tt.want {
				t.Errorf("checkWebEndpointFromArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkWebEndpoint(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test Success 200",
			args: args{
				url: "https://httpstat.us/200",
			},
			want: true,
		},
		{
			name: "Test Success 204",
			args: args{
				url: "https://httpstat.us/204",
			},
			want: true,
		},
		{
			name: "Test Failure 400",
			args: args{
				url: "https://httpstat.us/400",
			},
			want: false,
		},
		{
			name: "Test Failure 500",
			args: args{
				url: "https://httpstat.us/400",
			},
			want: false,
		},
		{
			name: "Test Redirect 301",
			args: args{
				url: "https://httpstat.us/301",
			},
			want: true,
		},
		{
			name: "Test Failure Timeout",
			args: args{
				url: "https://httpstat.us/524",
			},
			want: false,
		},
		{
			name: "Test Failure Nonexistent",
			args: args{
				url: "https://thisdomainshouldnot.exist/503",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkWebEndpoint(tt.args.url); got != tt.want {
				t.Errorf("checkWebEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testMain(t *testing.T) {
	// Only run the main part when a specific env variable is set
	if os.Getenv("TEST_MAIN") == "1" {
		// We manipulate the Args to set them up for the testcases
		// after this test we restore the initial args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{"hc", os.Getenv("TEST_MAIN_URL")}
		main()

		return
	}

	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Test Invalid Arguments",
			url:  "",
			want: 1,
		},
		{
			name: "Test Valid Arguments",
			url:  "https://httpstat.us/200",
			want: 0,
		},
		{
			name: "Test Valid Arguments Wrong Url",
			url:  "https://httpstat.us/500",
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start the actual test in a different subprocess
			cmd := exec.Command(os.Args[0], "-test.run=Test_testMain")
			cmd.Env = append(os.Environ(), "TEST_MAIN=1", "TEST_MAIN_URL="+tt.url)
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}

			// Check that the program exited
			err := cmd.Wait()

			if exiterr, ok := err.(*exec.ExitError); ok {
				// The program has exited with an exit code != 0

				// This works on both Unix and Windows. Although package
				// syscall is generally platform dependent, WaitStatus is
				// defined for both Unix and Windows and in both cases has
				// an ExitStatus() method with the same signature.
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					exitStatus := status.ExitStatus()
					if exitStatus != tt.want {
						t.Errorf("main() = %v, want %v", exitStatus, tt.want)
					}
				} else {
					t.Errorf("main() = unknown err: %v, want %v", err, tt.want)
				}
			} else {
				if err == nil && tt.want != 0 {
					t.Errorf("main() = 0, want %v", tt.want)
				}
			}
		})
	}
}
