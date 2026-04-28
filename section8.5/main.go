package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

// ANSI カラーコード
const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
)

// Jarvis AIアシスタント構造体
type Jarvis struct {
	name      string
	user      string
	bootTime  time.Time
	systems   map[string]bool
	powerLvl  int
}

// 新規Jarvisインスタンス生成
func NewJarvis(user string) *Jarvis {
	return &Jarvis{
		name:     "J.A.R.V.I.S.",
		user:     user,
		bootTime: time.Now(),
		systems: map[string]bool{
			"ARC REACTOR":      true,
			"REPULSORS":        true,
			"FLIGHT STABILIZE": true,
			"TARGETING":        true,
			"COMMS":            true,
			"HUD":              true,
			"LIFE SUPPORT":     true,
			"ARMOR INTEGRITY":  true,
		},
		powerLvl: 100,
	}
}

// タイピングアニメーション風出力
func typeOut(text string, color string) {
	fmt.Print(color)
	for _, c := range text {
		fmt.Printf("%c", c)
		time.Sleep(15 * time.Millisecond)
	}
	fmt.Println(colorReset)
}

// 起動シーケンス
func (j *Jarvis) Boot() {
	clearScreen()
	printLogo()
	time.Sleep(500 * time.Millisecond)

	bootSteps := []string{
		"BIOS check............................[OK]",
		"Loading kernel modules................[OK]",
		"Initializing arc reactor..............[OK]",
		"Calibrating repulsor outputs..........[OK]",
		"Establishing satellite uplink.........[OK]",
		"Synchronizing with STARK NETWORK......[OK]",
		"Loading neural interface..............[OK]",
		"Activating heads-up display...........[OK]",
		"Running diagnostics...................[OK]",
	}

	fmt.Println(colorCyan + "\n[ BOOT SEQUENCE INITIATED ]" + colorReset)
	fmt.Println(strings.Repeat("─", 50))

	for _, step := range bootSteps {
		fmt.Print(colorGreen + "  > " + colorReset)
		typeOut(step, colorBlue)
		time.Sleep(150 * time.Millisecond)
	}

	fmt.Println(strings.Repeat("─", 50))
	fmt.Println(colorGreen + colorBold + "\n[ ALL SYSTEMS NOMINAL ]" + colorReset)
	time.Sleep(700 * time.Millisecond)

	greeting := fmt.Sprintf("\nGood %s, %s. J.A.R.V.I.S. at your service.",
		timeOfDay(), j.user)
	typeOut(greeting, colorCyan+colorBold)
	time.Sleep(300 * time.Millisecond)
}

// ASCIIアートロゴ
func printLogo() {
	logo := `
    ╔═══════════════════════════════════════════════════╗
    ║                                                   ║
    ║      ░░░ J.A.R.V.I.S. MARK VII ░░░                ║
    ║                                                   ║
    ║   Just A Rather Very Intelligent System           ║
    ║   STARK INDUSTRIES © - All Rights Reserved        ║
    ║                                                   ║
    ╚═══════════════════════════════════════════════════╝
`
	fmt.Println(colorCyan + logo + colorReset)
}

// 画面クリア
func clearScreen() {
	if runtime.GOOS == "windows" {
		fmt.Print("\033[H\033[2J")
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

// 時間帯判定
func timeOfDay() string {
	h := time.Now().Hour()
	switch {
	case h < 12:
		return "morning"
	case h < 18:
		return "afternoon"
	default:
		return "evening"
	}
}

// システムステータス表示
func (j *Jarvis) Status() {
	fmt.Println(colorCyan + "\n┌─[ SYSTEM STATUS ]──────────────────────────┐" + colorReset)
	for name, ok := range j.systems {
		status := colorGreen + "ONLINE " + colorReset
		if !ok {
			status = colorRed + "OFFLINE" + colorReset
		}
		fmt.Printf(colorCyan+"│ "+colorReset+"%-25s [%s]"+colorCyan+"  │\n"+colorReset, name, status)
	}
	uptime := time.Since(j.bootTime).Round(time.Second)
	fmt.Printf(colorCyan+"│ "+colorReset+"%-25s %s"+colorCyan+"        │\n"+colorReset,
		"UPTIME:", uptime)
	fmt.Printf(colorCyan+"│ "+colorReset+"%-25s %d%%"+colorCyan+"               │\n"+colorReset,
		"REACTOR OUTPUT:", j.powerLvl)
	fmt.Println(colorCyan + "└────────────────────────────────────────────┘" + colorReset)
}

// コマンド応答生成
func (j *Jarvis) Respond(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))

	switch {
	case strings.Contains(input, "hello"), strings.Contains(input, "hi"), strings.Contains(input, "こんにちは"):
		return fmt.Sprintf("Hello, %s. How may I assist you today?", j.user)

	case strings.Contains(input, "time"), strings.Contains(input, "時間"):
		return fmt.Sprintf("The current time is %s.", time.Now().Format("15:04:05"))

	case strings.Contains(input, "date"), strings.Contains(input, "日付"):
		return fmt.Sprintf("Today is %s.", time.Now().Format("Monday, January 2, 2006"))

	case strings.Contains(input, "status"), strings.Contains(input, "ステータス"):
		j.Status()
		return "Diagnostics complete, sir."

	case strings.Contains(input, "weather"), strings.Contains(input, "天気"):
		return "I'm afraid I don't have access to live weather data in this build, sir."

	case strings.Contains(input, "joke"), strings.Contains(input, "ジョーク"):
		jokes := []string{
			"Why did the AI cross the road? To optimize the chicken's path.",
			"I would tell you a UDP joke, but you might not get it.",
			"There are 10 types of people: those who understand binary, and those who don't.",
		}
		return jokes[rand.Intn(len(jokes))]

	case strings.Contains(input, "suit up"), strings.Contains(input, "armor"):
		return "Initiating Mark VII deployment sequence. Stand clear, sir."

	case strings.Contains(input, "help"):
		return "Available commands: hello, time, date, status, weather, joke, suit up, shutdown, exit."

	case strings.Contains(input, "shutdown"), strings.Contains(input, "exit"), strings.Contains(input, "quit"):
		return "__EXIT__"

	default:
		return "I'm not sure how to respond to that, sir. Try 'help' for a list of commands."
	}
}

// シャットダウンシーケンス
func (j *Jarvis) Shutdown() {
	fmt.Println(colorYellow + "\n[ SHUTDOWN SEQUENCE INITIATED ]" + colorReset)
	steps := []string{
		"Saving session data...",
		"Powering down repulsors...",
		"Closing satellite uplink...",
		"Reactor entering standby mode...",
	}
	for _, s := range steps {
		fmt.Print(colorYellow + "  > " + colorReset)
		typeOut(s, colorBlue)
		time.Sleep(200 * time.Millisecond)
	}
	typeOut(fmt.Sprintf("\nGoodbye, %s. Until next time.", j.user), colorCyan+colorBold)
	fmt.Println()
}

// メイン対話ループ
func (j *Jarvis) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(colorBold + colorGreen + "\n" + j.user + " > " + colorReset)
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "" {
			continue
		}

		response := j.Respond(input)
		if response == "__EXIT__" {
			j.Shutdown()
			return
		}

		fmt.Print(colorCyan + colorBold + "JARVIS > " + colorReset)
		typeOut(response, colorCyan)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	user := "Sir"
	if len(os.Args) > 1 {
		user = os.Args[1]
	}

	jarvis := NewJarvis(user)
	jarvis.Boot()
	jarvis.Run()
}