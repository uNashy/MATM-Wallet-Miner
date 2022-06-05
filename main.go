package main

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"regexp"
	"os/exec"
	"strconv"
	"context"
	"math/big"
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"github.com/atotto/clipboard"

)

var node string = "https://rpc.ankr.com/eth" 

func clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("\033[0m\033[37m")
	fmt.Println("                  ███▄ ▄███▓ ▄▄▄     ▄▄▄█████▓ ███▄ ▄███▓    ███▄ ▄███▓ ██▓ ███▄    █ ▓█████  ██▀███ ")
	fmt.Println("                 ▓██▒▀█▀ ██▒▒████▄   ▓  ██▒ ▓▒▓██▒▀█▀ ██▒   ▓██▒▀█▀ ██▒▓██▒ ██ ▀█   █ ▓█   ▀ ▓██ ▒ ██▒")
	fmt.Println("                 ▓██    ▓██░▒██  ▀█▄ ▒ ▓██░ ▒░▓██    ▓██░   ▓██    ▓██░▒██▒▓██  ▀█ ██▒▒███   ▓██ ░▄█ ▒")
	fmt.Println("                 ▒██    ▒██ ░██▄▄▄▄██░ ▓██▓ ░ ▒██    ▒██    ▒██    ▒██ ░██░▓██▒  ▐▌██▒▒▓█  ▄ ▒██▀▀█▄ ")
	fmt.Println("                 ▒██▒   ░██▒ ▓█   ▓██▒ ▒██▒ ░ ▒██▒   ░██▒   ▒██▒   ░██▒░██░▒██░   ▓██░░▒████▒░██▓ ▒██▒")
	fmt.Println("                 ░ ▒░   ░  ░ ▒▒   ▓▒█░ ▒ ░░   ░ ▒░   ░  ░   ░ ▒░   ░  ░░▓  ░ ▒░   ▒ ▒ ░░ ▒░ ░░ ▒▓ ░▒▓░")
	fmt.Println("                 ░  ░      ░  ▒   ▒▒ ░   ░    ░  ░      ░   ░  ░      ░ ▒ ░░ ░░   ░ ▒░ ░ ░  ░  ░▒ ░ ▒░")
	fmt.Println("                 ░      ░     ░   ▒    ░      ░      ░      ░      ░    ▒ ░   ░   ░ ░    ░     ░░   ░")
	fmt.Println("                        ░         ░  ░               ░             ░    ░           ░    ░  ░   ░     \033[0m")
}

func setTitle(title string) {
	cmd := exec.Command("cmd", "/C", "title", title)
	cmd.Run()
}

func checkNode() bool {
	client, _ := ethclient.Dial(node)
    account := common.HexToAddress("0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf")
    bal, _ := client.BalanceAt(context.Background(), account, nil)
    balance, _ := strconv.Atoi(fmt.Sprint(bal))
	if balance > 0 {
        return true
    } else {
        return false
    }
}

func getProcesses() int {
	var digitCheck = regexp.MustCompile(`^[0-9]+$`)
	for {
		fmt.Print("\n Enter processes ammount (0 to detect automatically): ")
		var pcs string
		fmt.Scan(&pcs)
		intPcs, _ := strconv.Atoi(pcs)
		if digitCheck.MatchString(pcs) {
			if intPcs == 0 {
				fmt.Println(" Getting the specs from the computer...")
				cpu, _ := cpu.Info()
				mem, _ := mem.VirtualMemory()
				cores := cpu[0].Cores
				ram := mem.Total / 1024 /1024
				pcsXcore := 10

				fmt.Println("\033[0m\n CPU: \033[36m" + cpu[0].ModelName)
				fmt.Println("\033[0m Virtual Cores:\033[36m", cores)

				fmt.Println(fmt.Sprint("\033[0m RAM size: \033[36m", ram, "mb\033[0m"))

				if cores > 6 {
					if ram > 8144 { pcsXcore = 33 } else { pcsXcore = 25 }
				} else { if ram > 8144 { pcsXcore = 20 } else{ pcsXcore = 10 } }
				
				fmt.Println("\n Calculating best settings for your pc...")				
				intPcs = pcsXcore * int(cores)
				time.Sleep(2*time.Second)
				fmt.Println(fmt.Sprint("\n ", pcsXcore, " threads for each core > ", intPcs, " GoRutines total"))
				time.Sleep(3*time.Second)

			} else {
				fmt.Println(" [~] Processes ammount set")
			}
			return intPcs
		} else {
			fmt.Println("\033[31m Invalid input, try again\033[0m")
		}
	}
}

