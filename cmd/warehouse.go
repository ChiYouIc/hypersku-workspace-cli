package cmd

import (
	"fmt"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

var warehouseCmd = &cobra.Command{
	Use:   "warehouse",
	Short: "仓库管理",
	Long:  "提供仓库相关信息查询，支持仓库物流轨迹等子命令。",
}

// 查询仓库物流轨迹
var getWarehouseTracking = &cobra.Command{
	Use:   "tracking [trackingNumber]",
	Short: "查询仓库物流轨迹",
	Long:  "根据物流单号查询仓库物流轨迹（快递签收、仓库签收、入库、物流轨迹、仓库操作等）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		trackingNumber := args[0]
		trackList, err := apis.NewWarehouseApi().GetWarehouseTracking(trackingNumber)
		if err != nil {
			cmd.PrintErr("查询仓库物流轨迹发生错误", err)
			cmd.Help()
			return
		}

		if len(trackList) == 0 {
			cmd.Print("未查询到仓库物流轨迹")
			return
		}

		instoreStatusMap := map[int]string{
			0: "未入库",
			1: "已入库",
		}

		var sb strings.Builder
		for i, info := range trackList {
			if i > 0 {
				fmt.Fprintln(&sb)
			}
			fmt.Fprintf(&sb, "【仓库物流轨迹：%s】\n", info.FinalNo)
			fmt.Fprintf(&sb, "物流单号：%s\n", info.FinalNo)
			fmt.Fprintf(&sb, "快递签收：%s\n", info.ExpressSignInfo)
			fmt.Fprintf(&sb, "快递签收时间：%s\n", info.ExpressSignTime)
			fmt.Fprintf(&sb, "仓库：%s\n", info.StoreAddressName)
			fmt.Fprintf(&sb, "仓库签收信息：%s\n", info.SignInfo)
			if info.SignTime != "" {
				fmt.Fprintf(&sb, "仓库签收时间：%s\n", info.SignTime)
			}
			fmt.Fprintf(&sb, "入库状态：%s\n", instoreStatusMap[info.InstoreStatus])
			// fmt.Fprintf(&sb, "入库信息：%s\n", info.InputInfo)
			if info.InstoreTime != "" {
				fmt.Fprintf(&sb, "入库时间：%s\n", info.InstoreTime)
			}

			// 物流轨迹步骤
			fmt.Fprintln(&sb, "物流轨迹：")
			if len(info.LogisticsSteps) == 0 {
				fmt.Fprintln(&sb, "无轨迹")
			} else {
				fmt.Fprintln(&sb, "|时间|轨迹|")
				fmt.Fprintln(&sb, "|----|----|")
				for _, step := range info.LogisticsSteps {
					fmt.Fprintf(&sb, "|%s|%s|\n", step.AcceptTime, step.Remark)
				}
			}

			// 仓库操作列表
			fmt.Fprintln(&sb, "仓库操作：")
			if len(info.ActionList) == 0 {
				fmt.Fprintln(&sb, "无操作记录")
			} else {
				fmt.Fprintln(&sb, "|操作时间|内容|")
				fmt.Fprintln(&sb, "|----|----|")
				for _, action := range info.ActionList {
					fmt.Fprintf(&sb, "|%s|%s|\n", action.ActionTime, action.Content)
				}
			}
		}

		cmd.Print(sb.String())
	},
}

func init() {
	warehouseCmd.AddCommand(getWarehouseTracking)
	rootCmd.AddCommand(warehouseCmd)
}
