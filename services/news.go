package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var httpClient = &http.Client{}

type New struct {
	title string
	body string
	url string
}

func GetRandomNew() string {
	// TODO: get news dynamically from an API
	// Original: https://www.clarin.com/viste/dato-asusta-tabaco-mata-45-000-personas-argentina-ano_0_y6Y5yP6a2M.html
	return `A pesar de haber disminuido casi en un 50 porciento los fumadores en Argentina en los últimos 25 años, según datos oficiales, el tabaquismo sigue siendo una de las mayores amenazas para la salud pública pues 225.000 personas enferman y 45.000 mueren cada año a causa del cigarrillo.De esos 45 mil decesos ligados al tabaco, 6 mil nunca probaron fumar, por lo que el humo pasivo todavía sigue siendo un grave problema sanitario.Esos fallecimientos representan el 14 porciento de todas las muertes en el país, según afirma un informe de la Red de Hospitales Universitarios de la Universidad de Buenos Aires (UBA), dado a conocer por la conmemoración del Día Mundial sin Tabaco, celebrado el miércoles pasado.En las últimas décadas, las políticas públicas de salud en Argentina buscaron disminuir la prevalencia del consumo del tabaco con medidas como la prohibición de fumar en espacios públicos, el empaquetado neutro y con advertencias sobre las consecuencias que tiene fumar, y la restricción de publicidad, junto con el desarrollo de programas antitabaco.Argentina firmó el Convenio Marco de la Organización Mundial de la Salud (OMS) para el Control del Tabaco (CMCT) el 25 de septiembre de 2003 pero por diversos intereses y lobbys de la industria, no ha pasado a la ratificación legislativa.`
}

func GetNew(topic string) (*New, error) {
	// Source: https://newsdata.io/documentation

	var newsDataApiKey = os.Getenv("NEWS_DATA_API_KEY")

	url := fmt.Sprintf("https://newsdata.io/api/1/news?apikey=%s&category=%s&language=es", newsDataApiKey, topic)

	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	
	if err != nil {
		log.Println("GetNew error -> Cannot create HTTP request: " + err.Error())
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	
	resp, err := client.Do(req)

	if err != nil {
		log.Println("GetNew error -> Cannot make HTTP request: " + err.Error())
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("GetNew error -> Cannot read body of response: " + err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	var result map[string]any
	err = json.Unmarshal([]byte(string(body)), &result)

	if err != nil {
		log.Println("GetNew error -> Cannot Unmarshal body of response: " + err.Error())
		return nil, err
	}

	firstResult := result["results"].([]any)[0].(map[string]any)

	new := New {
		title: firstResult["title"].(string),
		body: firstResult["content"].(string),
		url: firstResult["link"].(string),
	}

	return &new, nil
}

func GetSummarizedMessage(new *New) string {
	summarizedBody := Summarize(new.body)
	return fmt.Sprintf("*%s*\n\n%s\n\n%s", new.title, summarizedBody, new.url)
}
