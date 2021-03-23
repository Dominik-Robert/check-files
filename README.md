## check-files

A nagios compliant check for different file operations

### Synopsis

With this check you can perform check operations with your files to prove if they are correct in age, size and count

### Options

```
  -c, --critical int       Specify the Critical value. All files which are older than 10 units (default 14)
  -d, --directory string   Specify the directory which contains the files (default "./")
  -h, --help               help for check-files
  -w, --warning int        Specify the Warning value. All files which are older than 10 units (default 10)
```

### SEE ALSO

* [check-files fileAge](check-files_fileAge.md)	 - Check the age of all files in a directory
* [check-files fileCount](check-files_fileCount.md)	 - Count your files in a directory
* [check-files fileSize](check-files_fileSize.md)	 - Check the size of a directory
* [check-files version](check-files_version.md)	 - Print the Version of this check

# TODO
- [ ] Add recursive option
- [ ] Add ignore field to all commands
- [ ] Add Binary and a Release Page
- [ ] Rework Documentation and this README file

