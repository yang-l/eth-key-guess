package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/exp/slices"
)

func main() {
	multiBalances()
}

type ethereumAccount struct {
	privateKeyHex string
	publicAddrHex string
}

func multiBalances() {
	token := os.Getenv("ETHERSCAN_TOKEN")
	if token == "" {
		log.Fatal("ETHERSCAN_TOKEN is not set")
	}

	limiter := time.Tick(275 * time.Millisecond)
	count := 0
	for true {
		<-limiter

		count++
		fmt.Printf("\033[2K\rCount %d", count)

		accountSet, keys := genRandKeySet()
		bals := getBalances(keys)

		idx := slices.IndexFunc(bals, func(r Result) bool { return r.Balance != "0" })
		if idx != -1 {
			text := fmt.Sprintf("%s\nFound address [%s] with balance [%s]\n", time.Now().Format(time.RFC3339), bals[idx].Account, bals[idx].Balance)
			fmt.Printf(text)
			writeToFile(text)

			writeToFile(fmt.Sprintf("Keyset (raw) - %s\n", accountSet))
			writeToFile(fmt.Sprintf("Result (raw) - %s\n\n", bals))
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

func genRandKeySet() ([]ethereumAccount, string) {
	var accountSet []ethereumAccount
	var pubKeys []string

	for i := 0; i < 20; i++ {
		privKeyHex, pubAddrHex := GetKeyHexValue(genRandomKey())
		accountSet = append(accountSet, ethereumAccount{privKeyHex, pubAddrHex})
		pubKeys = append(pubKeys, pubAddrHex)
	}
	keys := strings.Join(pubKeys[:], ",")

	return accountSet, keys
}

func writeToFile(output string) {
	f, e := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if e != nil {
		log.Fatal(e)
	}
	defer f.Close()

	if _, e = f.WriteString(output); e != nil {
		log.Fatal(e)
	}
}

type Response struct {
	Status  string
	Message string
	Result  []Result
}

type Result struct {
	Account string
	Balance string
}

func getResponse(addr string) (string, string, []Result) {
	token := os.Getenv("ETHERSCAN_TOKEN")
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=balancemulti&address=%s&tag=latest&apikey=%s", addr, token)

	resp, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer resp.Body.Close()

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Fatal(e)
	}

	var response Response
	json.Unmarshal(body, &response)

	return response.Status, response.Message, response.Result
}

func getBalances(addr string) []Result {
	stat, msg, res := getResponse(addr)

	if (stat != "1") || (msg != "OK") {
		log.Fatalf("Failed to get response from etherscan. Status: %s, Message: %s, Result: %s", stat, msg, res)
	}

	return res
}
