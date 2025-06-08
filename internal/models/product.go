package models

type Product struct {
	ID               string `dynamodbav:"id"`
	Name             string `dynamodbav:"name"`
	Price            string `dynamodbav:"price"`
	Category         string `dynamodbav:"category"`
	Domain           string `dynamodbav:"domain"`
	ImageURL         string `dynamodbav:"imageUrl"`
	PricePerUnit     string `dynamodbav:"pricePerUnit"`
	EntityType       string `dynamodbav:"entity_type"`
	Timestamp        string `dynamodbav:"timestamp"`
	URL              string `dynamodbav:"url"`
	RawTextContent   string `dynamodbav:"rawTextContent"`
	RawHTML          string `dynamodbav:"rawHtml"`
	RawHTMLExtracted string `json:"RawHTMLExtracted"`
	TTL              int64  `dynamodbav:"ttl"`
}
