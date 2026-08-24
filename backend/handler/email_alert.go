package handler

import "github.com/gofiber/fiber/v3"

func EmailAlert(c fiber.Ctx) error {
	return c.Render("email-alert", fiber.Map{
		"folio":              "1234",
		"date":               "03/May/26",
		"facilityLabel":      "Jerry's house",
		"vehicleType":        "Una Canionetita",
		"driver":             "Jerry",
		"step":               "7.- Piso Interior",
		"user":               "GSuarez",
		"plates":             "773SUZ",
		"boxNumber":          "s/n",
		"truckTransportline": "Transportes Terrestres Trejo",
		"boxTransportline":   "Transportes Terrestres Trejo",
		"notes":              "Canionetita con 'gomito' en los asientos",
	})
}
