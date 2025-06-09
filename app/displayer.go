package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

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

func (d *Displayer) ShowProducts(products []models.Product, randomize bool) {
	d.logger.Debug("Displaying products", "count", len(products), "randomize", randomize)

	// Randomize the order if requested
	if randomize {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(products), func(i, j int) {
			products[i], products[j] = products[j], products[i]
		})
		d.logger.Debug("Products randomized")
	}

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
