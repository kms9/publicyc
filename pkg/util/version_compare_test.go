package util

import "testing"

func TestVersionCompareV1EqualV2(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1==v2",
			args:args{
				v1: "1.1.1",
				v2: "1.1.1",
			},
			wantRet: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1>v2 major",
			args:args{
				v1: "2.1.1",
				v2: "1.1.1",
			},
			wantRet: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger2(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1>v2 minor",
			args:args{
				v1: "1.2.1",
				v2: "1.1.1",
			},
			wantRet: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger3(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1>v2 revision",
			args:args{
				v1: "1.1.2",
				v2: "1.1.1",
			},
			wantRet: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger4(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1>v2 longer",
			args:args{
				v1: "1.1.1.1",
				v2: "1.1.1",
			},
			wantRet: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger5(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1>v2 two digits",
			args:args{
				v1: "1.10.1",
				v2: "1.9.1",
			},
			wantRet: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger6(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1<v2 major",
			args:args{
				v1: "1.1.1",
				v2: "2.1.1",
			},
			wantRet: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger7(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1<v2 minor",
			args:args{
				v1: "1.1.1",
				v2: "1.2.1",
			},
			wantRet: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger8(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1<v2 revision",
			args:args{
				v1: "1.1.1",
				v2: "1.1.2",
			},
			wantRet: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger9(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1<v2 shorter",
			args:args{
				v1: "1.1.1",
				v2: "1.1.1.1",
			},
			wantRet: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger10(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1<v2 two digits",
			args:args{
				v1: "1.1.9",
				v2: "1.1.10",
			},
			wantRet: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger11(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1 == v2 all with letter v",
			args:args{
				v1: "v1.1.1",
				v2: "v1.1.1",
			},
			wantRet: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestVersionCompareV1Bigger12(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "版本比较，v1 == v2 one with letter v",
			args:args{
				v1: "v1.1.1",
				v2: "1.1.1",
			},
			wantRet: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := VersionCompare(tt.args.v1, tt.args.v2); gotRet != tt.wantRet {
				t.Errorf("VersionCompare() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}