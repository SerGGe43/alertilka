package controller

import (
	"fmt"
	"github.com/SerGGe43/alertilka/internal/domain"
)

func (c *Controller) addIndicator(alertID int64) error {
	_, err := c.IndicatorDB.Add(domain.Indicator{
		AlertID:     alertID,
		IndicatorID: "",
		Value:       -1,
	})
	if err != nil {
		return fmt.Errorf("can't add Indicator: %w", err)
	}
	return nil
}

func (c *Controller) addIndicatorId(e domain.Event) error {
	return nil
}
