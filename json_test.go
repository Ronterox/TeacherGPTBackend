package main

import (
	"encoding/json"
	"log"
	"math"
	"testing"
)

var jsonString string = `{
    "topic": "Periodismo como herramienta política",
    "content": "¿Quién lee noticias, periódicos, prensa digital, las ve o las escucha? ¿Y quién cree en ella? ¿A qué me refiero? ¿Creéis que es puramente objetiva?",
    "chunk": "En España, el porcentaje de personas que evitan las noticias 'a veces o a menudo' ha pasado del 26% en 2017 al 35% en 2022. En España, 3 de cada 10 personas (35%) evita consumir noticias 'a veces o a menudo'.",
    "answer": "Ramon se cree todo lo que lee en los periódicos y noticias",
    "correct": false,
    "reason": ""
}`

func TestParseJson(t *testing.T) {
	var questions QuestionOpen
	err := json.Unmarshal([]byte(jsonString), &questions)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", questions)
}

func TestOutofBounds(t *testing.T) {
    length := len(jsonString)
    log.Printf("Length: %d", length)
    log.Printf("Out of bounds: %s", jsonString[0:int(math.Min(float64(2 * length), float64(length)))])
}
