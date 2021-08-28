[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) [![GoDoc](https://godoc.org/github.com/alexj212/tldr?status.png)](http://godoc.org/github.com/alexj212/tldr)  [![travis](https://travis-ci.org/alexj212/tldr.svg?branch=master)](https://travis-ci.org/alexj212/tldr) [![Go Report Card](https://goreportcard.com/badge/github.com/alexj212/tldr)](https://goreportcard.com/report/github.com/alexj212/tldr)



# tldr
Lightweight tldr client.<br>

This tldr client, will download and install the initial see of tldr pages from https://codeload.github.com/tldr-pages/tldr/zip/refs/heads/main/tldr-main.zip. The file will be downloaded to  ```~/.cache/tldr/tldr-main.zip```. The directory is unpacked to ```~/.cache/tldr/tldr-main/```. This implementation allows for local pages to be stored in ```~/.cache/tldr/custom```. These pages will be available along with the office repo pages. Options are provided to allow for hiding the custom or official pages. Options are also provided to update the local cache from the official repo pages. 
<br>
The same colors and help file parsing mechanism as in the original python script have been used.<br>


# Hard Fork - Initial Code
https://github.com/HardDie/myTldr
https://github.com/free2k/myTldr



# How to install
```go get github.com/alexj212/tldr```
or download from releases binary package
