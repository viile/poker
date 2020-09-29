package landlord

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestNewCardsContain(t *testing.T) {
	var r bool
	r = NewCards("345667").Contain(NewCards("346"))
	assert.Equal(t,r,true,"Contain")

	r = NewCards("45667").Contain(NewCards("346"))
	assert.Equal(t,r,false,"Contain")

	r = NewCards("45667").Contain(NewCards("46667"))
	assert.Equal(t,r,false,"Contain")
}

func TestNewCardsRemove(t *testing.T) {
	var r string
	r = NewCards("345667").Remove(NewCards("346")).String()
	assert.Equal(t,r,"567","Remove")

}

func TestNewCardsCompare(t *testing.T) {
	var (
		i int
		err error
	)

	i,err = NewCards("a").Compare(NewCards("2"))
	assert.Equal(t,i,1,"Compare")
	assert.Equal(t,err,nil,"err")

	i,err = NewCards("a").Compare(NewCards("3"))
	assert.Equal(t,i,-1,"Compare")
	assert.Equal(t,err,nil,"err")

	i,err = NewCards("a").Compare(NewCards("a"))
	assert.Equal(t,i,0,"Compare")
	assert.Equal(t,err,nil,"err")

	i,err = NewCards("a").Compare(NewCards("3333"))
	assert.Equal(t,i,1,"Compare")
	assert.Equal(t,err,nil,"err")

	i,err = NewCards("a").Compare(NewCards("sb"))
	assert.Equal(t,i,1,"Compare")
	assert.Equal(t,err,nil,"err")

	i,err = NewCards("34445556").Compare(NewCards("36667778"))
	assert.Equal(t,i,1,"Compare")
	assert.Equal(t,err,nil,"err")

	i,err = NewCards("36667778").Compare(NewCards("9999"))
	assert.Equal(t,i,1,"Compare")
	assert.Equal(t,err,nil,"err")

	i,err = NewCards("34567").Compare(NewCards("56789"))
	assert.Equal(t,i,1,"Compare")
	assert.Equal(t,err,nil,"err")

	i,err = NewCards("4444").Compare(NewCards("sb"))
	assert.Equal(t,i,1,"Compare")
	assert.Equal(t,err,nil,"err")

	i,err = NewCards("344447").Compare(NewCards("588889"))
	assert.Equal(t,i,1,"Compare")
	assert.Equal(t,err,nil,"err")

}

func TestNewCardsParser(t *testing.T) {
	var (
		ct int
		err error
	)

	ct,err = NewCards("a").Parser()
	assert.Equal(t,ct, TypeSole,"TypeSole")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("33").Parser()
	assert.Equal(t,ct, TypePair,"TypePair")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("sb").Parser()
	assert.Equal(t,ct, TypeJokerBomb,"TypeJokerBomb")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("333").Parser()
	assert.Equal(t,ct, TypeTrio,"TypeTrio")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("333a").Parser()
	assert.Equal(t,ct, TypeTrioSole,"TypeTrioSole")
	assert.Equal(t,err,nil,"err")
	ct,err = NewCards("3444").Parser()
	assert.Equal(t,ct, TypeTrioSole,"TypeTrioSole")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("3333").Parser()
	assert.Equal(t,ct, TypeBomb,"TypeBomb")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("33344").Parser()
	assert.Equal(t,ct, TypeTrioPair,"TypeTrioPair")
	assert.Equal(t,err,nil,"err")
	ct,err = NewCards("33444").Parser()
	assert.Equal(t,ct, TypeTrioPair,"TypeTrioPair")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("334455").Parser()
	assert.Equal(t,ct, TypePairChain,"TypePairChain")
	assert.Equal(t,err,nil,"err")
	ct,err = NewCards("3344556677").Parser()
	assert.Equal(t,ct, TypePairChain,"TypePairChain")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("444555").Parser()
	assert.Equal(t,ct, TypeAirplane,"TypeAirplane")
	assert.Equal(t,err,nil,"err")
	ct,err = NewCards("444555666777").Parser()
	assert.Equal(t,ct, TypeAirplane,"TypeAirplane")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("3444555k").Parser()
	assert.Equal(t,ct, TypeAirplaneSole,"TypeAirplaneSole")
	assert.Equal(t,err,nil,"err")
	ct,err = NewCards("34555666").Parser()
	assert.Equal(t,ct, TypeAirplaneSole,"TypeAirplaneSole")
	ct,err = NewCards("3334445k").Parser()
	assert.Equal(t,ct, TypeAirplaneSole,"TypeAirplaneSole")
	assert.Equal(t,err,nil,"err")
	ct,err = NewCards("34445556667k").Parser()
	assert.Equal(t,ct, TypeAirplaneSole,"TypeAirplaneSole")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("33444555kk").Parser()
	assert.Equal(t,ct, TypeAirplanePair,"TypeAirplanePair")
	assert.Equal(t,err,nil,"err")
	ct,err = NewCards("33444555kk").Parser()
	assert.Equal(t,ct, TypeAirplanePair,"TypeAirplanePair")
	assert.Equal(t,err,nil,"err")
	ct,err = NewCards("33444555kk").Parser()
	assert.Equal(t,ct, TypeAirplanePair,"TypeAirplanePair")
	assert.Equal(t,err,nil,"err")
	ct,err = NewCards("3344455566677kk").Parser()
	assert.Equal(t,ct, TypeAirplanePair,"TypeAirplanePair")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("333345").Parser()
	assert.Equal(t,ct, TypeDualSole,"TypeDualSole")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("344445").Parser()
	assert.Equal(t,ct, TypeDualSole,"TypeDualSole")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("345555").Parser()
	assert.Equal(t,ct, TypeDualSole,"TypeDualSole")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("33334455").Parser()
	assert.Equal(t,ct, TypeDualPair,"TypeDualPair")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("33444455").Parser()
	assert.Equal(t,ct, TypeDualPair,"TypeDualPair")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("33445555").Parser()
	assert.Equal(t,ct, TypeDualPair,"TypeDualPair")
	assert.Equal(t,err,nil,"err")

	ct,err = NewCards("56789tj").Parser()
	assert.Equal(t,ct, TypeSoleChain,"TypeSoleChain")
	assert.Equal(t,err,nil,"err")

}
