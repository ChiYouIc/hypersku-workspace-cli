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

// 仓库列表查询
var warehouseListCmd = &cobra.Command{
	Use:   "list",
	Short: "按名称查询仓库",
	Long:  "按仓库名称分页查询仓库列表",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")

		result, err := apis.NewWarehouseApi().GetWarehousePage(name, page, limit)
		if err != nil {
			cmd.PrintErrf("查询仓库列表失败: %v\n", err)
			return
		}
		if result == nil || len(result.Rows) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "共 %d 条\n\n无数据", result.Total)
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", page, limit, result.Total)
		fmt.Fprintln(cmd.OutOrStdout(), "|ID|仓库名称|联系人|联系电话|地址|状态|")
		fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|----|----|")
		for _, row := range result.Rows {
			fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%s|%s|%s%s|%s|\n",
				row.ID, row.Name, row.ContactsName, row.ContactsPhone, row.ProvinceName, row.CityName, row.StatusStr)
		}
	},
}

// 物流列表查询
var logisticsListCmd = &cobra.Command{
	Use:   "logistics-list",
	Short: "按名称查询物流",
	Long:  "按物流名称分页查询物流列表",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")

		result, err := apis.NewWarehouseApi().GetLogisticsPage(name, 0, page, limit)
		if err != nil {
			cmd.PrintErrf("查询物流列表失败: %v\n", err)
			return
		}
		if result == nil || len(result.Rows) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "共 %d 条\n\n无数据", result.Total)
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", page, limit, result.Total)
		fmt.Fprintln(cmd.OutOrStdout(), "|ID|名称|别名|物流类型|报价方式|状态|")
		fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|----|----|")
		for _, row := range result.Rows {
			fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%s|%d|%d|%d|\n",
				row.ID, row.Name, row.Alias, row.ShippingType, row.PriceType, row.Status)
		}
	},
}

func init() {
	warehouseListCmd.Flags().StringP("name", "n", "", "仓库名称")
	warehouseListCmd.Flags().IntP("page", "p", 1, "页码")
	warehouseListCmd.Flags().IntP("limit", "l", 20, "每页条数")

	logisticsListCmd.Flags().StringP("name", "n", "", "物流名称")
	logisticsListCmd.Flags().IntP("page", "p", 1, "页码")
	logisticsListCmd.Flags().IntP("limit", "l", 20, "每页条数")

	warehouseCmd.AddCommand(getWarehouseTracking)
	warehouseCmd.AddCommand(warehouseListCmd)
	warehouseCmd.AddCommand(logisticsListCmd)
	rootCmd.AddCommand(warehouseCmd)
}
