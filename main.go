package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Person struct{
	Nome string
	Idade int
	Pontuação int
	NomeByte []byte
}

func main() {
	
	if len(os.Args) >= 3 {
		fileNameIn := os.Args[1] //indice 0 é o próprio nome do arquivo executado
		fileNameOut := os.Args[2]

		//valida se arquivo de entrada existe
		if _, err := os.Stat(fileNameIn); err == nil{
			//arquivo encontrado
			//valiad se nome do arquivo de saída é válido, se validar já é criado um arquivo vazio com o nome escolhido
			if ValidateFileOut(fileNameOut) {
				rows := GetData(fileNameIn)
				if rows == nil {
					println("Erro ao carregar registros: " + err.Error())
				} else {
					
					people := make([]Person,len(rows)-1) //retira 1 do cabeçalho
					for k, row := range rows {
						if k == 0 { continue } //pulando cabeçalho do arquivo
						
						idade, _ := strconv.Atoi(row[1])
						pontuacao, _ := strconv.Atoi(row[2])
						person := Person{
							Nome: row[0],
							Idade: idade,
							Pontuação: pontuacao,
							NomeByte: []byte(row[0]),
						}
						people[k-1] = person
					}

					// fmt.Println("Sem sort: ",people)

					sort.SliceStable(people,func(i, j int) bool { 
						if people[i].Nome != people[j].Nome	{
							var stri string = people[i].Nome
							var strj string = people[j].Nome
							var strilow string = strings.ToLower(stri)
							var strjlow string = strings.ToLower(strj)
							if strilow == strjlow {
								return people[i].Nome < people[j].Nome
							}
							return strilow < strjlow
						}
						return people[i].Idade < people[j].Idade
					})
					// fmt.Println("Sort por nome: ",people)

					if err := WriteOutputFile(people, fileNameOut); err != nil {
						println("Erro ao gravar aquivo da saída: " + err.Error())
					}
				}
			} else {
				println("Arquivo de saída inválido: " + fileNameOut)
			}
		} else if errors.Is(err, os.ErrNotExist){ 
			println("Arquivo de entrada '"+fileNameIn+"' não existe.")
		}
	} else {
		println("Argumentos de entrada não informados: \nUtilização: go run main.go <arquivo_entrada> <arquivo_saida>")
	}
}

func GetData(filename string) [][]string{
	fileIn, err := os.Open(filename)
	defer fileIn.Close()
	if err != nil {
		println("Erro ao abrir arquivo de entrada: " + err.Error())
		return nil
	} else {
		reader := csv.NewReader(fileIn)
		rows, err := reader.ReadAll()

		if err != nil {
			println("Erro ao carregar registros: " + err.Error())
			return nil
		}
		return rows
	}
}

func WriteOutputFile(people []Person, filename string) error {
	// println(filename)
	for _, row := range people {
		fmt.Println(row)
	}

	fileOut, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fileOut.Close()	

	writer := csv.NewWriter(fileOut)
	defer writer.Flush()

	if err := writer.Write([]string{"Nome","Idade","Pontuação"}); err != nil {
		return err
	}

	for _, row := range people {
		line := []string{row.Nome , strconv.Itoa(row.Idade) , strconv.Itoa(row.Pontuação)}
		if err := writer.Write(line); err != nil {
			return err
		}
	}
	return nil
}


func ValidateFileOut(filename string) bool{
	vReturn := true

	newFile, e := os.Create(filename)
	if e != nil {
		vReturn = false
		println("Erro ao validar arquivo de saída: ", e.Error())
	}
	newFile.Close()

	return vReturn
}