package model

type Movie struct {
	Title  string  `dynamodbav:"Title, omitempty"`
	Year   int32   `dynamodbav:"Year,omitempty"`
	Plot   string  `dynamodbav:"Plot,omitempty"`
	Rating float32 `dynamodaav:"Rating,omitempty"`
}
