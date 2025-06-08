package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gkwa/bouncingbeaver/internal/logger"
	"github.com/gkwa/bouncingbeaver/internal/models"
	"github.com/gkwa/bouncingbeaver/internal/processing"
)

type Client struct {
	htmlExtractor *processing.HTMLExtractor
	logger        *logger.Logger
}

func NewClient() *Client {
	return &Client{
		htmlExtractor: processing.NewHTMLExtractor(),
		logger:        logger.New(0), // Basic logger for debugging
	}
}

func (c *Client) UnmarshalProducts(items []map[string]types.AttributeValue) ([]models.Product, error) {
	var products []models.Product
	err := attributevalue.UnmarshalListOfMaps(items, &products)
	if err != nil {
		return nil, err
	}

	// Post-process to extract HTML
	for i := range products {
		c.logger.Debug("Processing product", "id", products[i].ID, "rawhtml_length", len(products[i].RawHTML))

		if products[i].RawHTML == "" {
			products[i].RawHTMLExtracted = "NO_RAW_HTML_DATA"
			c.logger.Debug("No raw HTML data for product", "id", products[i].ID)
			continue
		}

		extractedHTML, err := c.htmlExtractor.ExtractHTML(products[i].RawHTML)
		if err != nil {
			products[i].RawHTMLExtracted = "EXTRACTION_FAILED: " + err.Error()
			c.logger.Error("HTML extraction failed", "id", products[i].ID, "error", err)
		} else {
			products[i].RawHTMLExtracted = extractedHTML
			c.logger.Debug("HTML extraction successful", "id", products[i].ID, "extracted_length", len(extractedHTML))
		}
	}

	return products, nil
}
