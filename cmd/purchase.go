package cmd

import (
	"fmt"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

var (
	// 订单类型
	orderType = map[int]string{
		1:  "商品订单",
		2:  "会员收费服务订单",
		3:  "退款",
		4:  "草稿",
		5:  "充值",
		6:  "备货",
		9:  "平台线下充值",
		11: "平台冲正",
		12: "订单部分退款",
		13: "授权订单",
		14: "集运",
		16: "补款订单",
		17: "产品退款",
		18: "自购商品",
		19: "自购备货",
		20: "自购集运",
		22: "海外仓备货",
		24: "样品",
	}
	// 订单状态
	orderStatus = map[int]string{
		1:  "待审核",
		2:  "待采购",
		3:  "待付款",
		4:  "采购中",
		5:  "待发货",
		6:  "运输中",
		7:  "已完成",
		8:  "已作废",
		9:  "异常",
		10: "退款中",
	}

	// 订单状态日志
	orderStatusLog = map[int]string{
		1:   "待审核",
		2:   "待采购",
		3:   "待付款",
		4:   "采购中",
		5:   "待发货",
		6:   "运输中",
		7:   "已完成",
		8:   "已作废",
		9:   "异常",
		10:  "退款中",
		100: "客户付款",
		101: "包裹签收",
		102: "包裹称重",
		103: "包裹入库",
	}
)

// purchaseCmd 表示 purchase 子命令
var purchaseCmd = &cobra.Command{
	Use:   "purchase",
	Short: "采购订单管理",
	Long:  "采购订单管理命令集，提供了采购订单相关的查询和管理功能。",
}

// 采购订单信息
var getOrderInfoCmd = &cobra.Command{
	Use:   "get-order-info [orderId]",
	Short: "查询采购订单详情",
	Long:  "根据订单号查询采购订单的详细信息。包括商品明细、金额等。",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		orderId := args[0]
		orderInfo, err := apis.NewPurchaseApi().GetOrderInfo(orderId)
		if err != nil {
			cmd.PrintErrf("查询订单 %s 信息失败\n", orderId)
			cmd.Help()
			return
		}

		orderInfoContent := fmt.Sprintf(`【订单】
采购订单号：%s
客户订单号：%s
下单时间：%s
仓库：%s
金额：￥%.2f
运费：￥%.2f
订单类型：%s
订单状态：%s`,
			orderInfo.ID,
			orderInfo.CustomerOrderId,
			orderInfo.CrtTime,
			orderInfo.Warehouse,
			orderInfo.AmountActuallyPaid,
			orderInfo.Freight,
			orderType[orderInfo.Type],
			orderStatus[orderInfo.Status])

		goodsList := make([]string, len(orderInfo.GoodsList))
		for i, goods := range orderInfo.GoodsList {
			goodsList[i] = fmt.Sprintf(`【商品项 %d】
SPU：%d
SKU：%d
名称：%s
类别：%s %s
数量：%d
单价：￥%.2f
属性：%s
重量：%.2fg`,
				i+1,
				goods.ProductId,
				goods.GoodsId,
				goods.GoodsName,
				goods.CategoryName,
				goods.SubCategoryName,
				goods.Num,
				goods.UnitPrice,
				goods.GoodsAttr,
				goods.Weight)
		}

		result := fmt.Sprintf("%s\n---\n%s", orderInfoContent, strings.Join(goodsList, "\n\n"))
		cmd.Print(result)
	},
}

// 采购订单地址
var getOrderAddress = &cobra.Command{
	Use:   "get-order-address [orderId]",
	Short: "查询采购订单地址",
	Long:  "根据订单号查询采购订单的地址信息",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		orderId := args[0]
		orderInfo, err := apis.NewPurchaseApi().GetOrderInfo(orderId)
		if err != nil {
			cmd.PrintErrf("查询订单 %s 地址失败\n", orderId)
			cmd.Help()
			return
		}

		address := orderInfo.OrdersAddress
		orderAddressContent := fmt.Sprintf(`地址：%s
国家：%s
省份：%s
市：%s
区：%s
邮编：%s
		`, address.Address,
			address.CountryName,
			address.SecondRegionName,
			address.ThirdRegionName,
			address.FourthRegionName,
			address.ZipCode)

		cmd.Print(orderAddressContent)
	},
}

// 采购日志
var getPurchaseLog = &cobra.Command{
	Use:   "get-order-log [orderId]",
	Short: "查询采购日志",
	Long:  "根据订单号查询采购日志。",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		orderId := args[0]
		purchaseLogList, err := apis.NewPurchaseApi().GetPurchaseLog(orderId)
		if err != nil {
			cmd.PrintErrf("查询订单 %s 采购日志失败\n", orderId)
			cmd.Help()
			return
		}

		resultList := make([]string, len(*purchaseLogList))
		for i, item := range *purchaseLogList {
			resultList[i] = fmt.Sprintf("|%s|%s|", item.UpdTime, orderStatusLog[item.Status])
		}

		cmd.Print("|时间|状态|\n|----|----|\n" + strings.Join(resultList, "\n"))

	},
}

// 国际物流轨迹
var getInternationalLogistics = &cobra.Command{
	Use:   "get-purchase-international-logistics [orderId]",
	Short: "查询采购订单国际物流轨迹",
	Long:  "根据订单号查询订单包裹国际段物流轨迹",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		orderId := args[0]
		internationalLogisticsList, err := apis.NewPurchaseApi().GetInternationalLogistics(orderId)
		if err != nil {
			cmd.PrintErrf("查询订单 %s 采购国际物流轨迹失败\n", orderId)
			cmd.Help()
			return
		}

		result := make([]string, len(*internationalLogisticsList))
		for i, item := range *internationalLogisticsList {
			trackInfoItemList := make([]string, len(item.TruckInfo.List)+2)
			trackInfoItemList[0] = "|时间|内容|位置|"
			trackInfoItemList[1] = "|----|----|----|"
			for j, trackItem := range item.TruckInfo.List {
				trackInfoItemList[j+2] = fmt.Sprintf("|%s|%s|%s|", trackItem.ProcessDate, trackItem.ProcessContent, trackItem.ProcessLocation)
			}

			result[i] = fmt.Sprintf(`【包裹 %d】
物流单号：%s
承运商：%s
仓库：%s
轨迹：

%s`,
				i+1,
				item.TrackingNumber,
				item.ExpressDelivery,
				item.Warehouse,
				strings.Join(trackInfoItemList, "\n"),
			)
		}

		cmd.Print(strings.Join(result, "---\n"))
	},
}

func init() {
	purchaseCmd.AddCommand(
		getOrderInfoCmd,           // 订单
		getOrderAddress,           // 地址
		getPurchaseLog,            // 采购日志
		getInternationalLogistics, // 国际物流轨迹
	)
	rootCmd.AddCommand(purchaseCmd)
}
