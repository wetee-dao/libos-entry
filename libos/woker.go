package libos

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	// "github.com/edgelesssys/ego/attestation"
	// "github.com/edgelesssys/ego/attestation/tcbstatus"
	// "github.com/edgelesssys/ego/eclient"
)

// func verifyReport(reportBytes, certBytes, signer []byte) error {
// 	report, err := eclient.VerifyRemoteReport(reportBytes)
// 	if err == attestation.ErrTCBLevelInvalid {
// 		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
// 		fmt.Println("We'll ignore this issue in this sample. For an app that should run in production, you must decide which of the different TCBStatus values are acceptable for you to continue.")
// 	} else if err != nil {
// 		return err
// 	}

// 	hash := sha256.Sum256(certBytes)
// 	if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
// 		return errors.New("report data does not match the certificate's hash")
// 	}

// 	// You can either verify the UniqueID or the tuple (SignerID, ProductID, SecurityVersion, Debug).

// 	if report.SecurityVersion < 2 {
// 		return errors.New("invalid security version")
// 	}
// 	if binary.LittleEndian.Uint16(report.ProductID) != 1234 {
// 		return errors.New("invalid product")
// 	}
// 	if !bytes.Equal(report.SignerID, signer) {
// 		return errors.New("invalid signer")
// 	}

// 	// For production, you must also verify that report.Debug == false

// 	return nil
// }

func workerGet(tlsConfig *tls.Config, url string) []byte {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func workerPost(tlsConfig *tls.Config, url string, json string) ([]byte, error) {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	payload := strings.NewReader(json)
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
