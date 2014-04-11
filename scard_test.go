package scard

import (
	"testing"
)

// FIXME: all tests assume pcscd is running, and there is a single reader with a
// card present

type testCard struct {
	ctx  *Context
	card *Card
}

func setup(t *testing.T) *testCard {
	var ctx *Context
	var card *Card

	ctx, err := EstablishContext()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	defer func() {
		if card == nil {
			ctx.Release()
		}
	}()

	readers, err := ctx.ListReaders()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(readers) == 0 {
		t.Logf("No readers")
		t.FailNow()
	}

	card, err = ctx.Connect(readers[0], SHARE_EXCLUSIVE, PROTOCOL_ANY)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	var status *CardStatus
	status, err = card.Status()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("CardState: %+#v\n", status)
	if status.State.String() != PRESENT.String() {

		//t.Errorf("unexpected state: %s\n", status.State.String())
		//	t.FailNow()
	}

	return &testCard{ctx: ctx, card: card}
}

func teardown(c *testCard) {
	if c.card != nil {
		c.card.Disconnect(LEAVE_CARD)
	}

	if c.ctx != nil {
		c.ctx.Release()
	}
}

func TestListReaders(t *testing.T) {
	ctx, err := EstablishContext()
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Release()
	readers, err := ctx.ListReaders()
	if err != nil {
		t.Fatal(err)
	}
	for _, reader := range readers {
		t.Log(reader)
	}
}

func TestListReaderGroups(t *testing.T) {
	ctx, err := EstablishContext()
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Release()
	groups, err := ctx.ListReaderGroups()
	if err != nil {
		t.Fatal(err)
	}
	for _, group := range groups {
		t.Log(group)
	}
}

func TestGetAttrib(t *testing.T) {
	c := setup(t)
	defer teardown(c)

	atr, err := c.card.GetAttrib(ATTR_ATR_STRING)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("ATTR_ATR_STRING: % x\n", atr)
	}

	vendor, err := c.card.GetAttrib(ATTR_VENDOR_NAME)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("ATTR_VENDOR_NAME: %s [% x]\n", string(vendor), vendor)
	}
}
