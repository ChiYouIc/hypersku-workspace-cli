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
		fmt.Fprintf(&sb, "客户ID：%d\n", info.CustomerID)
		fmt.Fprintf(&sb, "是否有已支付订单：%s\n", boolText(info.HasOrder))
		fmt.Fprintf(&sb, "绑定店铺：%s\n", nonEmptyOrNone(info.Stores))
		fmt.Fprintln(&sb, "【问卷信息】")
		fmt.Fprintf(&sb, "DS 经验：%s\n", durationTypeText(info.EngagedTime))
		fmt.Fprintf(&sb, "周广告预算：%s\n", weeklyAdBudgetText(info.WeeklyAdBudget))
		fmt.Fprintf(&sb, "月订单量预期：%s\n", orderVolumeText(info.OrderVolume))
		fmt.Fprintf(&sb, "细分市场：%d\n", info.Niche)
		fmt.Fprintf(&sb, "意向服务：%s\n", nonEmptyOrNone(info.ServiceInterest))
		fmt.Fprintln(&sb, "【基础档案】")
		fmt.Fprintf(&sb, "公司：%s\n", nonEmptyOrNone(info.Company))
		fmt.Fprintf(&sb, "国家：%s\n", nonEmptyOrNone(info.CountryName))
		fmt.Fprintf(&sb, "注册区域：%s\n", nonEmptyOrNone(info.RegRegion))
		fmt.Fprintf(&sb, "注册时间：%s\n", nonEmptyOrNone(info.SignedUpAt))
		fmt.Fprintf(&sb, "最近登录时间：%s\n", nonEmptyOrNone(info.LastLoginTime))
		fmt.Fprintf(&sb, "用户等级：%s 订单量级别：%s\n", nonEmptyOrNone(info.Level), nonEmptyOrNone(info.OrderLevel))
		fmt.Fprintf(&sb, "客户标签：%s\n", customerTagText(info.Tag))
		fmt.Fprintln(&sb, "【渠道来源】")
		fmt.Fprintf(&sb, "渠道来源：%s\n", channelOrDash(info.ChannelSource))
		fmt.Fprintf(&sb, "渠道来源二级：%s\n", channelOrDash(info.ChannelSourceSub))
		fmt.Fprintf(&sb, "渠道 Medium：%s\n", channelOrDash(info.ChannelMedium))
		fmt.Fprintf(&sb, "渠道 Campaign：%s\n", channelOrDash(info.ChannelCampaign))
		fmt.Fprintf(&sb, "合作来源代码：%s\n", channelOrDash(info.PartnerSource))
		fmt.Fprintf(&sb, "渠道来源链接：%s", channelOrDash(info.ChannelURL))

		cmd.Print(sb.String())
	},
}

// 枚举值转文案，命中不了时展示原始值
func boolText(v bool) string {
	if v {
		return "是"
	}
	return "否"
}

func nonEmptyOrNone(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

func customerTagText(v int) string {
	switch v {
	case 0:
		return "默认"
	case 1:
		return "老用户"
	case 2:
		return "新用户"
	case 3:
		return "潜在用户"
	case 4:
		return "流失用户"
	}
	return fmt.Sprint(v)
}

// 渠道字段为空时展示 --
func channelOrDash(s string) string {
	if v := strings.TrimSpace(s); v != "" {
		return v
	}
	return "--"
}

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
		getCustomerDetailCmd, // 客户档案
	)

	rootCmd.AddCommand(customerCmd)
}
