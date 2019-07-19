package service

import (
	"bytes"
	"database/sql"
	"fmt"
	"ggfly/dao"
	"github.com/llairunhao/abiu"
	"github.com/sirupsen/logrus"
	"time"
)

type ProductService interface {
	SaveProduct(requestId string, product *dao.Product) error

	SaveProducts(requestId string, products []*dao.Product) error

	UpdateProduct(requestId string, product *dao.Product) error

	ListOfProducts(requestId string, request *dao.ProductListRequest) ([]*dao.Product, error)

	RemoveProduct(requestId string, product *dao.Product) error
}

type productService struct {
	db *sql.DB
}

func (service productService) SaveProduct(requestId string, product *dao.Product) error {
	sqlString := service.insetSqlString()

	stmt, err := service.db.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	res, err := service.execInset(stmt, product)
	if err != nil {
		return err
	}

	productID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	product.ProductID = abiu.NewInt64(productID)
	return nil
}

func (service productService) UpdateProduct(requestId string, product *dao.Product) error {
	sqlString := service.updateSqlString(product)
	fmt.Println(sqlString)

	stmt, err := service.db.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	_, err = service.execUpdate(stmt, product)
	if err != nil {
		return err
	}
	return nil
}

func (service productService) SaveProducts(requestId string, products []*dao.Product) error {
	tx, err := service.db.Begin()
	if err != nil {
		return err
	}
	sqlString := service.insetSqlString()
	stmt, err := tx.Prepare(sqlString)
	if err != nil {
		return err
	}
	for _, product := range products {
		_, err := service.execInset(stmt, product)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (service productService) RemoveProduct(requestId string, product *dao.Product) error {
	sqlString := `DELETE FROM product WHERE id = ?`
	stmt, err := service.db.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	_, err = stmt.Exec(product.ProductID.Value())
	return err
}

func (service productService) ListOfProducts(requestId string, request *dao.ProductListRequest) ([]*dao.Product,
	error) {
	sqlString := `SELECT
						id,
						product_name,
						max_price,
						min_price,
						shipping,
						rate,
						review_count,
						order_count,
       					origin_url,
						gmt_create,
						gmt_modified 
					FROM
						product 
						LIMIT ?,
						?`
	stmt, err := service.db.Prepare(sqlString)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	rows, err := stmt.Query(request.PageSize, request.Page)
	if err != nil {
		return nil, err
	}

	var products = make([]*dao.Product, request.PageSize)
	var i = 0
	for rows.Next() {
		p := dao.NewProduct()
		err = rows.Scan(
			p.ProductID.ValuePtr(),
			p.MaxPrice.ValuePtr(),
			p.MinPrice.ValuePtr(),
			p.Shipping.ValuePtr(),
			p.Rate.ValuePtr(),
			p.ReviewCount.ValuePtr(),
			p.OrderCount.ValuePtr(),
			p.OriginURL.ValuePtr(),
			p.CreateDate.ValuePtr(),
			p.ModifiedDate.ValuePtr(),
		)
		if err != nil {
			return nil, err
		}
		products[i] = p
		i += 1
	}
	return products[:i], nil
}

func (service productService) insetSqlString() string {
	return `INSERT INTO product 
			SET	
			product_name=?, 
			max_price=?, 
			min_price=?, 
			shipping=?, 
			rate=?, 
			review_count=?, 
			order_count=?,
			origin_url=?,
			gmt_create=?, 
			gmt_modified=?`
}

func (service productService) execInset(stmt *sql.Stmt, product *dao.Product) (sql.Result, error) {
	return stmt.Exec(product.Name.Value(),
		product.MaxPrice.Value(),
		product.MinPrice.Value(),
		product.Shipping.Value(),
		product.Rate.Value(),
		product.ReviewCount.Value(),
		product.OrderCount.Value(),
		product.OriginURL.Value(),
		time.Now(),
		time.Now())
}

func (service productService) updateSqlString(product *dao.Product) string {
	var buffer bytes.Buffer
	buffer.WriteString("UPDATE product SET ")
	if product.Name != nil {
		buffer.WriteString("product_name = ?, ")
	}
	if product.MaxPrice != nil {
		buffer.WriteString("max_price = ?, ")
	}
	if product.MinPrice != nil {
		buffer.WriteString("min_price = ?, ")
	}
	if product.Shipping != nil {
		buffer.WriteString("shipping = ?, ")
	}
	if product.Rate != nil {
		buffer.WriteString("rate = ?, ")
	}
	if product.ReviewCount != nil {
		buffer.WriteString("review_count = ?, ")
	}
	if product.OrderCount != nil {
		buffer.WriteString("order_count = ?, ")
	}
	if product.OriginURL != nil {
		buffer.WriteString("origin_url = ?, ")
	}
	b := buffer.Bytes()[:buffer.Len()-2]
	return fmt.Sprint(string(b), " WHERE id = ?")

}

func (service productService) execUpdate(stmt *sql.Stmt, product *dao.Product) (sql.Result, error) {
	var args = make([]interface{}, 9)
	var i = 0
	if product.Name != nil {
		args[i] = product.Name.Value()
		i += 1
	}
	if product.MaxPrice != nil {
		args[i] = product.MaxPrice.Value()
		i += 1
	}
	if product.MinPrice != nil {
		args[i] = product.MinPrice.Value()
		i += 1
	}
	if product.Shipping != nil {
		args[i] = product.Shipping.Value()
		i += 1
	}
	if product.Rate != nil {
		args[i] = product.Rate.Value()
		i += 1
	}
	if product.ReviewCount != nil {
		args[i] = product.ReviewCount.Value()
		i += 1
	}
	if product.OrderCount != nil {
		args[i] = product.OrderCount.Value()
		i += 1
	}
	args[i] = product.ProductID.Value()
	i += 1

	return stmt.Exec(args[:i]...)
}
