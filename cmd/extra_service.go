package cmd

import (
	"fmt"
	"strconv"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

var extraServiceCmd = &cobra.Command{
	Use:   "extra-service",
	Short: "增值服务查询",
	Long: `增值服务查询命令集（只读）。支持按名称、仓库查询增值服务信息。

可用子命令：
  list           查询增值服务名称列表
  warehouses     查询某增值服务关联的仓储
  goods          查询某增值服务关联的商品
  warehouse-add  查询某仓库的增值服务`,
}

// 增值服务名称列表
var extraServiceListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询增值服务名称列表",
	Long:  "获取所有可用的增值服务名称",
	Run: func(cmd *cobra.Command, args []string) {
		names, err := apis.NewExtraServiceApi().GetExtraServiceNames()
		if err != nil {
			cmd.PrintErrf("查询增值服务列表失败: %v\n", err)
			return
		}
		if len(names) == 0 {
			cmd.Println("暂无增值服务")
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), "|ID|名称|类型|")
		fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|")
		for _, n := range names {
			fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%d|\n", n.ID, n.Name, n.Type)
		}
	},
}

// 查询增值服务关联的仓储
var extraServiceWarehousesCmd = &cobra.Command{
	Use:   "warehouses [serviceId]",
	Short: "查询增值服务关联的仓储",
	Long:  "根据增值服务 ID 查询该服务关联的仓储列表",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceId, err := strconv.Atoi(args[0])
		if err != nil || serviceId == 0 {
			cmd.PrintErrf("无效的服务ID: %s\n", args[0])
			return
		}
		list, err := apis.NewExtraServiceApi().ListValueAddedByService(serviceId)
		if err != nil {
			cmd.PrintErrf("查询增值服务仓储失败: %v\n", err)
			return
		}
		if len(list) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "服务ID %d 无关联仓储\n", serviceId)
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), "|仓储ID|仓储名称|服务名称|状态|")
		fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|")
		for _, item := range list {
			statusStr := "启用"
			if !item.Status {
				statusStr = "停用"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%s|%s|\n",
				item.RepertoryID, item.RepertoryName, item.Name, statusStr)
		}
	},
}

// 查询增值服务关联的商品
var extraServiceGoodsCmd = &cobra.Command{
	Use:   "goods [serviceId]",
	Short: "查询增值服务关联的商品",
	Long:  "根据增值服务 ID 查询该服务关联的商品列表",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceId, err := strconv.Atoi(args[0])
		if err != nil || serviceId == 0 {
			cmd.PrintErrf("无效的服务ID: %s\n", args[0])
			return
		}
		list, err := apis.NewExtraServiceApi().ListGoodsByService(serviceId)
		if err != nil {
			cmd.PrintErrf("查询增值服务商品失败: %v\n", err)
			return
		}
		if len(list) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "服务ID %d 无关联商品\n", serviceId)
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), "|商品ID|商品名称|")
		fmt.Fprintln(cmd.OutOrStdout(), "|----|----|")
		for _, item := range list {
			fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|\n", item.GoodsID, item.GoodsName)
		}
	},
}

// 查询仓库的增值服务
var extraServiceWarehouseAddCmd = &cobra.Command{
	Use:   "warehouse-add [repertoryId]",
	Short: "查询仓库的增值服务",
	Long:  "根据仓库 ID 查询该仓库开通的增值服务列表",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repertoryId, err := strconv.Atoi(args[0])
		if err != nil || repertoryId == 0 {
			cmd.PrintErrf("无效的仓库ID: %s\n", args[0])
			return
		}
		names, err := apis.NewExtraServiceApi().ListWarehouseAddValue(repertoryId)
		if err != nil {
			cmd.PrintErrf("查询仓库增值服务失败: %v\n", err)
			return
		}
		if len(names) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "仓库 %d 暂无增值服务\n", repertoryId)
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "仓库 %d 开通的增值服务:\n\n", repertoryId)
		for i, name := range names {
			fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s\n", i+1, name)
		}
	},
}

func init() {
	extraServiceCmd.AddCommand(extraServiceListCmd)
	extraServiceCmd.AddCommand(extraServiceWarehousesCmd)
	extraServiceCmd.AddCommand(extraServiceGoodsCmd)
	extraServiceCmd.AddCommand(extraServiceWarehouseAddCmd)

	rootCmd.AddCommand(extraServiceCmd)
}
