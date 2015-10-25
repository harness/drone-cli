package main

import (
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-exec/yaml/secure"
	"github.com/drone/drone-go/drone"

	"github.com/square/go-jose"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

var SecureCmd = cli.Command{
	Name:  "secure",
	Usage: "creates a secure yaml file",
	Action: func(c *cli.Context) {
		handle(c, SecureYamlCmd)
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "in",
			Usage: "input path to the plaintext secret file",
			Value: ".drone.sec.yml",
		},
		cli.StringFlag{
			Name:  "out",
			Usage: "output path for the encrypted secret file",
			Value: ".drone.sec",
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "name of the repository",
		},
		cli.StringFlag{
			Name:  "yaml",
			Usage: "path to .drone.yml file",
			Value: ".drone.yml",
		},
	},
}

func SecureYamlCmd(c *cli.Context, client drone.Client) error {
	var (
		repo    = c.String("repo")
		inFile  = c.String("in")
		outFile = c.String("out")
		ymlFile = c.String("yaml")
	)

	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}

	keypair, err := client.RepoKey(owner, name)
	if err != nil {
		return err
	}

	key, err := toPublicKey(keypair.Public)
	if err != nil {
		return err
	}

	// read the .drone.sec.yml file (plain text)
	plaintext, err := ioutil.ReadFile(inFile)
	if err != nil {
		return err
	}

	// parse the .drone.sec.yml file
	sec := &secure.Secure{}
	err = yaml.Unmarshal(plaintext, sec)
	if err != nil {
		return err
	}

	// read the .drone.yml file and caclulate the
	// checksum. add to the .drone.sec.yml file.
	yml, err := ioutil.ReadFile(ymlFile)
	if err == nil {
		sec.Checksum = sha256sum(string(yml))
	}

	// re-marshal the .drone.sec.yml file since we've
	// added the checksum
	plaintext, err = yaml.Marshal(sec)
	if err != nil {
		return err
	}

	// encrypt the .drone.sec.yml file
	ciphertext, err := encrypt(plaintext, key)
	if err != nil {
		return err
	}

	// write the encrypted .drone.sec.yml file to .drone.sec
	return ioutil.WriteFile(outFile, []byte(ciphertext), 0664)
}

// toPublicKey parses a public key and returns an *rsa.PublicKey.
// credit to stackoverflow http://stackoverflow.com/q/31593329
func toPublicKey(key string) (*rsa.PublicKey, error) {
	raw := []byte(key)
	pub, _, _, _, err := ssh.ParseAuthorizedKey(raw)
	if err != nil {
		return nil, err
	}
	return reflect.ValueOf(pub).
		Convert(reflect.TypeOf(new(rsa.PublicKey))).Interface().(*rsa.PublicKey), nil
}

// encrypt encrypts a plaintext variable using JOSE with
// RSA_OAEP and A128GCM algorithms.
func encrypt(plaintext []byte, pubKey *rsa.PublicKey) (string, error) {
	var encrypted string

	// Creates a new encrypter using defaults
	encrypter, err := jose.NewEncrypter(jose.RSA_OAEP, jose.A128GCM, pubKey)
	if err != nil {
		return encrypted, err
	}
	// Encrypts the plaintext value and serializes
	// as a JOSE string.
	object, err := encrypter.Encrypt(plaintext)
	if err != nil {
		return encrypted, err
	}
	return object.CompactSerialize()
}

func sha256sum(in string) string {
	h := sha256.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}
