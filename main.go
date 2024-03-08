package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	// "time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Resp struct {
	UxpCode  string `json:"uxp-code"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Size     int    `json:"filesize"`
}

var resp = &Resp{
	UxpCode:  "0",
	Message:  "",
	Filename: "",
	Size:     0,
}

func upload(c echo.Context) error {
	filename := c.Param("name")
	err := writeFile(filename, c.Request().Body)
	if err != nil {
		resp.UxpCode = "1"
		resp.Message = fmt.Sprintf("%v", err)
		c.JSON(http.StatusInternalServerError, resp)
	}

	resp.Filename = filename
	resp.Size = int(c.Request().ContentLength)
	return c.JSON(http.StatusOK, resp)
}

func writeFile(filename string, src io.Reader) error {

	dst, err := os.Create(filepath.Join("uploads", filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

func download(c echo.Context) error {
	filename := c.Param("name")
	return c.File(filepath.Join("uploads", filename))
}

// func downloadT(c echo.Context) error {
// 	filename := c.Param("name")
// 	filePath := filepath.Join("uploads", filename)

// 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 		return echo.NewHTTPError(http.StatusNotFound, "File not found")
// 	}

// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 	}
// 	defer file.Close()

// 	c.Response().Header().Set("Content-Type", "application/octet-stream")

// 	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)

// 	return c.Stream(http.StatusOK, "application/octet-stream", file)
// }

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.PUT("/:name", upload)
	e.GET("/:name", download)
	// e.GET("/2/:name", downloadT)

	e.Logger.Fatal(e.Start(":1323"))
}
