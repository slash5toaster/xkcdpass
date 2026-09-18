# XKCD Style passwords

Go program to create xkcd style passwords aka
[https://xkcd.com/936/](https://xkcd.com/936/)

This will attempt to use the system dictionary, if not available it will fall back to using the [EFF large wordlist](https://www.eff.org/document/passphrase-wordlists)

If that is not available, it will generate a random HEX value in place of the words

```bash
xkcdpass -h
Usage of xkcdpass:
  -dict string
        path to dictionary file (default "/usr/share/dict/words")
  -eff
        use the EFF large wordlist (downloaded and cached on first use), overrides --dict
  -number int
        number of digits: max 16 (default 3)
  -special
        Add a special char "@,#,$,%,^,&,+,_,?,~"
  -words int
        number of dictionary words to use (default 2)

```

Example

```bash
$ xkcdpass -eff -special
bullseye302Acronym&

```

The original bash version has been included for history's sake, it has been updated to use the EFF long word list as a fallback, and has an option to use that as the source dictionary.

```bash
xkcdpass.sh -h
Usage: xkcdpass.sh
  -h help 
  -n Number of digits (e.g. -n 3 gives 000-999) defaults to 3, max 10
  -s Adds special characters '(@ # $ % ^ & + _)'
  -e Use the EFF long word list as the dictionary
```
