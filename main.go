package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	input_path := "/Data_source_dir"
	output_path := "Lotsize2018.csv"
	f, err := os.Create(output_path)
	if err != nil {
		log.Fatal(err)
	}
	output_writer := csv.NewWriter(f)
	output_writer.Write([]string{"date", "instrument", "expiry", "lotsize"})
	output_writer.Flush()
	paths := []string{}
	e := filepath.Walk(input_path, func(path string, info os.FileInfo, err error) error {
		if err == nil && strings.Contains(info.Name(), ".csv") {
			paths = append(paths, path)
		}
		return nil
	})
	log.Println("File Cache done, starting DB and loading...")
	if e != nil {
		log.Println(e)
	}
	LotSizes := make(map[string]float64)
	for _, filepath := range paths {
		date := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
		date = date[2 : len(date)-4]
		cur_file, err := os.Open(filepath)
		if err != nil {
			log.Println(err)
			continue
		}
		cur_reader := csv.NewReader(cur_file)
		contents, _ := cur_reader.ReadAll()
		for _, line := range contents {
			inst_name := line[0]
			inst_name = strings.ReplaceAll(inst_name, "FUTIDX", "")
			inst_name = strings.ReplaceAll(inst_name, "FUTSTK", "")
			inst_name = strings.ReplaceAll(inst_name, "OPTIDX", "")
			inst_name = strings.ReplaceAll(inst_name, "OPTSTK", "")
			if _, ok := LotSizes[fmt.Sprintf("%s_%s", date, inst_name)]; !ok {
				j, _ := strconv.ParseFloat(line[9], 64)
				k, _ := strconv.ParseFloat(line[10], 64)
				if j == 0 || k == 0 {
					continue
				}
				lotsize := j / k
				if !strings.Contains(inst_name, "CE") && !strings.Contains(inst_name, "PE") {
					LotSizes[fmt.Sprintf("%s_%s", date, inst_name)] = lotsize
				} else {
					pos := strings.LastIndex(inst_name, "-")
					pos = pos + 5
					LotSizes[fmt.Sprintf("%s_%s", date, inst_name[:pos])] = lotsize
				}

			} else {
				continue
			}
		}
	}

	for name, lotsize := range LotSizes {
		date := ConvDateToTS(strings.Split(name, "_")[0])
		name_split := (strings.Split(name, "_")[1])
		e := ConvDateToTS(name_split[len(name_split)-11:])
		inst := name_split[:len(name_split)-11]
		temp := []string{
			fmt.Sprint(date),
			inst,
			fmt.Sprint(e),
			fmt.Sprint(lotsize),
		}
		output_writer.Write(temp)
		output_writer.Flush()
	}
}

func ConvDateToTS(name string) int64 {
	if !strings.Contains(name, "-") {
		inter, _ := time.Parse("02-01-2006", fmt.Sprintf("%s-%s-20%s", name[:2], name[2:4], name[4:]))
		return inter.Unix()
	} else {
		month := strings.Split(name, "-")[1]
		switch month {
		case "JAN":
			name = strings.ReplaceAll(name, month, "01")
		case "FEB":
			name = strings.ReplaceAll(name, month, "02")
		case "MAR":
			name = strings.ReplaceAll(name, month, "03")
		case "APR":
			name = strings.ReplaceAll(name, month, "04")
		case "MAY":
			name = strings.ReplaceAll(name, month, "05")
		case "JUN":
			name = strings.ReplaceAll(name, month, "06")
		case "JUL":
			name = strings.ReplaceAll(name, month, "07")
		case "AUG":
			name = strings.ReplaceAll(name, month, "08")
		case "SEP":
			name = strings.ReplaceAll(name, month, "09")
		case "OCT":
			name = strings.ReplaceAll(name, month, "10")
		case "NOV":
			name = strings.ReplaceAll(name, month, "11")
		case "DEC":
			name = strings.ReplaceAll(name, month, "12")
		default:
			log.Fatal("Error", month)
		}
		inter, _ := time.Parse("02-01-2006", name)
		return inter.Unix()
	}

}
