package cmd

import (
	"fmt"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

// 客户画像
var customerProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "客户画像查询",
	Long:  "查询客户画像相关数据，支持订单、交易等子命令。",
}

// 客户画像-订单
var customerProfileOrderCmd = &cobra.Command{
	Use:   "order",
	Short: "客户画像-订单",
	Long:  "查询客户画像订单相关数据，支持订单统计、日订单数量等子命令。",
}

// 客户画像-交易
var customerProfileTransactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "客户画像-交易",
	Long:  "查询客户画像交易相关数据，支持交易统计、交易流水等子命令。",
}

// 查询客户订单统计
var getCustomerProfileOrderCountCmd = func() *cobra.Command {
	var startDate, endDate string

	cmd := &cobra.Command{
		Use:   "count [customerId]",
		Short: "查询客户订单统计",
		Long:  "通过客户 ID 查询指定时间范围内的订单统计信息，包含总订单数、日均订单数、日最大/最小订单数。",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			customerId := args[0]
			count, err := apis.NewCustomerProfileApi().GetCustomerProfileOrderCount(customerId, startDate+" 00:00:00", endDate+" 23:59:59")
			if err != nil {
				cmd.PrintErrf("查询客户 %s 订单统计失败：%v\n", customerId, err)
				return
			}
			if count == nil {
				cmd.Print("未查询到客户订单统计信息")
				return
			}

			var sb strings.Builder
			fmt.Fprintln(&sb, "【订单统计】")
			fmt.Fprintf(&sb, "统计区间：%s ~ %s\n", count.StartDate, count.EndDate)
			fmt.Fprintf(&sb, "总订单数：%d\n", count.Total)
			fmt.Fprintf(&sb, "日均订单数：%d\n", count.Avg)
			fmt.Fprintf(&sb, "日最大订单数：%d\n", count.Max)
			fmt.Fprintf(&sb, "日最小订单数：%d", count.Min)

			cmd.Print(sb.String())
		},
	}

	cmd.Flags().StringVarP(&startDate, "start", "", "", "开始日期，格式：yyyy-MM-dd")
	cmd.Flags().StringVarP(&endDate, "end", "", "", "结束日期，格式：yyyy-MM-dd")
	_ = cmd.MarkFlagRequired("start")
	_ = cmd.MarkFlagRequired("end")

	return cmd
}()

// 查询客户日订单数量
var getCustomerProfileDailyOrdersCmd = func() *cobra.Command {
	var startDate, endDate string

	cmd := &cobra.Command{
		Use:   "daily [customerId]",
		Short: "查询客户日订单数量",
		Long:  "通过客户 ID 查询指定时间范围内每日的订单数量，包含付款、履约、超期、退款订单数量。每次最多返回 90 天的数据。",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			customerId := args[0]
			pageData, err := apis.NewCustomerProfileApi().GetCustomerProfileReceiveVisitorsInfo(customerId, startDate+" 00:00:00", endDate+" 23:59:59", 1, 90)
			if err != nil {
				cmd.PrintErrf("查询客户 %s 日订单数量失败：%v\n", customerId, err)
				return
			}
			if pageData == nil || len(pageData.Rows) == 0 {
				cmd.Print("未查询到客户日订单数量信息")
				return
			}

			var sb strings.Builder
			fmt.Fprintf(&sb, "总数：%d\n", pageData.Total)
			fmt.Fprintln(&sb, "|日期|订单数|付款订单|履约订单|超期订单|退款订单|")
			fmt.Fprintln(&sb, "|----|----|----|----|----|----|")
			for _, row := range pageData.Rows {
				fmt.Fprintf(&sb, "|%s|%d|%d|%d|%d|%d|\n",
					row.CrtDate,
					row.DaysNum,
					row.FkOrderNum,
					row.FulfilledNum,
					row.OverdueNum,
					row.TkOrderNum)
			}

			cmd.Print(strings.TrimRight(sb.String(), "\n"))
		},
	}

	cmd.Flags().StringVarP(&startDate, "start", "", "", "开始日期，格式：yyyy-MM-dd")
	cmd.Flags().StringVarP(&endDate, "end", "", "", "结束日期，格式：yyyy-MM-dd")
	_ = cmd.MarkFlagRequired("start")
	_ = cmd.MarkFlagRequired("end")

	return cmd
}()

