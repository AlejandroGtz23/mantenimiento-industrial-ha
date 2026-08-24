package handler

import (
	"fmt"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
)

func Upload(c fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	if token := form.Value["token"]; len(token) > 0 {
		fmt.Println(token[0])
	}
	//
	files := form.File["docs"]
	for _, file := range files {
		dst := filepath.Join("./docs", file.Filename)
		if err := c.SaveFile(file, dst); err != nil {
			return err
		}
	}
	return c.SendString("success")
}

func Upload2(c fiber.Ctx) error {
	if form, err := c.MultipartForm(); err == nil {
		if token := form.Value["token"]; len(token) > 0 {
			fmt.Println(token[0])
		}

		files := form.File["docs"]
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

			if err := c.SaveFile(file, fmt.Sprintf("./docs/%s", file.Filename)); err != nil {
				return err
			}
		}
		return c.SendString("success")
	}
	return c.SendStatus(500)
}
