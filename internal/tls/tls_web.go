/*
© Copyright IBM Corporation 2019, 2024

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package tls

import (
	"fmt"
	"os"

	"github.com/ibm-messaging/mq-container/internal/keystore"
	"github.com/ibm-messaging/mq-container/internal/mqtemplate"
	"github.com/ibm-messaging/mq-container/internal/pathutils"
	"github.com/ibm-messaging/mq-container/pkg/logger"
)

// webKeystoreDefault is the name of the default web server Keystore
const webKeystoreDefault = "default.p12"

// ConfigureWebTLS configures TLS for the web server
func ConfigureWebTLS(keyLabel string, log *logger.Logger) error {

	// Return immediately if we have no certificate to use as identity
	if keyLabel == "" && os.Getenv("MQ_GENERATE_CERTIFICATE_HOSTNAME") == "" {
		return nil
	}

	tlsConfigLink := "/run/tls.xml"
	tlsConfigTemplate := "/etc/mqm/web/installations/Installation1/servers/mqweb/tls.xml.tpl"

	err := mqtemplate.ProcessTemplateFile(tlsConfigTemplate, tlsConfigLink, map[string]string{}, log)
	if err != nil {
		return err
	}

	return nil
}

// ConfigureWebKeyStore configures the Web Keystore
func ConfigureWebKeystore(p12Truststore KeyStoreData, keyLabel string) (string, error) {

	webKeystore := webKeystoreDefault
	if keyLabel != "" {
		webKeystore = keyLabel + ".p12"
	}
	webKeystoreFile := pathutils.CleanPath(keystoreDirDefault, webKeystore)

	// Check if a new self-signed certificate should be generated
	if keyLabel == "" {

		// Get hostname to use for self-signed certificate
		genHostName := os.Getenv("MQ_GENERATE_CERTIFICATE_HOSTNAME")

		// Create the Web Keystore
		newWebKeystore := keystore.NewPKCS12KeyStore(webKeystoreFile, p12Truststore.Password)
		err := newWebKeystore.Create()
		if err != nil {
			return "", fmt.Errorf("Failed to create Web Keystore %s: %v", webKeystoreFile, err)
		}

		// Generate a new self-signed certificate in the Web Keystore
		err = newWebKeystore.CreateSelfSignedCertificate("default", fmt.Sprintf("CN=%s", genHostName), genHostName)
		if err != nil {
			return "", fmt.Errorf("Failed to generate certificate in Web Keystore %s with DN of 'CN=%s': %v", webKeystoreFile, genHostName, err)
		}
	} else {
		// Check Web Keystore already exists
		_, err := os.Stat(webKeystoreFile)
		if err != nil {
			return "", fmt.Errorf("Failed to find existing Web Keystore %s: %v", webKeystoreFile, err)
		}
	}

	// Check Web Truststore already exists
	_, err := os.Stat(p12Truststore.Keystore.Filename)
	if err != nil {
		return "", fmt.Errorf("Failed to find existing Web Truststore %s: %v", p12Truststore.Keystore.Filename, err)
	}

	return webKeystore, nil
}
