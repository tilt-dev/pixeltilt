package client

import (
	"testing"
)

func TestSillyFunc(t *testing.T) {
	s := SillyFunc(1, 1)
	if s != 2 {
		t.Errorf("wrong silly; wanted 2 got: %d", s)
	}
}

func TestTheSillyFuncStrikesBack(t *testing.T) {
	s := TheSillyFuncStrikesBack(1, 1)
	if s != 2 {
		t.Errorf("wrong silly; wanted 2 got: %d", s)
	}
}

func TestReturnOfTheSillyFunc(t *testing.T) {
	s := ReturnOfTheSillyFunc(1, 1)
	if s != 2 {
		t.Errorf("wrong silly; wanted 2 got: %d", s)
	}
}

func TestThePhantomSillyFunc(t *testing.T) {
	s := ThePhantomSillyFunc(1, 1)
	if s != 2 {
		t.Errorf("wrong silly; wanted 2 got: %d", s)
	}
}

func TestAttackOfTheSilly(t *testing.T) {
	s := AttackOfTheSilly(1, 1)
	if s != 2 {
		t.Errorf("wrong silly; wanted 2 got: %d", s)
	}
}

func TestRevengeOfTheSilly(t *testing.T) {
	s := RevengeOfTheSilly(1, 1)
	if s != 2 {
		t.Errorf("wrong silly; wanted 2 got: %d", s)
	}
}

func TestSillyOne(t *testing.T) {
	s := SillyOne(1, 1)
	if s != 2 {
		t.Errorf("wrong silly; wanted 2 got: %d", s)
	}
}
