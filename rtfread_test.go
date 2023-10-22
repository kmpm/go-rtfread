package rtfread

import (
	"os"
	"testing"
)

func TestParseFile(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		compareFile string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		// {"infotext", args{"../testdata/infotext.rtf"}, want{"../testdata/infotext.txt"}, false},
		// {"infotext-2", args{"testdata/infotext-2.rtf"}, want{"testdata/infotext-2.txt"}, false},
		// {"ad", args{"testdata/ad.rtf"}, want{"testdata/ad.txt"}, false},
		// {"np.new", args{"testdata/np.new.rtf"}, want{"testdata/np.new.txt"}, false},
		{"sample", args{"./testdata/sample.rtf"}, want{"../testdata/sample.txt"}, false},
		// {"file-sample_100kB", args{"testdata/file-sample_100kB.rtf"}, want{"testdata/file-sample_100kB.txt"}, false},
		// TODO: Add more test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			data, err := os.ReadFile(tt.want.compareFile)
			if err != nil {
				t.Errorf("ToString() error = %v, wantErr %v", err, tt.wantErr)
			}
			want := string(data)

			if got != want {
				t.Errorf("\nToString() = %#v,\n      want = %#v", got, want)
			}
		})
	}
}
