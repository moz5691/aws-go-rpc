package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	rpc "github.com/moz5691/twirp-lambda/rpc/movies"
)

var (
	err           error
	yearInt       int
	ratingFloat   float64
	apigatewayURL string
)

func main() {

	// Enabled the following to test via local
	//client := rpc.NewMoviesProtobufClient("http://localhost:8080", &http.Client{})

	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apigatewayURL = os.Getenv("APIGATEWAY_URL")

	client := rpc.NewMoviesProtobufClient(apigatewayURL, &http.Client{})

	method := flag.String("m", "r", "c-create, r-read, u-update, d-delete")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	switch *method {
	case "c":
		fmt.Println("=== Create : create a new data set ===")
		fmt.Print("Title: ")
		title, _ := reader.ReadString('\n')
		title = strings.Replace(title, "\n", "", -1)

		fmt.Print("Year: ")
		year, _ := reader.ReadString('\n')
		year = strings.Replace(year, "\n", "", -1)
		if yearInt, err = strconv.Atoi(year); err != nil {
			fmt.Println("fail to convert to integer")
		}

		fmt.Print("Plot: ")
		plot, _ := reader.ReadString('\n')
		plot = strings.Replace(plot, "\n", "", -1)

		fmt.Print("Rating: ")
		rating, _ := reader.ReadString('\n')
		rating = strings.Replace(rating, "\n", "", -1)
		if ratingFloat, err = strconv.ParseFloat(rating, 64); err != nil {
			fmt.Println("fail to convert to float")
		}

		item := &rpc.Item{
			Title:  title,
			Year:   int32(yearInt),
			Plot:   plot,
			Rating: float32(ratingFloat),
		}

		res, err := client.Post(context.Background(), item)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("From server: %+v\n", res)

	case "r":
		fmt.Println("=== Read : read a data with Title & Year ===")
		fmt.Print("Title: ")
		title, _ := reader.ReadString('\n')
		title = strings.Replace(title, "\n", "", -1)

		fmt.Print("Year: ")
		year, _ := reader.ReadString('\n')
		year = strings.Replace(year, "\n", "", -1)
		if yearInt, err = strconv.Atoi(year); err != nil {
			fmt.Println("fail to convert to integer")
		}

		term := &rpc.SearchTerm{
			Title: title,
			Year:  int32(yearInt),
		}

		res, err := client.GetByTitle(context.Background(), term)

		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			os.Exit(1)
		}
		fmt.Printf("From server: %+v\n", res)

	case "u":
		fmt.Println("=== Update : update a data set with Title & Year, create a new data if no match found")
		fmt.Print("Title: ")
		title, _ := reader.ReadString('\n')
		title = strings.Replace(title, "\n", "", -1)

		fmt.Print("Year: ")
		year, _ := reader.ReadString('\n')
		year = strings.Replace(year, "\n", "", -1)
		if yearInt, err = strconv.Atoi(year); err != nil {
			fmt.Println("fail to convert to integer")
		}

		fmt.Print("Plot: ")
		plot, _ := reader.ReadString('\n')
		plot = strings.Replace(plot, "\n", "", -1)

		fmt.Print("Rating: ")
		rating, _ := reader.ReadString('\n')
		rating = strings.Replace(rating, "\n", "", -1)
		if ratingFloat, err = strconv.ParseFloat(rating, 64); err != nil {
			fmt.Println("fail to convert to float")
		}

		term := &rpc.Item{
			Title:  title,
			Year:   int32(yearInt),
			Plot:   plot,
			Rating: float32(ratingFloat),
		}

		res, err := client.Update(context.Background(), term)

		if err != nil {
			fmt.Printf("Error: %+v\n", err)
		}
		fmt.Printf("From server: %+v\n", res)

	case "d":
		fmt.Println("=== Delete : delete a data set with Title & Year")
		fmt.Print("Title: ")
		title, _ := reader.ReadString('\n')
		title = strings.Replace(title, "\n", "", -1)

		fmt.Print("Year: ")
		year, _ := reader.ReadString('\n')
		year = strings.Replace(year, "\n", "", -1)
		if yearInt, err = strconv.Atoi(year); err != nil {
			fmt.Println("fail to convert to integer")
		}

		term := &rpc.SearchTerm{
			Title: title,
			Year:  int32(yearInt),
		}

		res, err := client.DeleteByTitle(context.Background(), term)

		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("From server: %+v\n", res)
	}

}
