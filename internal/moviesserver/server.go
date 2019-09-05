package moviesserver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/moz5691/twirp-lambda/internal/awsutils"
	"github.com/moz5691/twirp-lambda/internal/model"
	rpc "github.com/moz5691/twirp-lambda/rpc/movies"
	log "github.com/sirupsen/logrus"
	"github.com/twitchtv/twirp"
)

// Server implements the Image service
type Server struct{}

func (s *Server) GetByTitle(ctx context.Context, req *rpc.SearchTerm) (*rpc.Item, error) {

	d := &model.Movie{}

	av, err := awsutils.MarshalMap(req)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		Key:       av,
		TableName: &awsutils.DynamoTableName,
	}

	res, err := awsutils.GetItem(input)

	if err != nil {
		fmt.Printf("error: %+v\n", err)
		return nil, twirp.NewError(twirp.NotFound, "Error")

	}
	if len(res.Item) == 0 {
		return nil, twirp.NewError(twirp.NotFound, "No movie, length is zero")
	}

	err = awsutils.UnmarshalMap(res.Item, d)
	if err != nil {
		return nil, err
	}
	return &rpc.Item{
		Title:  d.Title,
		Year:   d.Year,
		Plot:   d.Plot,
		Rating: d.Rating,
	}, nil

}

func (s *Server) DeleteByTitle(ctx context.Context, req *rpc.SearchTerm) (*rpc.SearchTerm, error) {
	av, err := awsutils.MarshalMap(req)

	if err != nil {
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		Key:       av,
		TableName: &awsutils.DynamoTableName,
	}
	_, err1 := awsutils.DeleteItem(input)

	if err1 != nil {
		fmt.Printf("error: %+v\n", err1)
		return nil, twirp.NewError(twirp.NotFound, "Error")
	}

	return &rpc.SearchTerm{
		Title: req.Title,
		Year:  req.Year,
	}, nil
}

func (s *Server) Update(ctx context.Context, req *rpc.Item) (*rpc.Item, error) {
	year := int(req.GetYear())
	rating := float64(req.GetRating())

	params := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {N: aws.String(strconv.FormatFloat(rating, 'E', -1, 32))},
			":p": {S: aws.String(req.GetPlot())},
		},
		TableName: &awsutils.DynamoTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Title": {
				S: aws.String(req.GetTitle()),
			},
			"Year": {
				N: aws.String(strconv.Itoa(year)),
			},
		},
		ReturnValues:     aws.String(dynamodb.ReturnValueUpdatedNew),
		UpdateExpression: aws.String("set Rating=:r, Plot=:p"),
	}

	resp, err := awsutils.UpdateItem(params)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, twirp.NewError(twirp.NotFound, "Error")
	}

	// set up interface to unmarshal resp.
	movie := &model.Movie{}
	err = awsutils.UnmarshalMap(resp.Attributes, &movie)

	if err != nil {
		fmt.Println("Unmarshal failed")
	}

	return &rpc.Item{
		Title:  req.Title,
		Year:   req.Year,
		Plot:   movie.Plot,
		Rating: movie.Rating,
	}, nil
}

func (s *Server) Post(ctx context.Context, req *rpc.Item) (*rpc.Item, error) {
	d := &model.Movie{
		Title:  req.GetTitle(),
		Year:   req.GetYear(),
		Plot:   req.GetPlot(),
		Rating: req.GetRating(),
	}
	av, err := awsutils.MarshalMap(d)
	if err != nil {
		return nil, err
	}

	putInput := &dynamodb.PutItemInput{
		Item: av,
	}
	if _, err = awsutils.PutItem(putInput); err != nil {
		log.WithFields(log.Fields{
			"cause":  err,
			"videos": fmt.Sprintf("%+v", d),
		}).Error("Error from Post Item")
		return nil, err
	}
	return &rpc.Item{
		Title:  req.GetTitle(),
		Year:   req.GetYear(),
		Plot:   req.GetPlot(),
		Rating: req.GetRating(),
	}, nil
}
