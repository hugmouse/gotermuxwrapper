package gotermux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Some vars. Gonna rid of them someday, but not today!
var (
	TD = "termux-dialog" // Just for saving some space
	RT = TResult{}       // Same thing as above
)

// ShareAction from termux-share
//
// Used for TShareView/TShareEdit/TShareSend constants
type ShareAction uint8

// Some constants for better readability when you going to use functions
const (
	TShareView ShareAction = iota // TermuxShare's action "View" flag
	TShareEdit                    // TermuxShare's action "Edit" flag
	TShareSend                    // TermuxShare's action "Send" flag
)

// TermuxDialog spawns new dialog with only title in it
func TermuxDialog(title string) TResult {
	executed := ExecAndListen(TD, []string{
		"-t", title})
	err := json.Unmarshal(executed, &RT)
	if err != nil {
		log.Println(err)
	}
	return RT
}

// TermuxDialogConfirm spawns new confirmation dialog
func TermuxDialogConfirm(td TDialogConfirm) TResult {
	executed := ExecAndListen(TD, []string{
		"confirm",
		"-i", td.Hint,
		"-t", td.Title,
	})
	err := json.Unmarshal(executed, &RT)
	if err != nil {
		log.Println(err)
	}
	return RT
}

// TermuxDialogCheckbox spawns new dialog with multiple values using checkboxes
func TermuxDialogCheckbox(td TDialogCheckbox) TResult {
	values := strings.Join(td.Values, ",")
	executed := ExecAndListen(TD, []string{
		"checkbox",
		"-v", values,
		"-t", td.Title,
	})
	err := json.Unmarshal(executed, &RT)
	if err != nil {
		log.Println(err)
	}
	return RT
}

// TermuxDialogCounter spawns new dialog with pick function in it
//
// User can pick a number in specified range
func TermuxDialogCounter(td TDialogCounter) TResult {
	values := fmt.Sprintf("%d,%d,%d", td.Min, td.Max, td.Start)
	executed := ExecAndListen(TD, []string{
		"counter",
		"-r", values,
		"-t", td.Title,
	})
	err := json.Unmarshal(executed, &RT)
	if err != nil {
		log.Println(err)
	}
	return RT
}

// TODO: Rewrite TermuxDialogDate

// TermuxDialogRadio spawns new dialog with pick function in it
//
// User can pick a single value from radio buttons
func TermuxDialogRadio(td TDialogRadio) TResult {
	values := strings.Join(td.Values, ",")
	executed := ExecAndListen(TD, []string{
		"radio",
		"-v", values,
		"-t", td.Title,
	})
	err := json.Unmarshal(executed, &RT)
	if err != nil {
		log.Println(err)
	}
	return RT
}

// TermuxDialogSheet spawns new dialog with pick function in it
//
// User can pick a value from sliding bottom sheet
//
// Be aware that this function returns "0" in the code result, not "-1" like others (Radio, Spinner)
func TermuxDialogSheet(td TDialogSheet) TResult {
	values := strings.Join(td.Values, ",")
	executed := ExecAndListen(TD, []string{
		"sheet",
		"-v", values,
		"-t", td.Title,
	})
	err := json.Unmarshal(executed, &RT)
	if err != nil {
		log.Println(err)
	}
	return RT
}

// TermuxDialogSpinner spawns new dialog with pick function in it
//
// User can pick a single value from a dropdown spinner
func TermuxDialogSpinner(td TDialogSpinner) TResult {
	values := strings.Join(td.Values, ",")
	executed := ExecAndListen(TD, []string{
		"spinner",
		"-v", values,
		"-t", td.Title,
	})
	err := json.Unmarshal(executed, &RT)
	if err != nil {
		log.Println(err)
	}
	return RT
}

// TermuxDialogSpeech spawns a new dialog that can obtain speech using device microphone
func TermuxDialogSpeech(td TDialogSpeech) TResult {
	executed := ExecAndListen(TD, []string{
		"speech",
		"-i", td.Hint,
		"-t", td.Title,
	})
	err := json.Unmarshal(executed, &RT)
	if err != nil {
		log.Println(err)
	}
	return RT
}

// TermuxDialogText spawns a new dialog with input text (default if no widget specified)
func TermuxDialogText(td TDialogText) TResult {
	if td.MultipleLine == true && td.NumberInput == true {
		log.Fatalln("Cannot use multilines with input numbers (see wiki.termux.com/wiki/Termux-dialog)")
	}

	command := []string{
		"text",
		"-i", td.Hint,
		"-t", td.Title,
	}

	if td.MultipleLine == true {
		command = append(command, "-m")
	}

	if td.NumberInput == true {
		command = append(command, "-n")
	}

	executed := ExecAndListen(TD, command)
	err := json.Unmarshal(executed, &RT)
	if err != nil {
		log.Println(err)
	}
	return RT
}

// TermuxDialogTime spawns new dialog with pick function in it
//
// User can pick a time value
func TermuxDialogTime(td TDialogTime) TResult {
	executed := ExecAndListen(TD, []string{
		"time",
		"-t", td.Title,
	})

	err := json.Unmarshal(executed, &RT)
	if err != nil {
		log.Println(err)
	}
	return RT
}

