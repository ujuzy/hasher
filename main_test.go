package main

import (
	"reflect"
	"testing"
)

func Test_splitUrls(t *testing.T) {
	type args struct {
		urls       []string
		groupCount int64
	}

	type testCase struct {
		name     string
		args     args
		expected [][]string
	}

	tests := []testCase{
		{
			name: "reqs count = group count",
			args: args{
				urls:       []string{"adjust.com", "google.com", "facebook.com", "yahoo.com", "yandex.com", "twitter.com"},
				groupCount: 6,
			},
			expected: [][]string{
				{"adjust.com"}, {"google.com"}, {"facebook.com"}, {"yahoo.com"}, {"yandex.com"}, {"twitter.com"},
			},
		},
		{
			name: "reqs count > group count (% != 0)",
			args: args{
				urls:       []string{"adjust.com", "google.com", "facebook.com", "yahoo.com", "yandex.com", "twitter.com"},
				groupCount: 4,
			},
			expected: [][]string{
				{"adjust.com", "yandex.com"}, {"google.com", "twitter.com"}, {"facebook.com"}, {"yahoo.com"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if actual := splitUrls(tc.args.urls, tc.args.groupCount); !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Actual %q not equal to expected %q", actual, tc.expected)
			}
		})
	}
}

func Test_sanitizeArgs(t *testing.T) {
	type testCase struct {
		name     string
		args     []string
		expected []string
	}

	tests := []testCase{
		{
			name:     "default",
			args:     []string{"adjust.com", "http://google.com", "youtube.com"},
			expected: []string{"http://adjust.com", "http://google.com", "http://youtube.com"},
		},
		{
			name:     "with parallel flag",
			args:     []string{"adjust.com", "http://google.com", "youtube.com", "-parallel", "10"},
			expected: []string{"http://adjust.com", "http://google.com", "http://youtube.com"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if actual := sanitizeArgs(tc.args); !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Actual %q not equal to expected %q", actual, tc.expected)
			}
		})
	}
}
