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

// 查询客户扩展信息
var getCustomerExtendsCmd = &cobra.Command{
	Use:   "extends [customerId]",
	Short: "查询客户扩展信息",
	Long:  "通过客户 ID 查询客户运营扩展信息，包含等级、最近 30 天订单数量、店铺数量、从事 Dropshipping 时长、周广告预算、月订单量、客户来源等。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		customerId := args[0]
		info, err := apis.NewCustomerInfoApi().GetCustomerExtendInfo(customerId)
		if err != nil {
			cmd.PrintErrf("查询客户 %s 扩展信息失败：%v\n", customerId, err)
			return
		}
		if info == nil {
			cmd.Print("未查询到客户信息")
			return
		}

		var sb strings.Builder
		fmt.Fprintln(&sb, "【客户信息】")
		fmt.Fprintf(&sb, "客户ID：%s\n", customerId)
		fmt.Fprintf(&sb, "等级：%s\n", info.Level)
		fmt.Fprintf(&sb, "订单等级：%s\n", info.OrderLevel)
		fmt.Fprintf(&sb, "最近30天订单数：%d\n", info.OrderNum)
		fmt.Fprintf(&sb, "店铺数量：%d\n", info.StoreNum)
		fmt.Fprintf(&sb, "从事Dropshipping时长：%s\n", durationTypeText(info.DurationType))
		fmt.Fprintf(&sb, "周广告预算：%s\n", weeklyAdBudgetText(info.WeeklyAdBudget))
		fmt.Fprintf(&sb, "月订单量：%s\n", orderVolumeText(info.OrderVolume))
		fmt.Fprintf(&sb, "客户来源：%s\n", info.CustomerSource)
		fmt.Fprintf(&sb, "资料更新时间：%s", info.AllocatedTime)

		cmd.Print(sb.String())
	},
}

// 查询客户档案
var getCustomerDetailCmd = &cobra.Command{
	Use:   "detail [customerId]",
	Short: "查询客户档案",
	Long:  "通过客户 ID 查询客户基础档案，包含注册时间、注册区域、最近登录时间、总订单数与总购买金额等。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		customerId := args[0]
		info, err := apis.NewCustomerInfoApi().GetCustomerInfo(customerId)
		if err != nil {
			cmd.PrintErrf("查询客户 %s 档案失败：%v\n", customerId, err)
			return
		}
		if info == nil {
			cmd.Print("未查询到客户档案")
			return
		}

		var sb strings.Builder
		fmt.Fprintln(&sb, "【客户档案】")
		fmt.Fprintf(&sb, "客户ID：%s\n", customerId)
		fmt.Fprintf(&sb, "注册时间：%s\n", info.RegTime)
		fmt.Fprintf(&sb, "注册区域：%s\n", info.RegRegion)
		fmt.Fprintf(&sb, "最近登录时间：%s\n", info.LastLoginTime)
		// fmt.Fprintf(&sb, "总订单数：%d\n", info.TotalOrder)
		// fmt.Fprintf(&sb, "总购买金额：%.2f", info.TotalAmount)

		cmd.Print(sb.String())
	},
}

// 枚举值转文案，命中不了时展示原始值
func durationTypeText(v int) string {
	if s, ok := apis.DurationTypeMap[v]; ok {
		return s
	}
	return fmt.Sprint(v)
}

func weeklyAdBudgetText(v int) string {
	if s, ok := apis.WeeklyAdBudgetMap[v]; ok {
		return s
	}
	return fmt.Sprint(v)
}

func orderVolumeText(v int) string {
	if s, ok := apis.OrderVolumeMap[v]; ok {
		return s
	}
	return fmt.Sprint(v)
}

func init() {
	customerCmd.AddCommand(
		getCustomerExtendsCmd, // 客户扩展信息
		getCustomerDetailCmd,  // 客户档案
	)

	rootCmd.AddCommand(customerCmd)
}
