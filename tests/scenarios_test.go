package scenarios_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infratographer/krakend-endpoints-tool/cmd"
)

func TestHappyPath(t *testing.T) {
	t.Parallel()

	endpointsDir := "scenarios/happy-path"
	cfg := "scenarios/happy-path/krakend.tmpl"
	outf := filepath.Join(t.TempDir(), "krakend.tmpl")
	err := cmd.Generate(endpointsDir, cfg, outf, "$ENDPOINTS$", false)
	assert.NoError(t, err, "should not fail")
}