// 查询客户交易统计
var getCustomerProfileTransactionCountCmd = func() *cobra.Command {
	var startDate, endDate string

	cmd := &cobra.Command{
		Use:   "count [customerId]",
		Short: "查询客户交易统计",
		Long:  "通过客户 ID 查询指定时间范围内的交易统计信息，包含总交易额、实际成交额、退款金额、总订单数、退款订单数、客单价。",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			customerId := args[0]
			data, err := apis.NewCustomerProfileApi().GetCustomerProfileTransactionCount(customerId, startDate+" 00:00:00", endDate+" 23:59:59")
			if err != nil {
				cmd.PrintErrf("查询客户 %s 交易统计失败：%v\n", customerId, err)
				return
			}
			if data == nil {
				cmd.Print("未查询到客户交易统计信息")
				return
			}

			var sb strings.Builder
			fmt.Fprintln(&sb, "【交易统计】")
			fmt.Fprintf(&sb, "总交易额：%.2f\n", data.TranAmount)
			fmt.Fprintf(&sb, "实际成交额：%.2f\n", data.ActualAmount)
			fmt.Fprintf(&sb, "退款金额：%.2f\n", data.RefundAmount)
			fmt.Fprintf(&sb, "总订单数：%d\n", data.AllOrderNum)
			fmt.Fprintf(&sb, "退款订单数：%d\n", data.RefundOrderNum)
			fmt.Fprintf(&sb, "客单价：%.2f", data.CustomerPrice)

			cmd.Print(sb.String())
		},
	}

	cmd.Flags().StringVarP(&startDate, "start", "", "", "开始日期，格式：yyyy-MM-dd")
	cmd.Flags().StringVarP(&endDate, "end", "", "", "结束日期，格式：yyyy-MM-dd")
	_ = cmd.MarkFlagRequired("start")
	_ = cmd.MarkFlagRequired("end")

	return cmd
}()

// 查询客户交易流水
var getCustomerProfileTransactionBillsCmd = func() *cobra.Command {
	var startDate, endDate string

	cmd := &cobra.Command{
		Use:   "bills [customerId]",
		Short: "查询客户交易流水",
		Long:  "通过客户 ID 查询指定时间范围内的每日交易流水，金额列依次为 CNY/USD/EUR。每次最多返回 90 天的数据。",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			customerId := args[0]
			pageData, err := apis.NewCustomerProfileApi().GetCustomerProfileTransactionBillRecords(customerId, startDate+" 00:00:00", endDate+" 23:59:59", 1, 90)
			if err != nil {
				cmd.PrintErrf("查询客户 %s 交易流水失败：%v\n", customerId, err)
				return
			}
			if pageData == nil || len(pageData.Rows) == 0 {
				cmd.Print("未查询到客户交易流水信息")
				return
			}

			var sb strings.Builder
			fmt.Fprintf(&sb, "总数：%d\n", pageData.Total)
			fmt.Fprintln(&sb, "|日期|订单数|退款订单数|交易额(￥/$/€)|实际成交额(￥/$/€)|退款金额(￥/$/€)|手续费(￥/$/€)|")
			fmt.Fprintln(&sb, "|----|----|----|----|----|----|----|")
			for _, row := range pageData.Rows {
				fmt.Fprintf(&sb, "|%s|%.0f|%.0f|%.2f/%.2f/%.2f|%.2f/%.2f/%.2f|%.2f/%.2f/%.2f|%.2f/%.2f/%.2f|\n",
					row.Day,
					row.OrderNum,
					row.RefundOrderNum,
					row.TranAmountCny, row.TranAmountUsd, row.TranAmountEur,
					row.ActualTranAmountCny, row.ActualTranAmountUsd, row.ActualTranAmountEur,
					row.RefundTranAmountCny, row.RefundTranAmountUsd, row.RefundTranAmountEur,
					row.TranFeeCny, row.TranFeeUsd, row.TranFeeEur)
			}

			cmd.Print(strings.TrimRight(sb.String(), "\n"))
		},
	}

	cmd.Flags().StringVarP(&startDate, "start", "", "", "开始日期，格式：yyyy-MM-dd")
	cmd.Flags().StringVarP(&endDate, "end", "", "", "结束日期，格式：yyyy-MM-dd")
	_ = cmd.MarkFlagRequired("start")
	_ = cmd.MarkFlagRequired("end")

	return cmd
}()

func init() {
	customerProfileCmd.AddCommand(
		customerProfileOrderCmd,       // 订单
		customerProfileTransactionCmd, // 交易
	)

	// 客户画像-订单
	customerProfileOrderCmd.AddCommand(
		getCustomerProfileOrderCountCmd, // 订单统计
		getCustomerProfileDailyOrdersCmd,
	)

	// 客户画像-交易
	customerProfileTransactionCmd.AddCommand(
		getCustomerProfileTransactionCountCmd, // 交易统计
		getCustomerProfileTransactionBillsCmd, // 交易流水
	)

	customerCmd.AddCommand(customerProfileCmd)
}
