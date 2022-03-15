
<p align="center">
  <img alt="golangci-lint logo" src="https://www3.blacklistalliance.com/static/media/login-logo.a109cbd2.png" height="75" />
  <h3 align="center">Scrub files from the command line.</h3>
</p>

## Badges

[//]: # (![Build Status]&#40;https://github.com/blacklistalliance/scrubfile/workflows/CI/badge.svg&#41;)
[![License](https://img.shields.io/github/license/blacklistalliance/scrubfile)](/LICENSE)
[![Release](https://img.shields.io/github/release/blacklistalliance/scrubfile.svg)](https://github.com/blacklistalliance/scrubfile/releases/latest)
[![GitHub Releases Stats of golangci-lint](https://img.shields.io/github/downloads/blacklistalliance/scrubfile/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=blacklistalliance&repository=scrubfile)

---

`scrubfile` is a utility to scrub files inplace against the Blacklist Litgation firewall.  You must have an account at [Blacklist Alliance.](https://www.blacklistalliance.com)

---
## Scrub file requirements
This application can take in a file with just numbers one per line.  Or you can have a csv formatted file for example
<pre>
Name,Phone,Notes
John Smith,2223334444,good lead
</pre>
You would scrub this file with the following options:
<pre>
-apikey [your api key] -hasheader=true -colnum=2 -file [your filename]
</pre>
The clean file will have the same format as the original file with all the bad leads 
taken out.  There are some additional options for downloading additional data.

---

## Downloading and Running

1. Click releases on the right 
2. Download the file that matches your computer architecture and operating system
   1. For windows choose scrubfile-windows-4.0-amd64.exe
   2. For Apple M1 and newer choose scrubfile-darwin-10.16-arm64
   3. For Apple Intel Macs choose scrubfile-darwin-10.16-amd64
   4. For linux choose scrubfile-linux-amd64
3. This is a command line application which means you need to use Terminal or CMD
4. Go to the directory where the executable is
5. If its Apple follow the steps below
6. Type the name of the file you downloaded and press enter
7. It will return a list of options you will need to append these to the name of the file
8. Example scrubfile-darwin-10.16-arm64 -apikey XXXX -file hotleads.csv -hasheader
9. After it is done you will have a file named hotleads_clean.csv these are all the clean numbers from your original file

---

## Apple MacOS Issues 

Apple blocks downloaded files from executing steps to fix:

1. Open a terminal session
2. Go to where you downloaded the file cd ~/Downloads
3. Run this command "chmod 700 [filename]" replacing the filename
4. Run this command "xattr -rd com.apple.quarantine [filename]" 
5. You should now be able to run the application without an error

---

##Command line options
<pre>
  -apikey string
    	API Key put between single quotes (REQUIRED)
  -carrier
    	Include carrier data
  -colnum string
    	Column number to use (default "1")
  -federaldnc
    	Include federal DNC
  -file string
    	File name to scrub (REQUIRED)
  -hasheader
    	Does the file have a header
  -includefeeds
    	Include carrier feeds
  -invalid
    	Include invalid data
  -noCarrier
    	Include no carrier
  -splitchar string
    	What character do we split on wrap in single quotes (default ",")
  -wireless
    	Include wireless
</pre>