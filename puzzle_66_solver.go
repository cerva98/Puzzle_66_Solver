package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
)

func main() {
	// Define o intervalo hexadecimal
	minHex := "20000000000000000"
	maxHex := "3ffffffffffffffff"

	// Converte os limites para big.Int
	min, success := new(big.Int).SetString(minHex, 16)
	if !success {
		fmt.Println("Erro ao converter o valor mínimo hexadecimal.")
		return
	}

	max, success := new(big.Int).SetString(maxHex, 16)
	if !success {
		fmt.Println("Erro ao converter o valor máximo hexadecimal.")
		return
	}

	// Alvo hash160
	targetHash := "20d45a6a762535700ce9e0b216e31994335db8a5"

	// Configuração de multiprocessamento
	numWorkers := runtime.NumCPU()
	fmt.Printf("Número de CPUs disponíveis: %d\n", numWorkers)

	// WaitGroup para sincronização
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Mutex para controlar acesso seguro à variável de contagem
	var mu sync.Mutex
	var checkCount int

	// Variáveis para medir a velocidade
	var lastCheckCount int
	var lastCheckTime time.Time

	// Função para processar chaves privadas em paralelo
	processPrivateKeys := func(workerID int) {
		defer wg.Done()

		for {
			// Gera uma chave privada inicial aleatória
			startPrivateKey, err := generateRandomPrivateKey(min, max)
			if err != nil {
				fmt.Println("Erro ao gerar número aleatório:", err)
				return
			}

			// Preenche a chave inicial com zeros à esquerda para garantir que tenha 64 caracteres
			startPrivateKeyHex := fmt.Sprintf("%064x", startPrivateKey)
			fmt.Printf("Worker %d iniciou com a chave: %s\n", workerID, startPrivateKeyHex)

			for i := 0; i < 10000000; i++ {
				// Incrementa a chave privada
				privateKey := new(big.Int).Add(startPrivateKey, big.NewInt(int64(i)))

				// Preenche a chave com zeros à esquerda para garantir que tenha 64 caracteres
				privateKeyHex := fmt.Sprintf("%064x", privateKey)

				// Deriva a chave pública comprimida da chave privada
				publicKeyCompressed, err := derivePublicKeyCompressed(privateKeyHex)
				if err != nil {
					fmt.Println("Erro ao derivar a chave pública:", err)
					return
				}

				// Converte a chave pública comprimida em hash160
				hash160 := hash160(publicKeyCompressed)

				// Incrementa a contagem de verificações de forma segura
				mu.Lock()
				checkCount++
				if checkCount%100000 == 0 {
					now := time.Now()
					elapsed := now.Sub(lastCheckTime).Seconds()
					fmt.Printf("Chaves geradas: %d, Velocidade: %.2f chaves/s\n", checkCount, float64(checkCount-lastCheckCount)/elapsed)
					lastCheckCount = checkCount
					lastCheckTime = now
				}
				mu.Unlock()

				// Verifica se o hash160 gerado corresponde ao alvo
				if hash160 == targetHash {
					// Calcula o tempo total de execução
					duration := time.Since(lastCheckTime)
					fmt.Printf("Chave privada correspondente encontrada pelo worker %d:\n", workerID)
					fmt.Printf("Chave privada: %s\n", privateKeyHex)
					fmt.Printf("Hash160 correspondente: %s\n", hash160)
					fmt.Printf("Alvo hash160: %s\n", targetHash)
					fmt.Printf("Tempo de execução: %s\n", duration)
					savePrivateKey(privateKeyHex)
					return
				}
			}
		}
	}

	// Inicia múltiplos workers
	for i := 0; i < numWorkers; i++ {
		go processPrivateKeys(i + 1)
	}

	// Marca o tempo de início
	lastCheckTime = time.Now()

	// Espera todos os workers terminarem
	wg.Wait()
}

func generateRandomPrivateKey(min, max *big.Int) (*big.Int, error) {
	privateKey, err := rand.Int(rand.Reader, new(big.Int).Sub(max, min))
	if err != nil {
		return nil, err
	}
	privateKey.Add(privateKey, min)
	return privateKey, nil
}

func derivePublicKeyCompressed(privateKeyHex string) ([]byte, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	return privateKey.PubKey().SerializeCompressed(), nil
}

func hash160(data []byte) string {
	return hex.EncodeToString(btcutil.Hash160(data))
}

func savePrivateKey(privateKeyHex string) {
	file, err := os.Create("matching_private_key.txt")
	if err != nil {
		log.Fatal("Erro ao criar arquivo:", err)
	}
	defer file.Close()

	_, err = file.WriteString(privateKeyHex)
	if err != nil {
		log.Fatal("Erro ao escrever no arquivo:", err)
	}

	fmt.Println("Chave privada salva em matching_private_key.txt")
}