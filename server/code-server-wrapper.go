package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/phayes/freeport"
)

var (
	cmd *exec.Cmd
	srv *http.Server
)

func launchCodeServer() (*exec.Cmd, int) {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	// Get cmd path
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dirPath := filepath.Dir(exePath)
	serverCmdPath := filepath.Join(dirPath, "code-server")

	// Run code-server
	cmd := exec.Command(serverCmdPath, "--no-auth", "-p", strconv.Itoa(port))
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
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
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	return srv
}

func initMenu() {
	// systray.SetIcon(getIcon("assets/clock.ico"))
	systray.SetTitle("VEDA")
	systray.SetTooltip("VEDA for VSCode Web Server")
	mQuit := systray.AddMenuItem("Quit", "Quit VSCode Web Server")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func cleanup() {
	// Kill HTTP server
	err := srv.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Kill code-server before quit
	cmd.Process.Kill()
}

func main() {
	runtime.LockOSThread()

	wrapperPort, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	cm, port := launchCodeServer()
	cmd = cm

	// Prepare cleanup
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-signalChan
		cleanup()
	}()

	srv = launchWebServer(wrapperPort, port)

	// systray.Run(initMenu, cleanup)
	systray.Run(func() {
		initMenu()
		systray.AddMenuItem(fmt.Sprintf("wp: %d", wrapperPort), "")
		systray.AddMenuItem(fmt.Sprintf("port: %d", port), "")
	}, cleanup)
}
