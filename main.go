package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type AtivoIbov struct {
	Codigo       string
	Acao         string
	Tipo         string
	QtdeTeorica  int64
	Participacao float64
}

type IbovRT struct {
	Data             string
	QtdeTeoricaTotal int64
	Redutor          float64
	Carteira         map[string]AtivoIbov
}

func main() {
	var headings, row []string
	var ibovRT IbovRT

	ibovRT.Carteira = make(map[string]AtivoIbov)

	doc, err := goquery.NewDocument("http://bvmf.bmfbovespa.com.br/indices/ResumoCarteiraTeorica.aspx?Indice=IBOV&amp;idioma=pt-br&idioma=pt-br") //  NewDocumentFromReader(strings.NewReader(resp.Body))
	if err != nil {
		log.Println("No url found")
		log.Fatal(err)
	}

	log.Println("Analisando Doc")

	inicio := time.Now()
	strDateIni := FormatDate(inicio, "20060102T150405") // Converter data inicial em string
	datestr := strDateIni[0:8]                          // Obter apenas data inicial
	// Find each table

	doc.Find("tfoot").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			rowhtml.Find("td").Each(func(indexth int, tableheading *goquery.Selection) {
				tableheading.Find("div").Each(func(indexth int, tablecont *goquery.Selection) {

					atext := formatNum(tablecont.Text())
					if atext != "" && atext != " " {

						log.Println(atext, tablecont.Attr)
						headings = append(headings, atext)
					}
				})
			})

			ibovRT.Data = datestr

			ibovRT.QtdeTeoricaTotal, _ = strconv.ParseInt(headings[2], 0, 0)

			headings[3] = strings.Replace(headings[3], ",", ".", -1)
			ibovRT.Redutor, _ = strconv.ParseFloat(headings[3], 64)
		})
	})
	doc.Find("tbody").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {

				atext := formatNum(tablecell.Text())
				if atext != "" && atext != " " {
					row = append(row, atext)
				}
			})
			if len(row) > 2 {

				ativo := AtivoIbov{Codigo: row[0], Acao: row[1], Tipo: row[2]}

				log.Println(row)

				row[3] = strings.Replace(row[3], ",", ".", -1)
				row[4] = strings.Replace(row[4], ",", ".", -1)

				ativo.QtdeTeorica, _ = strconv.ParseInt(row[3], 0, 0)
				ativo.Participacao, _ = strconv.ParseFloat(row[4], 64)
				ibovRT.Carteira[row[0]] = ativo

			}
			row = nil
		})
	})

	log.Println(ibovRT)

}

// FormatDate qq coisa
func FormatDate(t time.Time, layout string) string {
	const bufSize = 64
	var b []byte
	max := len(layout) + 10
	if max < bufSize {
		var buf [bufSize]byte
		b = buf[:0]
	} else {
		b = make([]byte, 0, max)
	}
	b = t.AppendFormat(b, layout)
	return string(b)
}

func formatNum(texto string) string {
	result := ""
	for _, ch := range texto {
		if ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' && ch != '.' {
			result = result + string(ch)
		}
	}
	return result
}
