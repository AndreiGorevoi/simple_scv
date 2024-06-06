package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	indexPath  = "./vcs/index.txt"
	logPath    = "./vcs/log.txt"
	configPath = "./vcs/config.txt"
)

func main() {
	prepareFolders()
	args := os.Args
	// No flags were sent
	if len(args) == 1 {
		printHelp()
		return
	}

	arg := args[1]

	switch arg {
	case "--help":
		printHelp()
	case "config":
		handleConfig(args[2:])
	case "add":
		handleAdd(args[2:])
	case "log":
		handleLog()
	case "commit":
		handleCommit(args[2:])
	case "checkout":
		handleCheckOut(args[2:])
	default:
		fmt.Printf("'%v' is not a SVCS command.\n", arg)
	}
}

func printHelp() {
	fmt.Println(`These are SVCS commands:
config     Get and set a username.
add        Add a file to the index.
log        Show commit logs.
commit     Save changes.
checkout   Restore a file.`)
}

// prepareFolders creates all required foldres if they have not been created yet
func prepareFolders() {
	if _, err := os.Stat("vcs"); os.IsNotExist(err) {
		err = os.MkdirAll("./vcs", os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		err = os.MkdirAll("./vcs/commits", os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range []string{"config.txt", "index.txt", "log.txt"} {
			file, err := os.Create("./vcs/" + f)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
		}
	}
}

func handleConfig(args []string) {
	data, err := os.ReadFile("./vcs/config.txt")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	if len(args) == 0 && len(data) == 0 { // name is not provided, config is empty. Print info
		fmt.Println("Please, tell me who you are.")
	} else if len(args) > 0 { // name is provided. Rewrite user
		os.WriteFile("./vcs/config.txt", []byte(args[0]), 0666)
		fmt.Printf("The username is %s.\n", args[0])
	} else { // name is not provided, but config is not empty. Print current user
		fmt.Printf("The username is %s.\n", string(data))
	}
}

func handleAdd(args []string) {
	file, err := os.OpenFile(indexPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)

	if len(args) == 0 && !sc.Scan() { // file name is not provided, index is empty. Print info
		fmt.Println("Add a file to the index.")
	} else if len(args) > 0 { // file name is provided. Add to index
		_, err = os.Stat(args[0])
		if os.IsNotExist(err) {
			fmt.Printf("Can't find '%s'.\n", args[0])
			return
		}
		fmt.Fprintln(file, args[0])
		fmt.Printf("The file '%s' is tracked.\n", args[0])
	} else { // file name is not provided, but index is not empty. Print current index
		fmt.Println("Tracked files:")
		fmt.Println(sc.Text())
		for sc.Scan() {
			fmt.Println(sc.Text())
		}
	}
}

func handleCommit(args []string) {
	if len(args) < 1 {
		fmt.Println("Message was not passed.")
		return
	}
	data := readIndexAsSlice()
	if len(data) > 0 && somehtingChanged(data) {
		activeUser := getUserName()
		commitHash := hashCommit(activeUser)
		commitPath := "./vcs/commits/" + commitHash
		createCommit(commitPath, data) // write everything from idex if no commits before
		addCommitToLog(args[0], activeUser, commitHash)
		fmt.Println("Changes are committed.")
	} else {
		fmt.Println("Nothing to commit.")
	}
}

// getLastCommitInfo returns the last commit info from log.txt file
func getLastCommitInfo() string {
	f, err := os.Open(logPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var lastCommit string
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		lastCommit = sc.Text()
	}

	return lastCommit
}

// hashCommit creates a hash ussing currect time + userName by sha1 algorithm and returns result
// as hexdecimal string
func hashCommit(userName string) string {
	sha := sha1.New()
	sha.Write([]byte(time.Now().String()))
	sha.Write([]byte(userName))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

// getUserName reads a current user from cofig.txt file
func getUserName() string {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

// addCommitToLog adds whole information into log.txt about a new commit
func addCommitToLog(msg, user, commitHash string) {
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_RDWR, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	fmt.Fprintf(logFile, "%s|%s|%s\n", commitHash, user, msg)
}

// createCommit creates a dir for a new commit and put files in it
func createCommit(commitPath string, files []string) {
	err := os.MkdirAll(commitPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		src, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()

		dst, err := os.Create(filepath.Join(commitPath, f))

		if err != nil {
			log.Fatal(err)
		}
		defer dst.Close()
		io.Copy(dst, src)
	}
}

// readIndexAsSlice reads index.txt file and return list of files in it
func readIndexAsSlice() []string {
	file, err := os.Open(indexPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	res := make([]string, 0)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		res = append(res, sc.Text())
	}
	return res
}

// somehtingChanged check if files were changed after last commit
func somehtingChanged(data []string) bool {
	lastCommit := getLastCommitInfo()
	if lastCommit == "" {
		return true
	}
	h1 := sha256.New()
	h2 := sha256.New()
	for _, fn := range data {
		src, err := os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()
		io.Copy(h1, src)

		commitHash := strings.Split(lastCommit, "|")[0]
		entries, err := os.ReadDir(filepath.Join("./vcs/commits/", commitHash))

		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {
			if entry.Name() == fn {
				f2, err := os.Open(filepath.Join("./vcs/commits/", commitHash, entry.Name()))
				if err != nil {
					log.Fatal(err)
				}

				defer f2.Close()
				io.Copy(h2, f2)
			}
		}

		if !bytes.Equal(h1.Sum(nil), h2.Sum(nil)) {
			return true
		}
		h1.Reset()
		h2.Reset()
	}

	return false
}

func handleLog() {
	fs, err := os.Stat(logPath)
	if err != nil {
		log.Fatal(err)
	}
	if fs.Size() <= 1 {
		fmt.Println("No commits yet.")
		return
	}

	file, err := os.Open(logPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	sc := bufio.NewScanner(file)
	logs := make([]string, 0)

	for sc.Scan() {
		logs = append(logs, sc.Text())
	}

	for i := len(logs) - 1; i >= 0; i-- {
		tokens := strings.Split(logs[i], "|")
		commitHash := tokens[0]
		author := tokens[1]
		msg := tokens[2]
		fmt.Printf("commit %s\nAuthor: %s\n%s\n\n", commitHash, author, msg)
	}

}

func handleCheckOut(args []string) {
	if len(args) < 1 {
		fmt.Print("Commit id was not passed.")
		return
	}

	folder := filepath.Join("./vcs/commits", args[0])
	entities, err := os.ReadDir(folder)

	if os.IsNotExist(err) {
		fmt.Println("Commit does not exist.")
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, entity := range entities {
		dst, err := os.Create(entity.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer dst.Close()

		src, err := os.Open(filepath.Join(folder, entity.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()

		io.Copy(dst, src)
	}
	fmt.Printf("Switched to commit %s.", args[0])
}
