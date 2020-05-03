## Tuenti Challenge 10

This repo has the code for all challenges from 1-16.
All that had coding are in go, except for problem 14 which was easier in Node.

Each problem is a module (either go or node), in order to set it up:

Go:

    go get

Node:
    
    yarn

All programs accept the input through stdin and output the result to stdout (plus the log to stderr). So executing them is just:

    cat testInput.txt | go run main.go > testOutput.txt
    

Parsing of the input is annotation based using this [library](https://github.com/furstenheim/challenge_encoding).


### Problems
### Problem 1. Rock paper scissors.

First problem was straightforward. Since there are only 9 cases we can just hardcode them.

### Problem 2. Strongest player
There is a list of players (we don't know how many) with a hidden strength. We are provided a list of matches and need to find out the strongest player.

Since they are completely ordered, and we want to find the maximum, we can just discard everybody that lost at least one match.
At the end, we pick up the one that hasn't lose any.
It could potentially mean that there are two people that have never lost, but that is just bad input, because in that case we don't have enough information.

### Problem 3. Galdos
Third problem provides us with a copy of "Fortunata y Jacinta", and we need to be able to know the number of occurrences of each word, and it's total ranking.

We proceed in three phases. 

First, we need to clean up the book of all punctuation. For example, if the word is "Hello!" we are interested in "hello". Since the book is not extremely large, we can simply load it as a string and apply a regex:

```go
func parseBook (in string) string {
	lowerString := strings.ToLower(in)
	re := regexp.MustCompile(`[^abcdefghijklmnñopqrstuvwxyzáéíóúü]`)
	spacedString := re.ReplaceAll([]byte(lowerString), []byte(" "))
	return string(spacedString)
}
```

Second, we count words. We are asked to ignore words that have less than 3 letters. The tricky part is that we are dealing with unicode. For example, `él` is encoded with three bytes `C3 A9 6C` but it is only length two. With go this is straightforward using runes:

```go
word := "él"
realLength := len([]rune("él"))
realLength == 2
```

Third we order them, taking into account that we have to use unicode order.

### Problem 4. Steam
In this problem we need to access the pre-production environment of a certain game company.
Pre-production is sitting inside a VPN, so we cannot access directly. However, the load balancer of production has access, so we need to trick the load balancer (NGINX) to allow us to access.

For me, this was definitely the hardest problem of the first part. After trying tons of different combinations, I finally came to the one that worked:

``` curl steam-origin.contest.tuenti.net:9876/games/cat_fight/get_key?url=http://pre.steam-origin.contest.tuenti.net:9876/games/cat_fight/get_key  -H "Host: pre.steam-origin.contest.tuenti.net:9876"
```

With the host header, we trick nginx (not really nginx but the faulty configuration) to direct us to the environment that we want. Validation is probably checking that the url matches, so what we do is adding a query param with the pre.steam-origin url. Name of the query param doesn't matter. It could've been:


``` curl steam-origin.contest.tuenti.net:9876/games/cat_fight/get_key?another-parameter=http://pre.steam-origin.contest.tuenti.net:9876/games/cat_fight/get_key  -H "Host: pre.steam-origin.contest.tuenti.net:9876"
```


### Problem 5. Tuenti numbers
