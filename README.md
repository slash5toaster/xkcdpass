# XKCD Style passwords 

Go program to create xkcd style passwords aka
[https://xkcd.com/936/](https://xkcd.com/936/)

This will attempt to use the system dictionary, if not available it will fall back to using the [EFF large wordlist](https://www.eff.org/document/passphrase-wordlists)
If that is not available, it will generate a random HEX value in place of the words

```
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
```
$ xkcdpass -eff -special
bullseye302Acronym&

```
