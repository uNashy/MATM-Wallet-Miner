package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"
	"os"
	"os/exec"
	"regexp"
	"strconv"


	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("")
	fmt.Println("                  ███▄ ▄███▓ ▄▄▄     ▄▄▄█████▓ ███▄ ▄███▓    ███▄ ▄███▓ ██▓ ███▄    █ ▓█████  ██▀███ ")
	fmt.Println("                 ▓██▒▀█▀ ██▒▒████▄   ▓  ██▒ ▓▒▓██▒▀█▀ ██▒   ▓██▒▀█▀ ██▒▓██▒ ██ ▀█   █ ▓█   ▀ ▓██ ▒ ██▒")
	fmt.Println("                 ▓██    ▓██░▒██  ▀█▄ ▒ ▓██░ ▒░▓██    ▓██░   ▓██    ▓██░▒██▒▓██  ▀█ ██▒▒███   ▓██ ░▄█ ▒")
	fmt.Println("                 ▒██    ▒██ ░██▄▄▄▄██░ ▓██▓ ░ ▒██    ▒██    ▒██    ▒██ ░██░▓██▒  ▐▌██▒▒▓█  ▄ ▒██▀▀█▄ ")
	fmt.Println("                 ▒██▒   ░██▒ ▓█   ▓██▒ ▒██▒ ░ ▒██▒   ░██▒   ▒██▒   ░██▒░██░▒██░   ▓██░░▒████▒░██▓ ▒██▒")
	fmt.Println("                 ░ ▒░   ░  ░ ▒▒   ▓▒█░ ▒ ░░   ░ ▒░   ░  ░   ░ ▒░   ░  ░░▓  ░ ▒░   ▒ ▒ ░░ ▒░ ░░ ▒▓ ░▒▓░")
	fmt.Println("                 ░  ░      ░  ▒   ▒▒ ░   ░    ░  ░      ░   ░  ░      ░ ▒ ░░ ░░   ░ ▒░ ░ ░  ░  ░▒ ░ ▒░")
	fmt.Println("                 ░      ░     ░   ▒    ░      ░      ░      ░      ░    ▒ ░   ░   ░ ░    ░     ░░   ░")
	fmt.Println("                        ░         ░  ░               ░             ░    ░           ░    ░  ░   ░     ")
}

func setTitle(title string) {
	cmd := exec.Command("cmd", "/C", "title", title)
	cmd.Run()
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

				fmt.Println("\n CPU: " + cpu[0].ModelName)
				fmt.Println(" Virtual Cores:", cores)

				fmt.Println(fmt.Sprint(" RAM size: ", ram, "mb"))

				pcsXcore := 10

				if cores > 6 {
					if ram > 8144 {
						pcsXcore = 33
					} else{
						pcsXcore = 25
					}
				} else {
					if ram > 8144 {
						pcsXcore = 20
					} else{
						pcsXcore = 10
					}
				}
				
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
			fmt.Println(" [!] Invalid input")
		}
	}
}

func getWallet(i int) {
	for {
		privateKey, _ := crypto.GenerateKey()
		privateKeyBytes := crypto.FromECDSA(privateKey)
		key := hexutil.Encode(privateKeyBytes)
		publicKey := privateKey.Public()
		publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		client, _ := ethclient.Dial("https://mainnet.infura.io/v3/555a4b9f11044e26a7d14afe70a45934")
		account := common.HexToAddress(address)
		bal, _ := client.BalanceAt(context.Background(), account, nil)
		balance, _ := strconv.Atoi(fmt.Sprint(bal))
		title := fmt.Sprint("MATM Wallet Miner | Balance ", balance, " ETH  | Cracking ", key)

		if balance > 0 {
			logHits(key, address, balance)
		}
		setTitle(title)
	}
}

func logHits(key string, address string, bal int) {
	fmt.Println("[+] New hit > " + key + " | Balance ", bal)
	f, _ := os.OpenFile("hits.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.Write([]byte(fmt.Sprint("NEW HIT [ ", time.Now(), " ] Private Key: ", key, " | Address: ", address, " | Balance: ", bal, " ETH\n")))
}

func main() {
	clear()
	fmt.Printf("\n                                                      Welcome back\n")
	var pcs int = getProcesses()

	clear()
	fmt.Println("\n GoRutines ammount:", pcs)

	for i := 0; i < pcs+1; i++ {
		go getWallet(i)
	}
	fmt.Println("\n Threads spawning completed\n Sleeping Main routine...\n\n")
	for { }
}
