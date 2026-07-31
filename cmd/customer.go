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

		orderContent := fmt.Sprintf(`【订单】
客户订单号：%s
订单金额：%s%.2f
实际金额：%s%.2f
运费：%s%.2f
增值服务费：%s%.2f
税费：%s%.2f
关税：%s%.2f
币种：%s
仓库：%s
采购状态：%s
支付时间：%s
创建时间：%s`,
			orderInfo.ID,
			orderInfo.CurrencySymbol, orderInfo.Amount,
			orderInfo.CurrencySymbol, orderInfo.ActualAmount,
			orderInfo.CurrencySymbol, orderInfo.Freight,
			orderInfo.CurrencySymbol, orderInfo.BrandingServiceAmount,
			orderInfo.CurrencySymbol, orderInfo.TaxAmount,
			orderInfo.CurrencySymbol, orderInfo.TariffAmount,
			orderInfo.CurrencyCode,
			orderInfo.WarehouseName,
			orderStatus[orderInfo.PurchaseStatus],
			orderInfo.PaymentTime,
			orderInfo.CrtTime)

		goodsList := make([]string, len(orderInfo.GoodsList))
		for i, goods := range orderInfo.GoodsList {
			goodsList[i] = fmt.Sprintf(`【商品项 %d】
SPU：%d
SKU：%s
名称：%s
属性：%s
图片：%s
数量：%d
销售价：%s%.2f
单价：%s%.2f
重量：%.2fg`,
				i+1,
				goods.ProductID,
				goods.GoodsID,
				goods.GoodsName,
				goods.AttrStr,
				goods.ImgURL,
				goods.Num,
				goods.CurrencySymbol, goods.SellingPrice,
				goods.CurrencySymbol, goods.UnitPrice,
				goods.Weight)
		}

		result := fmt.Sprintf("%s\n---\n%s", orderContent, strings.Join(goodsList, "\n\n"))
		cmd.Print(result)
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

func init() {
	customerCmd.AddCommand(customerOrderCmd)

	// 订单
	customerOrderCmd.AddCommand(
		getCustomerOrderInfo,      // 商品信息
		getCustomerOrderLogistics, // 物流
		getCustomerOrderAddress,   // 地址
	)

	rootCmd.AddCommand(customerCmd)
}
