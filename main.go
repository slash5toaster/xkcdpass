// Go program to create xkcd style passwords aka
// https://xkcd.com/936/

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const effWordlistURL = "https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt"

// ----------------------------------------------------------------------------\\
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

// ----------------------------------------------------------------------------\\
func countLinesInFile(fileName string) (int, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return 0, err
	}

	buf := make([]byte, 1024)
	lines := 0

	for {
		readBytes, err := file.Read(buf)

		if err != nil {
			if readBytes == 0 && err == io.EOF {
				err = nil
			}
			return lines, err
		}

		lines += bytes.Count(buf[:readBytes], []byte{'\n'})
	}

	return lines, nil
}

// ----------------------------------------------------------------------------\\
func getSpecialChar() string {
	s := "@,#,$,%,^,&,+,_,?,~"
	spchar := strings.Split(s, ",")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomIndex := r.Intn(len(spchar))

	pick := spchar[randomIndex]

	return pick
}

// ----------------------------------------------------------------------------\\
func getRandomNumber(numPwr int) int {
	// get a random number with digits specified by numPwr
	// e.g.  numPwr = 4 gives numbers between 1000 and 9999
	var pwNum int = 1
	var maxNum int = int(math.Pow10(numPwr))

	// get a pseudorandom seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// make sure we have a number with correct digits
	for true {
		pwNum = r.Intn(maxNum)
		if pwNum >= maxNum/10 && pwNum < maxNum {
			break
		}
	}
	return pwNum
}

// ----------------------------------------------------------------------------\\
func getRandomWord(dictFile string) string {
	// read the file
	var myWord, chosenWord = "nada", ""
	var maxLine, err = countLinesInFile(dictFile)
	var chosenLine int = 1 // default to first line

	// get a random number
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	f, err := os.Open(dictFile)

	// if the file open fails, generate a hash
	if err == nil {
		// pick a line
		chosenLine = r.Intn(maxLine)

		defer f.Close()

		// scan the file until we get to our desired line.
		scanner := bufio.NewScanner(f)
		for lineNumb := 0; lineNumb <= chosenLine; lineNumb++ {
			scanner.Scan()
			if lineNumb == chosenLine {
				chosenWord = scanner.Text()
				break
			}
		}
		// remove punctuation from string
		reg := regexp.MustCompile("[^a-zA-Z0-9]+")
		myWord = reg.ReplaceAllString(chosenWord, "")
	} else {
		// if the file open fails for any reason, generate a random number as hex
		myWord = strconv.FormatInt(int64(getRandomNumber(18)), 16)
	}
	return myWord
}

// ----------------------------------------------------------------------------\\
func downloadEFFWordlist(destPath string) error {
	resp, err := http.Get(effWordlistURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download EFF wordlist: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// the EFF wordlist is tab-separated "diceroll\tword" per line;
	// keep only the word so getRandomWord's punctuation stripping doesn't mangle it
	scanner := bufio.NewScanner(resp.Body)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		word := fields[len(fields)-1]
		if word == "" {
			continue
		}
		fmt.Fprintln(writer, word)
	}
	return scanner.Err()
}

// ----------------------------------------------------------------------------\\
func getEFFDictPath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}

	xkcdCacheDir := filepath.Join(cacheDir, "xkcdpass")
	if err := os.MkdirAll(xkcdCacheDir, 0o755); err != nil {
		return "", err
	}

	destPath := filepath.Join(xkcdCacheDir, "eff_large_wordlist.txt")

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		if err := downloadEFFWordlist(destPath); err != nil {
			return "", err
		}
	}

	return destPath, nil
}

// ----------------------------------------------------------------------------\\
func main() {
	var sChar string = ""
	var maxDigit int = 16
	var defaultDigit = 3
	var defaultWords = 2

	numPtr := flag.Int("number", defaultDigit, "number of digits: max "+strconv.Itoa(maxDigit))
	specialPtr := flag.Bool("special", false, "Add a special char")
	wordsPtr := flag.Int("words", defaultWords, "number of dictionary words to use")
	dictPtr := flag.String("dict", "/usr/share/dict/words", "path to dictionary file")
	effPtr := flag.Bool("eff", false, "use the EFF large wordlist (downloaded and cached on first use), overrides --dict")

	flag.Parse()

	if *effPtr {
		effDictPath, err := getEFFDictPath()
		checkError(err)
		*dictPtr = effDictPath
	}

	// catch if exponent overruns the int max (*must* be less than 18)
	if *numPtr > maxDigit {
		*numPtr = maxDigit
		fmt.Fprintln(os.Stderr, "--number set to max of "+strconv.Itoa(maxDigit))
	}

	if *numPtr <= 0 {
		*numPtr = defaultDigit
		fmt.Fprintln(os.Stderr, "--number set to default of "+strconv.Itoa(defaultDigit))
	}

	if *wordsPtr <= 0 {
		*wordsPtr = defaultWords
		fmt.Fprintln(os.Stderr, "--words set to default of "+strconv.Itoa(defaultWords))
	}

	if *specialPtr {
		sChar = getSpecialChar()
	}

	// first word stays lowercase, remaining words are title-cased
	var sb strings.Builder
	sb.WriteString(strings.ToLower(getRandomWord(*dictPtr)))
	sb.WriteString(strconv.Itoa(getRandomNumber(*numPtr)))
	for i := 1; i < *wordsPtr; i++ {
		sb.WriteString(strings.Title(getRandomWord(*dictPtr)))
	}
	sb.WriteString(sChar)

	// print a generated password
	fmt.Println(sb.String())
}

// End of file, if this is missing the file is truncated
///////////////////////////////////////////////////////////////////////////////
