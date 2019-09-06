resource "aws_dynamodb_table" "movies-dynamodb-table" {
  name           = "movies"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "Title"
  range_key      = "Year"

  attribute {
    name = "Title"
    type = "S"
  }
  attribute {
    name = "Year"
    type = "N"
  }
}
