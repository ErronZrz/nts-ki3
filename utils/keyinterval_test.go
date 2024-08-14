package utils

import (
	"fmt"
	"testing"
)

func TestAnalyzeInterval(t *testing.T) {
	prefix := "C:\\Corner\\TMP\\毕设\\NTP\\Ntage15\\0814-key_id_data\\2024-07-30_"
	t1 := "2024073012"
	t2 := "2024081412"

	m, err := AnalyzeInterval(prefix, t1, t2)
	if err != nil {
		t.Error(err)
	}
	for k, v := range m {
		fmt.Println(k, v)
	}
}

func TestSaveIntervalTo(t *testing.T) {
	prefix := "C:\\Corner\\TMP\\毕设\\NTP\\Ntage15\\0814-key_id_data\\2024-07-30_"
	dst := "C:\\Corner\\TMP\\毕设\\NTP\\Ntage15\\key-id-pairs.txt"
	t1 := "2024073012"
	t2 := "2024081412"

	err := SaveIntervalTo(prefix, dst, t1, t2)
	if err != nil {
		t.Error(err)
	}
}

func TestCrossCompare(t *testing.T) {
	kePath := "C:\\Corner\\TMP\\毕设\\NTP\\Ntage15\\2024-08-14_ntske_0-with-ucloud.txt"
	itvPath := "C:\\Corner\\TMP\\毕设\\NTP\\Ntage15\\0814-key-id-pairs.txt"
	res, err := CrossCompare(kePath, itvPath)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}

func Test_intervalRange(t *testing.T) {
	type args struct {
		ids []string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
	}{
		{
			name: "test-1",
			args: args{
				ids: []string{"A", "A", "A", "B", "B", "C", "C", "D"},
			},
			want:  2,
			want1: 3,
		},
		{
			name: "test-2",
			args: args{
				ids: []string{"A", "A", "A", "B", "B", "C", "C"},
			},
			want:  2,
			want1: 3,
		},
		{
			name: "test-3",
			args: args{
				ids: []string{},
			},
			want:  -1,
			want1: -1,
		},
		{
			name: "test-4",
			args: args{
				ids: []string{"A", "A", "A"},
			},
			want:  2,
			want1: -1,
		},
		{
			name: "test-5",
			args: args{
				ids: []string{"A", "A", "B", "B", "C"},
			},
			want:  1,
			want1: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := intervalRange(tt.args.ids)
			if got != tt.want {
				t.Errorf("intervalRange() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("intervalRange() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
