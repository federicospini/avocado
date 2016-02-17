package judge

import (
	"crypto/rand"
	"errors"
	"io"

	"github.com/agl/ed25519"
	"github.com/golang/protobuf/proto"
	"github.com/jtremback/usc-core/wire"
)

// Phases of a Tx
// Created
// Confirmed
// Verified

func sliceTo64Byte(slice []byte) *[64]byte {
	var array [64]byte
	copy(array[:], slice[:64])
	return &array
}

func sliceTo32Byte(slice []byte) *[32]byte {
	var array [32]byte
	copy(array[:], slice[:32])
	return &array
}

func randomBytes(c uint) ([]byte, error) {
	b := make([]byte, c)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		return nil, err
	}
	return b, nil
}

type Phase int

const (
	PENDING_OPEN   Phase = 1
	OPEN           Phase = 2
	PENDING_CLOSED Phase = 3
	CLOSED         Phase = 4
)

type Channel struct {
	ChannelId string
	Phase     Phase

	OpeningTx         *wire.OpeningTx
	OpeningTxEnvelope *wire.Envelope

	LastFullUpdateTx         *wire.UpdateTx
	LastFullUpdateTxEnvelope *wire.Envelope

	Judge    *Judge
	Accounts []*Account

	Fulfillments [][]byte
}

type Account struct {
	Name    string
	Pubkey  []byte
	Address string
	Judge   *Judge
}

type Judge struct {
	Name    string
	Pubkey  []byte
	Privkey []byte
}

// NewJudge makes a new judge
func NewJudge(name string, address string) (*Judge, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Judge{
		Name:    name,
		Pubkey:  pub[:],
		Privkey: priv[:],
	}, nil
}

// VerifyOpeningTx checks the signatures and state of a fully-signed OpeningTx,
// unmarshals it and returns it.
func (ep *Judge) VerifyOpeningTx(ev *wire.Envelope) (*wire.Envelope, *wire.OpeningTx, error) {
	otx := wire.OpeningTx{}
	err := proto.Unmarshal(ev.Payload, &otx)
	if err != nil {
		return nil, nil, err
	}

	// Check signatures
	if !ed25519.Verify(sliceTo32Byte(otx.Pubkeys[0]), ev.Payload, sliceTo64Byte(ev.Signatures[0])) {
		return nil, nil, errors.New("signature 0 invalid")
	}
	if !ed25519.Verify(sliceTo32Byte(otx.Pubkeys[1]), ev.Payload, sliceTo64Byte(ev.Signatures[1])) {
		return nil, nil, errors.New("signature 1 invalid")
	}

	// Sign envelope
	ev.Signatures = append(ev.Signatures, [][]byte{ed25519.Sign(sliceTo64Byte(ep.Privkey), ev.Payload)[:]}...)

	return ev, &otx, nil
}

// NewChannel creates a new Channel from an Envelope containing an opening transaction,
// an Account and a Peer.
func (ep *Judge) NewChannel(ev *wire.Envelope, accounts []*Account) (*Channel, error) {
	otx := &wire.OpeningTx{}
	err := proto.Unmarshal(ev.Payload, otx)
	if err != nil {
		return nil, err
	}

	ch := &Channel{
		ChannelId:         otx.ChannelId,
		OpeningTx:         otx,
		OpeningTxEnvelope: ev,
		Accounts:          accounts,
		Judge:             ep,
		Phase:             OPEN,
	}

	return ch, nil
}

// VerifyUpdateTx checks the signatures and the state of a fully-signed UpdateTx,
// unmarshals it, signs it, and returns it.
func (ch *Channel) VerifyUpdateTx(ev *wire.Envelope) (*wire.Envelope, *wire.UpdateTx, error) {
	// Check signatures
	if !ed25519.Verify(sliceTo32Byte(ch.OpeningTx.Pubkeys[0]), ev.Payload, sliceTo64Byte(ev.Signatures[0])) {
		return nil, nil, errors.New("signature 0 invalid")
	}
	if !ed25519.Verify(sliceTo32Byte(ch.OpeningTx.Pubkeys[1]), ev.Payload, sliceTo64Byte(ev.Signatures[1])) {
		return nil, nil, errors.New("signature 1 invalid")
	}

	utx := &wire.UpdateTx{}
	err := proto.Unmarshal(ev.Payload, utx)
	if err != nil {
		return nil, nil, err
	}

	// Check ChannelId
	if utx.ChannelId != ch.OpeningTx.ChannelId {
		return nil, nil, errors.New("ChannelId does not match")
	}

	// Sign envelope
	ev.Signatures = append(ev.Signatures, [][]byte{ed25519.Sign(sliceTo64Byte(ch.Judge.Privkey), ev.Payload)[:]}...)

	return ev, utx, nil
}

func (ch *Channel) StartHoldPeriod(utx *wire.UpdateTx) error {
	if ch.Phase == PENDING_CLOSED {
		if ch.LastFullUpdateTx.SequenceNumber > utx.SequenceNumber {
			return errors.New("update tx with higher sequence number exists")
		}
	}
	ch.Phase = PENDING_CLOSED
	ch.LastFullUpdateTx = utx
	return nil
}

// AddFulfillment verifies a fulfillment's signature and adds it to the Channel's
// Fulfillments array.
func (ch *Channel) AddFulfillment(ev *wire.Envelope) error {
	if ch.Phase != PENDING_CLOSED {
		return errors.New("channel must be pending closed")
	}

	if !ed25519.Verify(sliceTo32Byte(ch.OpeningTx.Pubkeys[0]), ev.Payload, sliceTo64Byte(ev.Signatures[0])) ||
		!ed25519.Verify(sliceTo32Byte(ch.OpeningTx.Pubkeys[1]), ev.Payload, sliceTo64Byte(ev.Signatures[1])) {
		return errors.New("signature invalid")
	}

	ful := wire.Fulfillment{}
	err := proto.Unmarshal(ev.Payload, &ful)
	if err != nil {
		return err
	}

	ch.Fulfillments = append(ch.Fulfillments, ful.State)
	return nil
}