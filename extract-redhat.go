package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "extract-redhat"
	app.Usage = "RedHatのセキュリティアップデート情報を抽出します。"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "infile,i",
			Value: "data/target.html",
			Usage: "RedHatのページから取得したHTMLファイル",
		},
		cli.StringFlag{
			Name:  "outdir,o",
			Value: "dist",
			Usage: "抽出した結果の保存先ディレクトリ",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Println("start:", app.Name)

		infile := c.String("i")
		outdir := c.String("o")
		p := regexp.MustCompile(`<[^>]+>`)

		b, err := ioutil.ReadFile(infile)
		if err != nil {
			log.Println(err)
			return err
		}
		r := strings.NewReader(string(b))
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			log.Println(err)
			return err
		}
		records := []string{
			"アドバイザリー\t概要\t深刻度\t公開日\t説明\tURL",
		}
		counter := 0
		elems := doc.Find("#DataTables_Table_0 > tbody > .ng-scope")
		denom := elems.Length()
		elems.Each(func(_ int, s *goquery.Selection) {
			fmt.Fprintf(os.Stderr, "\r[%3d/%3d]", counter, denom)

			id := s.Find(".td-cve > span > a").First().Text()
			overview := s.Find(".td-synopsis > span").First().Text()
			severity := s.Find(".td-impact > span").First().Text()
			startDate := s.Find(".td-date > span").First().Text()
			nd, err := time.Parse("02 Jan 2006", startDate)
			if err != nil {
				log.Println(err)
				return
			}
			startDate = nd.Format("2006/01/02")
			url := "https://access.redhat.com/errata/" + id
			var descline string

			d, err := goquery.NewDocument(url)
			if err != nil {
				log.Println(err)
				return
			}
			d.Find("body").Each(func(_ int, s *goquery.Selection) {
				htmltext, err := s.Html()
				if err != nil {
					log.Println(err)
					return
				}
				lines := strings.Split(htmltext, "\n")
				for i := 0; i < len(lines); i++ {
					line := strings.TrimSpace(strings.ToLower(lines[i]))
					if line == "<h2>description</h2>" {
						descline = lines[i+1]
						descline = p.ReplaceAllString(descline, "")
						descline = strings.TrimSpace(descline)
						break
					}
				}
			})
			text := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", id, overview, severity, startDate, descline, url)
			records = append(records, text)

			counter++
		})

		if err := os.MkdirAll(outdir, os.ModeDir); err != nil {
			log.Println(err)
			return err
		}

		fn := time.Now().Format("20060102_150405.csv")
		outfile := filepath.Join(outdir, fn)
		if err := ioutil.WriteFile(outfile, []byte(strings.Join(records, "\n")), os.ModePerm); err != nil {
			log.Println(err)
			return err
		}

		fmt.Println()
		log.Println("completed:", app.Name)
		return nil
	}

	app.Run(os.Args)
}
