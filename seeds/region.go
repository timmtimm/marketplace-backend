package seeds

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"marketplace-backend/business/regions"
	"marketplace-backend/util"
	"os"
	"path/filepath"
	"strings"
)

func SeedRegion(regionUC regions.UseCase) {
	path, _ := os.Getwd()

	files, err := ioutil.ReadDir(filepath.Join(path, "seeds/country"))
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".csv" {
			continue
		}

		data, err := os.Open(filepath.Join(path, "seeds/country/"+file.Name()))
		if err != nil {
			panic(err)
		}
		defer data.Close()

		fmt.Println("Seeding " + util.GetFilenameWithoutExtension(file.Name()) + "...")

		if util.GetFilenameWithoutExtension(file.Name()) == "indonesia" {
			domain := regions.Domain{
				Country: "Indonesia",
			}

			csvReader := csv.NewReader(data)

			for {
				rec, err := csvReader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					panic(err)
				}

				regionCode := strings.Split(rec[0], ".")
				regionName := rec[1]

				if len(regionCode) == 1 {
					domain.Province = strings.Title(strings.ToLower(regionName))
				} else if len(regionCode) == 2 {
					domain.Regency = strings.Title(strings.ToLower(regionName))
				} else if len(regionCode) == 3 {
					domain.District = strings.Title(strings.ToLower(regionName))
				} else if len(regionCode) == 4 {
					domain.Subdistrict = strings.Title(strings.ToLower(regionName))

					_, _, err = regionUC.Create(&domain)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
