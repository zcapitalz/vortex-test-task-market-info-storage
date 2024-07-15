package storages

import (
	"bytes"
	"fmt"
	"market-info-storage/internal/domain"
	"strconv"

	"github.com/pkg/errors"
)

type DepthOrders []domain.DepthOrder

func (a *DepthOrders) Scan(value interface{}) error {
	arrayBytes, ok := value.([]byte) // value has a form of {\"(0.010782342,24)\",\"(0.010765101,11.5)\"}
	if !ok {
		return fmt.Errorf("failed to convert value to []byte")
	}
	arrayBytes = arrayBytes[2 : len(arrayBytes)-2] // trim {\" and \"}

	splitArrayBytes := bytes.Split(arrayBytes, []byte{'"', ',', '"'})
	res := make([]domain.DepthOrder, 0, len(splitArrayBytes))
	for _, depthOrderBytes := range splitArrayBytes {
		depthOrderBytes = depthOrderBytes[1 : len(depthOrderBytes)-1] // trim ( and )
		splitDepthOrderBytes := bytes.Split(depthOrderBytes, []byte{','})

		var err error
		var depthOrder domain.DepthOrder
		depthOrder.Price, err = strconv.ParseFloat(string(splitDepthOrderBytes[0]), 64)
		if err != nil {
			return errors.Wrap(err, "parse DepthOrder.Price field")
		}
		depthOrder.BaseQty, err = strconv.ParseFloat(string(splitDepthOrderBytes[1]), 64)
		if err != nil {
			return errors.Wrap(err, "parse DepthOrder.BaseQty field")
		}

		res = append(res, depthOrder)
	}

	*a = res
	return nil
}
