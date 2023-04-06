package zmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCmp_New(t *testing.T) {
	check := assert.New(t)
	cmp := NewChMap()
	HelpBench(1000, cmp)
	check.Equal(1000, len(cmp.data))
}
