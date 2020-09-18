package card

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestParserCards(t *testing.T) {
	var (
		ct int
		err error
	)

	ct,err = ParserCards("a")
	assert.Equal(t,ct,0,"sole")
	assert.Equal(t,err,nil,"sole err")
}
