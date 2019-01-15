package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	keychain "github.com/keybase/go-keychain"
)

const airportBin = "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport"

func main() {
	auth, ssid := getWifiInfo()

	if ssid == "" {
		fmt.Fprintln(os.Stderr, "Could not retrieve SSID.")
		os.Exit(1)
	}
	fmt.Printf("\033[90m … getting password for \"%s\". \033[39m\n", ssid)
	fmt.Println("\033[90m … keychain prompt incoming. \033[39m")

	password, err := getWifiPassword(ssid)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not retrieve password. %v", err)
		os.Exit(1)
	}

	if password == "" {
		fmt.Fprintln(os.Stderr, "Could not retrieve password.")
		os.Exit(1)
	}

	fmt.Printf("\033[96m ✓ \"%s\" \033[39m\n\n", password)

	qr := composeWifiString(auth, ssid, password)

	displayQRCode(qr)
}

func getWifiPassword(ssid string) (password string, err error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetDescription("AirPort network password")
	query.SetAccount(ssid)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	results, err := keychain.QueryItem(query)
	if err != nil {
		return "", err
	} else if len(results) != 1 {
		return "", nil
	} else {
		password = string(results[0].Data)
	}

	return password, err
}

func getWifiInfo() (authenticationType string, ssid string) {
	output, _ := exec.Command(airportBin, "-I").Output()
	authenticationType = getWifiInfoFromRegex(string(output), "link auth: (.*)")

	if strings.Contains(authenticationType, "wpa") {
		authenticationType = "WPA"
	} else if strings.Contains(authenticationType, "wep") {
		authenticationType = "WEP"
	}

	ssid = getWifiInfoFromRegex(string(output), " SSID: (.*)")

	return
}

func getWifiInfoFromRegex(output string, expr string) string {
	regex, _ := regexp.Compile(expr)
	match := regex.FindStringSubmatch(output)
	return match[1]
}

func displayQRCode(qr string) {
	config := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
		QuietZone: 1,
	}
	qrterminal.GenerateWithConfig(qr, config)
}

func composeWifiString(authenticationType string, ssid string, password string) string {
	if authenticationType == "WPA" || authenticationType == "WEP" {
		authenticationType = "T:" + authenticationType
	} else {
		authenticationType = ""
	}

	ssid = escapeMECARDString(ssid)
	ssid = "S:" + ssid

	password = escapeMECARDString(password)
	if authenticationType != "" {
		password = "P:" + password
	}

	return fmt.Sprintf("WIFI:%s;%s;%s;", authenticationType, ssid, password)
}

func escapeMECARDString(s string) string {
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, `;`, `\;`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)
	return s
}
