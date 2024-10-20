package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["transaction_service/controllers:Order_itemsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:Order_itemsController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:Order_itemsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:Order_itemsController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:Order_itemsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:Order_itemsController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:Order_itemsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:Order_itemsController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:Order_itemsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:Order_itemsController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"],
        beego.ControllerComments{
            Method: "ConfirmOrder",
            Router: `/confirm-order`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:OrdersController"],
        beego.ControllerComments{
            Method: "GetUserOrders",
            Router: `/get-user-orders`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:Transaction_detailsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:Transaction_detailsController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:Transaction_detailsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:Transaction_detailsController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:Transaction_detailsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:Transaction_detailsController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:Transaction_detailsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:Transaction_detailsController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"],
        beego.ControllerComments{
            Method: "GetUserTransactions",
            Router: `/get-user-transactions`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"] = append(beego.GlobalControllerRouter["transaction_service/controllers:TransactionsController"],
        beego.ControllerComments{
            Method: "GetUserTransactionsByDate",
            Router: `/get-user-transactions-by-date`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
