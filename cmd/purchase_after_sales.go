package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

var purchaseAfterSalesCmd = &cobra.Command{
	Use:   "purchase-after-sales",
	Short: "采购售后查询",
	Long: `采购售后工单查询命令集（只读）。支持按交易号、订单号查询。

可用子命令：
  by-trade   按交易号查询采购售后工单
  by-order   按订单号查询采购售后工单`,
}

// 按交易号查询采购售后
var purchaseAfterSalesByTradeCmd = &cobra.Command{
	Use:   "by-trade [tradeId]",
	Short: "按交易号查询采购售后工单",
	Long:  "通过交易号（第三方订单号）查询采购售后工单列表",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tradeId, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			cmd.PrintErrf("无效的交易号: %s\n", args[0])
			return
		}

		res, err := apis.NewPurchaseAfterSalesApi().PageList(apis.PurchaseAfterSalesQuery{
			Page:         1,
			Limit:        20,
			ThirdOrderID: tradeId,
		})
		if err != nil {
			cmd.PrintErrf("查询采购售后失败: %v\n", err)
			return
		}
		if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "交易号 %s 无采购售后记录\n", args[0])
			return
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "交易号: %s   共 %d 条\n\n", args[0], res.Data.Total)
		for i, row := range res.Data.Rows {
			if i > 0 {
				fmt.Fprintln(&sb)
			}
			fmt.Fprintf(&sb, "【采购售后 %d】\n", i+1)
			fmt.Fprintf(&sb, "  工单号: %s\n", row.WorkOrderCode)
			fmt.Fprintf(&sb, "  交易号: %s\n", row.ThirdOrderIdStr)
			fmt.Fprintf(&sb, "  快递单号: %s\n", row.TrackingNumber)
			fmt.Fprintf(&sb, "  仓库: %s\n", row.WarehouseName)
			fmt.Fprintf(&sb, "  商品: %s (SKU: %s)\n", row.GoodsName, row.GoodsSku)
			fmt.Fprintf(&sb, "  数量: %d   金额: %.2f\n", row.Quantity, row.TotalAmount)
			fmt.Fprintf(&sb, "  异常类型: %s\n", row.AbnormalTypeStr)
			fmt.Fprintf(&sb, "  状态: %s\n", row.StatusStr)
			fmt.Fprintf(&sb, "  处理人: %s\n", row.HandlerName)
			fmt.Fprintf(&sb, "  客户: %s\n", row.CustomerUsername)
			fmt.Fprintf(&sb, "  创建时间: %s\n", row.CrtTime)
			if row.Remark != "" {
				fmt.Fprintf(&sb, "  备注: %s\n", row.Remark)
			}
		}
		cmd.Print(sb.String())
	},
}

// 按订单号查询采购售后
var purchaseAfterSalesByOrderCmd = &cobra.Command{
	Use:   "by-order [orderId]",
	Short: "按订单号查询采购售后工单",
	Long:  "通过订单号查询采购售后工单列表",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		orderId, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			cmd.PrintErrf("无效的订单号: %s\n", args[0])
			return
		}

		res, err := apis.NewPurchaseAfterSalesApi().PageList(apis.PurchaseAfterSalesQuery{
			Page:    1,
			Limit:   20,
			OrderID: orderId,
		})
		if err != nil {
			cmd.PrintErrf("查询采购售后失败: %v\n", err)
			return
		}
		if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "订单号 %s 无采购售后记录\n", args[0])
			return
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "订单号: %s   共 %d 条\n\n", args[0], res.Data.Total)
		for i, row := range res.Data.Rows {
			if i > 0 {
				fmt.Fprintln(&sb)
			}
			fmt.Fprintf(&sb, "【采购售后 %d】\n", i+1)
			fmt.Fprintf(&sb, "  工单号: %s\n", row.WorkOrderCode)
			fmt.Fprintf(&sb, "  交易号: %s\n", row.ThirdOrderIdStr)
			fmt.Fprintf(&sb, "  快递单号: %s\n", row.TrackingNumber)
			fmt.Fprintf(&sb, "  仓库: %s\n", row.WarehouseName)
			fmt.Fprintf(&sb, "  商品: %s (SKU: %s)\n", row.GoodsName, row.GoodsSku)
			fmt.Fprintf(&sb, "  数量: %d   金额: %.2f\n", row.Quantity, row.TotalAmount)
			fmt.Fprintf(&sb, "  异常类型: %s\n", row.AbnormalTypeStr)
			fmt.Fprintf(&sb, "  状态: %s\n", row.StatusStr)
			fmt.Fprintf(&sb, "  处理人: %s\n", row.HandlerName)
			fmt.Fprintf(&sb, "  客户: %s\n", row.CustomerUsername)
			fmt.Fprintf(&sb, "  创建时间: %s\n", row.CrtTime)
			if row.Remark != "" {
				fmt.Fprintf(&sb, "  备注: %s\n", row.Remark)
			}
		}
		cmd.Print(sb.String())
	},
}

func init() {
	purchaseAfterSalesCmd.AddCommand(purchaseAfterSalesByTradeCmd)
	purchaseAfterSalesCmd.AddCommand(purchaseAfterSalesByOrderCmd)

	rootCmd.AddCommand(purchaseAfterSalesCmd)
}
