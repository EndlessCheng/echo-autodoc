package autodoc

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_readTailComment(t *testing.T) {
	assert.Equal(t, "aa", _readTailComment(")// aa"))
	assert.Equal(t, "aa", _readTailComment("))// aa"))
	assert.Equal(t, "/ aa", _readTailComment(")/// aa"))
	assert.Equal(t, "aa", _readTailComment(")//aa"))
	assert.Equal(t, "aa", _readTailComment(") // aa"))
	assert.Equal(t, "aa", _readTailComment(")    // aa"))
	assert.Equal(t, "", _readTailComment(")    // "))
	assert.Equal(t, "", _readTailComment(") n // aa"))
	assert.Equal(t, "", _readTailComment(" // aa"))
	assert.Equal(t, "", _readTailComment(" aa"))
	assert.Equal(t, "", _readTailComment("//"))
	assert.Equal(t, "", _readTailComment("/"))
	assert.Equal(t, "", _readTailComment(")"))
	assert.Equal(t, "", _readTailComment(") "))
	assert.Equal(t, "", _readTailComment(") /"))
	assert.Equal(t, "", _readTailComment(""))
}
