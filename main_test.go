package main

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestGetKeyHexValue(t *testing.T) {
	var tests = []struct {
		privKeyHex, pubAddrHex string
	}{
		{"0x39bde15be9762ee36368564cb409b6a34a754dd7d0c52c511d18d5b734a1e5ad", "0x3874bB817D9e976155525DeEB6b86f74B08f9625"},
		{"0x510ebfdeed563fc6ab15dd09c5c4f491fcd4b97934f7ca04618c63c28a727a10", "0xF23598b829197b3967C9F2CABFf9B34c813eb8BE"},
		{"0x3db4ab827f283cac19a9ac65bab46344d9e75d1389c7196c8378f9def7a2533d", "0xBa29Ce183E1955748f0F02a0AF648085e1726A83"},
		{"0x9163e07dd770c4c7f92fcd4ef4ec5219db0ab3bf62ff81ad18bc3bb9d40fc7b5", "0x50d4924504329C1D928745953209FBE637e156B4"},
		{"0x2d1d878f639a396f47cc8502a95d7c57a1cd0aa3df5fb3eb2b42b9fa91636f0d", "0x0395C7Ce7AbC8493e373f65938E4E18cf5566611"},
		{"0xce2f6ca8df23e2bb0b4788804146921203733905c1b338884b3c1fe229961d2c", "0x0Bf0DAAfeE42Cde40a04FC3e91f6145117B3020a"},
		{"0x413dfe62e67f7ff2e430cb9d7a22a1d0656436672203015d31059ac830e988e7", "0xa615F8946c2542d7B1B2DF62ce39325A7DA06C36"},
		{"0x5adb81f15bc42bb58ebec114070d031ac900e65446e6f70f87c3535c34c1d3c5", "0xd0ca8d0f6517Ec56f278FbcC172fE48DA868cf01"},
		{"0x23e93a463da5401307b7afbcd88e062ed1de0a7af8b3abd5b9ba41f66dabe7e2", "0xA310aC70cca518264Ac67C1F18246C9C733b2692"},
		{"0xd1c8069fae5e05ec816ab4e54dc328dccd08453f352af9325ea2d23bf4e68352", "0x9328858E56c95c59d301C38E4237f4482684DF59"},
		{"0x2506bf1df301e3758fbdbe8db5eae21dc88ad520d85d4d2353ae12c139b21e34", "0xEd2cb22E77ed67A814143c0C7Ee8800a9Ee592Dc"},
		{"0x5bcd3a007faf8268992659050c421b31bf4a363619cd031f164b7746eee811c5", "0x1E967394657fF0D1938ec3A5202899d1637AfA6F"},
		{"0xe017e5e155603b93871079616d3bdb194eef113950f017f74445d873a4aa62b8", "0x148E31A7E9a899789072826e071013E2DFa282ce"},
	}

	for _, tc := range tests {
		t.Run(tc.pubAddrHex, func(t *testing.T) {
			key, _ := crypto.HexToECDSA((tc.privKeyHex)[2:])
			gotPubKeyHex := crypto.PubkeyToAddress(key.PublicKey).Hex()

			if gotPubKeyHex != tc.pubAddrHex {
				t.Errorf("got %s, want %s", gotPubKeyHex, tc.pubAddrHex)
			}
		})
	}
}
