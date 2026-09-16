package cmd

import (
	"fmt"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

var supplierCmd = &cobra.Command{
	Use:   "supplier",
	Short: "供应商管理查询",
	Long: `供应商管理查询命令集（只读）。提供供应商详情、列表、采购次数、排行、售后、履约等查询能力。

可用子命令：
  detail         查询供应商详情
  list           分页查询供应商列表
  pur-count      根据供应商登录ID查询采购次数
  ranking        供应商排行信息
  spu-list       供应商 SPU 信息
  after-sales    供应商售后信息
  fulfillment    供应商履约信息`,
}

// ============ supplier detail ============

var supplierDetailCmd = &cobra.Command{
	Use:   "detail [supplierId]",
	Short: "查询供应商详情",
	Long:  "根据供应商 ID 查询供应商完整信息（名称、链接、联系人、地址、采购员等）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := 0
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil || id == 0 {
			cmd.PrintErrf("无效的供应商 ID: %s\n", args[0])
			return
		}

		supplier, err := apis.NewSupplierApi().GetSupplier(id)
		if err != nil {
			cmd.PrintErrf("查询供应商失败: %v\n", err)
			return
		}

		platformMap := map[int]string{
			1: "1688", 17: "京东", 18: "淘宝", 19: "天猫", 25: "拼多多", 99: "HyperFrom",
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "供应商ID: %d\n", supplier.ID)
		fmt.Fprintf(&sb, "名称: %s\n", supplier.Name)
		fmt.Fprintf(&sb, "平台: %s\n", platformMap[supplier.Type])
		fmt.Fprintf(&sb, "平台登录ID: %s\n", supplier.PlatformLoginID)
		fmt.Fprintf(&sb, "平台会员ID: %s\n", supplier.PlatformMemberID)
		fmt.Fprintf(&sb, "链接: %s\n", supplier.URL)
		fmt.Fprintf(&sb, "邮箱: %s\n", supplier.Email)
		fmt.Fprintf(&sb, "旺旺号: %s\n", supplier.WangWangAccount)
		fmt.Fprintf(&sb, "联系人: %s\n", supplier.ContactsName)
		fmt.Fprintf(&sb, "联系人手机: %s\n", supplier.ContactsPhone)
		fmt.Fprintf(&sb, "微信: %s\n", supplier.Wechat)
		fmt.Fprintf(&sb, "到货周期: %d 天\n", supplier.ArrivalPeriod)
		fmt.Fprintf(&sb, "国家: %s\n", supplier.CountryName)
		fmt.Fprintf(&sb, "省份: %s\n", supplier.SecondRegionName)
		fmt.Fprintf(&sb, "城市: %s\n", supplier.ThirdRegionName)
		fmt.Fprintf(&sb, "地址: %s\n", supplier.AddressDetail)
		fmt.Fprintf(&sb, "采购员: %s\n", supplier.PurchaserName)
		fmt.Fprintf(&sb, "备注: %s\n", supplier.Remark)
		cmd.Print(sb.String())
	},
}

// ============ supplier list ============

var supplierListCmd = &cobra.Command{
	Use:   "list",
	Short: "分页查询供应商列表",
	Long:  "分页查询供应商列表，支持按名称筛选",
	Run: func(cmd *cobra.Command, args []string) {
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")
		name, _ := cmd.Flags().GetString("name")

		result, err := apis.NewSupplierApi().ListSupplier(apis.SupplierPageQuery{
			Page:  page,
			Limit: limit,
			Name:  name,
		})
		if err != nil {
			cmd.PrintErrf("查询供应商列表失败: %v\n", err)
			return
		}
		if result == nil || result.Data == nil || len(result.Data.Rows) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n无数据", page, limit, 0)
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", page, limit, result.Data.Total)
		// map[string]interface{} 类型，输出 key-value 格式
		for i, row := range result.Data.Rows {
			fmt.Fprintf(cmd.OutOrStdout(), "--- 供应商 %d ---\n", i+1)
			if v, ok := row["id"]; ok {
				fmt.Fprintf(cmd.OutOrStdout(), "ID: %v\n", v)
			}
			if v, ok := row["name"]; ok {
				fmt.Fprintf(cmd.OutOrStdout(), "名称: %v\n", v)
			}
			if v, ok := row["type"]; ok {
				fmt.Fprintf(cmd.OutOrStdout(), "类型: %v\n", v)
			}
			if v, ok := row["url"]; ok {
				fmt.Fprintf(cmd.OutOrStdout(), "链接: %v\n", v)
			}
			if v, ok := row["state"]; ok {
				fmt.Fprintf(cmd.OutOrStdout(), "状态: %v\n", v)
			}
		}
	},
}

