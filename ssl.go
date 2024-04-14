package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"strings"
)

func GenerateKeyPair(domain, filePath string, signInfo *pkix.Name) (csr, key string, err error) {
	keyBytes, key, err := GenerateKey(domain, filePath)
	if err != nil {
		return
	}

	csr, err = GenerateCSR(keyBytes, domain, filePath, signInfo)
	if err != nil {
		return
	}

	return
}

/*
CSR will be generated in
filePath+domain-com+".pub"
*/
func GenerateCSR(keyBytes *rsa.PrivateKey, domain, filePath string, signInfo *pkix.Name) (string, error) {

	if signInfo == nil {
		oidEmailAddress := asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}
		signInfo = &pkix.Name{
			CommonName:         domain,
			Country:            []string{"ID"},
			Province:           []string{"Jakarta"},
			Locality:           []string{"Jakarta Barat"},
			Organization:       []string{"Bendt Indonesia"},
			OrganizationalUnit: []string{"IT Services"},
			ExtraNames: []pkix.AttributeTypeAndValue{
				{
					Type: oidEmailAddress,
					Value: asn1.RawValue{
						Tag:   asn1.TagIA5String,
						Bytes: []byte("admin@" + domain),
					},
				},
			},
		}
	}

	template := x509.CertificateRequest{
		Subject:            *signInfo,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	//Check should be ended by /
	if filePath != "" {
		lastChar := string(filePath[len(filePath)-1])
		if lastChar != "/" {
			filePath += "/"
		}
	}

	domain = strings.ReplaceAll(domain, ".", "-")
	filePath += domain + ".pub"

	var publicCSR bytes.Buffer
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	if err != nil {
		return "", err
	}

	publicKeyBlock := &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}
	err = pem.Encode(&publicCSR, publicKeyBlock)
	if err != nil {
		return "", err
	}

	err = WriteFile("", filePath, publicCSR.String())
	if err != nil {
		return publicCSR.String(), err
	}

	return publicCSR.String(), nil
}

/*
Key will be generated in
filePath+domain-com
*/
func GenerateKey(domain, filePath string) (*rsa.PrivateKey, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	//Check should be ended by /
	if filePath != "" {
		lastChar := string(filePath[len(filePath)-1])
		if lastChar != "/" {
			filePath += "/"
		}
	}

	domain = strings.ReplaceAll(domain, ".", "-")
	filePath += domain

	var buff bytes.Buffer
	privateKeyBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	err = pem.Encode(&buff, privateKeyBlock)
	if err != nil {
		panic(err)
	}

	err = WriteFile("", filePath, buff.String())
	if err != nil {
		return nil, buff.String(), err
	}

	return privateKey, buff.String(), nil
}
