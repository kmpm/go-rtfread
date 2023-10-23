package rtfread

import (
	"log/slog"
	"os"
	"testing"
)

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}
	file, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	logger := slog.New(slog.NewTextHandler(file, opts))
	slog.SetDefault(logger)
}

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
		{"infotext", args{"./testdata/infotext.rtf"}, want{"./testdata/infotext.txt"}, false},
		{"infotext-2", args{"testdata/infotext-2.rtf"}, want{"testdata/infotext-2.txt"}, false},
		{"ad", args{"testdata/ad.rtf"}, want{"testdata/ad.txt"}, false},
		{"np.new", args{"testdata/np.new.rtf"}, want{"testdata/np.new.txt"}, false},
		// failing likeness tests
		// {"file-sample_100kB", args{"./testdata/file-sample_100kB.rtf"}, want{"./testdata/file-sample_100kB.txt"}, false},
		// {"sample", args{"./testdata/sample.rtf"}, want{"./testdata/sample.txt"}, false},
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
				err = os.WriteFile(tt.name+".got.txt", []byte(got), 0666)
				if err != nil {
					t.Errorf("write want error = %v", err)
				}
				t.Errorf("\nToString() = %#v,\n      want = %#v", got, want)
			}
		})
	}
}
