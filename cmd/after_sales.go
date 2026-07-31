package cmd

import (
	"fmt"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

var afterSalesCmd = &cobra.Command{
	Use:   "after-sales",
	Short: "售后管理",
	Long:  "",
}

// 查询1688售后
var get1688AfterSalesCmd = &cobra.Command{
	Use:   "1688 [thirdOrderId]",
	Short: "查询1688售后工单",
	Long:  "通过第三方订单号（1688订单号、交易号）在 Hypersku 平台查询 1688 售后工单。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		thirdOrderId := args[0]
		if res, err := apis.NewAfterSalesApi().Get1688AfterSales(thirdOrderId); err != nil {
			cmd.PrintErr("查询1688售后工单发生错误", err)
			return
		} else {
			if len(*res) == 0 {
				cmd.Print("未查询到1688售后工单信息")
				return
			}

			afterSalesList := make([]string, len(*res)+2)
			afterSalesList[0] = "|退款ID|交易号|退款状态|退款类型|纠纷类型|当前操作|退款金额|退款运费|总退款金额|退货运费|创建时间|更新时间|剩余处理时间|"
			afterSalesList[1] = "|----|----|----|----|----|----|----|----|----|----|----|----|----|"
			for i, item := range *res {
				afterSalesList[i+2] = fmt.Sprintf("|%s|%s|%s|%s|%s|%s|%.2f|%.2f|%.2f|%.2f|%s|%s|%s|",
					item.RefundID,
					item.ThirdOrderID,
					item.StatusShowDesc,
					item.DisputeType,
					item.DisputeRequestDesc,
					item.ActionDesc,
					item.RefundAmount,
					item.Freight,
					item.AllRefundAmount,
					item.BackFreight,
					item.CrtTime,
					item.UpdTime,
					item.TimeDesc,
				)
			}

			cmd.Print(strings.Join(afterSalesList, "\n"))
		}

	},
}

// 查询1688售后商品
var get1688AfterSalesGoodsCmd = &cobra.Command{
	Use:   "goods [thirdOrderId] [refundId]",
	Short: "查询1688售后商品",
	Long:  "通过第三方订单号（1688订单号）和退款ID查询1688售后工单中的商品项。",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		thirdOrderId, refundId := args[0], args[1]
		if res, err := apis.NewAfterSalesApi().Get1688AfterSalesGoods(thirdOrderId, refundId); err != nil {
			cmd.PrintErr("查询1688售后商品发生错误", err)
			return
		} else {
			if len(*res) == 0 {
				cmd.Print("未查询到1688售后商品信息")
				return
			}

			goodsList := make([]string, len(*res)+2)
			goodsList[0] = "|商品名称|SKU|属性|数量|图片|"
			goodsList[1] = "|----|----|----|----|----|"
			for i, goods := range *res {
				goodsList[i+2] = fmt.Sprintf("|%s|%d|%s|%d|%s|",
					goods.GoodsName,
					goods.GoodsID,
					goods.GoodsAttr,
					goods.Num,
					goods.ImgURL,
				)
			}

			cmd.Print(strings.Join(goodsList, "\n"))
		}
	},
}

// 查询1688售后详情
var get1688AfterSalesDetailCmd = &cobra.Command{
	Use:   "detail [refundId]",
	Short: "查询1688售后详情",
	Long:  "通过退款ID查询1688售后的详细信息，包含退款金额、退款原因、卖家信息等。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		refundId := args[0]
		if detail, err := apis.NewAfterSalesApi().Get1688AfterSalesDetail(refundId); err != nil {
			cmd.PrintErr("查询1688售后详情发生错误", err)
			return
		} else {
			if detail == nil {
				cmd.Print("未查询到1688售后详情")
				return
			}

			cmd.Print(fmt.Sprintf(`【售后详情】
退款ID：%s
交易号：%s
退款状态：%s
退款服务：%s
退款原因：%s
退款说明：%s
商品退款金额：%.2f
退款运费：%.2f
退款总金额：%.2f
卖家：%s
会员登录名：%s
手机：%s
电话：%s
申请退款时间：%s`,
				detail.RefundID,
				detail.ThirdOrderID,
				detail.StatusDesc,
				detail.RefundType,
				detail.RefundReason,
				detail.RefundDesc,
				detail.RefundAmount,
				detail.Freight,
				detail.AllRefundAmount,
				detail.CompanyName,
				detail.SellerLoginID,
				detail.Mobile,
				detail.Phone,
				detail.ApplyTime,
			))
		}
	},
}

// 查询1688售后留言
var get1688AfterSalesMessageCmd = &cobra.Command{
	Use:   "message [refundId]",
	Short: "查询1688售后留言",
	Long:  "通过退款ID查询1688售后的留言记录。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		refundId := args[0]
		if res, err := apis.NewAfterSalesApi().Get1688AfterSalesMessage(refundId); err != nil {
			cmd.PrintErr("查询1688售后留言发生错误", err)
			return
		} else {
			if len(*res) == 0 {
				cmd.Print("未查询到1688售后留言信息")
				return
			}

			messageList := make([]string, len(*res)+2)
			messageList[0] = "|留言时间|留言方|留言内容|操作留言|"
			messageList[1] = "|----|----|----|----|"
			for i, msg := range *res {
				messageList[i+2] = fmt.Sprintf("|%s|%s|%s|%s|",
					msg.GmtCreate,
					msg.OperatorLoginID,
					msg.Discription,
					msg.OperateRemark,
				)
			}

			cmd.Print(strings.Join(messageList, "\n"))
		}
	},
}

func init() {

	get1688AfterSalesCmd.AddCommand(
		get1688AfterSalesGoodsCmd,
		get1688AfterSalesDetailCmd,
		get1688AfterSalesMessageCmd,
	)

	afterSalesCmd.AddCommand(
		get1688AfterSalesCmd,
	)
	rootCmd.AddCommand(afterSalesCmd)
}
