package cli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/beeploop/aes-encrypt/encrypt"
)

type MODE string

const (
	ENCRYPT MODE = "encrypt"
	DECRYPT MODE = "decrypt"

	KEY_LEN          int = 32
	DEFAULT_FILENAME     = "output"
)

type cli struct {
	inputFile  string
	key        []byte
	mode       MODE
	outputFile string
	wd         string
}

func New() (*cli, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cli := &cli{
		mode:       ENCRYPT,
		outputFile: filepath.Join(wd, DEFAULT_FILENAME), // File extension automatically placed on save
		wd:         wd,
	}

	return cli, nil
}

func (c *cli) SetKey(keySource string) error {
	ext := filepath.Ext(keySource)

	if ext == ".txt" {
		if key, err := c.readContent(keySource); err != nil {
			return err
		} else {
			if len(key) != KEY_LEN {
				return errors.New("Invalid key. Key must be 32 character string for AES256")
			}

			c.key = key
			return nil
		}
	}

	key := []byte(keySource)
	if len(key) != KEY_LEN {
		return errors.New("Invalid key. Key must be 32 character string for AES256")
	}

	c.key = key
	return nil
}

func (c *cli) SetMode(mode MODE) {
	c.mode = mode
}

func (c *cli) SetInputFile(inputFile string) error {
	if inputFile == "" {
		return errors.New("empty file")
	}

	path, err := filepath.Abs(inputFile)
	if err != nil {
		return err
	}

	c.inputFile = path
	return nil
}

func (c *cli) SetOutput(out string) error {
	if out == "." {
		c.outputFile = filepath.Join(c.wd, DEFAULT_FILENAME)
		return nil
	}

	if out == "" {
		fmt.Println("Empty output location, saving to default location")
		return nil
	}

	abs, err := filepath.Abs(out)
	if err != nil {
		return err
	}

	c.outputFile = abs
	return nil
}

func (c *cli) Run() error {
	switch c.mode {
	case ENCRYPT:
		encrypted, err := c.encrypt()
		if err != nil {
			return err
		}

		return c.saveFile(encrypted)

	case DECRYPT:
		decrypted, err := c.decrypt()
		if err != nil {
			return err
		}

		return c.saveFile(decrypted)

	default:
		return errors.New("Unsupported mode")
	}
}

func (c *cli) encrypt() ([]byte, error) {
	input, err := c.readFile()
	if err != nil {
		return nil, err
	}

	encryptor, err := encrypt.New(c.key)
	if err != nil {
		return nil, err
	}

	encrypted, err := encryptor.Encrypt(input)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

func (c *cli) decrypt() ([]byte, error) {
	input, err := c.readFile()
	if err != nil {
		return nil, err
	}

	decryptor, err := encrypt.New(c.key)
	if err != nil {
		return nil, err
	}

	decrypted, err := decryptor.Decrypt(input)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func (c *cli) saveFile(data []byte) error {
	ext := filepath.Ext(c.inputFile)
	file := c.outputFile + ext

	if err := os.WriteFile(file, data, 0777); err != nil {
		return err
	}

	return nil
}

func (c *cli) readFile() ([]byte, error) {
	b, err := os.ReadFile(c.inputFile)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *cli) readContent(source string) ([]byte, error) {
	key, err := os.ReadFile(filepath.Join(c.wd, source))
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewReader(bytes.NewReader(key))
	line, _, err := scanner.ReadLine()
	if err != nil {
		return nil, err
	}

	return []byte(line), nil
}
