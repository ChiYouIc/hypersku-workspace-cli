package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

// 风险等级映射
var riskLevelMap = map[int]string{
	0: "无风险",
	1: "低风险",
	2: "中风险",
	3: "高风险",
}

// AI 分析类型
var aiTaskTypeMap = map[string]string{
	"international_logistics": "国际物流异常",
	"inventory":               "库存动销",
	"purchase":                "采购售后",
	"supplier":                "供应商",
	"customer_order":          "客户订单",
}

var aiAnalysisCmd = &cobra.Command{
	Use:   "ai-analysis",
	Short: "AI 分析数据查询",
	Long: `AI 智能分析查询命令集。提供 HyperSKU 平台各维度 AI 分析数据的只读查询能力。

可用子命令：
  intl-logistics          国际物流异常分析
  intl-logistics-count    国际物流异常风险等级统计
  inventory               库存动销分析
  inventory-count         库存动销风险等级统计
  purchase                采购售后分析
  purchase-count          采购售后风险等级统计
  supplier                供应商分析
  supplier-count          供应商等级统计
  customer-order          客户订单分析
  customer-order-count    客户订单风险等级统计
  task-progress           AI 任务风险等级汇总`,
}

// ============ 通用 flag 绑定 ============

func bindAnalysisQueryFlags(cmd *cobra.Command, query *apis.AiAnalysisQuery) {
	cmd.Flags().IntVarP(&query.Page, "page", "p", 1, "页码")
	cmd.Flags().IntVarP(&query.Limit, "limit", "l", 20, "每页条数")
	cmd.Flags().StringVarP(&query.Username, "username", "u", "", "客户账号")
	cmd.Flags().StringVarP(&query.SearchKey, "search", "s", "", "搜索关键字")
	cmd.Flags().IntVar(&query.ManagerID, "manager-id", 0, "客户经理 ID")
	cmd.Flags().StringVar(&query.Type, "type", "", "处理状态")
	cmd.Flags().StringVar(&query.Source, "source", "", "来源")
	cmd.Flags().IntVar(&query.RiskLevel, "risk-level", 0, "风险等级 (0=全部 1=低 2=中 3=高)")
	cmd.Flags().IntVar(&query.HyperskuStatus, "status", 0, "异常状态")
	cmd.Flags().IntVar(&query.SupplierLevel, "supplier-level", 0, "供应商等级")
	cmd.Flags().StringVar(&query.Day, "day", "", "日期 (YYYY-MM-DD)")
}

// ============ 辅助函数 ============

func printCountTable(cmd *cobra.Command, title string, counts []apis.AiBaseCount) {
	fmt.Fprintf(cmd.OutOrStdout(), "\n%s:\n\n", title)
	if len(counts) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "无数据")
		return
	}
	fmt.Fprintln(cmd.OutOrStdout(), "|等级|数量|")
	fmt.Fprintln(cmd.OutOrStdout(), "|----|----|")
	for _, c := range counts {
		target := c.Target
		if v, err := strconv.Atoi(target); err == nil {
			if name, ok := riskLevelMap[v]; ok {
				target = fmt.Sprintf("%s (%d)", name, v)
			}
		}
		fmt.Fprintf(cmd.OutOrStdout(), "|%s|%d|\n", target, c.Count)
	}
}

// ============ 国际物流异常分析 ============

