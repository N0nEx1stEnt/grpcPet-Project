package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Product struct {
	Name   string  `json:"product"`
	Price  float64 `json:"price"`
	Rating float64 `json:"rating"`
}

func main() {
	//filePath := "db.csv"
	filePath := os.Getenv("file_path") // Получаем путь к файлу с помощью переменного окружения

	// Открытие файла для чтения
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Не удалось открыть файл: %v\n", err)
		return
	}
	defer file.Close()

	// Получение расширения файла
	fileExt := getFileExtension(filePath)

	// Чтение данных из файла
	var products []Product
	switch fileExt {
	case "json":
		products, err = readJSON(file)
	case "csv":
		products, err = readCSV(file)
	default:
		fmt.Println("Неподдерживаемый формат файла.")
		return
	}

	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	// Нахождение самого дорогого продукта
	sort.SliceStable(products, func(i, j int) bool {
		return products[i].Price > products[j].Price
	})
	mostExpensive := products[0]

	// Нахождение продукта с самым высоким рейтингом
	sort.SliceStable(products, func(i, j int) bool {
		return products[i].Rating > products[j].Rating
	})
	highestRating := products[0]

	fmt.Printf("Самый дорогой продукт: %s ($%.2f)\n", mostExpensive.Name, mostExpensive.Price)
	fmt.Printf("Продукт с самым высоким рейтингом: %s (рейтинг %.2f)\n", highestRating.Name, highestRating.Rating)
}

// Чтение данных из JSON-файла
func readJSON(file io.Reader) ([]Product, error) {
	var products []Product
	err := json.NewDecoder(file).Decode(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Чтение данных из CSV-файла
func readCSV(file io.Reader) ([]Product, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var products []Product
	for i, _ := range records {
		if i == 0 {
			continue
		}
		price, err := strconv.ParseFloat(strings.TrimSpace(records[i][1]), 64)
		if err != nil {
			return nil, err
		}

		rating, err := strconv.ParseFloat(strings.TrimSpace(records[i][2]), 64)
		if err != nil {
			return nil, err
		}

		product := Product{
			Name:   records[i][0],
			Price:  price,
			Rating: rating,
		}

		products = append(products, product)
	}

	return products, nil
}

// Получение расширения файла
func getFileExtension(filePath string) string {
	return filepath.Ext(filePath)[1:]
}
