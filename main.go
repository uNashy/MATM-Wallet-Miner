package main

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"regexp"
	"errors"
	"strings"
	"os/exec"
	"strconv"
	"context"
	"net/http"
	"math/big"
	"io/ioutil"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	
	"github.com/atotto/clipboard"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	
	"github.com/admin100/util/console"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var currentVersion = 1.5
var updatesUrl = "https://raw.githubusercontent.com/uNashy/MATM-Wallet-Miner/main/VERSION_INFO"

var node string = "https://rpc.ankr.com/eth" 
var configFileName string = "config.json"

var numCheck = regexp.MustCompile(`^[0-9]+$`)
var addressCheck = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

var highMode string = "\n \033[31mHigh performance mode\033[0m"
var mediumMode string = "\n \033[33mMedium performance mode\033[0m"
var lowMode string = "\n \033[32mLow performance mode\033[0m"
var customMode string = "\n \033[33mCustom threads mode\033[0m"

var checkAddress string  = "0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf"

var gens int = 0

type UserConfig struct {
    Address string
    Threads int
	Mode string
	BotToken string
}

type Spinner struct {
    message string
    i int
}

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

func checkNode() bool {
	client, _ := ethclient.Dial(node)
    account := common.HexToAddress(checkAddress)
    bal, _ := client.BalanceAt(context.Background(), account, nil)
    balance, _ := strconv.Atoi(fmt.Sprint(bal))
	if balance > 0 {
        return true
    } else {
        return false
    }
}

func checkForUpdates() {
	fmt.Println(" Checking for updates...")
	time.Sleep(1*time.Second)
	resp, _ := http.Get(updatesUrl)
	body, _ := ioutil.ReadAll(resp.Body)
	floatNum, _ := strconv.ParseFloat(strings.TrimSpace(string(body)), 64)  
	if floatNum > currentVersion { fmt.Println(fmt.Sprint("\033[36m Update available\033[0m > New version: ", floatNum)) } else { fmt.Println("\033[32m No updates available\033[0m") }

}

func logHits(key string, address string, bal int) {
	fmt.Println("[+] New hit > " + key + " | Balance ", bal)

	f, _ := os.OpenFile("hits.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.Write([]byte(fmt.Sprint("NEW HIT [ ", time.Now(), " ] Private Key: ", key, " | Address: ", address, " | Balance: ", bal, " ETH\n")))
}

func setPriorityProcess() {
	pid := os.Getpid()
	fmt.Println("\n Setting up process", pid ,"priority to High...")
	cmd := exec.Command("cmd", "/c", fmt.Sprint("wmic process where processid='", pid, "' CALL setpriority '256'"))
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("\033[31m Failed to set high priority\033[0m")
	} else {
		fmt.Println("\033[32m Process high priority set\033[0m")
	}
	
}

func getProcesses() (int, string) {
	for {
		fmt.Print("\n Enter processes ammount (0 to detect automatically): ")
		var pcs string
		fmt.Scan(&pcs)
		intPcs, _ := strconv.Atoi(pcs)
		mode := customMode
		if numCheck.MatchString(pcs) {
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
					if ram > 8144 { pcsXcore = 45; mode = highMode } else { pcsXcore = 39; mode = mediumMode}
				} else { if ram > 8144 { pcsXcore = 32; mode = mediumMode } else { pcsXcore = 21; mode = lowMode} }
				
				fmt.Println("\n Calculating best settings for your pc...")				
				time.Sleep(1*time.Second)
				fmt.Println(mode)
				intPcs = pcsXcore * int(cores)
				time.Sleep(2*time.Second)
				fmt.Println(fmt.Sprint(" ", pcsXcore, " threads for each core > ", intPcs, " GoRutines total"))
				time.Sleep(3*time.Second)

			} else {
				fmt.Println("\033[32m Processes ammount set\033[0m")
			}
			return intPcs, mode
		} else {
			fmt.Println("\033[31m Invalid input, try again\033[0m")
		}
	}
}

func getAddress() string {
	clip, _ := clipboard.ReadAll()
	if addressCheck.MatchString(clip) {
		fmt.Println("\033[32m\n Address copied by clipboard:\033[0m", clip)
		return clip
	} else {
		for {
			fmt.Print("\n Enter your ETH address: ")
			var addr string
			fmt.Scan(&addr)
			if addressCheck.MatchString(addr) {
				fmt.Println("\033[32m Address accepted!\033[0m")
				return addr
			} else {
				fmt.Println("\033[31m Invalid address, try again\033[0m")
			}
		}
	}
}

func getChoice() bool {
	for {
		var choice string
		fmt.Scan(&choice)
		if strings.Contains("y", strings.ToLower(choice)) {
			return true
		} else if strings.Contains("n", strings.ToLower(choice)) {
			return false
		} else {
			fmt.Println("\033[31m Invalid input, try again\033[0m")
		}
	}
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
		
		if balance > 0 {
			createRawTransaction(key, addr, balance*1000000000000000000)
			logHits(key, address, balance*1000000000000000000)
		}
		gens = gens + 1
		title := fmt.Sprint("MATM Wallet Miner | Balance ", balance*1000000000000000000, " ETH | Generated ", gens, " | Cracking ", key)
		console.SetConsoleTitle(title)
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

func checkMode(str string) bool { for _, v := range []string{highMode, mediumMode, lowMode, customMode} { if v == str { return true } }
    return false
}

func readConfig() (string, int, string){
    content, _ := ioutil.ReadFile(configFileName)

    var config UserConfig
    json.Unmarshal(content, &config)

    return config.Address, config.Threads, config.Mode
}

func writeConfig(address string, threads int, mode string) {
	data := UserConfig{
        Address: address,
        Threads: threads,
		Mode: mode,
    }
	file, _ := json.MarshalIndent(data, "", "  ")
	_ = ioutil.WriteFile(configFileName, file, 0644)
}

func newConfig() (string, int, string) {
	addr := getAddress()
	pcs, mode := getProcesses()
	writeConfig(addr, pcs, mode)
	
	return addr, pcs, mode
}

func main() {
	console.SetConsoleTitle(fmt.Sprint("MATM Wallet Miner | Version: ", currentVersion))
	clear()
	fmt.Printf("\n\033[0m                                                      Welcome back\n")
	checkForUpdates()
	fmt.Println("\n Checking " + node + " for connection...")
	addr, pcs, mode := "", 0, ""
	if checkNode() {
		time.Sleep(1*time.Second)
		fmt.Println("\033[32m Connection established!\033[0m")
		setPriorityProcess()
		_, err := os.OpenFile(configFileName, os.O_WRONLY, 0664)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("\033[31m\n config.json not found!\033[0m")
			addr, pcs, mode = newConfig()		
		} else {
			addr, pcs, mode = readConfig()
			if ! addressCheck.MatchString(addr) || ! numCheck.MatchString(strconv.Itoa(pcs)) || ! checkMode(mode) {
				fmt.Println("\033[31m\n config.json found but corrupt!\033[0m")
				addr, pcs, mode = newConfig()
			} else {
				fmt.Print("\033[32m\n Valid config.json found!\033[0m Want use it? [y/n] ")
				if ! getChoice() { addr, pcs, mode = newConfig() }
			}
		}
		clear()
		fmt.Println("\n GoRutines ammount:", pcs)
		fmt.Println(mode)
		
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
