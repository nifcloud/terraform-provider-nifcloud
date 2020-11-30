package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

const (
	pemBlockTypeCertificate   = `CERTIFICATE`
	pemBlockTypeRsaPrivateKey = `RSA PRIVATE KEY`
)

var tlsX509CertificateSerialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)

// GeneratePrivateKey generates a RSA private key PEM string.
func GeneratePrivateKey(t *testing.T, bits int) string {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		t.Fatal(err)
	}

	return string(pem.EncodeToMemory(&pem.Block{
		Bytes: x509.MarshalPKCS1PrivateKey(key),
		Type:  pemBlockTypeRsaPrivateKey,
	}))
}

// GenerateCertificate generates a local CA x509 certificate PEM string.
func GenerateCertificate(t *testing.T, caKeyPem, caCertificatePem, keyPem, commonName string) string {
	caCertificateBlock, _ := pem.Decode([]byte(caCertificatePem))
	caCertificate, err := x509.ParseCertificate(caCertificateBlock.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	caKeyBlock, _ := pem.Decode([]byte(caKeyPem))
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	keyBlock, _ := pem.Decode([]byte(keyPem))
	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	serialNumber, err := rand.Int(rand.Reader, tlsX509CertificateSerialNumberLimit)
	if err != nil {
		t.Fatal(err)
	}

	certificate := &x509.Certificate{
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		NotAfter:              time.Now().Add(24 * time.Hour),
		NotBefore:             time.Now(),
		SerialNumber:          serialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"terraform-provider-nifcloud-testing"},
		},
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, certificate, caCertificate, &key.PublicKey, caKey)
	if err != nil {
		t.Fatal(err)
	}

	return string(pem.EncodeToMemory(&pem.Block{
		Bytes: certificateBytes,
		Type:  pemBlockTypeCertificate,
	}))
}

// GenerateSelfSignedCertificateAuthority generates a x509 CA certificate PEM string.
func GenerateSelfSignedCertificateAuthority(t *testing.T, keyPem string) string {
	keyBlock, _ := pem.Decode([]byte(keyPem))
	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	serialNumber, err := rand.Int(rand.Reader, tlsX509CertificateSerialNumberLimit)
	if err != nil {
		t.Fatal(err)
	}

	certificate := &x509.Certificate{
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		NotAfter:              time.Now().Add(24 * time.Hour),
		NotBefore:             time.Now(),
		SerialNumber:          serialNumber,
		Subject: pkix.Name{
			CommonName:   "terraform-provider-nifcloud-ca",
			Organization: []string{"terraform-provider-nifcloud-testing"},
		},
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	return string(pem.EncodeToMemory(&pem.Block{
		Bytes: certificateBytes,
		Type:  pemBlockTypeCertificate,
	}))
}
