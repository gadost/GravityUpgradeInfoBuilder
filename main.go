package main

import (
    "crypto/sha256"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "net/http"
    "time"
)

var releaseTag string
var currentHeight int64
var blockTime int64
var apxDate string
var binariesMap = map[string]string{
    "linux/amd64":   "gravity-linux-amd64",
    "darwin/amd64":  "gravity-darwin-amd64",
    "linux/arm64":   "gravity-linux-arm64",
    "windows/amd64": "gravity-windows-amd64.exe",
}

func init() {
    flag.StringVar(&releaseTag, "tag", "", "release tag \ne.g  -tag v1.3.0")
    flag.Int64Var(&currentHeight, "current-height", 0, "current height e.g. 400001")
    flag.Int64Var(&blockTime, "block-time", 6300, "current block time in ms e.g. 6300")
    flag.StringVar(&apxDate, "upgrade-date", "", "date&time for update e.g. 18 Jan 22 16:00 UTC   (RFC822)")
}

func main() {
    flag.Parse()
    if currentHeight != 0 && blockTime != 0 && len(apxDate) > 0 {
        tsX := parseDate(apxDate).UnixMilli()
        tsN := time.Now().UTC().UnixMilli()
        blockX := (tsX-tsN)/(blockTime) + currentHeight
        fmt.Printf("%s ~ calculated height %d\n", apxDate, blockX)
    }
    if len(releaseTag) > 0 {
        var binaryDownloadUrl = fmt.Sprintf("https://github.com/Gravity-Bridge/Gravity-Bridge/releases/download/%s/", releaseTag)
        var binary = map[string]string{}
        for k, v := range binariesMap {
            binary[k] = binaryDownloadUrl + v + "?checksum=sha256:" + getSHA256(getFile(binaryDownloadUrl+v))
        }

        buildJSON(binary)
    }
}

func getFile(s string) []byte {
    resp, err := http.Get(s)
    if err != nil {
        panic(err)
    }
    if resp.StatusCode != 200 {
        panic(fmt.Sprintf("Release Not Found"))
    }
    defer resp.Body.Close()
    bodyBytes, _ := io.ReadAll(resp.Body)
    return bodyBytes
}

func getSHA256(b []byte) string {
    return fmt.Sprintf("%x", sha256.Sum256(b))
}

func buildJSON(s map[string]string) {
    x := map[string]interface{}{"binaries": s}
    res, _ := json.Marshal(x)
    fmt.Println(string(res))
}

func parseDate(s string) time.Time {
    t, _ := time.Parse(time.RFC822, s)
    return t
}
