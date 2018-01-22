// Copyright (c) 2018 Lynn Root
package main

import (
    "bufio"
    "fmt"
    "flag"
    "math"
    "math/rand"
    "os"
    "sort"
    "strings"
    "time"

    // "github.com/faiface/beep"
    // "github.com/faiface/beep/mp3"
    // "github.com/faiface/beep/speaker"
)

type Note struct {
    Name        string
    SciName     string
    Hertz       float64
    filePath    string
}

type StartingFreq struct {
    NoteName    string
    SciName     int
    Hertz       float64
}

var AFreq       = StartingFreq{"A", 0, 27.5000}
var ASharpFreq  = StartingFreq{"A#", 0, 29.1352}
var BFlatFreq   = StartingFreq{"Bb", 0, 29.1352}
var BFreq       = StartingFreq{"B", 0, 30.8677}
var CFreq       = StartingFreq{"C", 1, 32.7032}
var CSharpFreq  = StartingFreq{"C#", 1, 34.6478}
var DFlatFreq   = StartingFreq{"Db", 1, 34.6478}
var DFreq       = StartingFreq{"D", 1, 36.7081}
var DSharpFreq  = StartingFreq{"D#", 1, 38.8909}
var EFlatFreq   = StartingFreq{"Eb", 1, 38.8909}
var EFreq       = StartingFreq{"E", 1, 41.2034}
var FFreq       = StartingFreq{"F", 1, 43.6535}
var FSharpFreq  = StartingFreq{"F#", 1, 46.2493}
var GFlatFreq   = StartingFreq{"Gb", 1, 46.2493}
var GFreq       = StartingFreq{"G", 1, 48.9994}
var GSharpFreq  = StartingFreq{"G#", 1, 51.9131}
var AFlatFreq   = StartingFreq{"Ab", 1, 51.9131}


func UserInput() string {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter guess: ")
    text, _ := reader.ReadString('\n')
    return strings.TrimRight(text, "\n")
}


func CreateAllOctaves(note string, sci_name int, starting_freq float64) []Note {
    var AllOctaves []Note
    var maxHertz = 4186.02  // C8
    for octave := 0; octave < 8; octave++ {
        var noteHertz = starting_freq * math.Pow(2, float64(octave))
        if noteHertz < maxHertz {
            var n = Note{}
            n.Name = note
            n.SciName = fmt.Sprintf("%v%v", note, sci_name + octave)
            n.Hertz = noteHertz
            n.filePath = fmt.Sprintf("data/%v/%v%v.mp3", note, note, octave)
            AllOctaves = append(AllOctaves, n)
        }
    }
    return AllOctaves
}

func CreatePiano() []Note {
    var noteFreqs = []StartingFreq{AFreq, ASharpFreq, BFreq, CFreq, CSharpFreq, DFreq, DSharpFreq, EFreq, FFreq, FSharpFreq, GFreq, GSharpFreq}
    var piano []Note
    for _, freq := range noteFreqs {
        allOctaves := CreateAllOctaves(freq.NoteName, freq.SciName, freq.Hertz)
        piano = append(piano, allOctaves...)
    }
    return piano
}

func main() {
    showCommand :=
    piano := CreatePiano()

    // why not.
    sort.Slice(piano, func(i, j int) bool {
      return piano[i].Hertz < piano[j].Hertz
    })

    // random note to play
    seedSource := rand.NewSource(time.Now().UnixNano())
    keySeed := rand.New(seedSource)
    key := keySeed.Intn(88)
    givenNote := piano[key]

    // once all data files are there
    // f, _ := os.Open(givenNote.filePath)
    f, _ := os.Open("data/C4.mp3")
    s, format, _ := mp3.Decode(f)

    // play note in a speaker - needs a channel otherwise blocks exit of program
    speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
    done := make(chan struct{})
    speaker.Play(beep.Seq(s, beep.Callback(func() {
        close(done)
    })))
    <-done

    guess := UserInput()

    if strings.Compare(strings.ToUpper(guess), givenNote.Name) == 0 {
        fmt.Println("Yay ", givenNote.Name)
    } else {
        fmt.Println("You're WRONG", givenNote.Name)
    }
}
