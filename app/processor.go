package app

import (
	"github.com/gkwa/bouncingbeaver/internal/dynamodb"
	"github.com/gkwa/bouncingbeaver/internal/logger"
)

type Processor struct {
	logger   *logger.Logger
	dynamodb *dynamodb.Client
}

func NewProcessor(verbosity int) *Processor {
	return &Processor{
		logger:   logger.New(verbosity),
		dynamodb: dynamodb.NewClient(),
	}
}

func (p *Processor) ProcessData(inputFile string) error {
	p.logger.Info("Processing DynamoDB data", "input", inputFile)

	sampleData, err := p.dynamodb.LoadData(inputFile)
	if err != nil {
		p.logger.Error("Failed to load data", "error", err, "input", inputFile)
		return err
	}

	products, err := p.dynamodb.UnmarshalProducts(sampleData)
	if err != nil {
		p.logger.Error("Failed to unmarshal products", "error", err)
		return err
	}

	p.logger.Debug("Successfully unmarshaled products", "count", len(products))

	displayer := NewDisplayer(p.logger)
	displayer.ShowProducts(products)

	return nil
}
