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
var orderInfoCmd = &cobra.Command{
	Use:       "info [orderId]",
	Short:     "查询采购订单详情",
	Long:      "根据订单号查询采购订单的详细信息。包括商品明细、金额等。",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"orderId"},
	Run: func(cmd *cobra.Command, args []string) {
		orderId := args[0]
		orderInfo, err := apis.NewPurchaseApi().GetOrderInfo(orderId)
		if err != nil {
			cmd.PrintErrf("查询订单 %s 信息失败\n", orderId)
			cmd.Help()
			return
		}

		if orderInfo == nil {
			cmd.Printf("查询到不到该订单信息")
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

// 分页查询订单
var orderInfoPageCmd = func() *cobra.Command {
	var opts apis.QueryPage

	cmd := &cobra.Command{
		Use:   "page",
		Short: "分页查询采购订单",
		Long:  "分页查询采购订单信息，支持日期、订单号、交易号、物流单号查询。仅返回订单的信息，不包含商品项、物流、地址等。",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.TrackingNumber == "" && opts.TransactionNo == "" {
				if opts.StartDate == "" || opts.EndDate == "" {
					cmd.PrintErr("查询日期、交易号、物流单号不能同时为空")
					cmd.Help()
					return
				}
			}

			if res, err := apis.NewPurchaseApi().PageList(opts); err != nil {
				cmd.PrintErr("分页查询订单信息失败", err)
				return
			} else {
				total := res.Data.Total
				rows := res.Data.Rows

				if len(rows) == 0 {
					cmd.Print(fmt.Sprintf("当前页码：%d，页大小：%d，总数：%d\n\n无数据", opts.Page, opts.Limit, total))
					return
				}

				contentList := make([]string, len(rows)+3)
				contentList[0] = fmt.Sprintf("当前页码：%d，页大小：%d，总数：%d\n", opts.Page, opts.Limit, total)
				contentList[1] = "|采购订单号|客户订单号|金额|下单时间|运费|类型|状态|仓库|"
				contentList[2] = "|----|----|----|----|----|----|----|----|"
				for i, row := range rows {
					contentList[i+3] = fmt.Sprintf("|%s|%s|%.2f|%s|%.2f|%s|%s|%s|",
						row.ID,
						row.CustomerOrderId,
						row.AmountActuallyPaid,
						row.CrtTime,
						row.Freight,
						orderType[row.Type],
						orderStatus[row.Status],
						row.Warehouse,
					)
				}

				cmd.Print(strings.Join(contentList, "\n"))
			}

		},
	}

	cmd.Flags().IntVarP(&opts.Page, "page", "p", 1, "页码(必填)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 10, "页大小(必填)")
	cmd.Flags().StringVarP(&opts.StartDate, "start", "", "", "开始时间，格式：yyyy-MM-dd HH:mm:ss")
	cmd.Flags().StringVarP(&opts.EndDate, "end", "", "", "结束时间，格式：yyyy-MM-dd HH:mm:ss")
	// cmd.Flags().StringVarP(&opts.Id, "orderId", "", "", "订单号")
	cmd.Flags().StringVarP(&opts.TransactionNo, "thirdOrderId", "", "", "交易号、第三方订单号")
	cmd.Flags().StringVarP(&opts.TrackingNumber, "trackingNumber", "", "", "物流单号")

	// 标记为必填
	rootCmd.MarkFlagRequired("page")
	rootCmd.MarkFlagRequired("limit")

	return cmd
}

// 采购日志
var purchaseLogCmd = &cobra.Command{
	Use:   "log [orderId]",
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

		if purchaseLogList == nil || len(*purchaseLogList) == 0 {
			cmd.Print("未查询到该订单的采购日志")
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
var internationalLogisticsCmd = &cobra.Command{
	Use:   "logistics [orderId]",
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

		if internationalLogisticsList == nil || len(*internationalLogisticsList) == 0 {
			cmd.Print("未查询到该订单的国际段物流信息")
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

	orderInfoCmd.AddCommand(orderInfoPageCmd())

	purchaseCmd.AddCommand(
		orderInfoCmd,              // 订单
		purchaseLogCmd,            // 采购日志
		internationalLogisticsCmd, // 国际物流轨迹
	)
	rootCmd.AddCommand(purchaseCmd)
}
