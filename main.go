package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	downMap = map[string]bool{
		"one":   false,
		"two":   false,
		"three": false,
		"four":  false,
		"five":  false,
		"six":   false,
		"seven": false,
		"eight": false,
	}
	downMapLock sync.Mutex
)

func main() {
	go setDown()

	// command-line flags
	upstreamAddr := flag.String("upstream-addr", "", "Upstream server address")
	localPort := flag.String("port", "8080", "Local server port")
	flag.Parse()

	// environment variables
	envUpstreamAddr := os.Getenv("UPSTREAM_ADDR")
	envLocalPort := os.Getenv("PORT")

	if  envLocalPort != "" {
		localPort = &envLocalPort
	}

	if *upstreamAddr == "" {
		upstreamAddr = &envUpstreamAddr
	}

	if *upstreamAddr == "" || *localPort == "" {
		log.Fatalf("invalid usage")
	}

	r := gin.Default()
	r.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"info": "brokenProxy",
		})
	})
	r.GET("/:host/*path", func(ctx *gin.Context) {
		h := ctx.Param("host")
		url := ctx.Param("path")
		if _, ok := downMap[h]; !ok {
			ctx.String(404, "invalid host")
			return
		}
		req, _ := http.NewRequest(ctx.Request.Method, *upstreamAddr+url, ctx.Request.Body)
		res, e := http.DefaultClient.Do(req)
		if e != nil {
			ctx.String(500, e.Error())
			return
		}
		down := downMap[h]
		if down {
			simulatedown(*ctx, res)
		} else {
			h := map[string]string{}
			for k, v := range res.Header {
				h[k] = strings.Join(v, ", ")
			}
			ctx.DataFromReader(res.StatusCode, res.ContentLength, "", res.Body, h)
		}

	})
	r.Run(":" + *localPort)
}

func simulatedown(ctx gin.Context, res *http.Response) {
	c := rand.Int() % 3
	if c == 0 {
		fmt.Println("sending error")
		ctx.String(500, "internal server error")
	}
	if c == 1 {
		// slow response
		fmt.Println("sending slow response")
		time.Sleep(15 * time.Second)
		h := map[string]string{}
		for k, v := range res.Header {
			h[k] = strings.Join(v, ", ")
		}
		ctx.DataFromReader(res.StatusCode, res.ContentLength, "", res.Body, h)
	}
	if c == 2 {
		// nonsensical response
		fmt.Println("sending nonsensical response")
		b,_ := ioutil.ReadAll(res.Body)
		messedUp := switchSegments(string(b))
		ctx.String(res.StatusCode, messedUp)
	}
}

func setDown() {
	for {
		downMapLock.Lock()
		numKeys := rand.Intn(4) + 1
		keys := make([]string, 0, len(downMap))
		for key := range downMap {
			keys = append(keys, key)
		}
		rand.Shuffle(len(keys), func(i, j int) {
			keys[i], keys[j] = keys[j], keys[i]
		})
		for i := 0; i < numKeys; i++ {
			downMap[keys[i]] = true
		}
		for i := numKeys; i < len(keys); i++ {
			downMap[keys[i]] = false
		}
		downMapLock.Unlock()
		time.Sleep(1 * time.Minute)
	}
}


func switchSegments(str string) string {
    minSegmentLength := 12
    strLen := len(str)

    if strLen < minSegmentLength*4 {
        return str
    }

    // Choose two random start indices for the segments
    rand1 := rand.Intn(strLen - minSegmentLength*2)
    rand2 := rand.Intn(strLen - minSegmentLength*2)
    // Ensure the two segments do not overlap
    for absDiff(rand2, rand1) < minSegmentLength {
        rand2 = rand.Intn(strLen - minSegmentLength*2)
    }

    // Swap the two segments
    if rand1 > rand2 {
        rand1, rand2 = rand2, rand1
    }

    seg1Start := rand1
    seg1End := rand1 + minSegmentLength
    seg2Start := rand2
    seg2End := rand2 + minSegmentLength

    result := str[:seg1Start] + str[seg2Start:seg2End] + str[seg1End:seg2Start] + str[seg1Start:seg1End] + str[seg2End:]

    return result
}


func absDiff(a, b int) int {
    if a > b {
        return a - b
    }
    return b - a
}
