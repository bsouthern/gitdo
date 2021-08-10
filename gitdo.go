package main

import (
	"fmt"
	"strings"

	// "os"
	"io/ioutil"
	"os/exec"
	"sync"

	// "log"

	"github.com/go-git/go-git"
)

var commands = []string{
	"gosec -fmt=json -out=gosec.json ./...",
	// "snyk test --print-deps --json",
}

func CheckIfError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func action(directory string) {
	// cd to directory and execute command
	for _, command := range commands {
		var command_args = strings.Split(command, " ")
		cmd := exec.Command(command_args[0], command_args[1:]...)
		cmd.Dir = directory
		out, _ := cmd.Output()
		fmt.Printf("%s\n", out)
	}
}

func clone(repo string, wg *sync.WaitGroup) {
	// TODO: check for directory. if exists, pull, else clone
	// TODO: rename this to clone/sync or something. Then define an operation function to call here

	defer wg.Done()
	// Check for empty string
	if repo == "" {
		fmt.Println("Empty string found...ignoring")
		return
	}

	directory := strings.Split(repo, "/")
	dir := directory[len(directory)-1]
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: repo,
		// Progress: os.Stdout,
		// Auth: &http.BasicAuth{
		// 	Username: username,
		// 	Password: password,
		// },
	})
	if err != nil {
		// do `git pull` if repo already exists
		if fmt.Sprint(err) == "repository already exists" {
			// // We instance a new repository targeting the given path (the .git folder)
			r, err := git.PlainOpen(dir)
			CheckIfError(err)

			// Get the working directory for the repository
			w, err := r.Worktree()
			CheckIfError(err)

			// Pull the latest changes from the origin remote and merge into the current branch
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			CheckIfError(err)
		} else {
			fmt.Println(err, repo)
		}
	}

	action(dir)
}

func main() {

	//eventually becomes func getRepos
	data, err := ioutil.ReadFile("list")
	if err != nil {
		fmt.Println("error reading file", err)
		return
	}

	repos := strings.Split(string(data), "\n")
	// fmt.Println(repos)
	// return

	wg := new(sync.WaitGroup)
	for _, repo := range repos {
		wg.Add(1)
		go clone(repo, wg)
	}
	wg.Wait()
	fmt.Println("done")
}
