package controllers

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/bruce-mig/go-admin/db"
	"github.com/bruce-mig/go-admin/middleware"
	"github.com/bruce-mig/go-admin/models"
	"github.com/gofiber/fiber/v2"
)

func ListOrders(c *fiber.Ctx) error {
	if err := middleware.IsAuthorized(c, "orders"); err != nil {
		return err
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))

	return c.JSON(models.Paginate(db.DB, &models.Order{}, page))
}

func Export(c *fiber.Ctx) error {
	if err := middleware.IsAuthorized(c, "orders"); err != nil {
		return err
	}
	filePath := "./csv/orders.csv"

	if err := CreateFile(filePath); err != nil {
		return err
	}

	return c.Download(filePath)
}

func CreateFile(filePath string) error {
	file, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var orders []models.Order

	db.DB.Preload("OrderItems").Find(&orders)

	writer.Write([]string{
		"ID", "Name", "Email", "Product Title", "Price", "Quantity",
	})

	for _, order := range orders {
		data := []string{
			strconv.Itoa(int(order.Id)),
			order.FirstName + " " + order.LastName,
			order.Email,
			"", "", "",
		}

		if err := writer.Write(data); err != nil {
			return err
		}

		for _, orderItem := range order.OrderItems {
			data := []string{
				"", "", "",
				orderItem.ProductTitle,
				strconv.Itoa(int(orderItem.Price)),
				strconv.Itoa(int(orderItem.Quantity)),
			}
			if err := writer.Write(data); err != nil {
				return err
			}
		}

	}
	return nil
}

type Sales struct {
	Date string `json:"date"`
	Sum  string `json:"sum"`
}

func Chart(c *fiber.Ctx) error {
	if err := middleware.IsAuthorized(c, "orders"); err != nil {
		return err
	}
	var sales []Sales
	db.DB.Raw(`
		SELECT DATE_FORMAT(o.created_at, '%Y-%m-%d') as date, SUM(oi.price * oi.quantity) as sum
		FROM orders o
		JOIN order_items oi on o.id = oi.order_id
		GROUP BY date
	`).Scan(&sales)

	return c.JSON(sales)
}
