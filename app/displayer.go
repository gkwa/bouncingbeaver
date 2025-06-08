package app

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gkwa/bouncingbeaver/internal/logger"
	"github.com/gkwa/bouncingbeaver/internal/models"
)

type Displayer struct {
	logger *logger.Logger
}

func NewDisplayer(logger *logger.Logger) *Displayer {
	return &Displayer{
		logger: logger,
	}
}

func (d *Displayer) ShowProducts(products []models.Product) {
	d.logger.Debug("Displaying products", "count", len(products))

	// Create an encoder that doesn't escape HTML
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	err := encoder.Encode(products)
	if err != nil {
		d.logger.Error("Failed to marshal products to JSON", "error", err)
		return
	}

	fmt.Print(buf.String())
}
