package proto

import (
	"context"
	"testing"
)

func TestParseFullMethod_Valid(t *testing.T) {
	tests := []struct {
		input       string
		wantService string
		wantMethod  string
	}{
		{"/helloworld.Greeter/SayHello", "helloworld.Greeter", "SayHello"},
		{"/pkg.v1.MyService/DoThing", "pkg.v1.MyService", "DoThing"},
		{"noSlashPrefix/Method", "noSlashPrefix", "Method"},
	}
	for _, tc := range tests {
		svc, meth, err := ParseFullMethod(tc.input)
		if err != nil {
			t.Errorf("ParseFullMethod(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if svc != tc.wantService {
			t.Errorf("ParseFullMethod(%q) service = %q, want %q", tc.input, svc, tc.wantService)
		}
		if meth != tc.wantMethod {
			t.Errorf("ParseFullMethod(%q) method = %q, want %q", tc.input, meth, tc.wantMethod)
		}
	}
}

func TestParseFullMethod_Invalid(t *testing.T) {
	invalid := []string{
		"",
		"/",
		"/OnlyService",
		"//EmptyService",
		"/Service/",
	}
	for _, tc := range invalid {
		_, _, err := ParseFullMethod(tc)
		if err == nil {
			t.Errorf("ParseFullMethod(%q) expected error, got nil", tc)
		}
	}
}

func TestResolveMethod_NilConn_ReturnsError(t *testing.T) {
	// ResolveMethod should propagate the invalid method format error
	// before attempting to use the connection.
	_, err := ResolveMethod(context.Background(), nil, "bad-format")
	if err == nil {
		t.Fatal("expected error for invalid method format, got nil")
	}
}

func TestShortName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"helloworld.Greeter", "Greeter"},
		{"Greeter", "Greeter"},
		{"a.b.c.Service", "Service"},
	}
	for _, tc := range cases {
		got := shortName(tc.input)
		if got != tc.want {
			t.Errorf("shortName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
