package main

import (
	"fmt"
	"os"
	"io"
	"io/ioutil"
	"bufio"
	"strings"
	"flag"
	"encoding/json"
)

type langModel struct {
	Bigram, Trigram map[string]float64
	BiMin, TriMin float64
}

//-------------------------
// Ngramm-Funktionen
//-------------------------

//.Lerne Ngramme
func learnNgrams(sentences []string, n int) (map[string]float64) {
	var ngrams = make(map[string]int)
	for _,a := range sentences {
		ngrams = countNgrams(ngrams, a, n)
	}
	return getNgramProb(ngrams)
}

// Zähle Ngramme
func countNgrams(ngrams map[string]int, sentence string, n int) (map[string]int) {
	words := append([]string{"^"}, strings.Split(sentence, " ")...)
	words = append(words, "$")
	if len(words) < n {
		return ngrams
	}
	for i:=0; i<len(words)-n+1; i++ {
		if noEmpty(words[i:i+n]) {
			ngrams[strings.Join(words[i:i+n], " ")]++
		}
	}
	return ngrams
}

// Bestimme Ngram-Wahrscheinlichkeit
func getNgramProb(ngram map[string]int) (map[string]float64) {
	var count = make(map[string]int)
	for ng, ngc := range ngram {
		count[ng[:strings.LastIndex(ng, " ")]] += ngc
	}
	var out = make(map[string]float64)
	for ng, ngc := range ngram {
		out[ng] = float64(ngc) / float64(count[ng[:strings.LastIndex(ng, " ")]])
	}
	return out
}

//-------------------------
// Haupt-Lernfunktion
//-------------------------

func learnEverything() (langModel) {
	var bi, tri map[string]float64
	sent := readSentences(os.Stdin)
	fmt.Printf("Lerne %d Sätze.\n", len(sent))
	bi = learnNgrams(sent, 2)
	tri = learnNgrams(sent, 3)
	return langModel{bi, tri, getMin(bi), getMin(tri)}
}

//-------------------------
// Auswertung
//-------------------------

// Bestimme die Wahrscheinlichkeit eines Satzes, für gegebenes languageModel
func (model langModel) getSentProb(sent string) (float64) {
	out := 1.0
	sent = strings.ToLower(sent)
	words := append([]string{"^"}, strings.Split(sent, " ")...)
	words = append(words, "$")
	var bigrList []string
	for i:=0; i<len(words)-1; i++ {
		bigrList = append(bigrList, words[i] + " " + words[i+1])
	}
	for _, b := range bigrList {
		if model.Bigram[b] == 0.0 {
			out = out * model.BiMin
		} else {
			out = out * model.Bigram[b]
		}
	}
	return out
}

func sentProbCheck(model langModel) {
	scanner := bufio.NewScanner(os.Stdin)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if line != "" {
			fmt.Printf("%e %s\n", model.getSentProb(line), line)
		}
	}
}

func mlsChecker(model langModel) {
	scanner := bufio.NewScanner(os.Stdin)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if line != "" {
			fmt.Println(model.mostLikelySentence(line))
		}
	}
}

func (model langModel) mostLikelySentence(inp string) (string) {
	var jline, topSent string
	var lineP, topSentP float64
	permutations := HeapsAlg(strings.Split(inp, " "))
	for _,line := range permutations {
		jline = strings.Join(line, " ")
		lineP = model.getSentProb(jline)
		if lineP > topSentP {
			topSent = jline
			topSentP = lineP
		}
	}
	return topSent
}

//-------------------------
// Hilfsfunktionen
//-------------------------

func check(e error) {
	if e != nil {
		fmt.Println("Error: ",e)
		os.Exit(1)
	}
}

func getMin(b map[string]float64) (float64) {
	curMin := 1.0
	for _,i := range b {
		if i < curMin {
			curMin = i
		}
	}
	return curMin
}

func noEmpty(words []string) (bool) {
	for _,w := range words {
		if w == "" {
			return false
		}
	}
	return true
}

// Lese Sätze ein und entferne Sonderzeichen
func readSentences(r io.Reader) ([]string) {
	scanner := bufio.NewScanner(r)
	var out []string
	var tmp string
	for scanner.Scan() {
		tmp = scanner.Text()
		if tmp != "" {
			out = append(out, tmp)
		}
	}
	return out
}

//-------------------------
// Heap's Algorithm
//-------------------------

func generate(n int, a []string, o *[][]string) {
	if n == 1 {
		// Rückgabe
		c := make([]string, len(a))
		copy(c, a)
		*o = append(*o, c)
	} else {
		for i := 0; i < n-1; i++ {
			generate(n-1, a, o)
			if n%2 == 0 {
				a[i], a[n-1] = a[n-1], a[i]
			} else {
				a[0], a[n-1] = a[n-1], a[0]
			}
		}
		generate(n-1, a, o)
	}
}

func HeapsAlg(words []string) ([][]string) {
	var tmp [][]string
	generate(len(words), words, &tmp)
	return tmp
}

//-------------------------
// Main
//-------------------------

func main() {
	// Optionen definieren und einlesen
	var learn = flag.Bool("learn", false, "")
	var verbose = flag.Bool("verbose", false, "")
	flag.Parse()
	var filename = flag.Arg(0)
	if filename == "" {
		fmt.Println("Filename required")
		os.Exit(1)
	}

	if *learn {
		// Bigramme lernen
		model := learnEverything()
		fmt.Printf("%d Bigramme und %d Trigramme eingelesen.\n", len(model.Bigram), len(model.Trigram))
		fmt.Printf("biMin: %e\ntriMin: %e\n", model.BiMin, model.TriMin)

		// Bigramme in json umwandeln
		b, err := json.Marshal(model)
		check(err)

		// json-Data in Datei speichern
		err = ioutil.WriteFile(filename, b, 0644)
		check(err)
	} else {
		// json aus Datei einlesen
		b, err := ioutil.ReadFile(filename)
		check(err)

		// json-Data in Bigramme umwandeln
		var model langModel
		err = json.Unmarshal(b, &model)
		check(err)

		// Wahrscheinlichkeit(en) bestimmen
		if *verbose {
			sentProbCheck(model)
		} else {
			mlsChecker(model)
		}
	}
}