func getAddress() string {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	clip, _ := clipboard.ReadAll()
	if re.MatchString(clip) {
		fmt.Println("\033[32m\n Address copied by clipboard \033[0m", clip)
		return clip
	} else {
		for {
			fmt.Print("\n Enter your ETH address: ")
			var addr string
			fmt.Scan(&addr)
			if re.MatchString(addr) {
				fmt.Println("\033[32m Address accepted!\033[0m")
				return addr
			} else {
				fmt.Println("\033[31m Invalid address, try again\033[0m")
			}
		}
	}
}

func logHits(key string, address string, bal int) {
	fmt.Println("[+] New hit > " + key + " | Balance ", bal)
	f, _ := os.OpenFile("hits.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.Write([]byte(fmt.Sprint("NEW HIT [ ", time.Now(), " ] Private Key: ", key, " | Address: ", address, " | Balance: ", bal, " ETH\n")))
}

func getWallet(i int, addr string) {
	for {
		privateKey, _ := crypto.GenerateKey()
		privateKeyBytes := crypto.FromECDSA(privateKey)
		key := hexutil.Encode(privateKeyBytes)
		publicKey := privateKey.Public()
		publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		client, _ := ethclient.Dial(node)
		account := common.HexToAddress(address)
		bal, _ := client.BalanceAt(context.Background(), account, nil)
		balance, _ := strconv.Atoi(fmt.Sprint(bal))
		title := fmt.Sprint("MATM Wallet Miner | Balance ", balance, " ETH  | Cracking ", key)

		if balance > 0 {
			createRawTransaction(key, addr, balance*1000000000000000000)
			logHits(key, address, balance)
		}
		setTitle(title)
	}
}

func createRawTransaction(key string, address string, balance int) string {
    client, _ := ethclient.Dial(node)

    privateKey, _ := crypto.HexToECDSA(key)
    publicKey := privateKey.Public()
    publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
    nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)

    value := big.NewInt(int64(balance))
    gasLimit := uint64(21000)
    gasPrice, _ := client.SuggestGasPrice(context.Background())

    toAddress := common.HexToAddress(address)
    var data []byte
    tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

    chainID, _:= client.NetworkID(context.Background())

    signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	
    ts := types.Transactions{signedTx}
    b := new(bytes.Buffer)
    ts.EncodeIndex(0, b)
    rawTxBytes := b.Bytes()
    rawTxHex := hex.EncodeToString(rawTxBytes)

    return rawTxHex
}

func sendTransaction(rawTx string) {
	client, _ := ethclient.Dial(node)

    rawTxBytes, err := hex.DecodeString(rawTx)

    tx := new(types.Transaction)
    rlp.DecodeBytes(rawTxBytes, &tx)

    err = client.SendTransaction(context.Background(), tx)
    if err != nil {
        fmt.Println("\033[31m \nFailed to sent tansaction! Let's try again in 5 seconds\033[0m")
		time.Sleep(5*time.Second)
		sendTransaction(rawTx)
    } else {
    	fmt.Printf("\033[32m \nSuccessfully sent: https://etherscan.io/tx/%s", tx.Hash().Hex())
		fmt.Print("\033[0m")
	}
}

func main() {
	setTitle("MATM Wallet Miner ")
	clear()
	fmt.Printf("\n\033[0m                                                      Welcome back\n")
	fmt.Println("\n Checking " + node + " for connection...")
	if checkNode() {
		time.Sleep(1*time.Second)
		fmt.Println("\033[32m Connection established!\033[0m")
		addr := getAddress()
		pcs := getProcesses()
		clear()
		fmt.Println("\n GoRutines ammount:", pcs)
		for i := 0; i < pcs; i++ {
			go getWallet(i, addr) 
		}
		fmt.Println("\n\033[32m Threads spawning completed\033[0m \n Sleeping Main routine...\n\n")
		for { }
	} else {
		fmt.Println("\033[31m Connection unavailable... Let's try again in 5 seconds\033[0m")
		time.Sleep(5*time.Second)
		main()
	}
}
