package cmd

import (
	"fmt"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

var customerCmd = &cobra.Command{
	Use:   "customer",
	Short: "客户管理",
	Long:  "提供客户相关信息的查询，支持订单详情、物流单号、收货地址等子命令。",
}

var customerOrderCmd = &cobra.Command{
	Use:   "order",
	Short: "客户订单查询",
	Long:  "通过客户订单号查询订单的完整信息，包含订单详情（商品/金额/币种）、物流单号与承运商、收货地址与税号。",
}

// 查询订单和商品信息
var getCustomerOrderInfo = &cobra.Command{
	Use:   "info [orderId]",
	Short: "查询订单信息",
	Long:  "通过客户订单号查询客户订单信息，包含订单和订单商品项内容。",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		orderId := args[0]
		orderInfo, err := apis.NewCustomerOrderApi().GetOrderInfo(orderId)
		if err != nil {
			cmd.PrintErrf("查询订单 %s 信息失败\n", orderId)
			cmd.Help()
			return
		}
		if orderInfo == nil {
			cmd.Print("未查询到订单信息")
			return
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "【订单】\n")
		fmt.Fprintf(&sb, "客户订单号：%s\n", orderInfo.ID)
		fmt.Fprintf(&sb, "订单金额：%s%.2f\n", orderInfo.CurrencySymbol, orderInfo.Amount)
		fmt.Fprintf(&sb, "实际金额：%s%.2f\n", orderInfo.CurrencySymbol, orderInfo.ActualAmount)
		fmt.Fprintf(&sb, "运费：%s%.2f\n", orderInfo.CurrencySymbol, orderInfo.Freight)
		fmt.Fprintf(&sb, "增值服务费：%s%.2f\n", orderInfo.CurrencySymbol, orderInfo.BrandingServiceAmount)
		fmt.Fprintf(&sb, "税费：%s%.2f\n", orderInfo.CurrencySymbol, orderInfo.TaxAmount)
		fmt.Fprintf(&sb, "关税：%s%.2f\n", orderInfo.CurrencySymbol, orderInfo.TariffAmount)
		fmt.Fprintf(&sb, "币种：%s\n", orderInfo.CurrencyCode)
		fmt.Fprintf(&sb, "仓库：%s\n", orderInfo.WarehouseName)
		fmt.Fprintf(&sb, "状态：%s\n", apis.CustomerOrderStatus[orderInfo.Status])
		fmt.Fprintf(&sb, "采购状态：%s\n", orderStatus[orderInfo.PurchaseStatus])
		fmt.Fprintf(&sb, "支付时间：%s\n", orderInfo.PaymentTime)
		fmt.Fprintf(&sb, "创建时间：%s\n", orderInfo.CrtTime)
		fmt.Fprintln(&sb, "---")

		for i, goods := range orderInfo.GoodsList {

			fmt.Fprintf(&sb, "【商品项 %d】\n", i+1)
			fmt.Fprintf(&sb, "SPU：%d\n", goods.ProductID)
			fmt.Fprintf(&sb, "SKU：%s\n", goods.GoodsID)
			fmt.Fprintf(&sb, "名称：%s\n", goods.GoodsName)
			fmt.Fprintf(&sb, "属性：%s\n", goods.AttrStr)
			fmt.Fprintf(&sb, "图片：%s\n", goods.ImgURL)
			fmt.Fprintf(&sb, "数量：%d\n", goods.Num)
			fmt.Fprintf(&sb, "销售价：%s%.2f\n", goods.CurrencySymbol, goods.SellingPrice)
			fmt.Fprintf(&sb, "单价：%s%.2f\n", goods.CurrencySymbol, goods.UnitPrice)
			fmt.Fprintf(&sb, "重量：%.2fg\n", goods.Weight)
			fmt.Fprintln(&sb, "---")
		}

		cmd.Print(sb.String())
	},
}

