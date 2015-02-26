package runner

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
)

// createUID is a helper function that will
// create a random, unique identifier.
func createUID() string {
	c := sha1.New()
	r := createRandom()
	io.WriteString(c, string(r))
	s := fmt.Sprintf("%x", c.Sum(nil))
	return s[0:10]
}

// createRandom creates a random block of bytes
// that we can use to generate unique identifiers.
func createRandom() []byte {
	k := make([]byte, sha1.BlockSize)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}

/*
func BashToBase64(script string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(script))
	return fmt.Sprintf("echo %s | base64 --decode | bash", encoded)
}

func CmdToBase64(cmd string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(cmd))
	return fmt.Sprintf("echo %s | base64 --decode", encoded)
}

func ConfigToBash(conf config.Config) string {
	var buf bytes.Buffer
	buf.WriteString("#!/bin/bash")
	for _, cmd := range conf.Script {
		echo := CmdToBase64(cmd)
		buf.WriteString(echo)
		buf.WriteString(cmd)
	}
	buf.WriteString("exit 0")
	return buf.String()
}

func CloneToBash(clone cloner.Cloner) string {
	var buf bytes.Buffer
	buf.WriteString("#!/bin/bash")
	buf.WriteString(NetrcToBash(clone.GetNetrc()))
	buf.WriteString(KeypairToBash(clone.GetKeypair()))
	for _, cmd := range clone.GetCmds() {
		echo := CmdToBase64(cmd)
		buf.WriteString(echo)
		buf.WriteString(cmd)
	}
	buf.WriteString("exit 0")
	return buf.String()
}

func NetrcToBash(netrc *cloner.Netrc) string {
	var buf bytes.Buffer
	if netrc == nil || netrc.Machine == "" {
		return buf.String()
	}
	return buf.String()
}

func KeypairToBash(keypair *cloner.Keypair) string {
	var buf bytes.Buffer
	if keypair == nil || keypair.Private == "" {
		return buf.String()
	}
	buf.WriteString("mkdir -p /root/.ssh")
	buf.WriteString("chmod 700 /root/.ssh")
	buf.WriteString(fmt.Sprintf("echo %s > ~/.ssh/id_rsa", keypair.Private))
	buf.WriteString("chmod 600 /root/.ssh/id_rsa")
	buf.WriteString("echo 'StrictHostKeyChecking no' > /root/.ssh/config")
	return buf.String()
}
*/
