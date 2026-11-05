package collector

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// CollectSSH executes one read-only diagnostic command over a verified SSH connection.
func CollectSSH(ctx context.Context, options SSHOptions) ([]string, error) {
	if options.Host == "" || options.User == "" || options.KeyPath == "" || options.KnownHosts == "" {
		return nil, fmt.Errorf("ssh-host, ssh-user, ssh-key and ssh-known-hosts are required")
	}
	if err := ValidatePlatform(options.Platform); err != nil {
		return nil, err
	}
	command := strings.TrimSpace(options.Command)
	if command == "" {
		// Les profils fournissent des commandes sûres tout en autorisant une commande explicite.
		// Profiles provide safe defaults while allowing an explicit command.
		command = defaultCommand(options.Platform)
	}

	signer, err := loadPrivateKey(options.KeyPath)
	if err != nil {
		return nil, err
	}
	hostKeyCallback, err := knownhosts.New(options.KnownHosts)
	if err != nil {
		return nil, fmt.Errorf("load known_hosts: %w", err)
	}

	clientConfig := &ssh.ClientConfig{
		User:            options.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hostKeyCallback,
		Timeout:         10 * time.Second,
	}

	dialer := net.Dialer{Timeout: 10 * time.Second}
	// La vérification de clé évite une connexion silencieuse à un équipement inattendu.
	// Host-key verification prevents silently connecting to an unexpected device.
	rawConnection, err := dialer.DialContext(ctx, "tcp", options.Host)
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", options.Host, err)
	}
	defer rawConnection.Close()

	connection, channels, requests, err := ssh.NewClientConn(rawConnection, options.Host, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh handshake: %w", err)
	}
	client := ssh.NewClient(connection, channels, requests)
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("open ssh session: %w", err)
	}
	defer session.Close()

	// La commande est exécutée une fois et sa sortie devient le snapshot normalisé.
	// The command runs once and its output becomes the normalized snapshot.
	output, err := session.CombinedOutput(command)
	if err != nil {
		return nil, fmt.Errorf("run remote command: %w: %s", err, strings.TrimSpace(string(output)))
	}

	var lines []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func defaultCommand(platform Platform) string {
	switch platform {
	case PlatformCiscoIOS:
		return "show ip interface brief; show ip route; show running-config | include ^hostname|^ip route|^ntp|^username|^line vty"
	default:
		return "ip -brief address; ip route"
	}
}

func loadPrivateKey(path string) (ssh.Signer, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read ssh private key: %w", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse ssh private key: %w", err)
	}
	return signer, nil
}
