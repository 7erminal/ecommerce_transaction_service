package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"transaction_service/controllers/functions"
	"transaction_service/models"
	"transaction_service/structs/requests"
	"transaction_service/structs/responses"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// OrdersController operations for Orders
type OrdersController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrdersController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("ConfirmOrder", c.ConfirmOrder)
	c.Mapping("GetAllByBranch", c.GetAllByBranch)
	c.Mapping("GetOrderCount", c.GetOrderCount)
	c.Mapping("ReturnOrder", c.ReturnOrder)
}

// Post ...
// @Title Post
// @Description create Orders
// @Param	body		body 	models.OrdersRequestDTO	true		"body for Orders content"
// @Success 201 {int} models.Orders
// @Failure 403 body is empty
// @router / [post]
func (c *OrdersController) Post() {
	serviceName := "ORDER"
	var v requests.OrdersRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	// Request Recieved. Processing ....
	logs.Info("Request recieved is::: ", v)

	// quantity_, _ := strconv.Atoi(v.TotalQuantity)
	quantity_ := 0
	// cost_, _ := strconv.ParseFloat(v.Cost, 32)
	cost_ := 0
	currency_id := v.Currency
	created_by := v.CreatedBy

	logs.Info("Total quantity is ")

	statusCode := 608
	message := "Error processing order"
	proceed := false

	if user, err := models.GetUsersById(created_by); err == nil {
		if cur, cur_err := models.GetCurrenciesById(currency_id); cur_err == nil {
			logs.Info("Currency found")
			// logs.Info("Time is ", time.Now().Day())
			// logs.Info("Time is ", int(time.Now().Month()))
			// logs.Info("Time is ", time.Now().Year())
			// logs.Info("Time is ", time.Now().Format("20060102"))
			var orderDate time.Time = time.Now()

			var allowedDateList [6]string = [6]string{"2006-01-02", "2006/01/02", "2006-01-02 15:04:05.000", "2006/01/02 15:04:05.000", "2006-01-02T15:04:05.000Z", "2006-01-02 15:04:05.000000 -0700 MST"}

			if v.OrderDate != "" {
				for _, date_ := range allowedDateList {
					logs.Debug("About to convert ", v.OrderDate)
					// Convert dob string to date
					tOrderDate, error := time.Parse(date_, v.OrderDate)

					if error != nil {
						logs.Error("Error parsing date", error)
						orderDate = time.Now()
					} else {
						logs.Error("Date converted to time successfully", tOrderDate)
						orderDate = tOrderDate

						break
					}
				}
			}

			var orderEndDate time.Time

			for _, date_ := range allowedDateList {
				logs.Debug("About to convert ", v.OrderEndDate)
				// Convert dob string to date
				tOrderDate, error := time.Parse(date_, v.OrderEndDate)

				if error != nil {
					logs.Error("Error parsing date", error)
					orderEndDate = time.Now()
				} else {
					logs.Error("Date converted to time successfully", tOrderDate)
					orderEndDate = tOrderDate

					break
				}
			}

			var customer models.Customers
			customerId := v.Customer

			logs.Info("Customer ID is ", customerId)

			if cust, err := models.GetCustomerById(customerId); err != nil {
				logs.Error("Customer not found ", err.Error())
			} else {
				customer = *cust
			}

			var order_ = models.Orders{OrderDesc: v.RequestType, Customer: &customer, OrderLocation: v.OrderLocation, Quantity: quantity_, Cost: float32(cost_), Currency: cur.CurrencyId, OrderDate: orderDate, OrderEndDate: orderEndDate, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: user, ModifiedBy: created_by}

			// Add order
			if _, err := models.AddOrders(&order_); err == nil {
				custId := strconv.FormatInt(customer.CustomerId, 10)
				go functions.UpdateCustomer(&c.Controller, custId, time.Now().String())
				orderNumber := time.Now().Format("20060102") + strconv.FormatInt(order_.OrderId, 10)
				logs.Info("Order number is ", orderNumber)
				onum, err := strconv.ParseInt(orderNumber, 10, 64)
				if err != nil {
					logs.Error("Unable to convert order number to int")
					panic(err)
				}
				order_.OrderNumber = strconv.FormatInt(onum, 10)
				cart_items := v.Items

				logs.Info("Cart items are ", cart_items)

				amount_ := float32(0.0)
				quantity_ = 0

				orderItems := make([]*models.Order_items, len(cart_items))
				if err := models.UpdateOrdersById(&order_); err == nil {
				forLoop:
					for q, r := range cart_items {
						logs.Info("q is ", q)
						logs.Info("and r is ", r.ItemId)
						// item_id, _ := strconv.ParseInt(r.ItemId, 0, 64)
						if item, err := models.GetItemsById(r.ItemId); err == nil {
							// each_quantity_, _ := strconv.Atoi(r.Quantity)
							each_quantity_ := r.Quantity
							finalQuantity := item.Quantity

							logs.Info("Quantity is ", each_quantity_)

							if iq, err := models.GetItem_quantityByItemId(item.ItemId); err == nil {
								// if item_, item_err := models.GetItemsById(item_id); item_err == nil {
								tempQuantity := iq.Quantity
								tempQuantity = tempQuantity - int(each_quantity_)
								finalQuantity = tempQuantity

								if tempQuantity < 0 {
									logs.Error("Quantity is less ", tempQuantity)
									statusCode = 609
									message = "Item quantity is less than ordered quantity."
									proceed = false

									continue forLoop
								}

								if status, err := models.GetStatusByName("PENDING"); err == nil {
									var order_items = models.Order_items{Order: &order_, Item: item, Quantity: each_quantity_, OrderDate: time.Now(), DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: created_by, Status: status}

									logs.Info("About to add order items")
									// Add order item
									if _, err := models.AddOrder_items(&order_items); err != nil {
										logs.Error("Error adding order item::: ", err.Error())
									} else {
										// amount_ = float32(amount_) + float32(item_.ItemPrice.ItemPrice)
										logs.Info("Performing order calculations")
										iq.Quantity = tempQuantity
										if err := models.UpdateItem_quantityById(iq); err != nil {
											logs.Error("Error upating item quantity")
											message = "Error updating the item quantity"
										}
										proceed = true
										logs.Info("Quantity is ", r.Quantity)
										logs.Info("Item price is ", item.ItemPrice.ItemPrice)
										logs.Info("Amount is ", amount_)
										amount_ = float32(amount_) + (float32(item.ItemPrice.ItemPrice) * float32(r.Quantity))
										quantity_ = quantity_ + int(each_quantity_)
										// each_quantity_ = int64(quantity_)
										logs.Info("Calculations completed. Amount is ", amount_, " and quantity is ", quantity_)
										orderItems[q] = &order_items
									}
								} else {
									logs.Error("Error adding order item. Could not find status::: ", err.Error())
								}
							} else {
								logs.Error("Error adding order item. Could not find quantity::: ", err.Error())
							}
							item.Quantity = finalQuantity
							item.LastOrderDate = orderDate

							if itemPrice, err := models.GetItem_pricesById(item.ItemPrice.ItemPriceId); err == nil {
								logs.Info("Update item amount paid ", float32(item.ItemPrice.ItemPrice)*float32(r.Quantity))
								itemPrice.AmountPaid = float32(item.ItemPrice.ItemPrice) * float32(r.Quantity)
								if itemPrice.AmountPaid >= itemPrice.ItemPrice {
									itemPrice.AmountPaid = itemPrice.ItemPrice
								}

								if err := models.UpdateItem_pricesById(itemPrice); err != nil {
									logs.Error("Error updating item::: ", err.Error())
									message = "Error updating the item price"
								}
							} else {
								logs.Error("Error adding order item. Could not find quantity::: ", err.Error())
							}

							if err := models.UpdateItemsById(item); err != nil {
								logs.Error("Error updating item::: ", err.Error())
								message = "Error updating the item quantity"
							}
						} else {
							logs.Error("Could not find this item ", err.Error())
						}
					}
				} else {
					logs.Error("There was an error adding the order number ", err.Error())
					panic(err)
				}

				logs.Info("Proceed state is ", proceed)

				if proceed {
					if amount_ == 0 || quantity_ == 0 {
						// amount_ = float32(cost_)
						statusCode = 400
						message = "Invalid request. Amount or quantity provided is invalid"
						var resp = responses.OrderResponseDTO{StatusCode: statusCode, Order: nil, StatusDesc: message}
						logs.Error("Error thrown when adding transaction details::: ")
						c.Ctx.Output.SetStatus(400)
						c.Data["json"] = resp
					} else {
						order_.Cost = amount_
						order_.Quantity = quantity_

						if err := models.UpdateOrdersById(&order_); err != nil {
							logs.Info("An error occurred when updating order")
						}

						logs.Info("About to move to transactions")

						branch := models.Branches{}

						if v.Branch != 0 {
							if branch_, err := models.GetBranchesById(v.Branch); err == nil {
								branch = *branch_
							} else {
								logs.Error("Error getting branch. Continue.")
							}

						}

						if service, err := models.GetServicesByName(serviceName); err == nil {
							status_ := "PENDING"
							if status, err := models.GetStatusByName(status_); err == nil {
								var transaction_ = models.Transactions{Order: &order_, Branch: &branch, Amount: amount_, TransactingCurrency: cur.CurrencyId, Status: status, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: int(created_by), ModifiedBy: int(created_by), Services: service}
								logs.Info("About to add transaction")
								if _, txn_err := models.AddTransactions(&transaction_); txn_err == nil {
									logs.Info("NO error adding transaction")
									status_code := "1022"
									var txn_details = models.Transaction_details{TransactionId: &transaction_, Amount: amount_, Comment: v.Comment, StatusCode: status_code, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: 1, ModifiedBy: 1}

									if _, txn_d_err := models.AddTransaction_details((&txn_details)); txn_d_err == nil {
										var customOrder responses.OrdersCustom = responses.OrdersCustom{
											OrderId:      order_.OrderId,
											OrderNumber:  order_.OrderNumber,
											Quantity:     order_.Quantity,
											Cost:         order_.Cost,
											CurrencyId:   order_.Currency,
											OrderDate:    order_.OrderDate,
											DateCreated:  order_.DateCreated,
											DateModified: order_.DateModified,
											OrderEndDate: order_.OrderEndDate,
											Customer:     order_.Customer,
											OrderDetails: orderItems,
											ReturnedDate: order_.ReturnedDate,
										}
										var customTxn responses.TransactionsCustom = responses.TransactionsCustom{
											TransactionId:       transaction_.TransactionId,
											Order:               &customOrder,
											Amount:              transaction_.Amount,
											TransactingCurrency: transaction_.TransactingCurrency,
											Status:              transaction_.Status.Status,
											DateCreated:         transaction_.DateCreated,
											DateModified:        transaction_.DateModified,
											CreatedBy:           transaction_.CreatedBy,
											ModifiedBy:          transaction_.ModifiedBy,
											Active:              transaction_.Active,
										}

										fmt.Printf("custom transaction of v: %+v\n", customTxn)
										statusCode = 200
										message = "Order successfully placed"
										var resp = responses.TransactionCustomResponseDTO{StatusCode: 200, Transaction: &customTxn, StatusDesc: message}
										c.Ctx.Output.SetStatus(200)
										c.Data["json"] = resp

									} else {
										var resp = responses.OrderResponseDTO{StatusCode: 808, Order: nil, StatusDesc: "Transaction details error!"}
										logs.Error("Error thrown when adding transaction details::: ", txn_d_err.Error())
										c.Data["json"] = resp
									}
								} else {
									var resp = responses.OrderResponseDTO{StatusCode: 807, Order: nil, StatusDesc: "Transaction error!"}
									logs.Error("Error thrown when adding transaction::: ", txn_err.Error())
									c.Data["json"] = resp
								}
							} else {
								var resp = responses.OrderResponseDTO{StatusCode: 807, Order: nil, StatusDesc: "Transaction error!"}
								logs.Error("Error thrown when adding transaction::: ", err.Error())
								c.Data["json"] = resp
							}
						} else {
							var resp = responses.OrderResponseDTO{StatusCode: 807, Order: nil, StatusDesc: "Transaction error: service"}
							logs.Error("Error thrown when adding transaction::: ", err.Error())
							c.Data["json"] = resp
						}
					}
				} else {
					logs.Info("Message and code are ", message, " :: ", statusCode)
					var resp = responses.OrderResponseDTO{StatusCode: statusCode, Order: nil, StatusDesc: message}
					c.Data["json"] = resp
				}

			} else {
				var resp = responses.OrderResponseDTO{StatusCode: 806, Order: nil, StatusDesc: "Order error!"}
				logs.Error("Error thrown when adding order::: ", err.Error())
				c.Ctx.Output.SetStatus(304)
				c.Data["json"] = resp
			}
		} else {
			logs.Info("Currency received is ", v.Currency)
			logs.Info("Currency NOT found ", cur_err.Error())
		}
	} else {
		var resp = responses.OrderResponseDTO{StatusCode: 809, Order: nil, StatusDesc: "Error adding order"}
		logs.Error("Error thrown when adding order::: ", err.Error())
		c.Ctx.Output.SetStatus(304)
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// Post ...
// @Title ConfirmOrder
// @Description Confirm Orders
// @Param	body		body 	models.ConfirmOrderDTO	true		"body for Orders content"
// @Success 201 {int} models.Orders
// @Failure 403 body is empty
// @router /confirm-order [post]
func (c *OrdersController) ConfirmOrder() {
	var v requests.ConfirmOrderDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	txn_id, _ := strconv.ParseInt(v.TransactionId, 0, 64)

	logs.Info("Transaction ID is ", v.TransactionId)

	if txn, txn_err := models.GetTransactionsById(txn_id); txn_err == nil {
		// status_ := "SUCCESS"
		if status, err := models.GetStatusByName(v.Status); err == nil {
			txn.CreatedBy, _ = strconv.Atoi(v.Confirmedby)
			txn.Status = status
			txn.Active = 1
			if utxn_err := models.UpdateTransactionsById(txn); utxn_err == nil {
				var customOrder responses.OrdersCustom = responses.OrdersCustom{OrderId: txn.Order.OrderId, OrderNumber: txn.Order.OrderNumber, Quantity: txn.Order.Quantity, Cost: txn.Order.Cost, CurrencyId: txn.Order.Currency, OrderDate: txn.Order.OrderDate, DateCreated: txn.Order.DateCreated, DateModified: txn.Order.DateModified, Customer: txn.Order.Customer, OrderDetails: txn.Order.OrderDetails}
				var customTxn responses.TransactionsCustom = responses.TransactionsCustom{TransactionId: txn.TransactionId, Order: &customOrder, Amount: txn.Amount, TransactingCurrency: txn.TransactingCurrency, Status: txn.Status.Status, DateCreated: txn.DateCreated, DateModified: txn.DateModified, CreatedBy: txn.CreatedBy, ModifiedBy: txn.ModifiedBy, Active: txn.Active}

				if order, err := models.GetOrdersById(txn.Order.OrderId); err == nil {
					if order_items, err := models.GetOrder_itemsByOrder(*order); err == nil {
						logs.Info("Order items are ", order_items)
						if order_items != nil {
							for _, item := range *order_items {
								item.Status = status
								if updateOrderId := models.UpdateOrder_itemsById(&item); updateOrderId == nil {
									logs.Info("Order item updated successfully")
								} else {
									logs.Error("Error updating order item::: ", updateOrderId.Error())
								}
							}
						}
					} else {
						var resp responses.TransactionResponseDTO = responses.TransactionResponseDTO{StatusCode: 806, Transaction: nil, StatusDesc: "Order error. Unable to find order!"}
						logs.Error("Error thrown when updating transaction::: ", err.Error())
						c.Ctx.Output.SetStatus(200)
						c.Data["json"] = resp
					}

					var resp = responses.TransactionCustomResponseDTO{StatusCode: 200, Transaction: &customTxn, StatusDesc: "Order successfully placed"}
					c.Ctx.Output.SetStatus(200)
					c.Data["json"] = resp
				} else {
					var resp responses.TransactionResponseDTO = responses.TransactionResponseDTO{StatusCode: 806, Transaction: nil, StatusDesc: "Order error. Unable to find order!"}
					logs.Error("Error thrown when updating transaction::: ", utxn_err.Error())
					c.Ctx.Output.SetStatus(200)
					c.Data["json"] = resp
				}
			} else {
				var resp responses.TransactionResponseDTO = responses.TransactionResponseDTO{StatusCode: 806, Transaction: nil, StatusDesc: "Order error!"}
				logs.Error("Error thrown when updating transaction::: ", utxn_err.Error())
				c.Ctx.Output.SetStatus(200)
				c.Data["json"] = resp
			}
		}
	} else {
		var resp responses.TransactionResponseDTO = responses.TransactionResponseDTO{StatusCode: 806, Transaction: nil, StatusDesc: "Order error!"}
		logs.Error("Error fetching transaction::: ", txn_err.Error())
		c.Ctx.Output.SetStatus(304)
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// ReturnOrder ...
// @Title ReturnOrder
// @Description Return Order
// @Param	body		body 	requests.ConfirmOrderDTO	true		"body for Orders content"
// @Success 201 {int} models.Orders
// @Failure 403 body is empty
// @router /return-order [post]
func (c *OrdersController) ReturnOrder() {
	var v requests.ConfirmOrderDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	txn_id, _ := strconv.ParseInt(v.TransactionId, 0, 64)

	logs.Info("Transaction ID is ", v.TransactionId)

	if txn, txn_err := models.GetTransactionsById(txn_id); txn_err == nil {
		logs.Info("Returned date is ", txn.Order.ReturnedDate)
		// status_ := "SUCCESS"
		if status, err := models.GetStatusByName(v.Status); err == nil {
			txn.CreatedBy, _ = strconv.Atoi(v.Confirmedby)
			txn.Status = status
			txn.Active = 1
			if utxn_err := models.UpdateTransactionsById(txn); utxn_err == nil {
				var customOrder responses.OrdersCustom = responses.OrdersCustom{OrderId: txn.Order.OrderId, OrderNumber: txn.Order.OrderNumber, Quantity: txn.Order.Quantity, Cost: txn.Order.Cost, CurrencyId: txn.Order.Currency, OrderDate: txn.Order.OrderDate, DateCreated: txn.Order.DateCreated, DateModified: txn.Order.DateModified, Customer: txn.Order.Customer, OrderDetails: txn.Order.OrderDetails, ReturnedDate: txn.Order.ReturnedDate, OrderEndDate: txn.Order.OrderEndDate}
				var customTxn responses.TransactionsCustom = responses.TransactionsCustom{TransactionId: txn.TransactionId, Order: &customOrder, Amount: txn.Amount, TransactingCurrency: txn.TransactingCurrency, Status: txn.Status.Status, DateCreated: txn.DateCreated, DateModified: txn.DateModified, CreatedBy: txn.CreatedBy, ModifiedBy: txn.ModifiedBy, Active: txn.Active}

				if order, err := models.GetOrdersById(txn.Order.OrderId); err == nil {

					order.ReturnedDate = time.Now()
					customOrder.ReturnedDate = order.ReturnedDate
					logs.Info("Updating return date to ", order.ReturnedDate)
					if err := models.UpdateOrdersById(order); err != nil {
						logs.Error("Error updating order::: ", err.Error())
					}
					logs.Info("Order items are ", order)
					if order_items, err := models.GetOrder_itemsByOrder(*order); err == nil {
						logs.Info("Order items are ", order_items)
						if order_items != nil {
							for _, item := range *order_items {
								item.Status = status
								if oItem, err := models.GetItem_quantityByItemId(item.Item.ItemId); err == nil {
									logs.Info("Item quantity is ", oItem.Quantity)
									logs.Info("Item quantity is ", item.Quantity)
									oItem.Quantity = oItem.Quantity + int(item.Quantity)
									if err := models.UpdateItem_quantityById(oItem); err != nil {
										logs.Error("Error updating item::: ", err.Error())
										message := "Error updating the item quantity"
										var resp responses.TransactionResponseDTO = responses.TransactionResponseDTO{StatusCode: 806, Transaction: nil, StatusDesc: message}
										logs.Error("Error thrown when updating transaction::: ", err.Error())
										c.Ctx.Output.SetStatus(200)
										c.Data["json"] = resp
										c.ServeJSON()
									} else {
										logs.Info("Item quantity updated successfully")
										item.Item.Quantity = int(item.Quantity) + item.Item.Quantity
										if err := models.UpdateItemsById(item.Item); err != nil {
											logs.Error("Error updating item::: ", err.Error())
											message := "Error updating the item quantity"
											var resp responses.TransactionResponseDTO = responses.TransactionResponseDTO{StatusCode: 806, Transaction: nil, StatusDesc: message}
											logs.Error("Error thrown when updating transaction::: ", err.Error())
											c.Ctx.Output.SetStatus(200)
											c.Data["json"] = resp
											c.ServeJSON()
										} else {
											logs.Info("Item updated successfully")
										}
									}
								}
								item.Item.Quantity = int(item.Quantity) + item.Item.Quantity
								if updateOrderId := models.UpdateOrder_itemsById(&item); updateOrderId == nil {
									logs.Info("Order item updated successfully")
								} else {
									logs.Error("Error updating order item::: ", updateOrderId.Error())
								}
							}
						}
					} else {
						var resp responses.TransactionResponseDTO = responses.TransactionResponseDTO{StatusCode: 806, Transaction: nil, StatusDesc: "Order error. Unable to find order!"}
						logs.Error("Error thrown when updating transaction::: ", err.Error())
						c.Ctx.Output.SetStatus(200)
						c.Data["json"] = resp
					}

					logs.Info("Order returned date is ", order.ReturnedDate)

					var resp = responses.TransactionCustomResponseDTO{StatusCode: 200, Transaction: &customTxn, StatusDesc: "Order successfully placed"}
					c.Ctx.Output.SetStatus(200)
					c.Data["json"] = resp
				} else {
					var resp responses.TransactionResponseDTO = responses.TransactionResponseDTO{StatusCode: 806, Transaction: nil, StatusDesc: "Order error. Unable to find order!"}
					logs.Error("Error thrown when updating transaction::: ", utxn_err.Error())
					c.Ctx.Output.SetStatus(200)
					c.Data["json"] = resp
				}
			} else {
				var resp responses.TransactionResponseDTO = responses.TransactionResponseDTO{StatusCode: 806, Transaction: nil, StatusDesc: "Order error!"}
				logs.Error("Error thrown when updating transaction::: ", utxn_err.Error())
				c.Ctx.Output.SetStatus(200)
				c.Data["json"] = resp
			}
		}
	} else {
		var resp responses.TransactionResponseDTO = responses.TransactionResponseDTO{StatusCode: 806, Transaction: nil, StatusDesc: "Order error!"}
		logs.Error("Error fetching transaction::: ", txn_err.Error())
		c.Ctx.Output.SetStatus(304)
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// Post ...
// @Title Get User Orders
// @Description get user orders
// @Param	body		body 	requests.GetUserOrdersRequest	true		"body for Transactions content"
// @Success 201 {int} models.OrdersResponseDTO
// @Failure 403 body is empty
// @router /get-user-orders [post]
func (c *OrdersController) GetUserOrders() {
	var v requests.GetUserOrdersRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if orders, err := models.GetOrdersByUser(v.Id); err == nil {
		logs.Debug("Item ID to get quantity is ", orders)
		var resp = responses.OrdersResponseDTO{StatusCode: 200, Orders: orders, StatusDesc: "Order details fetched successfully"}
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = resp
	} else {
		logs.Debug("An error occurred getting orders ", err.Error())
		var resp = responses.OrdersResponseDTO{StatusCode: 608, Orders: nil, StatusDesc: "Failed to fetch orders"}
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Orders by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Orders
// @Failure 403 :id is empty
// @router /:id [get]
func (c *OrdersController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetOrdersById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Orders
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Orders
// @Failure 403
// @router / [get]
func (c *OrdersController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 100
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllOrders(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// GetAllByBranch ...
// @Title Get All Orders by Branch
// @Description get Orders by branch
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Orders
// @Failure 403
// @router /branch/:id [get]
func (c *OrdersController) GetAllByBranch() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllOrders(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Orders
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Orders	true		"body for Orders content"
// @Success 200 {object} models.Orders
// @Failure 403 :id is not int
// @router /:id [put]
func (c *OrdersController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Orders{OrderId: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateOrdersById(&v); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Orders
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *OrdersController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteOrders(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetItemCount ...
// @Title Get Item Quantity
// @Description get Item_quantity by Item id
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Success 200 {object} responses.StringResponseDTO
// @Failure 403 :id is empty
// @router /count/ [get]
func (c *OrdersController) GetOrderCount() {
	// q, err := models.GetItemsById(id)

	var query = make(map[string]string)

	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	v, err := models.GetOrderCount(query)
	count := strconv.FormatInt(v, 10)

	if err != nil {
		logs.Error("Error fetching count of customers ... ", err.Error())
		resp := responses.StringResponseDTO{StatusCode: 301, Value: "", StatusDesc: err.Error()}
		c.Data["json"] = resp
	} else {
		resp := responses.StringResponseDTO{StatusCode: 200, Value: count, StatusDesc: "Count fetched successfully"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}
