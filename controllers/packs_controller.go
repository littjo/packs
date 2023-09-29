package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PacksController struct {
}

func NewPacksController() *PacksController {
	return &PacksController{}
}

const PACKS_FILENAME = "./packs.json"

func calculatePacks(packSizes []int, items int) []int {
	sum := 0
	// Packing items
	for i := len(packSizes) - 1; i >= 0 && items > 0; i-- {
		packSize := packSizes[i]
		packs := items / packSize
		if packs > 0 {
			items %= packSize
			sum += packs * packSize
		}
	}
	if items > 0 {
		sum += packSizes[0]
	}

	numberOfPacks := make([]int, len(packSizes))

	// Optimizing number of packs
	for i := len(packSizes) - 1; i >= 0 && sum > 0; i-- {
		packSize := packSizes[i]
		packs := sum / packSize
		if packs > 0 {
			sum %= packSize
			numberOfPacks[i] = packs
		}
	}

	return numberOfPacks
}

func writePacksToFile(filename string, packs []int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(packs)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func readPacksFromFile(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var arr []int
	err = json.Unmarshal(data, &arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func (uc *PacksController) WritePacksHandler(c *gin.Context) {
	var arr []int
	if err := c.BindJSON(&arr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := writePacksToFile(PACKS_FILENAME, arr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Packs written to file successfully"})
}

func (uc *PacksController) ReadPacksHandler(c *gin.Context) {
	arr, err := readPacksFromFile(PACKS_FILENAME)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"packs": arr})
}

func (uc *PacksController) CalculatePacksHandler(c *gin.Context) {
	items, err := strconv.Atoi(c.Param("items"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter"})
	} else {
		if packSizes, err := readPacksFromFile(PACKS_FILENAME); err != nil {
			log.Println("Error when reading the pack sizes from disk.")
		} else {
			numOfPacks := calculatePacks(packSizes, items)
			var builder strings.Builder
			builder.WriteString("You will receive:\n")
			for i := len(packSizes) - 1; i >= 0; i-- {
				if numOfPacks[i] > 0 {
					builder.WriteString(fmt.Sprintf("%d x %d\n", numOfPacks[i], packSizes[i]))
				}
			}
			c.JSON(http.StatusOK, gin.H{"result": builder.String()})
		}
	}
}
