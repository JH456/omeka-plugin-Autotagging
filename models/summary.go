package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/jdkato/prose/chunk"

	"github.com/jdkato/prose/tag"
	"github.com/jdkato/prose/tokenize"
)

func GetStringCoOccurrences(sentences [][]string, n int) map[string]map[string]int {
	result := make(map[string]map[string]int)
	for _, sentence := range sentences {
		for i, word := range sentence {
			for j := i - n; j < i+n+1; j++ {
				if j >= 0 && j != i && j < len(sentence) {
					if _, ok := result[word]; ok {
						result[word][sentence[j]] += 1
					} else {
						result[word] = make(map[string]int)
						result[word][sentence[j]] = 1
					}
				}
			}
		}
	}
	return result
}

func GetCoOccurrences(sentences [][]int, n int) map[int]map[int]int {
	result := make(map[int]map[int]int)
	for _, sentence := range sentences {
		for i, word := range sentence {
			for j := i - n; j < i+n+1; j++ {
				if j >= 0 && j != i && j < len(sentence) {
					if _, ok := result[word]; ok {
						result[word][sentence[j]] += 1
					} else {
						result[word] = make(map[int]int)
						result[word][sentence[j]] = 1
					}
				}
			}
		}
	}
	return result
}

func SentencesToIndices(sentences [][]string) ([][]int, map[string]int, map[int]string) {
	converted := make([][]int, len(sentences))
	for i, sentence := range sentences {
		converted[i] = make([]int, len(sentence))
	}
	forwardMapping := make(map[string]int)
	reverseMapping := make(map[int]string)
	count := 0
	for i, sentence := range sentences {
		for j, word := range sentence {
			if index, ok := forwardMapping[word]; ok {
				converted[i][j] = index
			} else {
				converted[i][j] = count
				reverseMapping[count] = word
				forwardMapping[word] = count
				count += 1
			}
		}
	}
	return converted, forwardMapping, reverseMapping
}

func IndicesToSentences(sentences [][]int, reverse map[int]string) [][]string {
	converted := make([][]string, len(sentences))
	for i, sentence := range sentences {
		converted[i] = make([]string, len(sentence))
	}
	for i, sentence := range sentences {
		for j, word := range sentence {
			converted[i][j] = reverse[word]
		}
	}
	return converted
}

func GetSparseApproxPageRank(coOccurences map[string]map[string]int, damping, threshold float64) map[string]float64 {
	total := len(coOccurences)
	result := make(map[string]float64)
	for i := range coOccurences {
		result[i] = 1 / float64(total)
	}
	prev := make(map[string]float64)
	for i := range coOccurences {
		prev[i] = 0
	}

	shouldContinue := true
	for shouldContinue {
		shouldContinue = false
		for i := range result {
			if math.Abs(result[i]-prev[i]) > threshold {
				shouldContinue = true
				continue
			}
		}

		for i := range result {
			prev[i] = result[i]
		}

		if !shouldContinue {
			return result
		}

		for i, neighbors := range coOccurences {
			cur := 0.0
			// For directed, neighbors should be incoming, not outgoing
			for j, _ := range neighbors {
				cur += prev[j] / float64(len(coOccurences[j]))
			}
			result[i] = 1 - damping + damping*cur
		}
	}

	return result
}

func GetKeywords(text string, threshold float64) []string {
	sentences := getCleanSentences(text)
	tokenizer := tokenize.NewTreebankWordTokenizer()
	tagger := tag.NewPerceptronTagger()
	words := make([][]string, len(sentences))
	for i, s := range sentences {
		tokenized := tokenizer.Tokenize(s)
		cur := make([]string, 0)
		for _, tok := range tagger.Tag(tokenized) {
			switch tok.Tag {
			case "NN", "NNS", "NNP", "NNPS", "JJ", "JJR", "JJS":
				cur = append(cur, tok.Text)
			}
		}
		words[i] = cur
	}

	coOccurrences := GetStringCoOccurrences(words, 1)
	result := GetSparseApproxPageRank(coOccurrences, 0.85, 0.000001)

	keywords := make([]string, 0)
	for word, confidence := range result {
		if confidence > threshold {
			keywords = append(keywords, word)
			// fmt.Printf("Selected keyword %v with confidence %v\n", word, confidence)
		}
	}

	return keywords
}

func GetNamedEntities(text string) []string {
	entities := make([]string, 0)
	words := tokenize.TextToWords(text)
	regex := chunk.TreebankNamedEntities

	tagger := tag.NewPerceptronTagger()
	for _, entity := range chunk.Chunk(tagger.Tag(words), regex) {
		entities = append(entities, entity)
	}
	return entities
}

type Pair struct {
	Key   string
	Value float64
}
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func sorted(m map[string]float64) PairList {
	pairs := make(PairList, len(m))
	i := 0
	for k, v := range m {
		pairs[i] = Pair{k, v}
		i++
	}

	sort.Sort(pairs)
	result := sort.Reverse(pairs)
	fmt.Println(result)

	return pairs
}

func getText() string {
	bytes, err := ioutil.ReadFile("/home/kpberry/Desktop/iada_pdfs/test/ahc_CAR_015_008_025_011-ocr.txt")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func getCleanText(text string, lower bool) string {
	space := regexp.MustCompile("\\s+")
	text = space.ReplaceAllString(text, " ")
	if lower {
		text = strings.ToLower(text)
	}
	nonLexical := regexp.MustCompile(`[^a-zA-Z0-9.?!'" ]+`)
	text = nonLexical.ReplaceAllString(text, "")
	text = stopwords.CleanString(text, "en", true)
	return text
}

func getCleanSentences(text string) []string {
	text = getCleanText(text, true)
	words := tokenize.NewPunktSentenceTokenizer().Tokenize(text)
	for i, sentence := range words {
		words[i] = strings.Trim(sentence, " ")
	}
	return words
}

func __main() {
	// TODO use tfidf for stuff
	sentences := getCleanSentences(getText())
	tokenizer := tokenize.NewTreebankWordTokenizer()
	tagger := tag.NewPerceptronTagger()
	words := make([][]string, len(sentences))
	count := 0
	for i, s := range sentences {
		tokenized := tokenizer.Tokenize(s)
		cur := make([]string, 0)
		for _, tok := range tagger.Tag(tokenized) {
			switch tok.Tag {
			case "NN", "NNS", "NNP", "NNPS", "JJ", "JJR", "JJS":
				cur = append(cur, tok.Text)
			}
		}
		words[i] = cur
	}

	coOccurrences := GetStringCoOccurrences(words, 1)
	result := GetSparseApproxPageRank(coOccurrences, 0.85, 0.000001)
	fmt.Println(sorted(result))

	fmt.Println(count)
}
