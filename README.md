# FileFetch
An alternative to ls written in golang. It's basically a directory reader with a simple function: **Get all file names.** But with some customizability, you can do *more* than get file names.

## Extra Functions
```
  -d, --directoriesfirst   List Directories before Files.
  -f, --format format      Date format for Last Modified, if enabled. (default "02/01/2006 15:04:05.000")
  -m, --lastmodified       Enable the Last Modified Section on Long Mode.
  -l, --long               Use Long Mode.
  -p, --permissions        Enable the Perms Section on Long Mode.
```