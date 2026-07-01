package termutils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

const (
	Bold      = "\033[1m"
	Italics   = "\033[3m"
	Underline = "\033[4m"
	Dim       = "\033[2m"
	None      = ""
)

const (
	IconDomain       = "🌍"
	IconServer       = "🖥"
	IconDestination  = "📍"
	IconLoadBalancer = "⚖️"
	IconRedirect     = "🔄"
	IconIssuer       = "🤵"
	IconCommonName   = "📜"
	IconOrganization = "💼"
	IconExpires      = "⏰"
	IconSans         = "✨"
	IconTls          = "🔒"
)

func YesNo(yesNo string) string {
	var textColor string
	switch yesNo {
	case "yes":
		textColor = Green
	case "no":
		textColor = Red
	default:
		textColor = Magenta
	}
	return textColor
}

func PrintDomain(sentence string) {
	PrintSentence(Dim, White, IconDomain, "Domain", sentence)
}

func PrintServerApp(sentence string) {
	PrintSentence(None, Blue, IconServer, "ServerApp", sentence)
}

func PrintDestination(sentence string) {
	PrintSentence(None, Blue, IconDestination, "Destination", sentence+" (Resolved)")
}

func PrintLoadBalancer(sentence string) {
	PrintSentence(None, YesNo(sentence), IconLoadBalancer, "LoadBalancer", sentence)
}

func PrintRedirectOnly(sentence string) {
	PrintSentence(Dim, YesNo(sentence), IconRedirect, "Redirect", sentence)
}

func PrintIssuedTo(sentence string) {
	PrintSentence(Bold, Cyan, IconOrganization, "Issued to", sentence)
}
func PrintCommonName(sentence string) {
	PrintSentence(Bold, White, IconCommonName, "CN", sentence)
}

func PrintIssuer(sentence string) {
	PrintSentence(Bold, Magenta, IconIssuer, "Issuer", sentence)
}

func PrintVersion(version int) {
	sentence := strconv.Itoa(version)
	PrintSentence(Dim, White, IconTls, "X.509 ver", sentence)
}

func PrintAlert(sentence string) {
	fmt.Println(Red + sentence + Reset)
}

func ExpirationColor(expiresIn int) string {
	switch {
	case expiresIn > 30:
		return Green
	case expiresIn >= 15:
		return Yellow
	default:
		return Red
	}
}

func PrintExpiration(expiration string, expiresIn int) {
	sentence := strconv.Itoa(expiresIn) + " days (" + expiration + ")"
	PrintSentence(Bold, ExpirationColor(expiresIn), IconExpires, "Exp. in", sentence)
}

func PrintSans(sans []string) {
	sentence := Dim + "[ " + Reset
	iterations := 0
	for _, value := range sans {
		sentence += Yellow + (value)
		iterations += 1
		if len(sans) > iterations {
			sentence += Dim + " | " + Reset
		}
	}
	sentence += Dim + " ]"
	PrintSentence(None, Yellow, IconSans, "SANs", sentence)
}

func PrintLine() {
	fmt.Println(Dim + strings.Repeat("_", 60) + Reset)
}
func PrintDoubleLine() {
	fmt.Println(Dim + Blue + strings.Repeat("=", 80) + Reset)
}

func Ask(question string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(question)
	scanner.Scan()
	return scanner.Text()
}

func ClearScreen() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error cleaning screen: ", err)
	}
}

func PrintDebug(text string) {
	fmt.Println(Dim + text + Reset)
}

func PrintSentence(style string, color string, icon string, label string, text string) {
	sentence := ""
	if icon != None {
		sentence += icon + Dim + " | " + Reset
	}
	if label != None {
		sentence += Italics + label + ": " + strings.Repeat(Dim+"_"+Reset, (12-len(label))) + " "
	}
	//sentence += fmt.Sprintf("%-10s", label)
	if style != None {
		sentence += style
	}
	if color != None {
		sentence += color
	}
	sentence += text + Reset
	fmt.Println(sentence)
}

func Quit() {
	PrintDebug("Have a nice day!.")
	os.Exit(0)
}