// 查询订单物流信息
var getCustomerOrderLogistics = &cobra.Command{
	Use:   "logistics [orderId]",
	Short: "查询订单物流信息",
	Long:  "通过客户订单号查询订单的物流信息，包含物流单号、承运商（不返回具体的物流轨迹）。",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		orderId := args[0]
		orderInfo, err := apis.NewCustomerOrderApi().GetOrderInfo(orderId)
		if err != nil {
			cmd.PrintErrf("查询订单 %s 物流信息失败\n", orderId)
			cmd.Help()
			return
		}
		if orderInfo == nil || len(orderInfo.LogisticsList) == 0 {
			cmd.Print("未查询到订单物流信息")
			return
		}

		logisticsList := make([]string, len(orderInfo.LogisticsList)+2)
		logisticsList[0] = "|物流单号|承运商|"
		logisticsList[1] = "|----|----|"
		for i, item := range orderInfo.LogisticsList {
			logisticsList[i+2] = fmt.Sprintf("|%s|%s|", item.TrackingNumber, item.ExpressDelivery)
		}

		cmd.Print(strings.Join(logisticsList, "\n"))
	},
}

// 查询订单地址信息
var getCustomerOrderAddress = &cobra.Command{
	Use:   "address [orderId]",
	Short: "查询订单地址信息",
	Long:  "通过客户订单号查询订单的送达地址。",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		orderId := args[0]
		orderInfo, err := apis.NewCustomerOrderApi().GetOrderInfo(orderId)
		if err != nil {
			cmd.PrintErrf("查询订单 %s 地址信息失败\n", orderId)
			cmd.Help()
			return
		}
		if orderInfo == nil {
			cmd.Print("未查询到订单地址信息")
			return
		}

		address := orderInfo.OrdersAddress
		cmd.Print(fmt.Sprintf(`【收货地址】
收货人：%s %s
国家：%s
省份：%s
城市：%s
区域：%s
详细地址：%s
邮编：%s
欧盟税号：%s
VAT：%s`,
			address.FirstName,
			address.LastName,
			address.CountryName,
			address.SecondRegionName,
			address.ThirdRegionName,
			address.FourthRegionName,
			address.Address,
			address.ZipCode,
			address.TaxNo,
			address.VatNo))
	},
}

// 查询订单退件
var getCustomerOrderReturnInfo = &cobra.Command{
	Use:   "return [customerOrderId]",
	Short: "查询客户订单退件工单",
	Long:  "通过客户订单号查询订单退件工单",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		customerOrderId := args[0]
		if returnInfo, err := apis.NewCustomerOrderReturnApi().GetOrderReturnInfo(customerOrderId); err != nil {
			cmd.PrintErr(err)
			return
		} else {
			if returnInfo == nil {
				cmd.Print("未查询到退件记录\n")
				return
			}

			var sb strings.Builder
			fmt.Fprintln(&sb, "【退件工单详情】")
			fmt.Fprintf(&sb, "客户订单号:%s\n", returnInfo.CustomerOrderID)
			fmt.Fprintf(&sb, "运单号:%s\n", returnInfo.WaybillNumber)
			fmt.Fprintf(&sb, "快递单号:%s\n", returnInfo.TrackingNumber)
			fmt.Fprintf(&sb, "客户单号:%s\n", returnInfo.CustomerOrderNumber)
			fmt.Fprintf(&sb, "国家:%s\n", returnInfo.CountryCode)
			fmt.Fprintf(&sb, "入库时间:%s\n", returnInfo.InstorageCreatedOn)
			fmt.Fprintf(&sb, "创建时间:%s\n", returnInfo.CrtTime)
			fmt.Fprintf(&sb, "工单类型:%s\n", apis.WorkOrderType[returnInfo.WorkOrderType])
			fmt.Fprintf(&sb, "工单状态:%s\n", apis.WorkOrderType[returnInfo.WorkOrderState])
			fmt.Fprintf(&sb, "退件状态:%s\n", apis.ReturnStatus[returnInfo.Status])
			fmt.Fprintf(&sb, "说明:%s\n", returnInfo.Describing)
			fmt.Fprintf(&sb, "留言:%s\n", returnInfo.Remark)
			fmt.Fprintf(&sb, "留言时间:%s\n", returnInfo.RemarkUpdateTime)
			cmd.Println(sb.String())
		}
	},
}

func init() {
	customerCmd.AddCommand(customerOrderCmd)

	// 订单
	customerOrderCmd.AddCommand(
		getCustomerOrderInfo,       // 商品信息
		getCustomerOrderLogistics,  // 物流
		getCustomerOrderAddress,    // 地址
		getCustomerOrderReturnInfo, // 退件
	)

	rootCmd.AddCommand(customerCmd)
}
