package main

import (
	"reflect"
	"testing"
)

func Test_sortUtility(t *testing.T) {
	type args struct {
		args   []string
		toSort []string
	}
	arg1 := args{
		args:   []string{},
		toSort: []string{"a", "c", "b"},
	}
	arg2 := args{
		args:   []string{"-c"},
		toSort: []string{"a", "c", "b"},
	}

	arg3 := args{
		args:   []string{"-n"},
		toSort: []string{"a34", "b3", "c8", "22d8"},
	}

	arg4 := args{
		args:   []string{"-n"},
		toSort: []string{"lvl34", "lvl8", "lvl3", "lvl22"},
	}

	arg5 := args{
		args:   []string{"-nc"},
		toSort: []string{"a34", "b3", "c8", "22d8"},
	}

	arg6 := args{
		args:   []string{"-ncr"},
		toSort: []string{"a34", "b3", "c8", "22d8"},
	}

	arg7 := args{
		args:   []string{"-nr"},
		toSort: []string{"a34", "b3", "c8", "22d8"},
	}

	arg8 := args{
		args:   []string{"-nr"},
		toSort: []string{"lvl34", "lvl8", "lvl3", "lvl22"},
	}

	arg9 := args{
		args:   []string{"-n", "-k=2"},
		toSort: []string{"34 54", "22 41", "1 88", "89 15"},
	}

	arg10 := args{
		args:   []string{"-n", "-k=1"},
		toSort: []string{"34 54", "22 41", "1 88", "89 15"},
	}

	arg11 := args{
		args:   []string{"-nc", "-k=1"},
		toSort: []string{"1 88", "22 41", "34 54", "89 15"},
	}

	arg12 := args{
		args:   []string{"-u"},
		toSort: []string{"aaaa", "aaa", "b", "bbb", "aaa", "b"},
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"Without flags", arg1, []string{"a", "b", "c"}, false},
		{"Check is sorted without another flags should return error", arg2, nil, true},
		{"Numeric sort in disorganized strings", arg3, []string{"b3", "c8", "a34", "22d8"}, false},
		{"Numeric sort in organized strings", arg4, []string{"lvl3", "lvl8", "lvl22", "lvl34"}, false},
		{"Check numeric sort in disorganized strings should return error", arg5, nil, true},
		{"Check numeric sort in organized strings should return error", arg6, nil, true},
		{"Numeric sort in disorganized strings reversed", arg7, []string{"22d8", "a34", "c8", "b3"}, false},
		{"Numeric sort in organized strings reversed", arg8, []string{"lvl34", "lvl22", "lvl8", "lvl3"}, false},
		{"Numeric sort by column 2", arg9, []string{"89 15", "22 41", "34 54", "1 88"}, false},
		{"Numeric sort by column 1", arg10, []string{"1 88", "22 41", "34 54", "89 15"}, false},
		{"Check Numeric sort by column 1 shouldn't return error", arg11, nil, false},
		{"Unique simple sort", arg12, []string{"aaa", "aaaa", "b", "bbb"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sortUtility(tt.args.args, tt.args.toSort)
			if (err != nil) != tt.wantErr {
				t.Errorf("sortUtility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortUtility() got = %v, want %v", got, tt.want)
			}
		})
	}
}
