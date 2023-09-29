package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

type testCase struct {
	name      string
	packSizes []int
	items     int
	expected  []int
}

var testCases = []testCase{
	{"Test case 1", []int{250, 500, 1000, 2000, 5000}, 1, []int{1, 0, 0, 0, 0}},
	{"Test case 2", []int{250, 500, 1000, 2000, 5000}, 250, []int{1, 0, 0, 0, 0}},
	{"Test case 3", []int{250, 500, 1000, 2000, 5000}, 251, []int{0, 1, 0, 0, 0}},
	{"Test case 4", []int{250, 500, 1000, 2000, 5000}, 501, []int{1, 1, 0, 0, 0}},
	{"Test case 5", []int{250, 500, 1000, 2000, 5000}, 12001, []int{1, 0, 0, 1, 2}},
	{"Test case 6", []int{250, 500, 1000, 2000, 5000}, 751, []int{0, 0, 1, 0, 0}},
	//{"Test case 7", []int{75, 250, 500, 1000, 2000, 5000}, 300, []int{4, 0, 0, 0, 0, 0}},
}

func TestCalculatePacks(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := calculatePacks(tc.packSizes, tc.items)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestWritePacksToFile(t *testing.T) {
	filename := "test.json"
	arr := []int{250, 500, 1000, 2000, 5000}

	err := writePacksToFile(filename, arr)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	arrRead, err := readPacksFromFile(filename)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(arr, arrRead) {
		t.Errorf("Expected %v, got %v", arr, arrRead)
	}

	err = os.Remove(filename)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestWritePacksHandler(t *testing.T) {
	uc := &PacksController{}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(`[1,2,3]`)))
	uc.WritePacksHandler(c)
	err := os.Remove(PACKS_FILENAME)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	assert.Equal(t, http.StatusOK, c.Writer.Status())
}

func TestReadPacksHandler(t *testing.T) {
	uc := &PacksController{}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	data := []int{1, 2, 3}
	err := writePacksToFile(PACKS_FILENAME, data)
	if err != nil {
		t.Fatal(err)
	}

	uc.ReadPacksHandler(c)

	assert.Equal(t, http.StatusOK, c.Writer.Status())
}

func TestCalculatePacksHandler(t *testing.T) {
	uc := &PacksController{}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/calculatePacks/10", nil)
	c.AddParam("items", "10")
	uc.CalculatePacksHandler(c)
	assert.Equal(t, http.StatusOK, c.Writer.Status())
}