// ============ supplier pur-count ============

var supplierPurCountCmd = &cobra.Command{
	Use:   "pur-count [loginId1,loginId2,...]",
	Short: "查询供应商采购次数",
	Long:  "根据供应商登录ID（逗号分隔）查询各供应商的采购次数",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		loginIds := strings.Split(args[0], ",")
		result, err := apis.NewSupplierApi().SupplierPurNumCount(loginIds)
		if err != nil {
			cmd.PrintErrf("查询采购次数失败: %v\n", err)
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), "|供应商登录ID|采购次数|")
		fmt.Fprintln(cmd.OutOrStdout(), "|----|----|")
		for k, v := range result {
			fmt.Fprintf(cmd.OutOrStdout(), "|%s|%d|\n", k, v)
		}
	},
}

// ============ supplier ranking ============

var supplierRankingCmd = func() *cobra.Command {
	var query apis.SupplierRankingQuery

	cmd := &cobra.Command{
		Use:   "ranking",
		Short: "供应商排行信息",
		Long:  "查询供应商排行数据（采购量、交易额、发货率、异常率等）",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := apis.NewSupplierApi().SupplyRankingList(query)
			if err != nil {
				cmd.PrintErrf("查询供应商排行失败: %v\n", err)
				return
			}
			if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n无数据", query.Page, query.Limit, 0)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", query.Page, query.Limit, res.Data.Total)
			fmt.Fprintln(cmd.OutOrStdout(), "|ID|名称|城市|采购量|订单数|交易额|销售额|发货率|异常率|利润率|")
			fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|----|----|----|----|----|----|")
			for _, row := range res.Data.Rows {
				fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%s|%d|%d|%.2f|%.2f|%.1f%%|%.1f%%|%.1f%%|\n",
					row.SupplyID, row.Name, row.SecondRegion,
					row.PurTotal, row.OrdNum, row.TotalAmount, row.SalesAmount,
					row.DeliveryRatio*100, row.AbnormalRatio*100, row.ProfitMargin*100)
			}
		},
	}
	cmd.Flags().IntVarP(&query.Page, "page", "p", 1, "页码")
	cmd.Flags().IntVarP(&query.Limit, "limit", "l", 20, "每页条数")
	cmd.Flags().StringVarP(&query.Name, "name", "n", "", "供应商名称")
	cmd.Flags().IntVar(&query.SearchDay, "days", 0, "近N天")
	cmd.Flags().StringVar(&query.Sort, "sort", "purTotal desc", "排序字段")
	cmd.Flags().StringVar(&query.Type, "type", "", "类型")
	return cmd
}()

// ============ supplier spu-list ============

var supplierSPUCmd = func() *cobra.Command {
	var query apis.SupplierRankingQuery

	cmd := &cobra.Command{
		Use:   "spu-list",
		Short: "供应商 SPU 信息",
		Long:  "查询供应商 SPU 维度数据",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := apis.NewSupplierApi().SupplySPUList(query)
			if err != nil {
				cmd.PrintErrf("查询供应商 SPU 失败: %v\n", err)
				return
			}
			if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n无数据", query.Page, query.Limit, 0)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", query.Page, query.Limit, res.Data.Total)
			fmt.Fprintln(cmd.OutOrStdout(), "|ID|名称|SPU数|SKU数|类目|采购均价|")
			fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|----|----|")
			for _, row := range res.Data.Rows {
				fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%d|%d|%s|%.2f|\n",
					row.SupplyID, row.Name, row.SpuNum, row.SkuNum, row.CategoryName, row.AvgPurPrice)
			}
		},
	}
	cmd.Flags().IntVarP(&query.Page, "page", "p", 1, "页码")
	cmd.Flags().IntVarP(&query.Limit, "limit", "l", 20, "每页条数")
	cmd.Flags().StringVarP(&query.Name, "name", "n", "", "供应商名称")
	return cmd
}()