var aiIntlLogisticsCmd = func() *cobra.Command {
	var query apis.AiAnalysisQuery

	cmd := &cobra.Command{
		Use:   "intl-logistics",
		Short: "国际物流异常分析",
		Long:  "分页查询国际物流异常 AI 分析结果，含物流信息、AI 摘要、风险等级",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := apis.NewAiAnalysisApi().InternationalLogisticsTable(query)
			if err != nil {
				cmd.PrintErrf("查询国际物流异常分析失败: %v\n", err)
				return
			}
			if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "总数：%d\n\n无数据", 0)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", query.Page, query.Limit, res.Data.Total)
			fmt.Fprintln(cmd.OutOrStdout(), "|ID|客户|订单号|仓库|物流单号|物流公司|最新轨迹|状态|AI风险|AI摘要|")
			fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|----|----|----|----|----|----|")
			for _, row := range res.Data.Rows {
				summary := row.AiSummary
				if len(summary) > 50 {
					summary = summary[:50] + "..."
				}
				fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%d|%s|%s|%s|%s|%s|%s|%s|\n",
					row.ID, row.CustomerUsername, row.OrderID, row.WarehouseName,
					row.TrackingNumber, row.LogisticsCompanyName,
					row.LatestLogisticsTrack,
					row.LatestLogisticsStatusText,
					riskLevelMap[row.AiRiskLevel], summary)
			}
		},
	}
	bindAnalysisQueryFlags(cmd, &query)
	return cmd
}()

var aiIntlLogisticsCountCmd = func() *cobra.Command {
	var query apis.AiAnalysisQuery

	cmd := &cobra.Command{
		Use:   "intl-logistics-count",
		Short: "国际物流异常风险等级统计",
		Run: func(cmd *cobra.Command, args []string) {
			counts, err := apis.NewAiAnalysisApi().InternationalLogisticsCountRiskLevel(query)
			if err != nil {
				cmd.PrintErrf("统计失败: %v\n", err)
				return
			}
			printCountTable(cmd, "国际物流异常风险等级分布", counts)
		},
	}
	bindAnalysisQueryFlags(cmd, &query)
	return cmd
}()

// ============ 库存动销分析 ============

var aiInventoryCmd = func() *cobra.Command {
	var query apis.AiAnalysisQuery

	cmd := &cobra.Command{
		Use:   "inventory",
		Short: "库存动销分析",
		Long:  "分页查询库存动销 AI 分析结果，含库存数量、消耗速度、AI 摘要、风险等级",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := apis.NewAiAnalysisApi().InventoryTable(query)
			if err != nil {
				cmd.PrintErrf("查询库存动销分析失败: %v\n", err)
				return
			}
			if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "总数：%d\n\n无数据", 0)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", query.Page, query.Limit, res.Data.Total)
			fmt.Fprintln(cmd.OutOrStdout(), "|客户|商品|仓库|库存|15天消耗|预计可用天数|90天销量|库龄(天)|AI风险|AI摘要|")
			fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|----|----|----|----|----|----|")
			for _, row := range res.Data.Rows {
				summary := row.AiSummary
				if len(summary) > 40 {
					summary = summary[:40] + "..."
				}
				fmt.Fprintf(cmd.OutOrStdout(), "|%s|%s|%s|%d|%d|%d|%d|%d|%s|%s|\n",
					row.CustomerName, row.GoodsName, row.WarehouseName,
					row.Quantity, row.AlertInventory, row.PredictUsageDays,
					row.SalesQuantity, row.InventoryAge,
					riskLevelMap[row.AiRiskLevel], summary)
			}
		},
	}
	bindAnalysisQueryFlags(cmd, &query)
	return cmd
}()

var aiInventoryCountCmd = func() *cobra.Command {
	var query apis.AiAnalysisQuery

	cmd := &cobra.Command{
		Use:   "inventory-count",
		Short: "库存动销风险等级统计",
		Run: func(cmd *cobra.Command, args []string) {
			counts, err := apis.NewAiAnalysisApi().InventoryCountRiskLevel(query)
			if err != nil {
				cmd.PrintErrf("统计失败: %v\n", err)
				return
			}
			printCountTable(cmd, "库存动销风险等级分布", counts)
		},
	}
	bindAnalysisQueryFlags(cmd, &query)
	return cmd
}()

// ============ 采购售后分析 ============

