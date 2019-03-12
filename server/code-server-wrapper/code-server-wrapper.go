package codeServerWrapper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	. "github.com/fand/veda-vscode-web/server/logger"
	"github.com/getlantern/systray"
	"github.com/phayes/freeport"
	"github.com/skratchdot/open-golang/open"
)

var (
	command     *exec.Cmd
	server      *http.Server
	wrapperPort int
	serverPort  int
)

func getCmdPath(p string) string {
	exePath, err := os.Executable()
	if err != nil {
		LogFatal(err)
	}

	dirPath := filepath.Dir(exePath)
	return filepath.Join(dirPath, p)
}

func launchCodeServer() (*exec.Cmd, int) {
	port, err := freeport.GetFreePort()
	if err != nil {
		LogFatal(err)
	}

	// Run code-server
	serverCmdPath := getCmdPath("code-server")
	cmd := exec.Command(serverCmdPath, "--no-auth", "-p", strconv.Itoa(port))

	err = cmd.Start()
	if err != nil {
		LogFatal(err)
	}

	return cmd, port
}

func launchWebServer(wrapperPort, port int) *http.Server {
	// Start HTTP server for IPC
	server := &http.Server{Addr: fmt.Sprintf(":%d", wrapperPort)}

	http.HandleFunc("/port", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, strconv.Itoa(port))
	})

	// Run server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			LogFatal(fmt.Sprintf("ListenAndServe(): %s", err))
		}
	}()

	return server
}

func initMenu() {
	iconPath := getCmdPath("../Resources/icon.icns")
	icon, err := ioutil.ReadFile(iconPath)
	if err != nil {
		LogFatal(err)
	}
	systray.SetIcon(icon)
	systray.SetTitle("VEDA")
	systray.SetTooltip("VEDA for VSCode Web Server")

	mOpen := systray.AddMenuItem("Open VSCode for Web", "Open VSCode in the browser")
	go func() {
		for {
			<-mOpen.ClickedCh
			open.Start(fmt.Sprintf("http://localhost:%d", serverPort))
		}
	}()

	mQuit := systray.AddMenuItem("Quit", "Quit VSCode Web Server")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func cleanup() {
	// Kill HTTP server
	err := server.Close()
	if err != nil {
		LogFatal(err)
	}

	// Kill code-server before quit
	command.Process.Kill()

	FlushLogger()
}

func main() {
	runtime.LockOSThread()
	InitLogger("/tmp/gl.veda.vscode.web.server/log-code-server-wrapper.txt")

	wrapperPort, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	command, serverPort = launchCodeServer()

	// Prepare cleanup
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-signalChan
		cleanup()
	}()

	server = launchWebServer(wrapperPort, serverPort)

	systray.Run(initMenu, cleanup)
}
