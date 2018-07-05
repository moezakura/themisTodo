package utils

import (
	"testing"
)

func TestPasswordSha512_Hash(t *testing.T) {
	password := "password"
	// `expected` はPython3.6を使って生成した。生成したコードは以下の通り。
	// >>> import hashlib
	// >>> hashlib.sha512(b"password").hexdigest()
	expected := "b109f3bbbc244eb82441917ed06d618b9008dd09b3befd1b5e07394c706a8bb980b1d7785e5976ec049b46df5f1326af5a2ea6d103fd07c95385ffab0cacbc86"
	pw := PasswordSha512{}
	actual := pw.Hash(password)
	if expected != actual {
		t.Errorf("Expected: %v, Actual: %v", expected, actual)
	}
}

func TestPasswordSha512_Equal(t *testing.T) {
	pw := PasswordSha512{}
	// `hash` はPython3.6を使って生成した。生成したコードは以下の通り。
	// >>> import hashlib
	// >>> hashlib.sha512(b"password").hexdigest()
	hash := "b109f3bbbc244eb82441917ed06d618b9008dd09b3befd1b5e07394c706a8bb980b1d7785e5976ec049b46df5f1326af5a2ea6d103fd07c95385ffab0cacbc86"
	actual := pw.Equal(hash, hash)
	if true != actual {
		t.Errorf("Expected: %v, Actual: %v", true, actual)
	}
	// `incorrectHash` はPython3.6を使って生成した。生成したコードは以下の通り。
	// >>> import hashlib
	// >>> hashlib.sha512(b"incorrect password").hexdigest()
	incorrectHash := "00577d986cbb0d230f331d5e19f8ed4941f96b657363b5055dde47df24118f01701c45c85820e8633f150f7f9e00431d8fbe301120840321486caa4ead3793e7"
	actual = pw.Equal(hash, incorrectHash)
	if false != actual {
		t.Errorf("Expected: %v, Actual: %v", false, actual)
	}
	unexpectedHash := "password"
	actual = pw.Equal(hash, unexpectedHash)
	if false != actual {
		t.Errorf("Expected: %v, Actual: %v", false, actual)
	}
}

func TestPasswordHmacSha512_Hash(t *testing.T) {
	password := "password"
	secret := []byte("secretKey")
	pw := PasswordHmacSha512{secret}
	// `expected` はPython3.6を使って生成した。生成したコードは以下の通り。
	// >>> import hmac, hashlib
	// >>> hmac.new(b"secretKey", b"password", digestmod=hashlib.sha512).hexdigest()
	expected := "62475363a39382ec3be1d3001406a1b298faf688471ee78fbbb45ebe097def9601e585465482d3ae16390b806572d6712b7150b004b5755e018ca2ce0e247535"
	actual := pw.Hash(password)
	if expected != actual {
		t.Errorf("Expected: %v, Actual: %v", expected, actual)
	}
}

func TestPasswordHmacSha512_Equal(t *testing.T) {
	secret := []byte("secretKey")
	pw := PasswordHmacSha512{secret}
	// `hash` はPython3.6を使って生成した。生成したコードは以下の通り。
	// >>> import hmac, hashlib
	// >>> hmac.new(b"secretKey", b"password", digestmod=hashlib.sha512).hexdigest()
	hash := "62475363a39382ec3be1d3001406a1b298faf688471ee78fbbb45ebe097def9601e585465482d3ae16390b806572d6712b7150b004b5755e018ca2ce0e247535"
	actual := pw.Equal(hash, hash)
	if true != actual {
		t.Errorf("Expected: %v, Actual: %v", true, actual)
	}
	// `incorrectHash` はPython3.6を使って生成した。生成したコードは以下の通り。
	// >>> import hmac, hashlib
	// >>> hmac.new(b"secretKey", b"incorrect password", digestmod=hashlib.sha512).hexdigest()
	incorrectHash := "fe4f379747bd92c40ec3cafac68503cd4026a3193fc5a553b9be1eac2cedf0cec3c18621beeda2e086e4f60f8db10c05fc54ce019176386a9e800224b5026d2b"
	actual = pw.Equal(hash, incorrectHash)
	if false != actual {
		t.Errorf("Expected: %v, Actual: %v", false, actual)
	}
	unexpectedHash := "password"
	actual = pw.Equal(hash, unexpectedHash)
	if false != actual {
		t.Errorf("Expected: %v, Actual: %v", false, actual)
	}
}
