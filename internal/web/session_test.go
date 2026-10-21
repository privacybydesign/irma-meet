package web

import (
	"bytes"
	"testing"
	"time"
)

func TestSealedCodecRoundTripsAndHidesTheData(t *testing.T) {
	codec, err := newSealedCodec(bytes.Repeat([]byte{1}, 32))
	if err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(time.Hour).Round(0)

	data, err := codec.Encode(deadline, map[string]any{"verified": []byte(`{"email":"alice@example.com"}`)})
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Contains(data, []byte("alice@example.com")) {
		t.Error("the encoded session contains the email address in plain text")
	}

	gotDeadline, values, err := codec.Decode(data)
	if err != nil || !gotDeadline.Equal(deadline) || string(values["verified"].([]byte)) != `{"email":"alice@example.com"}` {
		t.Errorf("round trip: %v %v %v", gotDeadline, values, err)
	}

	again, _ := codec.Encode(deadline, map[string]any{"verified": []byte("x")})
	other, _ := codec.Encode(deadline, map[string]any{"verified": []byte("x")})

	if bytes.Equal(again, other) {
		t.Error("two encodings of the same session are identical: the nonce is not random")
	}
}

func TestSealedCodecTreatsUnreadableDataAsAnEmptySession(t *testing.T) {
	codec, _ := newSealedCodec(bytes.Repeat([]byte{1}, 32))
	otherKey, _ := newSealedCodec(bytes.Repeat([]byte{2}, 32))

	sealed, _ := codec.Encode(time.Now().Add(time.Hour), map[string]any{"k": []byte("v")})
	tampered := append([]byte(nil), sealed...)
	tampered[len(tampered)-1] ^= 0xff

	for name, data := range map[string]([]byte){"other key": sealed, "tampered": tampered, "short": {1, 2, 3}, "empty": nil} {
		decoder := codec
		if name == "other key" {
			decoder = otherKey
		}

		deadline, values, err := decoder.Decode(data)
		if err != nil || len(values) != 0 || !deadline.After(time.Now()) {
			t.Errorf("%s: deadline %v, values %v, err %v: want an empty, live session", name, deadline, values, err)
		}
	}
}

func TestSealedCodecNeedsAnAESKey(t *testing.T) {
	if _, err := newSealedCodec([]byte("short")); err == nil {
		t.Error("a 5 byte key was accepted")
	}
}
