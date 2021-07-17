package writer

import (
	"os"
	"testing"

	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

func TestNewWriter(t *testing.T) {
	tmp := t.TempDir()
	f, err := os.CreateTemp(tmp, "")
	if err != nil {
		t.Errorf("failed to create temp file for new Writer, %v", err)
	}
	w := NewWriter(f)
	if w == nil {
		t.Errorf("the return of the NewWriter is nil")
	}
}

// TODO CheckFilePath

func TestWriteHeader(t *testing.T) {
	tests := []struct {
		sz       int
		expected string
	}{
		{sz: 5, expected: "header;A;B;C;D;E\n"},
		{sz: 27, expected: "header;A;B;C;D;E;F;G;H;I;J;K;L;M;N;O;P;Q;R;S;T;U;V;W;X;Y;Z;AA\n"},
		{sz: 53, expected: "header;A;B;C;D;E;F;G;H;I;J;K;L;M;N;O;P;Q;R;S;T;U;V;W;X;Y;Z;AA;AB;AC;AD;AE;AF;AG;AH;AI;AJ;AK;AL;AM;AN;AO;AP;AQ;AR;AS;AT;AU;AV;AW;AX;AY;AZ;BA\n"},
	}

	for _, tt := range tests {
		tmp := t.TempDir()
		f, _ := os.CreateTemp(tmp, "")
		defer func() { f.Close() }()

		w := NewWriter(f)
		w.Comma = ';'

		err := w.WriteHeader(tt.sz)
		// checks if write was sucessful
		if err != nil {
			t.Errorf(
				"failed to write header of size %d to file %s",
				tt.sz,
				f.Name(),
			)
		}

		// checks content of file
		b, err := os.ReadFile(f.Name())
		if err != nil {
			t.Errorf(
				"Failed to read file %s after write",
				f.Name(),
			)
		}

		if string(b) != tt.expected {
			t.Errorf(
				"error expected %v, received %v",
				tt.expected,
				string(b),
			)
		}
	}
}

func TestElementsObjs(t *testing.T) {

	tests := []struct {
		elems     models.Elements
		separator rune
		expected  string
		err       string
	}{
		{
			elems: models.Elements{
				{Objs: []float64{1.0, 2.0, 3.0}},
			},
			expected: "elem[0],1.000,2.000,3.000\n",
		},
		{
			elems: models.Elements{
				{Objs: []float64{0.01, 0.02, 0.03}},
				{Objs: []float64{0.004, 0.005, 0.006}},
			},
			separator: ';',
			expected:  "elem[0];0.010;0.020;0.030\nelem[1];0.004;0.005;0.006\n",
		},
		{
			err: "empty slice of elements",
		},
	}

	for _, tt := range tests {
		tmp := t.TempDir()
		f, _ := os.CreateTemp(tmp, "")
		defer func() { f.Close() }()

		w := NewWriter(f)
		if tt.separator != 0 {
			w.Comma = tt.separator
		}

		err := w.ElementsObjs(tt.elems)
		if err != nil && err.Error() != tt.err {
			t.Errorf(
				"failed ElementsObjs, got %v and expected %v",
				err,
				tt.err,
			)
		}

		b, _ := os.ReadFile(f.Name())

		if string(b) != tt.expected {
			t.Errorf(
				"error expected %v, got %v",
				tt.expected,
				string(b),
			)
		}
	}
}

func TestElementsVectors(t *testing.T) {

	tests := []struct {
		elems     models.Elements
		separator rune
		expected  string
		err       string
	}{
		{
			elems: models.Elements{
				{X: []float64{1.0, 2.0, 3.0}},
			},
			expected: "elem[0],1.000,2.000,3.000\n",
		},
		{
			elems: models.Elements{
				{X: []float64{0.01, 0.02, 0.03}},
				{X: []float64{0.004, 0.005, 0.006}},
			},
			separator: ';',
			expected:  "elem[0];0.010;0.020;0.030\nelem[1];0.004;0.005;0.006\n",
		},
		{
			err: "empty slice of elements",
		},
	}

	for _, tt := range tests {
		tmp := t.TempDir()
		f, _ := os.CreateTemp(tmp, "")
		defer func() { f.Close() }()

		w := NewWriter(f)
		if tt.separator != 0 {
			w.Comma = tt.separator
		}

		err := w.ElementsVectors(tt.elems)
		if err != nil && err.Error() != tt.err {
			t.Errorf(
				"failed ElementsObjs, got %v and expected %v",
				err,
				tt.err,
			)
		}

		b, _ := os.ReadFile(f.Name())

		if string(b) != tt.expected {
			t.Errorf(
				"error expected %v, got %v",
				tt.expected,
				string(b),
			)
		}
	}
}