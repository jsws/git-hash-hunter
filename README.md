# Git Hash Hunter 🏹  
![Tests](https://github.com/jsws/git-hash-hunter/workflows/Tests/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/jsws/git-hash-hunter)](https://goreportcard.com/report/github.com/jsws/git-hash-hunter)

A tool that allows you to set your commit hashes to whatever you like, inspired by [lucky-commit](https://github.com/not-an-aardvark/lucky-commit).

## Install
```console
$ # Ensure $GOBIN is set and in $PATH.
$ go install 
```


## Usage
Run `githh` inside the git repo you want to change the commit hash.

```console

$ githh
Usage: githh <HASH> 
e.g.
     githh 0000000
```


```console

$ githh 0000000
Original hash: 6d1cf69c87ab8a317b80ba599486aa73946edbfe
Original commit message: "Add tests"
New hash: 000000091177ebe2f4159542289e6d7aa601083b
New Message: "Add tests 
 
                                 
          
"
Took 11.035s @ 2370584 H/s
```

## Explained
Git Hash Hunter will add permutations of the whitespace characters `\n`, `\t` and `<space>` to your commit message and hash the commit until the desired hash is found. Once a sequence of whitespaces is found the commit will be amended.

The given hash can be of any length but the longer the hash the longer `githh` will take to find a message. 