var aiPurchaseCmd = func() *cobra.Command {
	var query apis.AiAnalysisQuery

	cmd := &cobra.Command{
		Use:   "purchase",
		Short: "采购售后分析",
		Long:  "分页查询采购售后 AI 分析结果，含物流/退款信息、AI 摘要、风险等级",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := apis.NewAiAnalysisApi().PurchaseTable(query)
			if err != nil {
				cmd.PrintErrf("查询采购售后分析失败: %v\n", err)
				return
			}
			if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "总数：%d\n\n无数据", 0)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", query.Page, query.Limit, res.Data.Total)
			fmt.Fprintln(cmd.OutOrStdout(), "|ID|来源|数据ID|交易号|物流单号|类型|风险|AI摘要|")
			fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|----|----|----|----|")
			for _, row := range res.Data.Rows {
				summary := row.Summary
				if len(summary) > 50 {
					summary = summary[:50] + "..."
				}
				sourceName := "物流异常"
				if row.Source == "after_sales_1688" {
					sourceName = "1688售后"
				}
				trackingNumber := ""
				if row.Logistics != nil {
					trackingNumber = row.Logistics.TrackingNumber
				}
				fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%s|%d|%s|%s|%s|%s|\n",
					row.ID, sourceName, row.DataID, row.ThirdOrderID,
					trackingNumber, row.Type, riskLevelMap[row.RiskLevel], summary)
			}
		},
	}
	bindAnalysisQueryFlags(cmd, &query)
	return cmd
}()

var aiPurchaseCountCmd = func() *cobra.Command {
	var query apis.AiAnalysisQuery

	cmd := &cobra.Command{
		Use:   "purchase-count",
		Short: "采购售后风险等级统计",
		Run: func(cmd *cobra.Command, args []string) {
			counts, err := apis.NewAiAnalysisApi().PurchaseCountRiskLevel(query)
			if err != nil {
				cmd.PrintErrf("统计失败: %v\n", err)
				return
			}
			printCountTable(cmd, "采购售后风险等级分布", counts)
		},
	}
	bindAnalysisQueryFlags(cmd, &query)
	return cmd
}()

// ============ 供应商分析 ============

var aiSupplierCmd = func() *cobra.Command {
	var query apis.AiAnalysisQuery

	cmd := &cobra.Command{
		Use:   "supplier",
		Short: "供应商 AI 分析",
		Long:  "分页查询供应商 AI 分析结果，含供应商信息、AI 摘要、等级",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := apis.NewAiAnalysisApi().SupplierTable(query)
			if err != nil {
				cmd.PrintErrf("查询供应商分析失败: %v\n", err)
				return
			}
			if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "总数：%d\n\n无数据", 0)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", query.Page, query.Limit, res.Data.Total)
			fmt.Fprintln(cmd.OutOrStdout(), "|供应商ID|名称|类型|等级|AI摘要|更新时间|")
			fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|----|----|")
			for _, row := range res.Data.Rows {
				summary := row.AiSummary
				if len(summary) > 50 {
					summary = summary[:50] + "..."
				}
				fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%s|%d|%s|%s|\n",
					row.SupplierID, row.SupplierName, row.SupplierTypeName,
					row.SupplierLevel, summary, row.AiSummaryUpdTime)
			}
		},
	}
	bindAnalysisQueryFlags(cmd, &query)
	return cmd
}()

var aiSupplierCountCmd = func() *cobra.Command {
	var query apis.AiAnalysisQuery

	cmd := &cobra.Command{
		Use:   "supplier-count",
		Short: "供应商等级统计",
		Run: func(cmd *cobra.Command, args []string) {
			counts, err := apis.NewAiAnalysisApi().SupplierCountSupplierLevel(query)
			if err != nil {
				cmd.PrintErrf("统计失败: %v\n", err)
				return
			}
			printCountTable(cmd, "供应商等级分布", counts)
		},
	}
	bindAnalysisQueryFlags(cmd, &query)
	return cmd
}()

// ============ 客户订单分析 ============

