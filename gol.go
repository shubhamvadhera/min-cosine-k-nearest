package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"os"

	"github.com/shopspring/decimal"
)

/* Tuple for document number pairs */
type DocPair struct {
	doc1 int
	doc2 int
}

/* Structs and functions related to sorting a map:
* adopted as is from https://gist.github.com/ikbear/4038654
 */

type sortedMap struct {
	m map[int]float64
	s []int
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[int]float64) []int {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]int, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}

/* End of sort functions and structs */

/* Global Variables */
var records int
var ptr []int
var ind []int
var val []int
var numCosineSim int
var cosims map[DocPair]float64
var normValues map[int]float64
var normVectors map[int]map[int]float64
var nbrslist []int

/* Function to print file statistics like number of rows */
func printFileStats(filename string) {
	if file, err := os.Open(filename); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		numLines := 0
		maxCol := 0
		maxi := 0
		for scanner.Scan() {
			numLines = numLines + 1
			line := scanner.Text()
			arr := strings.Split(line, " ")
			for _, str := range arr {
				if maxCol < len(arr) {
					maxCol = len(arr)
				}
				if i, _ := strconv.Atoi(str); i > maxi {
					maxi, _ = strconv.Atoi(str)
				}
			}
		}
		fmt.Println("number of lines: ", numLines)
		fmt.Println("max Value: ", maxi)
		fmt.Println("max Columns: ", maxCol)
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}

/* Reads the file and creates ptr, ind, val slices */
func readCreateCSR(filename string) {
	if file, err := os.Open(filename); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		ptr = append(ptr, 0)
		for scanner.Scan() {
			records += 1
			arr := strings.Split(strings.Trim(scanner.Text(), " "), " ")
			//fmt.Println(arr)
			l := len(arr)
			//println(l/2)
			ptr = append(ptr, ptr[len(ptr)-1]+l/2)
			for i := 0; i < l; i++ {
				x, _ := strconv.Atoi(arr[i])
				if i%2 == 0 {
					ind = append(ind, x)
				} else {
					val = append(val, x)
				}
			}
		}
		//fmt.Println(ptr)
		//fmt.Println(ind)
		//fmt.Println(val)
		fmt.Println("Number of records: ", records)
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}

func truncate(f float64) float64 {
	f_d := decimal.NewFromFloat(f)
	f_t := f_d.Truncate(5)
	f_r,_ := f_t.Float64()
	return f_r
}

/* Calculates Cosine Similarity by dot product
* algorithm adopted from findsim by Prof. David C. Anastasiu
 */
func cosineSimDot(docNum1 int, docNum2 int) float64 {
	numCosineSim += 1
	if docNum1 == docNum2 {
		return 1.0
	}
	// caching of cosine similarities
	docPair := DocPair{doc1: docNum1, doc2: docNum2}
	if v, ok := cosims[docPair]; ok {
		return v
	}
	docPairInv := DocPair{doc1: docNum2, doc2: docNum1}
	if v, ok := cosims[docPairInv]; ok {
		return v
	}

	fr1 := ptr[docNum1-1]
	to1 := ptr[docNum1]
	fr2 := ptr[docNum2-1]
	to2 := ptr[docNum2]
	i := fr1
	j := fr2
	dp := 0.0
	normDoc1 := 0.0
	normDoc2 := 0.0
	simi := 0.0
	for i < to1 && j < to2 {
		if i == to1 { //i reaches end
			normDoc2 += float64(val[j]*val[j])
			j += 1
		} else if j == to2 { // j reaches end
			normDoc1 += float64(val[i]*val[i])
			i += 1
		} else if ind[i] > ind[j] {
			normDoc2 += float64(val[j]*val[j])
			j += 1
		} else if ind[i] < ind[j] {
			normDoc1 += float64(val[i]*val[i])
			i += 1
		} else {
			dp += float64(val[i]*val[j])
			normDoc1 += float64(val[i]*val[i])
			normDoc2 += float64(val[j]*val[j])
			i += 1
			j += 1
		}
	}
	if dp > 0.0 {
		simi = dp / math.Sqrt(normDoc1*normDoc2)
	}
	/*cosims[docPair] = simi*/
	docPair = DocPair{doc1: docNum2, doc2:docNum1}
	cosims[docPair] = simi
	return simi
}

/* Calculates and caches norm of a document */
func normValue(docnumber int) float64 {
	if v, ok := normValues[docnumber]; ok {
		return v
	}
	fr := ptr[docnumber-1]
	to := ptr[docnumber]
	nval := 0.0
	for i := fr; i < to; i++ {
		nval += truncate(float64(val[i] * val[i]))
	}
	nval = truncate(math.Sqrt(nval))
	normValues[docnumber] = nval
	return nval
}

