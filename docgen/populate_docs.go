package docgen

import (
	"fmt"
	"github.com/couchbase/indexing/secondary/common"
	"github.com/couchbase/indexing/secondary/logging"
	"github.com/couchbase/indexing/secondary/tests/framework/kvutility"
	"github.com/couchbase/indexing/secondary/tools/ProjectorPerf/docgen/stats"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"
)

var Username, Password, Hostaddress string

func Init(username, password, hostaddress string) {
	Username = username
	Password = password
	Hostaddress = hostaddress
}

var totalThreads = runtime.NumCPU()

/*
toBucket => true means that all the documents will be
pushed to bucket but the type of documents will depend
on the scope and collections
*/

func PopulateDocs(bucket string, toBucket bool,
	scopes, collections []string, numDocs int64) {

	diff := populateDocs(bucket, toBucket, scopes, collections, numDocs)
	if diff != nil {
		for key, value := range diff[Hostaddress] {
			// Get the scope and collection from key
			split := strings.Split(key, ":")
			scope := split[0]
			coll := split[1]
			PopulateDocs(bucket, toBucket, []string{scope}, []string{coll}, value)
		}
	} else {
		return
	}

}

func populateDocs(bucket string, toBucket bool,
	scopes, collections []string, numDocs int64) map[string]map[string]int64 {

	beforeStats := stats.GetStats(bucket)
	numDocsPerType := numDocs / (int64)(len(collections))
	maxDocsPerIter := (int64)(5000000)

	collectionIds := make([]string, 0)
	if !toBucket {
		for i, coll := range collections {
			retry := 0
			for {
				cid := kvutility.GetCollectionID(bucket, scopes[i], coll, Username, Password, Hostaddress)
				if cid == "" {
					logging.Infof("Retrying because of UNKNOWN_COLLECTION error")
					time.Sleep(5 * time.Second)
					retry++
					if retry > 100 {
						panic(fmt.Errorf("Empty collectionID"))
					}
				} else {
					collectionIds = append(collectionIds, cid)
					break
				}
			}
		}
	}

	numDocsPending := numDocsPerType
	for {
		if numDocsPending > maxDocsPerIter {
			numDocsPending = numDocsPending - maxDocsPerIter
			pushDocs(bucket, toBucket,
				collectionIds, scopes, collections,
				maxDocsPerIter)

		} else {
			pushDocs(bucket, toBucket,
				collectionIds, scopes, collections,
				numDocsPending)
			break
		}
	}

	time.Sleep(5 * time.Second)
	// Validate the number of docs in the bucket/collections
	for {
		retry := 0
		afterStats := stats.GetStats(bucket)
		if !stats.Validate(beforeStats, afterStats, scopes, collections, numDocsPerType) {
			logging.Fatalf("Docs were not populated in the bucket\n before: %v\n, after: %v\n, retrying...%v", beforeStats, afterStats, retry)
			time.Sleep(30 * time.Second)
			if retry > 20 {
				// Generate the diff and return the diff
				return stats.GetDiff(beforeStats, afterStats, (int64)(numDocsPerType))
			}
		} else {
			return nil
		}
	}

}

func pushDocs(bucket string, toBucket bool,
	scopes, collections, collectionIds []string,
	numDocs int64) {

	// a. Generate docs for populating to bucket
	var wg1 sync.WaitGroup

	before := time.Now().UnixNano()
	docCh := make([]chan map[string]interface{}, totalThreads)
	stepPerThread := make([]int64, totalThreads)

	// b. Generate slice of steps. last thread will push remaining docs
	step := numDocs / (int64)(totalThreads)
	sum := (int64)(0)
	for tid := 0; tid < totalThreads; tid++ {
		sum += step
		if tid == totalThreads-1 {
			stepPerThread[tid] = numDocs - sum
		} else {
			stepPerThread[tid] = step
		}
	}

	for tid := 0; tid < totalThreads; tid++ {
		docCh[tid] = make(chan map[string]interface{}, stepPerThread[tid])
		wg1.Add(1)
		go func(index int) {
			seed := rand.New(rand.NewSource(time.Now().UnixNano()))
			defer wg1.Done()
			//defer close(docCh[index])
			var i int64
			for i = 0; i < step; i++ {
				doc := generateJson(String(100, seed), seed)
				docCh[index] <- doc
			}
			return
		}(tid)
	}
	wg1.Wait()
	after := time.Now().UnixNano()
	logging.Infof("Time taken to generate docs for numDocs: %v is: %v", numDocs, after-before)

	docChForNextColl := make([]chan map[string]interface{}, totalThreads)
	for i := 0; i < totalThreads; i++ {
		docChForNextColl[i] = make(chan map[string]interface{}, len(docCh[i]))
	}

	repopulateDocCh := func(docChForNextColl []chan map[string]interface{}) {
		docCh = make([]chan map[string]interface{}, totalThreads)
		var wg3 sync.WaitGroup
		for i := 0; i < totalThreads; i++ {
			wg3.Add(1)
			go func(index int) {
				defer wg3.Done()
				docCh[index] = make(chan map[string]interface{}, len(docChForNextColl[index]))
				for {
					select {
					case doc := <-docChForNextColl[index]:
						docCh[index] <- doc
					default:
						return
					}
				}
			}(i)
		}
		wg3.Wait()
		logging.Infof("Re-populated docCh, len(docCh[0]): %v, len(docChForNextColl): %v", len(docCh[0]), len(docChForNextColl[0]))
	}

	// c. Push docs to collections
	for i, collection := range collections {
		logging.Infof("Pushing docs for collection: %v", collection)
		var wg2 sync.WaitGroup
		for tid := 0; tid < totalThreads; tid++ {
			wg2.Add(1)
			go func(tid int) {
				defer wg2.Done()
				seed := rand.New(rand.NewSource(time.Now().UnixNano()))

				// Connect to a bukcet
				url := "http://" + bucket + ":" + Password + "@" + Hostaddress
				b, err := common.ConnectBucket(url, "default", bucket)
				if err != nil {
					panic(err)
				}
				defer b.Close()

				pushDoc := func(index int, doc map[string]interface{}) {
					doc["type"] = scopes[i] + ":" + collection
					retry := 0
					for {
						key := "Users-" + String(20, seed)
						if !toBucket {
							err = b.SetC(key, collectionIds[i], 0, doc)
							if err != nil {
								logging.Warnf("Reveived error: %v", err)
								time.Sleep(30 * time.Second)
								retry++
								if retry > 5 {
									return
								}
							} else {
								docChForNextColl[index] <- doc
								return
							}
						} else {
							err = b.Set(key, 0, doc)
							if err != nil {
								logging.Warnf("Reveived error: %v", err)
								time.Sleep(30 * time.Second)
								retry++
								if retry > 5 {
									return
								}
							} else {
								docChForNextColl[index] <- doc
								return
							}
						}
					}
				}

				// Push docs
				for {
					select {
					case doc := <-docCh[tid]:
						pushDoc(tid, doc)
					default:
						return
					}
				}
			}(tid)
			wg2.Wait()
			logging.Infof("Done pushing document for collection: %v, len(docChForNextColl[0]): %v", collection, len(docChForNextColl[0]))
			repopulateDocCh(docChForNextColl)
		}
		logging.Infof("Time taken to generate and push docs to collections: %v is: %v", collectionIds, time.Now().UnixNano()-before)
	}
}
