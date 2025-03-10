package main

import (
	"fmt"
	"os"
	"strings"
	"golang.org/x/text/message"
	"golang.org/x/text/language"
	"github.com/cornfeedhobo/pflag"
)

var (
	hlsize   = 4
	hlmodi   = 13
	hlperms  = 11
	dirsize  int64
	dircount int
)

func sep(num int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", num)
}

func printPA(str string, pad int, ansi string) {
	if pad > 0 {
		fmt.Print(ansi + str + "\033[0m" + strings.Repeat(" ", pad-len(str)))
	} else {
		fmt.Print(ansi + str + "\033[0m")
	}
}

func printFiles(dirOnly, longMode, lastModifiedEnable, permsEnable bool, dateFormat string) {
	files, _ := os.ReadDir(".")
	for _, file := range files {
		fileinfo, _ := os.Stat(file.Name())
		if dirOnly && file.IsDir() || !dirOnly && !file.IsDir() {
			if longMode {
				if dirOnly {
					printPA("Directory", hlsize, "\033[36m")
				} else {
					printPA(sep(fileinfo.Size())+"B", hlsize, "\033[36m")
				}
				if lastModifiedEnable {
					printPA(fileinfo.ModTime().Format(dateFormat), hlmodi, "\033[33m")
				}
				if permsEnable {
					if dirOnly {
						printPA(strings.Replace(fileinfo.Mode().Perm().String(), "-", "d", 1), hlperms, "\033[31m")
					} else {
						printPA(fileinfo.Mode().Perm().String(), hlperms, "\033[31m")
					}
				}
				if dirOnly {
					printPA(file.Name(), 0, "\033[34m")
				} else {
					printPA(file.Name(), 0, "\033[32m")
				}
				fmt.Println("")
			} else {
				if dirOnly {
					printPA(file.Name(), 0, "\033[34m")
					fmt.Print(" ")
				} else if !file.IsDir() {
					printPA(file.Name(), 0, "\033[32m")
					fmt.Print(" ")
				}
			}
		}
	}
}

func main() {
	longMode := pflag.BoolP("long", "l", false, "Use Long Mode.")
	dirFirst := pflag.BoolP("directoriesfirst", "d", false, "List Directories before Files.")
	lastModifiedEnable := pflag.BoolP("lastmodified", "m", false, "Enable the Last Modified Section on Long Mode.")
	permsEnable := pflag.BoolP("permissions", "p", false, "Enable the Perms Section on Long Mode.")
	dateFormat := pflag.StringP("format", "f", "02/01/2006 15:04:05.000", "Date `format` for Last Modified, if enabled.")
	pflag.Parse()

	files, _ := os.ReadDir(".")
	for _, file := range files {
		fileinfo, _ := os.Stat(file.Name())
		if len(sep(fileinfo.Size())+"B") > hlsize {
			hlsize = len(sep(fileinfo.Size()) + "B")
		}
		if len(fileinfo.ModTime().Format(*dateFormat)) > hlmodi {
			hlmodi = len(fileinfo.ModTime().Format(*dateFormat))
		}
		if len(fileinfo.Mode().Perm().String()) > hlperms {
			hlperms = len(fileinfo.Mode().Perm().String())
		}

		if file.IsDir() {
			dircount++
		} else {
			dirsize += fileinfo.Size()
		}
	}

	hlsize++
	hlmodi++
	hlperms++

	if *longMode {
		printPA("Size", hlsize, "\033[1;4m")
		if *lastModifiedEnable {
			printPA("Last Modified", hlmodi, "\033[1;4m")
		}
		if *permsEnable {
			printPA("Permissions", hlperms, "\033[1;4m")
		}
		printPA("Name", hlsize, "\033[1;4m")
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
	if dircount == 1 {
		plural = "y"
	}
	if !*longMode {
		fmt.Println()
	}
	fmt.Printf("Fetched \033[36;1m%sB \033[0mof Files and \033[34;1m%s \033[0mDirector%s.\n", sep(dirsize), sep(int64(dircount)), plural)
}