// TermuxBatteryStatus returns the status of the device battery
func TermuxBatteryStatus() TBattery {
	t := TBattery{}
	status := ExecAndListen("termux-battery-status", nil)

	err := json.Unmarshal(status, &t)
	if err != nil {
		log.Fatalln(err)
	}

	return t
}

// TermuxBrightness sets the display brightness.
//
// Note that this may not work if automatic brightness control is enabled.
func TermuxBrightness(val uint8) []byte {
	u := strconv.FormatUint(uint64(val), 10)
	executed := ExecAndListen("termux-brightness", []string{
		u})
	return executed
}

// TermuxClipboardGet gets the system clipboard text
func TermuxClipboardGet() string {
	executed := ExecAndListen("termux-clipboard-get", nil)
	return string(executed)
}

// TermuxClipboardSet sets the system clipboard text
func TermuxClipboardSet(clip TClipboard) {
	if len(clip.Text) > 0 {
		ExecAndListen("termux-clipboard-set", []string{
			clip.Text,
		})
	} else {
		log.Println("Clipboard is empty!")
	}
}

// TermuxContactList returns list of contacts
func TermuxContactList() TCalls {
	c := TCalls{}
	executed := ExecAndListen("termux-contact-list", nil)
	err := json.Unmarshal(executed, &c)
	if err != nil {
		log.Println(err)
	}
	return c
}

// TermuxDownload downloads a resource using the system download manager
//
// Returns nothing. See: wiki.termux.com/wiki/Termux-download
func TermuxDownload(description, title string) {
	ExecAndListen("termux-download", []string{
		"-d", description,
		"-t", title,
	})
}

// TermuxInfraredFrequencies query the infrared transmitter's supported carrier frequencies
func TermuxInfraredFrequencies() string {
	executed := ExecAndListen("termux-infrared-frequencies", nil)
	return string(executed)
}

// TermuxInfraredTransmit transmits an infrared pattern
func TermuxInfraredTransmit(timings []uint) string {

	// this is cool, but readability is shit and performance too
	// uint -> string
	// [1 2 3 4 5]
	// [1,2,3,4,5]
	//  1,2,3,4,5
	values := strings.Trim(strings.Replace(fmt.Sprint(timings), " ", ",", -1), "[]")

	executed := ExecAndListen("termux-infrared-transmit", []string{
		"-f", values,
	})
	return string(executed)
}

// TermuxLocation gets device location
func TermuxLocation(location TLocation) TLocationResult {
	result := TLocationResult{}
	executed := ExecAndListen("termux-location", []string{
		"-p", location.Provider,
		"-r", location.Request,
	})
	err := json.Unmarshal(executed, &result)
	if err != nil {
		log.Println(err)
	}
	return result
}

// TermuxMediaPlayerPlayFile plays specified media file
func TermuxMediaPlayerPlayFile(path string) {
	ExecAndListen("termux-media-player", []string{
		"play", path,
	})
}

// TermuxMediaPlayerResume resumes playback if paused
func TermuxMediaPlayerResume() {
	ExecAndListen("termux-media-player", []string{"play"})
}

// TermuxMediaPlayerStop quits playback
func TermuxMediaPlayerStop() {
	ExecAndListen("termux-media-player", []string{"stop"})
}

// TermuxMediaPlayerPause pauses playback
func TermuxMediaPlayerPause() {
	ExecAndListen("termux-media-player", []string{"pause"})
}

// TermuxMediaPlayerInfo displays current playback information
func TermuxMediaPlayerInfo() string {
	executed := ExecAndListen("termux-media-player", []string{"info"})
	return string(executed)
}

// TermuxMediaPlayerScan scans the specified file(s) and add to the media content provider
//
// recur - scans directories recursively [set True to use]
//
// verbose - verbose mode [set True to use]
func TermuxMediaPlayerScan(recur, verbose bool) string {

	var command []string

	if recur == true {
		command = append(command, "-r")
	}
	if verbose == true {
		command = append(command, "-v")
	}

	executed := ExecAndListen("termux-media-scan", command)
	return string(executed)
}

// TermuxShare shares a file specified as argument or the text received on stdin
func TermuxShare(t TShare) string {
	command := []string{"-a"}

	switch t.Action {
	default:
	case TShareView:
		command = append(command, "view")
		break
	case TShareEdit:
		command = append(command, "edit")
		break
	case TShareSend:
		command = append(command, "send")
	}

	if t.Default == true {
		command = append(command, "-d")
	}

	if t.Title != "" {
		command = append(command, "-t", t.Title)
	}

	return string(ExecAndListen("termux-share", command))
}

// TermuxVibrate vibrate the device
func TermuxVibrate(t TVibrate) {
	command := []string{"-t", string(t.Duration)}

	if t.SilentModeIgnore == true {
		command = append(command, "-f")
	}

	ExecAndListen("termux-vibrate", command)
}

// ExecAndListen is a function, that build around "exec.Command()"
//
// returns cmd output
func ExecAndListen(command string, args []string) []byte {
	//log.Printf("Arguments: %+v\n", args)
	cmd := exec.Command(command, args...)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		_, err := os.Stderr.WriteString(err.Error())
		if err != nil {
			log.Fatalln("I really don't know how you done this. But you did.", err)
		}
	}
	return cmdOutput.Bytes()
}
