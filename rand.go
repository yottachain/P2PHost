package host

import (
	"crypto/rand"

	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

type pr struct {
	ID   peer.ID
	Priv ci.PrivKey
	pub  ci.PubKey
}

// RandomPeer 生成随机peerid 密钥对
func RandomPeer() (*pr, error) {
	// First, select a source of entropy. We're using the stdlib's crypto reader here
	src := rand.Reader

	// Now create a 2048 bit RSA key using that
	priv, pub, err := ci.GenerateKeyPairWithReader(ci.Secp256k1, 256, src)
	if err != nil {
		return nil, err
	}

	// Now that we have a keypair, lets create our identity from it
	pid, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return nil, err
	}
	return &pr{pid, priv, pub}, nil
}
