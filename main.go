package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type instagramUser struct {
	fullname string
	username string
}

func compareFollowers(followers []instagramUser, following []instagramUser) {
	file, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	followerMap := make(map[string]bool)
	for _, follower := range followers {
		followerMap[follower.username] = true
	}
	for _, followingUser := range following {
		if !followerMap[followingUser.username] {
			fmt.Fprintf(file, "Follower not found. Username: %s, Fullname: %s\n", followingUser.username, followingUser.fullname)
		}
	}
	fmt.Println("Output written to output.txt")
}

func inputHeadersMapper(input string) http.Header {

	h := make(http.Header)
	lines := strings.SplitSeq(input, "\n")

	for line := range lines {
		line = strings.TrimSpace(line)

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		h.Set(key, value)
	}
	return h
}

func fetchFollowers(userId string, headers http.Header, kind string) []instagramUser {
	var followers []instagramUser
	var url string
	first := true
	var maxId string
	for {
		if first {
			url = fmt.Sprintf("https://www.instagram.com/api/v1/friendships/%s/%s/?count=25", userId, kind)
			first = false
		} else {
			url = fmt.Sprintf("https://www.instagram.com/api/v1/friendships/%s/%s/?count=25&max_id=%s&search_surface=follow_list_page", userId, kind, maxId)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		req.Header = headers

		client := &http.Client{}

		resp, err := client.Do(req)
		time.Sleep(10 * time.Millisecond)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var jsonData map[string]any

		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			panic(err)
		}
		users := jsonData["users"].([]any)
		var ok bool
		maxId, ok = jsonData["next_max_id"].(string)

		for _, user := range users {
			userMap := user.(map[string]any)
			followers = append(followers, instagramUser{fullname: userMap["full_name"].(string), username: userMap["username"].(string)})
		}
		fmt.Println("Processing: ", kind, " . Current max Id: ", maxId)
		if ok {
			continue
		} else {
			break
		}
	}

	return followers
}

func main() {
	userId := "xx"
	headers := inputHeadersMapper(`Host: www.instagram.com
User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:149.0) Gecko/20100101 Firefox/149.0
Accept: */*
Accept-Language: en-US,en;q=0.9
X-CSRFToken: xxx
X-IG-App-ID: xxx
X-ASBD-ID: xxx
X-IG-WWW-Claim: xxx
X-Web-Session-ID: xxx
X-IG-Max-Touch-Points: 0
X-Requested-With: XMLHttpRequest
Alt-Used: www.instagram.com
Connection: keep-alive
Referer: https://www.instagram.com/xxx/
Cookie: csrftoken=xxx; datr=xxx; ig_did=xxxx; mid=xx; ig_nrcb=1; ds_user_id=xx; sessionid=xx; ps_l=1; ps_n=1; rur=xxx; dpr=1.2; wd=1600x899
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-origin
Pragma: no-cache
Cache-Control: no-cache
TE: trailers`)

	followers := fetchFollowers(userId, headers, "followers")
	following := fetchFollowers(userId, headers, "following")
	compareFollowers(followers, following)
}
