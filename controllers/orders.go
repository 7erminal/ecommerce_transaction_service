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

	if user, err := functions.GetUser(&c.Controller, requests.GetUserRequest{UserId: created_by}); err == nil {
		if cur, cur_err := functions.GetCurrency(&c.Controller, requests.GetCurrencyRequest{CurrencyId: currency_id}); cur_err == nil {
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

			var customer responses.Customers
			customerId := v.Customer

			logs.Info("Customer ID is ", customerId)

			if cust, err := functions.GetCustomer(&c.Controller, requests.GetCustomerRequest{CustomerId: customerId}); err != nil {
				logs.Error("Customer not found ", err.Error())
			} else {
				customer = *cust.Result
			}

			customerIdStr := strconv.FormatInt(customer.CustomerId, 10)
			userIdStr := strconv.FormatInt(user.User.UserId, 10)

			var order_ = models.Orders{
				OrderDesc:     v.RequestType,
				CustomerId:    customerIdStr,
				CustomerName:  customer.FullName,
				CustomerEmail: customer.Email,
				CustomerPhone: customer.PhoneNumber,
				OrderLocation: v.OrderLocation,
				Quantity:      quantity_,
				Cost:          float32(cost_),
				Currency:      cur.Result.Symbol,
				OrderDate:     orderDate,
				OrderEndDate:  orderEndDate,
				DateCreated:   time.Now(),
				DateModified:  time.Now(),
				CreatedBy:     userIdStr,
				ModifiedBy:    userIdStr}

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
						if item, err := functions.GetItem(&c.Controller, requests.GetItemRequest{ItemId: r.ItemId}); err == nil {
							// each_quantity_, _ := strconv.Atoi(r.Quantity)
							each_quantity_ := r.Quantity
							finalQuantity := item.Item.Quantity

							logs.Info("Quantity is ", each_quantity_)

							// if item_, item_err := models.GetItemsById(item_id); item_err == nil {
							tempQuantity := item.Item.ItemQuantity.Quantity
							tempQuantity = tempQuantity - int(each_quantity_)
							finalQuantity = tempQuantity
							totalPrice := item.Item.ItemPrice.ItemPrice * float32(each_quantity_)

							if tempQuantity < 0 {
								logs.Error("Quantity is less ", tempQuantity)
								statusCode = 609
								message = "Item quantity is less than ordered quantity."
								proceed = false

								continue forLoop
							}

							if status, err := models.GetStatusByName("PENDING"); err == nil {
								itemIdStr := strconv.FormatInt(item.Item.ItemId, 10)
								var order_items = models.Order_items{
									Order:        &order_,
									Item:         itemIdStr,
									ItemName:     item.Item.ItemName,
									Category:     item.Item.Category.CategoryName,
									UnitPrice:    item.Item.ItemPrice.ItemPrice,
									Quantity:     each_quantity_,
									TotalPrice:   totalPrice,
									OrderDate:    time.Now(),
									DateCreated:  time.Now(),
									DateModified: time.Now(),
									CreatedBy:    userIdStr,
									Status:       status}

								logs.Info("About to add order items")
								// Add order item
								if _, err := models.AddOrder_items(&order_items); err != nil {
									logs.Error("Error adding order item::: ", err.Error())
								} else {
									// amount_ = float32(amount_) + float32(item_.ItemPrice.ItemPrice)
									logs.Info("Performing order calculations")
									item.Item.ItemQuantity.Quantity = tempQuantity
									if itemuq, err := functions.UpdateItemQuantity(&c.Controller, requests.UpdateItemQuantityRequest{ItemId: strconv.FormatInt(item.Item.ItemId, 10), Quantity: each_quantity_}); err != nil {
										logs.Error("Error upating item quantity")
										message = "Error updating the item quantity"
									} else {
										proceed = true
										logs.Info("Quantity is ", itemuq.Item.Quantity)
										logs.Info("Item price is ", item.Item.ItemPrice)
										logs.Info("Amount is ", amount_)
										amount_ = float32(amount_) + (float32(item.Item.ItemPrice.ItemPrice) * float32(r.Quantity))
										quantity_ = quantity_ + int(each_quantity_)
										// each_quantity_ = int64(quantity_)
										logs.Info("Calculations completed. Amount is ", amount_, " and quantity is ", quantity_)
										orderItems[q] = &order_items
									}
								}
							} else {
								message = "Error adding order item. Could not find status"
								statusCode = 612
								logs.Error("Error adding order item. Could not find status::: ", err.Error())
							}

							item.Item.Quantity = finalQuantity
							item.Item.LastOrderDate = orderDate

							logs.Info("Update item amount paid ", float32(item.Item.ItemPrice.ItemPrice)*float32(r.Quantity))
							// itemPrice.AmountPaid = itemPrice.AmountPaid + (float32(item.Item.ItemPrice.ItemPrice) * float32(r.Quantity))
							// if itemPrice.AmountPaid >= itemPrice.AltItemPrice {
							// 	itemPrice.AmountPaid = itemPrice.ItemPrice
							// }

							// item.Item.ItemPrice.Discount = itemPrice.Discount

							// if err := models.UpdateItem_pricesById(itemPrice); err != nil {
							// 	logs.Error("Error updating item::: ", err.Error())
							// 	message = "Error updating the item price"
							// }

							// if err := models.UpdateItemsById(item); err != nil {
							// 	logs.Error("Error updating item::: ", err.Error())
							// 	message = "Error updating the item quantity"
							// }
						} else {
							message = "Error adding order item. Could not find item"
							statusCode = 610
							logs.Error("Could not find this item ", err.Error())
						}
					}
				} else {
					message = "Error updating order"
					statusCode = 611
					logs.Error("There was an error adding the order number ", err.Error())
					panic(err)
				}

				logs.Info("Proceed state is ", proceed)

				if proceed {
					order_.Cost = amount_
					order_.Quantity = quantity_

					if err := models.UpdateOrdersById(&order_); err != nil {
						logs.Info("An error occurred when updating order")
					}

					logs.Info("About to move to transactions")

					branch := responses.Branches{}

					if v.Branch != "" {
						if branch_, err := functions.GetBranch(&c.Controller, requests.GetBranchRequest{BranchId: v.Branch}); err == nil {
							branch = *branch_.Branch
						} else {
							logs.Error("Error getting branch. Continue.")
						}
					}

					if service_, err := functions.GetService(&c.Controller, requests.GetServiceRequest{ServiceId: serviceName}); err == nil {
						service := service_.Service
						status_ := "PENDING"
						if status, err := models.GetStatusByName(status_); err == nil {
							branchStr := strconv.FormatInt(branch.BranchId, 10)
							currencyStr := strconv.FormatInt(cur.Result.CurrencyId, 10)
							// convert created_by to int
							serviceIdStr := strconv.FormatInt(service.ServiceId, 10)
							var transaction_ = models.Transactions{
								Order:          &order_,
								BranchId:       branchStr,
								BranchName:     branch.Branch,
								Amount:         amount_,
								CurrencyId:     currencyStr,
								CurrencySymbol: cur.Result.Symbol,
								Status:         status,
								DateCreated:    time.Now(),
								DateModified:   time.Now(),
								CreatedBy:      created_by,
								ModifiedBy:     created_by,
								ServiceId:      serviceIdStr,
								ServiceName:    service.ServiceName}
							logs.Info("About to add transaction")
							if _, txn_err := models.AddTransactions(&transaction_); txn_err == nil {
								logs.Info("NO error adding transaction")
								status_code := "1022"
								var txn_details = models.Transaction_details{TransactionId: &transaction_, Amount: amount_, Comment: v.Comment, StatusCode: status_code, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: created_by, ModifiedBy: created_by}

								if _, txn_d_err := models.AddTransaction_details((&txn_details)); txn_d_err == nil {
									customerIdStr := strconv.FormatInt(customer.CustomerId, 10)
									customerData := responses.CustomersAlt{
										CustomerId:    customerIdStr,
										CustomerName:  customer.FullName,
										CustomerEmail: customer.Email,
										CustomerPhone: customer.PhoneNumber}

									orderDetails := []responses.OrderItemsCustom{}
									if transaction_.Order.OrderDetails != nil {
										for _, orderItem := range transaction_.Order.OrderDetails {
											itemData := responses.ItemAlt{
												ItemId:       orderItem.Item,
												ItemName:     orderItem.ItemName,
												Description:  "",
												Price:        orderItem.TotalPrice,
												Currency:     orderItem.Order.Currency,
												Category:     orderItem.Category,
												DateCreated:  orderItem.DateCreated,
												DateModified: orderItem.DateModified}

											orderIdStr := strconv.FormatInt(orderItem.Order.OrderId, 10)
											orderItemsCustom := responses.OrderItemsCustom{
												OrderItemId: orderItem.OrderItemId,
												OrderId:     orderIdStr,
												Item:        &itemData,
												Quantity:    orderItem.Quantity,
												OrderDate:   orderItem.OrderDate,
												Status:      orderItem.Status.Status}

											orderDetails = append(orderDetails, orderItemsCustom)
										}
									}
									var customOrder responses.OrdersCustom = responses.OrdersCustom{
										OrderId:      order_.OrderId,
										OrderNumber:  order_.OrderNumber,
										Quantity:     order_.Quantity,
										Cost:         order_.Cost,
										Currency:     order_.Currency,
										OrderDate:    order_.OrderDate,
										DateCreated:  order_.DateCreated,
										DateModified: order_.DateModified,
										OrderEndDate: order_.OrderEndDate,
										Customer:     &customerData,
										OrderDetails: orderDetails,
										ReturnedDate: order_.ReturnedDate,
									}
									var customTxn responses.TransactionsCustom = responses.TransactionsCustom{
										TransactionId:       transaction_.TransactionId,
										Order:               &customOrder,
										Amount:              transaction_.Amount,
										TransactingCurrency: transaction_.CurrencySymbol,
										Status:              transaction_.Status.Status,
										DateCreated:         transaction_.DateCreated,
										DateModified:        transaction_.DateModified,
										CreatedBy:           transaction_.CreatedBy,
										ModifiedBy:          transaction_.ModifiedBy,
										Active:              transaction_.Active,
										BranchName:          transaction_.BranchName,
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

	logs.Info("Transaction ID is ", v.TransactionId)

	if txn, txn_err := models.GetTransactionsById(v.TransactionId); txn_err == nil {
		// status_ := "SUCCESS"
		if status, err := models.GetStatusByName(v.Status); err == nil {
			txn.CreatedBy = v.Confirmedby
			txn.Status = status
			txn.Active = 1
			if utxn_err := models.UpdateTransactionsById(txn); utxn_err == nil {
				customerData := responses.CustomersAlt{
					CustomerId:    txn.Order.CustomerId,
					CustomerName:  txn.Order.CustomerName,
					CustomerEmail: txn.Order.CustomerEmail,
					CustomerPhone: txn.Order.CustomerPhone}

				orderDetails := []responses.OrderItemsCustom{}
				if txn.Order.OrderDetails != nil {
					for _, orderItem := range txn.Order.OrderDetails {
						itemData := responses.ItemAlt{
							ItemId:       orderItem.Item,
							ItemName:     orderItem.ItemName,
							Description:  "",
							Price:        orderItem.TotalPrice,
							Currency:     orderItem.Order.Currency,
							Category:     orderItem.Category,
							DateCreated:  orderItem.DateCreated,
							DateModified: orderItem.DateModified}

						orderIdStr := strconv.FormatInt(orderItem.Order.OrderId, 10)
						orderItemsCustom := responses.OrderItemsCustom{
							OrderItemId: orderItem.OrderItemId,
							OrderId:     orderIdStr,
							Item:        &itemData,
							Quantity:    orderItem.Quantity,
							OrderDate:   orderItem.OrderDate,
							Status:      orderItem.Status.Status}

						orderDetails = append(orderDetails, orderItemsCustom)
					}
				}
				var customOrder responses.OrdersCustom = responses.OrdersCustom{
					OrderId:      txn.Order.OrderId,
					OrderNumber:  txn.Order.OrderNumber,
					Quantity:     txn.Order.Quantity,
					Cost:         txn.Order.Cost,
					Currency:     txn.Order.Currency,
					OrderDate:    txn.Order.OrderDate,
					DateCreated:  txn.Order.DateCreated,
					DateModified: txn.Order.DateModified,
					Customer:     &customerData,
					OrderDetails: orderDetails}
				var customTxn responses.TransactionsCustom = responses.TransactionsCustom{
					TransactionId:       txn.TransactionId,
					Order:               &customOrder,
					Amount:              txn.Amount,
					TransactingCurrency: txn.CurrencySymbol,
					Status:              txn.Status.Status,
					DateCreated:         txn.DateCreated,
					DateModified:        txn.DateModified,
					CreatedBy:           txn.CreatedBy,
					ModifiedBy:          txn.ModifiedBy,
					Active:              txn.Active,
					BranchName:          txn.BranchName}

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

	// txn_id, _ := strconv.ParseInt(v.TransactionId, 0, 64)

	logs.Info("Transaction ID is ", v.TransactionId)

	if txn, txn_err := models.GetTransactionsById(v.TransactionId); txn_err == nil {
		logs.Info("Returned date is ", txn.Order.ReturnedDate)
		// status_ := "SUCCESS"
		if status, err := models.GetStatusByName(v.Status); err == nil {
			txn.CreatedBy = v.Confirmedby
			txn.Status = status
			txn.Active = 1
			if utxn_err := models.UpdateTransactionsById(txn); utxn_err == nil {
				customerData := responses.CustomersAlt{
					CustomerId:    txn.Order.CustomerId,
					CustomerName:  txn.Order.CustomerName,
					CustomerEmail: txn.Order.CustomerEmail,
					CustomerPhone: txn.Order.CustomerPhone}

				orderDetails := []responses.OrderItemsCustom{}
				if txn.Order.OrderDetails != nil {
					for _, orderItem := range txn.Order.OrderDetails {
						itemData := responses.ItemAlt{
							ItemId:       orderItem.Item,
							ItemName:     orderItem.ItemName,
							Description:  "",
							Price:        orderItem.TotalPrice,
							Currency:     orderItem.Order.Currency,
							Category:     orderItem.Category,
							DateCreated:  orderItem.DateCreated,
							DateModified: orderItem.DateModified}

						orderIdStr := strconv.FormatInt(orderItem.Order.OrderId, 10)
						orderItemsCustom := responses.OrderItemsCustom{
							OrderItemId: orderItem.OrderItemId,
							OrderId:     orderIdStr,
							Item:        &itemData,
							Quantity:    orderItem.Quantity,
							OrderDate:   orderItem.OrderDate,
							Status:      orderItem.Status.Status}

						orderDetails = append(orderDetails, orderItemsCustom)
					}
				}
				var customOrder responses.OrdersCustom = responses.OrdersCustom{
					OrderId:      txn.Order.OrderId,
					OrderNumber:  txn.Order.OrderNumber,
					Quantity:     txn.Order.Quantity,
					Cost:         txn.Order.Cost,
					Currency:     txn.Order.Currency,
					OrderDate:    txn.Order.OrderDate,
					DateCreated:  txn.Order.DateCreated,
					DateModified: txn.Order.DateModified,
					Customer:     &customerData,
					OrderDetails: orderDetails,
					ReturnedDate: txn.Order.ReturnedDate,
					OrderEndDate: txn.Order.OrderEndDate}
				var customTxn responses.TransactionsCustom = responses.TransactionsCustom{
					TransactionId:       txn.TransactionId,
					Order:               &customOrder,
					Amount:              txn.Amount,
					TransactingCurrency: txn.CurrencySymbol,
					Status:              txn.Status.Status,
					DateCreated:         txn.DateCreated,
					DateModified:        txn.DateModified,
					CreatedBy:           txn.CreatedBy,
					ModifiedBy:          txn.ModifiedBy,
					Active:              txn.Active,
					BranchName:          txn.BranchName}

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
								if oItem, err := functions.GetItemQuantity(&c.Controller, requests.GetItemRequest{ItemId: item.Item}); err == nil {
									logs.Info("Item quantity is ", oItem.Quantity)
									logs.Info("Item quantity is ", item.Quantity)
									oItem.Quantity.Quantity = int(oItem.Quantity.Quantity) + int(item.Quantity)
									if _, err := functions.UpdateItemQuantity(&c.Controller, requests.UpdateItemQuantityRequest{ItemId: item.Item, Quantity: int(oItem.Quantity.Quantity)}); err != nil {
										logs.Error("Error updating item::: ", err.Error())
										message := "Error updating the item quantity"
										var resp responses.TransactionResponseDTO = responses.TransactionResponseDTO{StatusCode: 806, Transaction: nil, StatusDesc: message}
										logs.Error("Error thrown when updating transaction::: ", err.Error())
										c.Ctx.Output.SetStatus(200)
										c.Data["json"] = resp
										c.ServeJSON()
									} else {
										logs.Info("Item quantity updated successfully")
									}
								}
								item.Quantity = int(item.Quantity) + int(item.Quantity)
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
		ordersResponse := []responses.OrdersCustom{}
		for _, order := range *orders {
			customerData := responses.CustomersAlt{
				CustomerId:    order.CustomerId,
				CustomerName:  order.CustomerName,
				CustomerEmail: order.CustomerEmail,
				CustomerPhone: order.CustomerPhone}

			orderDetails := []responses.OrderItemsCustom{}
			if order.OrderDetails != nil {
				for _, orderItem := range order.OrderDetails {
					itemData := responses.ItemAlt{
						ItemId:       orderItem.Item,
						ItemName:     orderItem.ItemName,
						Description:  "",
						Price:        orderItem.TotalPrice,
						Currency:     orderItem.Order.Currency,
						Category:     orderItem.Category,
						DateCreated:  orderItem.DateCreated,
						DateModified: orderItem.DateModified}

					orderIdStr := strconv.FormatInt(orderItem.Order.OrderId, 10)
					orderItemsCustom := responses.OrderItemsCustom{
						OrderItemId: orderItem.OrderItemId,
						OrderId:     orderIdStr,
						Item:        &itemData,
						Quantity:    orderItem.Quantity,
						OrderDate:   orderItem.OrderDate,
						Status:      orderItem.Status.Status}

					orderDetails = append(orderDetails, orderItemsCustom)
				}
			}
			var customOrder responses.OrdersCustom = responses.OrdersCustom{
				OrderId:      order.OrderId,
				OrderNumber:  order.OrderNumber,
				Quantity:     order.Quantity,
				Cost:         order.Cost,
				Currency:     order.Currency,
				OrderDate:    order.OrderDate,
				DateCreated:  order.DateCreated,
				DateModified: order.DateModified,
				OrderEndDate: order.OrderEndDate,
				Customer:     &customerData,
				OrderDetails: orderDetails,
				ReturnedDate: order.ReturnedDate}
			ordersResponse = append(ordersResponse, customOrder)
		}
		var resp = responses.OrdersResponseDTO{StatusCode: 200, Orders: &ordersResponse, StatusDesc: "Order details fetched successfully"}
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
		logs.Error("Error fetching order by ID::: ", err.Error())
		var resp = responses.OrderResponseDTO{StatusCode: 608, Order: nil, StatusDesc: "Failed to fetch order details"}
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = resp
	} else {
		var customerData = responses.CustomersAlt{
			CustomerId:    v.CustomerId,
			CustomerName:  v.CustomerName,
			CustomerEmail: v.CustomerEmail,
			CustomerPhone: v.CustomerPhone}

		orderDetails := []responses.OrderItemsCustom{}
		if v.OrderDetails != nil {
			for _, orderItem := range v.OrderDetails {
				itemData := responses.ItemAlt{
					ItemId:       orderItem.Item,
					ItemName:     orderItem.ItemName,
					Description:  "",
					Price:        orderItem.TotalPrice,
					Currency:     orderItem.Order.Currency,
					Category:     orderItem.Category,
					DateCreated:  orderItem.DateCreated,
					DateModified: orderItem.DateModified}

				orderIdStr := strconv.FormatInt(orderItem.Order.OrderId, 10)
				orderItemsCustom := responses.OrderItemsCustom{
					OrderItemId: orderItem.OrderItemId,
					OrderId:     orderIdStr,
					Item:        &itemData,
					Quantity:    orderItem.Quantity,
					OrderDate:   orderItem.OrderDate,
					Status:      orderItem.Status.Status}

				orderDetails = append(orderDetails, orderItemsCustom)
			}
		}
		var orderData = responses.OrdersCustom{
			OrderId:      v.OrderId,
			OrderNumber:  v.OrderNumber,
			Quantity:     v.Quantity,
			Cost:         v.Cost,
			Currency:     v.Currency,
			OrderDate:    v.OrderDate,
			DateCreated:  v.DateCreated,
			DateModified: v.DateModified,
			OrderEndDate: v.OrderEndDate,
			Customer:     &customerData,
			OrderDetails: orderDetails,
			ReturnedDate: v.ReturnedDate,
		}
		v.OrderDetails = nil
		var resp = responses.OrderResponseDTO{StatusCode: 200, Order: &orderData, StatusDesc: "Order details fetched successfully"}
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = resp
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
	var limit int64 = 30
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
		logs.Error("Error fetching all orders::: ", err.Error())
		var resp = responses.OrdersResponseDTO{StatusCode: 608, Orders: nil, StatusDesc: "Failed to fetch order details"}
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = resp
	} else {
		ordersResponse := []responses.OrdersCustom{}
		for _, order := range l {
			m := order.(models.Orders)
			var customerData = responses.CustomersAlt{
				CustomerId:    m.CustomerId,
				CustomerName:  m.CustomerName,
				CustomerEmail: m.CustomerEmail,
				CustomerPhone: m.CustomerPhone}

			orderDetails := []responses.OrderItemsCustom{}
			if m.OrderDetails != nil {
				for _, orderItem := range m.OrderDetails {
					itemData := responses.ItemAlt{
						ItemId:       orderItem.Item,
						ItemName:     orderItem.ItemName,
						Description:  "",
						Price:        orderItem.TotalPrice,
						Currency:     orderItem.Order.Currency,
						Category:     orderItem.Category,
						DateCreated:  orderItem.DateCreated,
						DateModified: orderItem.DateModified}

					orderIdStr := strconv.FormatInt(orderItem.Order.OrderId, 10)
					orderItemsCustom := responses.OrderItemsCustom{
						OrderItemId: orderItem.OrderItemId,
						OrderId:     orderIdStr,
						Item:        &itemData,
						Quantity:    orderItem.Quantity,
						OrderDate:   orderItem.OrderDate,
						Status:      orderItem.Status.Status}

					orderDetails = append(orderDetails, orderItemsCustom)
				}
			}
			var orderData = responses.OrdersCustom{
				OrderId:      m.OrderId,
				OrderNumber:  m.OrderNumber,
				Quantity:     m.Quantity,
				Cost:         m.Cost,
				Currency:     m.Currency,
				OrderDate:    m.OrderDate,
				DateCreated:  m.DateCreated,
				DateModified: m.DateModified,
				OrderEndDate: m.OrderEndDate,
				Customer:     &customerData,
				OrderDetails: orderDetails,
				ReturnedDate: m.ReturnedDate,
			}
			ordersResponse = append(ordersResponse, orderData)
		}
		var resp = responses.OrdersResponseDTO{StatusCode: 200, Orders: &ordersResponse, StatusDesc: "Order details fetched successfully"}
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = resp
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
