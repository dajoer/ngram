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
	"math"
)

type langModel struct {
	Unigram, Bigram, Trigram map[string]int
	WordCount int
	Lambda1, Lambda2, Lambda3 float64
	Alpha float64
}

//-------------------------
// Ngramm-Funktionen
//-------------------------

//.Lerne Ngramme
func learnNgrams(sentences []string, n int) (map[string]int) {
	var ngrams = make(map[string]int)
	for _,a := range sentences {
		ngrams = countNgrams(ngrams, a, n)
	}
	return ngrams
}

// Zähle Ngramme
func countNgrams(ngrams map[string]int, sentence string, n int) (map[string]int) {
	words := strings.Split(sentence, " ")
	if n > 1 {
		words = append([]string{"^"}, words...)
		words = append(words, "$")
	}
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
		out[ng] = math.Log(float64(ngc) / float64(count[ng[:strings.LastIndex(ng, " ")]]))
	}
	return out
}

func getUnigramProb(model langModel, unigr string) float64 {
	ugc := float64(model.Unigram[unigr])
	wc := float64(model.WordCount)
	return math.Log((ugc + model.Alpha) / (wc + model.Alpha))
}

func getBigramProb(model langModel, bigr string) float64 {
	ugc := float64(model.Unigram[bigr[:strings.LastIndex(bigr, " ")]])
	bgc := float64(model.Bigram[bigr])
	return math.Log((bgc + model.Alpha) / (ugc + model.Alpha))
}

func getTrigramProb(model langModel, trigr string) float64 {
	bgc := float64(model.Bigram[trigr[:strings.LastIndex(trigr, " ")]])
	tgc := float64(model.Trigram[trigr])
	return math.Log((tgc + model.Alpha) / (bgc + model.Alpha))
}

//-------------------------
// Haupt-Lernfunktion
//-------------------------

func learnEverything() (langModel) {
	var uni, bi, tri map[string]int
	sent := readSentences(os.Stdin)
	fmt.Printf("Lerne %d Sätze.\n", len(sent))
	uni = learnNgrams(sent, 1)
	bi = learnNgrams(sent, 2)
	tri = learnNgrams(sent, 3)
	count := 0
	for _,i := range uni {
		count += i
	}
	return langModel{uni, bi, tri, count, 0.0, 1.0, 0.0, 0.001}
}

//-------------------------
// Auswertung
//-------------------------

// Bestimme die Wahrscheinlichkeit eines Satzes, für gegebenes languageModel
func (model langModel) getSentProb(sent string) (float64) {
	out := 1.0
	sent = strings.ToLower(sent)
	words := append([]string{"^", "^"}, strings.Split(sent, " ")...)
	words = append(words, "$")
	var trigrList [][]string
	for i:=0; i<len(words)-2; i++ {
		trigrList = append(trigrList, []string{words[i], words[i+1], words[i+1]})
	}
	for _, b := range trigrList {
		out = out + getInterpTrigramProb(model, b)
	}
	return out
}

func getInterpTrigramProb(model langModel, trig []string) float64 {
	if len(trig) < 3 {
		return 1.0
	}
	ugp := model.Lambda1 * getUnigramProb(model, trig[0])
	bgp := model.Lambda2 * getBigramProb(model, strings.Join(trig[0:2], " "))
	tgp := model.Lambda3 * getTrigramProb(model, strings.Join(trig[0:3], " "))
	return ugp + bgp + tgp
}

func sentProbCheck(model langModel, v bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if line != "" {
			output(model.getSentProb(line), line, v)
		}
	}
}

func mostLikelySentence(model langModel, v bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var line, topSent string
	var lineP, topSentP float64
	for scanner.Scan() {
		line = scanner.Text()
		if line != "" {
			lineP = model.getSentProb(line)
			if lineP > topSentP || topSentP == 0.0 {
				topSent = line
				topSentP = lineP
			}
		}
	}
	if v {
		fmt.Printf("%e\t%s\n", topSentP, topSent)
	} else {
		fmt.Printf("%s\n", topSent)
	}
}

//-------------------------
// Hilfsfunktionen
//-------------------------

func output(p float64, s string, v bool) {
	if v {
		fmt.Printf("%e\t%s\n", p, s)
	} else {
		fmt.Printf("%e\n", p)
	}
}

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
// Main
//-------------------------

func main() {
	// Optionen definieren und einlesen
	var learn = flag.Bool("learn", false, "")
	var verbose = flag.Bool("v", false, "")
	var best = flag.Bool("b",false,"")
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
		if *best {
			mostLikelySentence(model, *verbose)
		} else {
			sentProbCheck(model, *verbose)
		}
	}
}