/* Calculates and caches norm vector of a document */
func normVector(docnumber int) map[int]float64 {
	if v, ok := normVectors[docnumber]; ok {
		return v
	}
	fr := ptr[docnumber-1]
	to := ptr[docnumber]
	dic := make(map[int]float64)
	for i := fr; i < to; i++ {
		dic[ind[i]] = truncate(float64(val[i]) / normValue(docnumber))
	}
	normVectors[docnumber] = dic
	return dic

}

/* Calculates Cosine Similarity by normalized vectors
* slower than cosineSimDot in Golang
 */
func cosineSimNorm(docNum1 int, docNum2 int) float64 {
	numCosineSim += 1
	if docNum1 == docNum2 {
		return 1.0
	}
	docPair := DocPair{doc1: docNum1, doc2: docNum2}
	if v, ok := cosims[docPair]; ok {
		return v
	}
	docPairInv := DocPair{doc1: docNum2, doc2: docNum1}
	if v, ok := cosims[docPairInv]; ok {
		return v
	}
	doc1norm := normVector(docNum1)
	doc2norm := normVector(docNum2)
	simi := 0.0
	for k, _ := range doc1norm {
		if _, ok := doc2norm[k]; ok {
			simi += truncate(doc1norm[k] * doc2norm[k])
		}
	}
	cosims[docPair] = simi
	return simi
}

/* Finds the probable neighbours as per IdxJoin */
func knnIdx() {
	nbrs := make(map[int]bool)

	for docnum := 1; docnum < records+1; docnum++ {
		if _, ok := nbrs[docnum]; ok {
			//print("nbr hit")
			continue
		}
		fr := ptr[docnum-1]
		to := ptr[docnum]
		l := len(ptr)
		for d := fr; d < to; d++ {
			for x := 1; x < l; x++ {
				if _, ok := nbrs[x]; ok {
					continue
				}
				n := ptr[x-1]
				if n == fr {
					continue
				}
				for ind[n] <= ind[d] && n < ptr[x] {
					if ind[n] == ind[d] {
						nbrs[x] = true
						nbrs[docnum] = true
						break
					}
					n += 1
				}
			}
		}
	}

	for k := range nbrs {
		nbrslist = append(nbrslist, k)
	}
	println("Neighbours created")
	l := len(nbrslist)
	println("Number of probable neighbours: ", l)
}

/* Finds similarity within probable neighbours */
func knn(docnum int, eps float64, k int) {
	klist := make(map[int]float64)
	//fmt.Println("nbrslist: ", nbrslist)
	l := len(nbrslist)
	//println("Number of probable neighbours: ", l)
	for i := 0; i < l; i++ {
		if nbrslist[i] != docnum {
			simi := cosineSimDot(docnum, nbrslist[i])
			if simi >= eps {
				// docPair := DocPair{doc1: nbrslist[i], doc2: nbrslist[j]}
				klist[nbrslist[i]] = simi
			}
		}
	}
	//println("Number of computed similarities: ", numCosineSim)
	if len(klist) < k {
		k = len(klist)
	}
	sklist := sortedKeys(klist)
	print(docnum, ": ")
	for i:=0; i<k; i++ {
		fmt.Print(sklist[i], " ", klist[sklist[i]], " ")
	}
	println()
	//fmt.Println("klist: ", sortedKeys(klist)[:k])
}

func knnAll(eps float64, k int) {
	knnIdx()
	for i:=1; i<records+1; i++ {
		knn(i,eps,k)
	}
}

func main() {
	println("********************************************************************************")
	println("findsim-golang (0.0.1), vInfo: [initial version]")
	println("mode: ij, iFile: , oFile: , k: , eps: ")
	println("********************************************************************************")
	tStart := time.Now()
	//printFileStats("wiki10.csr")
	readCreateCSR("wiki1k.csr")
	//knnIdx()
	cosims = make(map[DocPair]float64)
	normValues = make(map[int]float64)
	normVectors = make(map[int]map[int]float64)
	for i := 1; i <= records; i++ {
		println(i)
		for j := i; j <= records; j++ {
			cosineSimDot(i, j)
		}
	}
	//knnIdx()
	//knnAll(0.8,10)
	//fmt.Println(len(nbrs))
	//fmt.Println(cosims)
	//fmt.Println(cosineSimDot(258,259))
	//fmt.Println(cosineSimDot(258,259))
	// fmt.Println(normValue(3))
	// fmt.Println(normValue(3))
	// fmt.Println(normVector(1))
	// fmt.Println(normVector(1))

	fmt.Println("Program time: ", time.Since(tStart))
}