// ============ supplier after-sales ============

var supplierAfterSalesCmd = func() *cobra.Command {
	var query apis.SupplierRankingQuery

	cmd := &cobra.Command{
		Use:   "after-sales",
		Short: "供应商售后信息",
		Long:  "查询供应商售后数据（成功/超时/关闭/退款数量）",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := apis.NewSupplierApi().SupplierAfterSalesList(query)
			if err != nil {
				cmd.PrintErrf("查询供应商售后失败: %v\n", err)
				return
			}
			if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n无数据", query.Page, query.Limit, 0)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", query.Page, query.Limit, res.Data.Total)
			fmt.Fprintln(cmd.OutOrStdout(), "|ID|供应商|交易总数|成功|超时关闭|关闭|退款成功|退款失败|退款总数|")
			fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|----|----|----|----|----|")
			for _, row := range res.Data.Rows {
				fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%d|%d|%d|%d|%d|%d|%d|\n",
					row.SupplyID, row.SupplierName, row.TradeTotalNum,
					row.SuccessNum, row.OverTimeCloseNum, row.CloseNum,
					row.ReturnSuccessNum, row.ReturnFailureNum, row.ReturnTotalNum)
			}
		},
	}
	cmd.Flags().IntVarP(&query.Page, "page", "p", 1, "页码")
	cmd.Flags().IntVarP(&query.Limit, "limit", "l", 20, "每页条数")
	cmd.Flags().StringVarP(&query.Name, "name", "n", "", "供应商名称")
	return cmd
}()

// ============ supplier fulfillment ============

var supplierFulfillmentCmd = func() *cobra.Command {
	var query apis.SupplierRankingQuery

	cmd := &cobra.Command{
		Use:   "fulfillment",
		Short: "供应商履约信息",
		Long:  "查询供应商履约数据（发货率、时效、错发少发等）",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := apis.NewSupplierApi().SupplierFulfillmentList(query)
			if err != nil {
				cmd.PrintErrf("查询供应商履约失败: %v\n", err)
				return
			}
			if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n无数据", query.Page, query.Limit, 0)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "当前页码：%d，页大小：%d，总数：%d\n\n", query.Page, query.Limit, res.Data.Total)
			fmt.Fprintln(cmd.OutOrStdout(), "|ID|供应商|已发货|未发货|平均时效|发货率|异常率|错发|少发|交易额|")
			fmt.Fprintln(cmd.OutOrStdout(), "|----|----|----|----|----|----|----|----|----|----|")
			for _, row := range res.Data.Rows {
				fmt.Fprintf(cmd.OutOrStdout(), "|%d|%s|%d|%d|%.1f|%.1f%%|%.1f%%|%d|%d|%.2f|\n",
					row.SupplyID, row.SupplierName,
					row.ShipmentsNum, row.UnshippedNum, row.AvgTime,
					row.DeliveryRatio*100, row.AbnormalRatio*100,
					row.WdNum, row.LessNum, row.TotalAmount)
			}
		},
	}
	cmd.Flags().IntVarP(&query.Page, "page", "p", 1, "页码")
	cmd.Flags().IntVarP(&query.Limit, "limit", "l", 20, "每页条数")
	cmd.Flags().StringVarP(&query.Name, "name", "n", "", "供应商名称")
	return cmd
}()

func init() {
	supplierListCmd.Flags().IntP("page", "p", 1, "页码")
	supplierListCmd.Flags().IntP("limit", "l", 20, "每页条数")
	supplierListCmd.Flags().StringP("name", "n", "", "供应商名称")

	supplierCmd.AddCommand(supplierDetailCmd)
	supplierCmd.AddCommand(supplierListCmd)
	supplierCmd.AddCommand(supplierPurCountCmd)
	supplierCmd.AddCommand(supplierRankingCmd)
	supplierCmd.AddCommand(supplierSPUCmd)
	supplierCmd.AddCommand(supplierAfterSalesCmd)
	supplierCmd.AddCommand(supplierFulfillmentCmd)

	rootCmd.AddCommand(supplierCmd)
}
