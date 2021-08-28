[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) [![GoDoc](https://godoc.org/github.com/alexj212/tldr?status.png)](http://godoc.org/github.com/alexj212/tldr)  [![travis](https://travis-ci.org/alexj212/tldr.svg?branch=master)](https://travis-ci.org/alexj212/tldr) [![Go Report Card](https://goreportcard.com/badge/github.com/alexj212/tldr)](https://goreportcard.com/report/github.com/alexj212/tldr)



# tldr
A lightweight tldr client that allows for custom pages.
<br>

# Description
This tldr client, will download and install the initial cache of tldr pages from https://codeload.github.com/tldr-pages/tldr/zip/refs/heads/main/tldr-main.zip.
<br>

 The file will be downloaded and saved locally to  ```~/.cache/tldr/tldr-main.zip```. 
 
 The archive is unpacked to ```~/.cache/tldr/tldr-main/```. 
 
 This implementation allows for local pages to be stored in ```~/.cache/tldr/custom```. 
 <br>
 
 These pages will be available to be displayed along with the official repo pages. Options are provided to allow for hiding the custom or official pages. Options are also provided to update the local cache from the official repo pages. The binary will let you know of the cache of tldr pages is older than 7 days. 

<br>
The same colors and help file parsing mechanism as in the original python script have been replicated.
<br>
<br>


# Hard Fork - Initial Code
The initial code was lifted from the following repos. 
<br>
https://github.com/HardDie/myTldr
<br>
https://github.com/free2k/myTldr



# How to install
Clone repository run ```make install```

This will install the binary in ~/bin/




