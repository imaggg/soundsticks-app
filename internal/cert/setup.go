package cert

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/imaggg/soundsticks/internal/config"
	"golang.org/x/crypto/pkcs12"
)

const aliceAssetPath = "assets/alice.p12"

func configDir() (string, error) {
	return config.Dir()
}

// Forget removes the saved certificate files, forcing re-setup on next launch.
func Forget() {
	d, _ := configDir()
	os.Remove(filepath.Join(d, "cert.pem"))
	os.Remove(filepath.Join(d, "key.pem"))
}

// IsSetupDone returns true when cert.pem + key.pem exist in the config dir.
func IsSetupDone() bool {
	d, err := configDir()
	if err != nil {
		return false
	}
	_, e1 := os.Stat(filepath.Join(d, "cert.pem"))
	_, e2 := os.Stat(filepath.Join(d, "key.pem"))
	return e1 == nil && e2 == nil
}

// LoadTLSCert returns the saved mTLS certificate.
func LoadTLSCert() (tls.Certificate, error) {
	d, err := configDir()
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.LoadX509KeyPair(
		filepath.Join(d, "cert.pem"),
		filepath.Join(d, "key.pem"),
	)
}

// ExtractFromAPK finds assets/alice.p12 inside an APK (which is a ZIP),
// auto-discovers its password, and saves PEM files.
func ExtractFromAPK(apkPath string) error {
	p12Data, err := fileFromZip(apkPath, aliceAssetPath)
	if err != nil {
		return fmt.Errorf("alice.p12 not found in APK: %w", err)
	}
	password, err := findP12Password(p12Data)
	if err != nil {
		return err
	}
	return decodeAndSave(p12Data, password)
}

// ExtractFromXAPK handles XAPK (ZIP containing an APK).
func ExtractFromXAPK(xapkPath string) error {
	rc, err := zip.OpenReader(xapkPath)
	if err != nil {
		return fmt.Errorf("open xapk: %w", err)
	}
	defer rc.Close()

	for _, f := range rc.File {
		if !strings.HasSuffix(f.Name, ".apk") {
			continue
		}
		r, err := f.Open()
		if err != nil {
			return err
		}
		apkData, err := io.ReadAll(r)
		r.Close()
		if err != nil {
			return err
		}
		p12Data, err := fileFromZipBytes(apkData, aliceAssetPath)
		if err != nil {
			continue
		}
		password, err := findP12Password(p12Data)
		if err != nil {
			continue
		}
		return decodeAndSave(p12Data, password)
	}
	return errors.New("no APK inside XAPK contains alice.p12")
}

// ImportP12 imports a .p12 directly. Tries the provided password first,
// then falls back to auto-discovering the password from known encrypted strings.
func ImportP12(p12Path, password string) error {
	data, err := os.ReadFile(p12Path)
	if err != nil {
		return err
	}
	if password != "" {
		if err := decodeAndSave(data, password); err == nil {
			return nil
		}
	}
	p, err := findP12Password(data)
	if err != nil {
		return fmt.Errorf("pkcs12: wrong password and no known password matched")
	}
	return decodeAndSave(data, p)
}

func decodeAndSave(p12Data []byte, password string) error {
	// ToPEM handles any number of safe bags (alice.p12 has leaf + CA cert + key),
	// unlike Decode which requires exactly 2 bags and panics on chains.
	blocks, err := pkcs12.ToPEM(p12Data, password)
	if err != nil {
		return fmt.Errorf("pkcs12 decode: %w", err)
	}

	// ToPEM order is not guaranteed — CA cert may precede leaf cert.
	// Select leaf cert: the one without IsCA flag in BasicConstraints.
	var leafBlock, keyBlock *pem.Block
	for _, b := range blocks {
		switch b.Type {
		case "CERTIFICATE":
			if leafBlock == nil {
				cert, err := x509.ParseCertificate(b.Bytes)
				if err == nil && (!cert.BasicConstraintsValid || !cert.IsCA) {
					leafBlock = b
				}
			}
		case "PRIVATE KEY", "RSA PRIVATE KEY", "EC PRIVATE KEY":
			if keyBlock == nil {
				keyBlock = b
			}
		}
	}
	// Fallback: alice.p12 might have no BasicConstraints — take any CERTIFICATE.
	if leafBlock == nil {
		for _, b := range blocks {
			if b.Type == "CERTIFICATE" {
				leafBlock = b
				break
			}
		}
	}
	if leafBlock == nil {
		return errors.New("no certificate found in p12")
	}
	if keyBlock == nil {
		return errors.New("no private key found in p12")
	}

	d, err := configDir()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(d, "cert.pem"), pem.EncodeToMemory(leafBlock), 0600); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d, "key.pem"), pem.EncodeToMemory(keyBlock), 0600)
}

func fileFromZip(zipPath, target string) ([]byte, error) {
	rc, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return findInFiles(rc.File, target)
}

func fileFromZipBytes(data []byte, target string) ([]byte, error) {
	rc, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	return findInFiles(rc.File, target)
}

func findInFiles(files []*zip.File, target string) ([]byte, error) {
	for _, f := range files {
		if f.Name == target {
			r, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer r.Close()
			return io.ReadAll(r)
		}
	}
	return nil, fmt.Errorf("%s not found in archive", target)
}
