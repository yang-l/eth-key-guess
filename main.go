package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	token := os.Getenv("ETHERSCAN_TOKEN")
	if token == "" {
		log.Fatal("ETHERSCAN_TOKEN is not set")
	}

	limiter := time.Tick(275 * time.Millisecond)
	count := 0
	for true {
		<-limiter

		privKeyHex, pubAddrHex := GetKeyHexValue(genRandomKey())
		bal := getBalance(pubAddrHex)

		count++
		fmt.Printf("\033[2K\rCount %d", count)

		if bal != "0" {
			fmt.Printf("bal: %s, key: %s, addr: %s\n", bal, privKeyHex, pubAddrHex)
		}
	}
}

func genRandomKey() *ecdsa.PrivateKey {
	key, _ := crypto.GenerateKey()
	return key
}

func GetKeyHexValue(randKey *ecdsa.PrivateKey) (string, string) {
	privKeyHex := hexutil.Encode(crypto.FromECDSA(randKey))
	pubAddrHex := crypto.PubkeyToAddress(randKey.PublicKey).Hex()

	return privKeyHex, pubAddrHex
}

type Response struct {
	Status  string
	Message string
	Result  string
}

func getResponse(addr string) (string, string, string) {
	token := os.Getenv("ETHERSCAN_TOKEN")
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=balance&address=%s&tag=latest&apikey=%s", addr, token)

	resp, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer resp.Body.Close()

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Fatal(e)
	}

	var respponse Response
	json.Unmarshal(body, &respponse)

	return respponse.Status, respponse.Message, respponse.Result
}

func getBalance(addr string) string {
	stat, msg, res := getResponse(addr)

	if (stat != "1") || (msg != "OK") {
		log.Fatalf("Failed to get response from etherscan. Status: %s, Message: %s, Result: %s", stat, msg, res)
	}

	return res
}
