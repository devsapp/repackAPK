package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/mozilla-services/pkcs7"
)

// consts ...
const (
	CertValidYears = 30
)

// sha1Sum ...
func sha1Sum(msg []byte) string {
	sha := sha1.Sum(msg)
	return base64.StdEncoding.EncodeToString(sha[:])
}

func signSF(fcCtx *FCContext) ([]byte, error) {
	sfFile := fmt.Sprintf("%s/%s.SF", fcCtx.WorkDir, fcCtx.SigFileName)
	cmd := exec.Command("openssl", "smime", "-sign", "-binary", "-noattr", "-outform", "DER",
		"-in", sfFile, "-inkey", PrivateKeyPEM_PATH, "-signer", CertPEM_PATH)
	buf, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cmd error: %v, output: %s", err, string(buf))
	}
	return buf, nil
}

func SignAndDetach(content []byte, cert *x509.Certificate, privkey *rsa.PrivateKey) (signed []byte, err error) {
	toBeSigned, err := pkcs7.NewSignedData(content)
	if err != nil {
		err = fmt.Errorf("Cannot initialize signed data: %s", err)
		return
	}
	if err = toBeSigned.AddSigner(cert, privkey, pkcs7.SignerInfoConfig{}); err != nil {
		err = fmt.Errorf("Cannot add signer: %s", err)
		return
	}

	// Detach signature, omit if you want an embedded signature
	toBeSigned.Detach()

	signed, err = toBeSigned.Finish()
	if err != nil {
		err = fmt.Errorf("Cannot finish signing data: %s", err)
		return
	}

	// Verify the signature
	pem.Encode(os.Stdout, &pem.Block{Type: "PKCS7", Bytes: signed})
	p7, err := pkcs7.Parse(signed)
	if err != nil {
		err = fmt.Errorf("Cannot parse our signed data: %s", err)
		return
	}

	// since the signature was detached, reattach the content here
	p7.Content = content

	if bytes.Compare(content, p7.Content) != 0 {
		err = fmt.Errorf("Our content was not in the parsed data:\n\tExpected: %s\n\tActual: %s", content, p7.Content)
		return
	}
	if err = p7.Verify(); err != nil {
		err = fmt.Errorf("Cannot verify our signed data: %s", err)
		return
	}

	return signed, nil
}

// signPKCS7 does the minimal amount of work necessary to embed an RSA
// signature into a PKCS#7 certificate.
//
// We prepare the certificate using the x509 package, read it back in
// to our custom data type and then write it back out with the signature.
func signPKCS7(rand io.Reader, priv *rsa.PrivateKey, msg []byte) ([]byte, error) {
	buf, err := ioutil.ReadFile(CertPEM_PATH)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(buf)
	if block == nil {
		return nil, fmt.Errorf("failed to decode pem: %s", CertPEM_PATH)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	b, err := x509.CreateCertificate(rand, cert, cert, priv.Public(), priv)
	if err != nil {
		return nil, err
	}

	c := certificate{}
	if _, err := asn1.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("asn1 unmarshal: %v", err)
	}

	h := sha1.New()
	h.Write(msg)
	hashed := h.Sum(nil)

	signed, err := rsa.SignPKCS1v15(rand, priv, crypto.SHA1, hashed)
	if err != nil {
		return nil, err
	}

	content := pkcs7SignedData{
		ContentType: oidSignedData,
		Content: signedData{
			Version: 1,
			DigestAlgorithms: []pkix.AlgorithmIdentifier{{
				Algorithm:  oidSHA1,
				Parameters: asn1.RawValue{Tag: 5},
			}},
			ContentInfo:  contentInfo{Type: oidData},
			Certificates: c,
			SignerInfos: []signerInfo{{
				Version: 1,
				IssuerAndSerialNumber: issuerAndSerialNumber{
					Issuer:       cert.Issuer.ToRDNSequence(),
					SerialNumber: int(cert.SerialNumber.Int64()),
				},
				DigestAlgorithm: pkix.AlgorithmIdentifier{
					Algorithm:  oidSHA1,
					Parameters: asn1.RawValue{Tag: 5},
				},
				DigestEncryptionAlgorithm: pkix.AlgorithmIdentifier{
					Algorithm:  oidRSAEncryption,
					Parameters: asn1.RawValue{Tag: 5},
				},
				EncryptedDigest: signed,
			}},
		},
	}

	return asn1.Marshal(content)
}

type pkcs7SignedData struct {
	ContentType asn1.ObjectIdentifier
	Content     signedData `asn1:"tag:0,explicit"`
}

// signedData is defined in rfc2315, section 9.1.
type signedData struct {
	Version          int
	DigestAlgorithms []pkix.AlgorithmIdentifier `asn1:"set"`
	ContentInfo      contentInfo
	Certificates     certificate  `asn1:"tag0,explicit"`
	SignerInfos      []signerInfo `asn1:"set"`
}

type contentInfo struct {
	Type asn1.ObjectIdentifier
	// Content is optional in PKCS#7 and not provided here.
}

// certificate is defined in rfc2459, section 4.1.
type certificate struct {
	TBSCertificate     tbsCertificate
	SignatureAlgorithm pkix.AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

// tbsCertificate is defined in rfc2459, section 4.1.
type tbsCertificate struct {
	Version      int `asn1:"tag:0,default:2,explicit"`
	SerialNumber int
	Signature    pkix.AlgorithmIdentifier
	Issuer       pkix.RDNSequence // pkix.Name
	Validity     validity
	Subject      pkix.RDNSequence // pkix.Name
	SubjectPKI   subjectPublicKeyInfo
	Extensions   []pkix.Extension `asn1:"tag:3,explicit"`
}

// validity is defined in rfc2459, section 4.1.
type validity struct {
	NotBefore time.Time
	NotAfter  time.Time
}

// subjectPublicKeyInfo is defined in rfc2459, section 4.1.
type subjectPublicKeyInfo struct {
	Algorithm        pkix.AlgorithmIdentifier
	SubjectPublicKey asn1.BitString
}

type signerInfo struct {
	Version                   int
	IssuerAndSerialNumber     issuerAndSerialNumber
	DigestAlgorithm           pkix.AlgorithmIdentifier
	DigestEncryptionAlgorithm pkix.AlgorithmIdentifier
	EncryptedDigest           []byte
}

type issuerAndSerialNumber struct {
	Issuer       pkix.RDNSequence // pkix.Name
	SerialNumber int
}

// Various ASN.1 Object Identifies, mostly from rfc3852.
var (
	oidPKCS7         = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7}
	oidData          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 1}
	oidSignedData    = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 2}
	oidSHA1          = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 26}
	oidRSAEncryption = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
)
