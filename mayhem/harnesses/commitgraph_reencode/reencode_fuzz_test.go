package commitgraphreencode

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-git/go-git/v6/plumbing/format/commitgraph"
)

// FuzzCommitgraphReencode opens the input as a commit-graph file with the
// public commitgraph.OpenFileIndex and writes the resulting Index back out
// with commitgraph.NewEncoder(w).Encode, the way a tool that rewrites or
// copies a repository's .git/objects/info/commit-graph does. Every error is
// a legitimate rejection of a malformed file; only a panic is a finding.
func FuzzCommitgraphReencode(f *testing.F) {
	f.Add([]byte("CGPH\x01\x01\x00\x00"))

	f.Fuzz(func(_ *testing.T, data []byte) {
		idx, err := commitgraph.OpenFileIndex(struct {
			io.ReaderAt
			io.Closer
		}{
			bytes.NewReader(data),
			io.NopCloser(nil),
		})
		if err != nil {
			return
		}
		defer idx.Close()

		var out bytes.Buffer
		_ = commitgraph.NewEncoder(&out).Encode(idx)
	})
}
