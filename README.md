
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

##Apple MacOS 

Apple blocks downloaded files from executing steps to fix:

1. Open a terminal session
2. Go to where you downloaded the file cd ~/Downloads
3. Run this command "chmod 700 [filename]" replacing the filename
4. Run this command "xattr -rd com.apple.quarantine [filename]" 
5. You should now be able to run the application without an error