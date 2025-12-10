// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"transaction_service/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/orders",
			beego.NSInclude(
				&controllers.OrdersController{},
			),
		),
		beego.NSNamespace("/transactions",
			beego.NSInclude(
				&controllers.TransactionsController{},
			),
		),
		beego.NSNamespace("/order-items",
			beego.NSInclude(
				&controllers.Order_itemsController{},
			),
		),
	)

	ns2 := beego.NewNamespace("/v2",

		beego.NSNamespace("/transactions",
			beego.NSInclude(
				&controllers.TransactionsV2Controller{},
			),
		),
	)

	beego.AddNamespace(ns2)

	beego.AddNamespace(ns)
}
