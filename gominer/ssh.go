package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

var tailLogOutput string

func sshClient() {
	hostname := "116.202.104.188"
	username := "root"
	privateKeyPath := "/Users/stepansalikov/.ssh/id_rsa" // Укажите путь к вашему ключу
	logFile := "/app/log-proxy.txt"

	// Читаем приватный ключ
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}

	// Запрашиваем пароль
	fmt.Print("Enter passphrase for private key: ")
	passphrase, err := term.ReadPassword(0)
	if err != nil {
		log.Fatalf("Unable to read passphrase: %v", err)
	}

	// Создаем ssh.Signer с использованием пароля
	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, passphrase)
	if err != nil {
		log.Fatalf("Unable to parse private key with passphrase: %v", err)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", hostname), config)
	if err != nil {
		log.Fatalf("failed to connect to SSH server: %v", err)
	}
	defer client.Close()

	output, err := tailLog(client, logFile)
	if err != nil {
		log.Fatalf("failed to tail log: %v", err)
	}
	fmt.Printf("Log output: %s\n", output)

}

func tailLog(client *ssh.Client, logFile string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Выполнение команды tail -f
	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := session.Start(fmt.Sprintf("tail -f %s", logFile)); err != nil {
		return "", fmt.Errorf("failed to start command: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	// Объявляем переменную для хранения вывода

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "Accepted") {
			// Меняем tailLogOutput если решение принято
			tailLogOutput = "Accepted"
			fmt.Printf("Accepted: %s\n", line)
		} else if strings.Contains(line, `"Above target"`) {
			// Меняем tailLogOutput если решение с ошибкой
			tailLogOutput = "Above target"
			fmt.Printf("Above Target: %s\n", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading output: %w", err)
	}

	if err := session.Wait(); err != nil {
		return "", fmt.Errorf("session wait error: %w", err)
	}

	return tailLogOutput, nil
}