var aiCustomerOrderCmd = func() *cobra.Command {
	var query apis.AiAnalysisQuery

	cmd := &cobra.Command{
		Use:   "customer-order",
		Short: "客户订单 AI 分析",
		Long:  "分页查询客户订单 AI 分析结果，含客户信息、AI 摘要、风险等级",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := apis.NewAiAnalysisApi().CustomerOrderTable(query)
			if err != nil {
				cmd.PrintErrf("查询客户订单分析失败: %v\n", err)
				return
			}
			if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "总数：%d\n\n无数据", 0)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", query.Page, query.Limit, res.Data.Total)
			fmt.Fprintln(cmd.OutOrStdout(), "|客户ID|账号|标签|AI风险|AI摘要|更新时间|")
			fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|----|----|")
			for _, row := range res.Data.Rows {
				summary := row.AiSummary
				if len(summary) > 50 {
					summary = summary[:50] + "..."
				}
				fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%d|%s|%s|%s|\n",
					row.CustomerID, row.Username, row.Tag,
					riskLevelMap[row.AiRiskLevel], summary, row.AiSummaryUpdTime)
			}
		},
	}
	bindAnalysisQueryFlags(cmd, &query)
	return cmd
}()

var aiCustomerOrderCountCmd = func() *cobra.Command {
	var query apis.AiAnalysisQuery

	cmd := &cobra.Command{
		Use:   "customer-order-count",
		Short: "客户订单风险等级统计",
		Run: func(cmd *cobra.Command, args []string) {
			counts, err := apis.NewAiAnalysisApi().CustomerOrderCountRiskLevel(query)
			if err != nil {
				cmd.PrintErrf("统计失败: %v\n", err)
				return
			}
			printCountTable(cmd, "客户订单风险等级分布", counts)
		},
	}
	bindAnalysisQueryFlags(cmd, &query)
	return cmd
}()

// ============ 任务进度汇总 ============

var aiTaskProgressCmd = &cobra.Command{
	Use:   "task-progress [type]",
	Short: "AI 任务风险等级汇总",
	Long: `按分析类型统计 AI 任务的风险等级分布。

可用类型：
  international_logistics  国际物流异常
  inventory                库存动销
  purchase                 采购售后
  supplier                 供应商
  customer_order           客户订单`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskType := args[0]
		if _, ok := aiTaskTypeMap[taskType]; !ok {
			cmd.PrintErrf("未知的任务类型: %s\n可用类型: %s\n", taskType, strings.Join([]string{
				"international_logistics", "inventory", "purchase", "supplier", "customer_order",
			}, ", "))
			return
		}

		counts, err := apis.NewAiAnalysisApi().TaskProgressCountRiskLevelByType(taskType)
		if err != nil {
			cmd.PrintErrf("统计失败: %v\n", err)
			return
		}
		printCountTable(cmd, aiTaskTypeMap[taskType]+"风险等级汇总", counts)
	},
}

func init() {
	aiAnalysisCmd.AddCommand(aiIntlLogisticsCmd)
	aiAnalysisCmd.AddCommand(aiIntlLogisticsCountCmd)
	aiAnalysisCmd.AddCommand(aiInventoryCmd)
	aiAnalysisCmd.AddCommand(aiInventoryCountCmd)
	aiAnalysisCmd.AddCommand(aiPurchaseCmd)
	aiAnalysisCmd.AddCommand(aiPurchaseCountCmd)
	aiAnalysisCmd.AddCommand(aiSupplierCmd)
	aiAnalysisCmd.AddCommand(aiSupplierCountCmd)
	aiAnalysisCmd.AddCommand(aiCustomerOrderCmd)
	aiAnalysisCmd.AddCommand(aiCustomerOrderCountCmd)
	aiAnalysisCmd.AddCommand(aiTaskProgressCmd)

	rootCmd.AddCommand(aiAnalysisCmd)
}
