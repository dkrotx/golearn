package main

import (
    "fmt"
    "time"
    "bufio"
    "strings"
)

var lyrics = `Risin' up, back on the street
Did my time, took my chances
Went the distance, now I'm back on my feet
Just a man and his will to survive

So many times, it happens too fast
You trade your passion for glory
Don't lose your grip on the dreams of the past
You must fight just to keep them alive

It's the eye of the tiger, it's the thrill of the fight
Risin' up to the challenge of our rival
And the last known survivor stalks his prey in the night
And he's watchin' us all with the eye of the tiger`

func printNewLine(sc *bufio.Scanner) bool {
    for {
        if !sc.Scan() {
            return false
        }
        // print short lines fast
        if n, _ := fmt.Println(sc.Text()); n > 3 {
            break
        }
    }
    return true
}

func main() {
    ticker := time.NewTicker(2 * time.Second)
    sc := bufio.NewScanner(strings.NewReader(lyrics))

    printNewLine(sc) // print first line immediately
    for range ticker.C {
        if !printNewLine(sc) {
            ticker.Stop()
        }
    }
}
