package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var initialData [][]float64
var center [][]float64
var numOfClusterCenter int
var numOfDimension int
var dataIndexMapToCluster []int

func main() {
	numOfClusterCenter = 4
	numOfDimension = 7
	data, err := ioutil.ReadFile("./data.txt")
	if nil != err {
		return
	}
	stringData := strings.Split(string(data), ",")
	stringData = stringData[:len(stringData)-1]
	initialData = make([][]float64, (len(stringData) / numOfDimension))
	dataIndexMapToCluster = make([]int, len(initialData))
	for index := range stringData {
		if 0 == index%numOfDimension {
			initialData[index/numOfDimension] = make([]float64, numOfDimension)
		}
		floatValue, _ := strconv.ParseFloat(stringData[index], 64)
		initialData[index/numOfDimension][index%numOfDimension] = floatValue
	}
	initClusterCenter()
	totalError1 := classify()
	setCenter()
	totalError2 := classify()
	setCenter()
	count:= 2
	for (0 != math.Abs(totalError1 - totalError2)) {
		totalError1 = totalError2
		totalError2 = classify()
		setCenter()
		count++
	}
}

func initClusterCenter() {
	center = make([][]float64, numOfClusterCenter)
	for centerIndex := range center {
		center[centerIndex] = make([]float64, numOfDimension)
		for valueIndex := range center[centerIndex] {
			center[centerIndex][valueIndex] = initialData[centerIndex*len(initialData)/numOfClusterCenter][valueIndex]
		}
	}
}

func getODistance(data1 []float64, data2 []float64) float64 {
	sum := 0.0
	for index := range data1 {
		sum += math.Pow(data1[index]-data2[index], 2)
	}
	return math.Sqrt(sum)
}

func classify() float64 {
	totalDistance := 0.0
	// distanceBetweenCenter := make([][]float64, len(initialData))
	for dataIndex := range initialData {
		// distanceBetweenCenter[dataIndex] = make([]float64, numOfClusterCenter)
		minDistance := 9999.99
		for centerIndex := range center {
			distance := getODistance(initialData[dataIndex], center[centerIndex])
			if distance < minDistance {
				minDistance = distance
				dataIndexMapToCluster[dataIndex] = centerIndex
			}
		}
	}
	for index := range dataIndexMapToCluster {
		totalDistance += getODistance(initialData[index], center[dataIndexMapToCluster[index]])
	}
	fmt.Println( "total distance ->" ,totalDistance)	
	return totalDistance
}

func setCenter() {
	for centerIndex := range center {
		count := 0.0
		sum := make([]float64, numOfDimension)
		for dataIndex := range initialData {
			if centerIndex == dataIndexMapToCluster[dataIndex] {
				count++
				for propertyIndex := range sum {
					sum[propertyIndex] += initialData[dataIndex][propertyIndex]
				}
			}
		}
		for propertyIndex := range sum {
			center[centerIndex][propertyIndex] = round(sum[propertyIndex]/count, 3)
		}
	}
	fmt.Println(center)	
}
func round(f float64, n int) float64 {
	pow10N := math.Pow10(n)
	return math.Trunc((f+0.5/pow10N)*pow10N) / pow10N
}
