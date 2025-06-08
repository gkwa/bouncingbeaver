package dynamodb

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBResponse struct {
	Items            []map[string]json.RawMessage `json:"Items"`
	Count            int                          `json:"Count"`
	ScannedCount     int                          `json:"ScannedCount"`
	ConsumedCapacity interface{}                  `json:"ConsumedCapacity"`
}

func (c *Client) LoadData(input string) ([]map[string]types.AttributeValue, error) {
	var reader io.Reader

	if input == "-" {
		reader = os.Stdin
	} else {
		file, err := os.Open(input)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", input, err)
		}
		defer file.Close()
		reader = file
	}

	return c.LoadDataFromReader(reader)
}

func (c *Client) LoadDataFromReader(reader io.Reader) ([]map[string]types.AttributeValue, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	var response DynamoDBResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	items := make([]map[string]types.AttributeValue, 0, len(response.Items))

	for _, item := range response.Items {
		attributeMap := make(map[string]types.AttributeValue)

		for key, rawValue := range item {
			attributeValue, err := c.parseAttributeValue(rawValue)
			if err != nil {
				return nil, fmt.Errorf("failed to parse attribute %s: %w", key, err)
			}
			attributeMap[key] = attributeValue
		}

		items = append(items, attributeMap)
	}

	return items, nil
}

func (c *Client) parseAttributeValue(raw json.RawMessage) (types.AttributeValue, error) {
	var temp map[string]interface{}
	if err := json.Unmarshal(raw, &temp); err != nil {
		return nil, err
	}

	if s, exists := temp["S"]; exists {
		return &types.AttributeValueMemberS{Value: s.(string)}, nil
	}

	if n, exists := temp["N"]; exists {
		return &types.AttributeValueMemberN{Value: n.(string)}, nil
	}

	if b, exists := temp["BOOL"]; exists {
		return &types.AttributeValueMemberBOOL{Value: b.(bool)}, nil
	}

	if _, exists := temp["NULL"]; exists {
		return &types.AttributeValueMemberNULL{Value: true}, nil
	}

	return nil, fmt.Errorf("unsupported attribute value type: %v", temp)
}
