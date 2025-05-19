package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type result struct {
	Resultados []struct {
		CantVotos     int    `json:"cant_votos"`
		IDCandidatura string `json:"id_candidatura"`
	} `json:"resultados"`
}

func main() {
	cant_de_bancas := 30
	url := "https://caba.datosoficiales.com/resultados/CABA/DIP.json"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-AR,es;q=0.9")
	req.Header.Set("Referer", "https://caba.datosoficiales.com/")
	req.Header.Set("Origin", "https://caba.datosoficiales.com")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response body - err=%s", err.Error())
	}

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		log.Fatalf("Unexpected content type: %s\nBody:\n%s", ct, body)
	}

	var resultados result

	err = json.Unmarshal(body, &resultados)
	if err != nil {
		log.Fatalf("error unmarshalling body - err=%s", err.Error())
	}

	votos := make(map[string]int)

	ADORNI := "Adorni"
	SANTORO := "Santoro"
	PRO := "Pro"
	ALIEN := "Alien"
	GORILA := "Gorila"

	for _, r := range resultados.Resultados {
		switch r.IDCandidatura {
		case "14.14.14":
			votos[ADORNI] = r.CantVotos
		case "4.4.4":
			votos[SANTORO] = r.CantVotos
		case "2.2.2":
			votos[PRO] = r.CantVotos
		case "3.3.3":
			votos[ALIEN] = r.CantVotos
		case "6.6.6":
			votos[GORILA] = r.CantVotos
		}
	}

	var ado, san, pro, ali, gor int

	currentSlot := map[string]int{
		"Adorni":  0,
		"Santoro": 0,
		"Pro":     0,
		"Alien":   0,
		"Gorila":  0,
	}

	for i := 1; i <= cant_de_bancas; i++ {
		valor_mas_grande := 0.0
		candidato_actual := ""

		for candidato, totalVotos := range votos {
			slot := currentSlot[candidato] + 1
			if slot <= cant_de_bancas {
				quotient := float64(totalVotos) / float64(slot)
				if quotient > valor_mas_grande {
					valor_mas_grande = quotient
					candidato_actual = candidato
				}
			}
		}

		switch candidato_actual {
		case "Adorni":
			ado++
		case "Santoro":
			san++
		case "Pro":
			pro++
		case "Alien":
			ali++
		case "Gorila":
			gor++
		}

		currentSlot[candidato_actual]++
	}

	fmt.Println("Numero de bancas por candidato:")
	fmt.Printf("Adorni: %d bancas\n", ado)
	fmt.Printf("Santoro: %d bancas\n", san)
	fmt.Printf("Pro: %d bancas\n", pro)
	fmt.Printf("Alien: %d bancas\n", ali)
	fmt.Printf("Gorila: %d bancas\n", gor)
}
