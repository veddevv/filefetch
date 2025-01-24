package main

import (
	"fmt"
	"os"
	"strings"
	"golang.org/x/text/message"
	"golang.org/x/text/language"
)

import flag "github.com/cornfeedhobo/pflag"

var hlsize = 4
var hlmodi = 13
var hlperms = 11
var dirsize int64 = 0
var dircount = 0

func sep(num int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", num)
}

func PrintPA(str string, pad int, ansi string) {
	if pad > 0 {
		fmt.Print(ansi + str + "\033[0m" + strings.Repeat(" ", pad - len(str)))
	} else {
		fmt.Print(ansi + str + "\033[0m")
	}
}

func printFiles(dirOnly bool, longMode bool, lastModifiedEnable bool, permsEnable bool, dateFormat string) {
	files, _ := os.ReadDir(".")
	for _, file := range files {
		fileinfo, _ := os.Stat(file.Name())
		if dirOnly && file.IsDir() || !dirOnly && !file.IsDir() {
			if longMode {
				if dirOnly { PrintPA("Directory", hlsize, "\033[36m") } else { PrintPA(sep(fileinfo.Size()) + "B", hlsize, "\033[36m") } 
				if lastModifiedEnable { PrintPA(fileinfo.ModTime().Format(dateFormat), hlmodi, "\033[33m") }
				if permsEnable { if dirOnly { PrintPA(strings.Replace(fileinfo.Mode().Perm().String(), "-", "d", 1), hlperms, "\033[31m") } else { PrintPA(fileinfo.Mode().Perm().String(), hlperms, "\033[31m") } } 
				if dirOnly { PrintPA(file.Name(), 0, "\033[34m") } else { PrintPA(file.Name(), 0, "\033[32m") }  
				fmt.Println("")
			} else {
				if dirOnly {
					PrintPA(file.Name(), 0, "\033[34m")
					fmt.Print(" ")
				} else if !file.IsDir() {
					PrintPA(file.Name(), 0, "\033[32m")
					fmt.Print(" ")
				}
			}
		}
	}	
}

func main() {
	longMode := flag.BoolP("long", "l", false, "Use Long Mode.")
	dirFirst := flag.BoolP("directoriesfirst", "d", false, "List Directories before Files.")

	lastModifiedEnable := flag.BoolP("lastmodified", "m", false, "Enable the Last Modified Section on Long Mode.")
	permsEnable := flag.BoolP("permissions", "p", false, "Enable the Perms Section on Long Mode.")
	
	dateFormat := flag.StringP("format", "f", "02/01/2006 15:04:05.000", "Date `format` for Last Modified, if enabled.")

	flag.Parse()

	files, _ := os.ReadDir(".")
	for _, file := range files {
		fileinfo, _ := os.Stat(file.Name())
		if len(sep(fileinfo.Size()) + "B") > hlsize { hlsize = len(sep(fileinfo.Size()) + "B") }
		if len(fileinfo.ModTime().Format(*dateFormat)) > hlmodi { hlmodi = len(fileinfo.ModTime().Format(*dateFormat)) }
		if len(fileinfo.Mode().Perm().String()) > hlperms { hlperms = len(fileinfo.Mode().Perm().String()) }
	
		if file.IsDir() { dircount += 1 } else { dirsize += fileinfo.Size() } 
	}

	hlsize += 1
	hlmodi += 1
	hlperms += 1

	if *longMode {
		PrintPA("Size", hlsize, "\033[1;4m")
		if *lastModifiedEnable { PrintPA("Last Modified", hlmodi, "\033[1;4m") }
		if *permsEnable { PrintPA("Permissions", hlperms, "\033[1;4m") }
		PrintPA("Name", hlsize, "\033[1;4m")
		fmt.Println("\033[0m")
	}

	if *dirFirst {
		printFiles(true, *longMode, *lastModifiedEnable, *permsEnable, *dateFormat)	
		printFiles(false, *longMode, *lastModifiedEnable, *permsEnable, *dateFormat)	
	} else {
		printFiles(false, *longMode, *lastModifiedEnable, *permsEnable, *dateFormat)	
		printFiles(true, *longMode, *lastModifiedEnable, *permsEnable, *dateFormat)	
	}

	plural := "ies"
	if dircount == 1 { plural = "y" }
	if !*longMode { fmt.Println() }
	fmt.Println("Fetched \033[36;1m" + sep(int64(dirsize)) + "B \033[0mof Files and \033[34;1m" + sep(int64(dircount))+ " \033[0mDirector" + plural + ".")
}
