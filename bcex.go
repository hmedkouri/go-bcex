package bcex

import (
	"errors"
	"fmt"
)

const (
	API_BASE = "https://api.blockchain.com/v3/exchange" // BCEX API endpoint
)

type Bcex struct {
	client *client
}

// New returns an instantiated HitBTC struct
func New(apiKey, apiSecret string) *Bcex {
	client := NewClient(apiKey, apiSecret)
	return &Bcex{client}
}

// handleErr gets JSON response from API and deal with error
func handleErr(r interface{}) error {
	switch v := r.(type) {
	case map[string]interface{}:
		error := r.(map[string]interface{})["error"]
		if error != nil {
			switch v := error.(type) {
			case map[string]interface{}:
				errorMessage := error.(map[string]interface{})["message"]
				return errors.New(errorMessage.(string))
			default:
				return fmt.Errorf("I don't know about type %T!\n", v)
			}
		}
	case []interface{}:
		return nil
	default:
		return fmt.Errorf("I don't know about type %T!\n", v)
	}

	return nil
}