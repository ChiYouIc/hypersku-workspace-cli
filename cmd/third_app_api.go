package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

var thirdAppApiCmd = &cobra.Command{
	Use:   "third-app-api",
	Short: "第三方平台数据查询",
	Long: `第三方平台（1688、Shopify 等）聚合查询命令集。仅提供数据查询能力，不包含写入/修改操作。

可用子命令：
  ali-logistics          查询 1688 订单物流信息
  ali-logistics-trace    查询 1688 订单物流轨迹
  ali-order-detail       查询 1688 订单详情
  ali-refund-detail      查询 1688 退款单详情（按退款单号）
  ali-refund-list        查询 1688 退款单列表（按交易号）
  ali-refund-operations  查询 1688 退款单操作记录
  ali-product            查询 1688 产品信息（按ID或URL）
  ali-supplier           查询 1688 供应商信息
  ali-mix-config         查询 1688 卖家混批设置
  ali-sub-accounts       查询 1688 子账号列表
  shop-product           查询店铺产品
  shop-order-list        查询店铺订单列表
  shop-inventory         查询店铺库存`,
}

// ============ 1688 子命令 ============

// ali-logistics 查询 1688 订单物流信息
var aliLogisticsCmd = &cobra.Command{
	Use:   "ali-logistics [orderId]",
	Short: "查询 1688 订单物流信息",
	Long:  "根据 1688 交易号查询订单物流信息（运单号、物流公司、发货状态）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := apis.NewThirdAppApi().GetAliLogisticsInfo(args[0])
		if err != nil {
			cmd.PrintErrf("查询 1688 物流信息失败: %v\n", err)
			return
		}
		if !result.Success {
			cmd.Printf("查询失败: %s (%s)\n", result.ErrorMessage, result.ErrorCode)
			return
		}
		if len(result.Result) == 0 {
			cmd.Println("未查询到物流信息")
			return
		}
		var sb strings.Builder
		for _, lo := range result.Result {
			fmt.Fprintf(&sb, "物流单号: %s\n", lo.LogisticsBillNo)
			fmt.Fprintf(&sb, "物流公司: %s\n", lo.LogisticsCompanyName)
			fmt.Fprintf(&sb, "发货状态: %s\n", lo.Status)
			fmt.Fprintf(&sb, "订单号列表: %s\n", lo.OrderEntryIds)
			if len(lo.SendGoods) > 0 {
				fmt.Fprintln(&sb, "商品:")
				for _, g := range lo.SendGoods {
					fmt.Fprintf(&sb, "  - %s (数量: %s, 重量: %s)\n", g.Name, g.SendGoodsAmount, g.SendGoodsWeight)
				}
			}
			fmt.Fprintln(&sb, "---")
		}
		cmd.Print(sb.String())
	},
}

// ali-logistics-trace 查询 1688 订单物流轨迹
var aliLogisticsTraceCmd = &cobra.Command{
	Use:   "ali-logistics-trace [orderId]",
	Short: "查询 1688 订单物流轨迹",
	Long:  "根据 1688 交易号查询订单的完整物流轨迹（包括仓库状态）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := apis.NewThirdAppApi().GetAliLogisticsTraceInfo(args[0])
		if err != nil {
			cmd.PrintErrf("查询 1688 物流轨迹失败: %v\n", err)
			return
		}
		if !result.Success {
			cmd.Printf("查询失败: %s (%s)\n", result.ErrorMessage, result.ErrorCode)
			return
		}
		if len(result.Result) == 0 {
			cmd.Println("未查询到物流轨迹")
			return
		}
		var sb strings.Builder
		for _, trace := range result.Result {
			fmt.Fprintf(&sb, "物流编号: %s\n", trace.LogisticsID)
			fmt.Fprintf(&sb, "物流单号: %s\n", trace.LogisticsBillNo)
			fmt.Fprintf(&sb, "物流公司: %s\n", trace.LogisticsCompany)
			fmt.Fprintf(&sb, "阿里状态: %s\n", trace.AliStatusStr)
			fmt.Fprintf(&sb, "仓库状态: %s\n", trace.WarehouseStatusStr)
			fmt.Fprintf(&sb, "仓库名称: %s\n", trace.WarehouseName)

			if len(trace.LogisticsSteps) > 0 {
				fmt.Fprintln(&sb, "\n物流轨迹:")
				fmt.Fprintln(&sb, "|时间|轨迹|")
				fmt.Fprintln(&sb, "|----|----|")
				for _, step := range trace.LogisticsSteps {
					fmt.Fprintf(&sb, "|%s|%s|\n", step.AcceptTime, step.Remark)
				}
			}
			if len(trace.WarehouseLogisticsSteps) > 0 {
				fmt.Fprintln(&sb, "\n仓库物流轨迹:")
				fmt.Fprintln(&sb, "|时间|轨迹|")
				fmt.Fprintln(&sb, "|----|----|")
				for _, step := range trace.WarehouseLogisticsSteps {
					fmt.Fprintf(&sb, "|%s|%s|\n", step.AcceptTime, step.Remark)
				}
			}
			fmt.Fprintln(&sb, "---")
		}
		cmd.Print(sb.String())
	},
}

// ali-order-detail 查询 1688 订单详情
var aliOrderDetailCmd = &cobra.Command{
	Use:   "ali-order-detail [orderId]",
	Short: "查询 1688 订单详情",
	Long:  "根据 1688 交易号查询订单详情（含商品明细、物流、交易条款），请求时间较长",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := apis.NewThirdAppApi().GetAliOrderDetail(args[0])
		if err != nil {
			cmd.PrintErrf("查询 1688 订单详情失败: %v\n", err)
			return
		}
		if !result.Success {
			cmd.Printf("查询失败: %s (%s)\n", result.ErrorMessage, result.ErrorCode)
			return
		}
		if result.Result == nil {
			cmd.Println("未查询到订单详情")
			return
		}
		d := result.Result
		if d.BaseInfo != nil {
			b := d.BaseInfo
			var sb strings.Builder
			fmt.Fprintf(&sb, "订单ID: %d\n", b.ID)
			fmt.Fprintf(&sb, "交易状态: %s\n", b.Status)
			fmt.Fprintf(&sb, "总金额: %.2f\n", b.TotalAmount)
			fmt.Fprintf(&sb, "产品付款金额: %.2f\n", b.SumProductPayment)
			fmt.Fprintf(&sb, "运费: %.2f\n", b.ShippingFee)
			fmt.Fprintf(&sb, "退款金额: %.2f\n", b.Refund)
			fmt.Fprintf(&sb, "创建时间: %s\n", b.CreateTime)
			fmt.Fprintf(&sb, "付款时间: %s\n", b.PayTime)
			fmt.Fprintf(&sb, "收货时间: %s\n", b.ReceivingTime)
			fmt.Fprintf(&sb, "买家: %s\n", b.BuyerLoginID)
			fmt.Fprintf(&sb, "卖家: %s\n", b.SellerLoginID)

			if len(d.ProductItems) > 0 {
				fmt.Fprintln(&sb, "\n商品明细:")
				fmt.Fprintln(&sb, "|名称|单价|数量|金额|子订单状态|")
				fmt.Fprintln(&sb, "|----|----|----|----|----|")
				for _, item := range d.ProductItems {
					fmt.Fprintf(&sb, "|%s|%.2f|%.0f|%.2f|%s|\n",
						item.Name, item.Price, item.Quantity, item.ItemAmount, item.Status)
				}
			}
			cmd.Print(sb.String())
		}
	},
}

// ali-refund-detail 查询 1688 退款单详情
var aliRefundDetailCmd = &cobra.Command{
	Use:   "ali-refund-detail [refundId]",
	Short: "查询 1688 退款单详情",
	Long:  "根据退款单逻辑主键（TQ+ID）查询 1688 退款单详情",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := apis.NewThirdAppApi().QueryAliRefundDetail(apis.AliRefundDetailQuery{
			RefundID:                 args[0],
			NeedTimeOutInfo:          true,
			NeedOrderRefundOperation: true,
		})
		if err != nil {
			cmd.PrintErrf("查询 1688 退款单详情失败: %v\n", err)
			return
		}
		if result.ErrorMessage != "" {
			cmd.Printf("查询失败: %s (%s)\n", result.ErrorMessage, result.ErrorCode)
			return
		}
		if result.Result == nil || result.Result.OpOrderRefundModelDetail == nil {
			cmd.Println("未查询到退款单信息")
			return
		}
		d := result.Result.OpOrderRefundModelDetail
		var sb strings.Builder
		fmt.Fprintf(&sb, "退款单编号: %d\n", d.ID)
		fmt.Fprintf(&sb, "退款单逻辑主键: %s\n", d.RefundID)
		fmt.Fprintf(&sb, "退款状态: %s\n", d.Status)
		fmt.Fprintf(&sb, "订单编号: %d\n", d.OrderID)
		fmt.Fprintf(&sb, "产品名称: %s\n", d.ProductName)
		fmt.Fprintf(&sb, "申请退款金额: %d 分 (%.2f 元)\n", d.ApplyPayment, float64(d.ApplyPayment)/100)
		fmt.Fprintf(&sb, "实际退款金额: %d 分 (%.2f 元)\n", d.RefundPayment, float64(d.RefundPayment)/100)
		fmt.Fprintf(&sb, "申请原因: %s\n", d.ApplyReason)
		fmt.Fprintf(&sb, "买家: %s\n", d.BuyerLoginID)
		fmt.Fprintf(&sb, "卖家: %s\n", d.SellerLoginID)
		fmt.Fprintf(&sb, "申请时间: %s\n", d.GmtApply)
		fmt.Fprintf(&sb, "完成时间: %s\n", d.GmtCompleted)
		if d.OnlyRefund {
			fmt.Fprintf(&sb, "类型: 仅退款\n")
		} else {
			fmt.Fprintf(&sb, "类型: 退货退款\n")
		}
		if d.RejectReason != "" {
			fmt.Fprintf(&sb, "拒绝原因: %s\n", d.RejectReason)
		}
		cmd.Print(sb.String())
	},
}

// ali-refund-list 查询 1688 退款单列表
var aliRefundListCmd = &cobra.Command{
	Use:   "ali-refund-list [orderId]",
	Short: "查询 1688 退款单列表",
	Long:  "根据交易号查询该订单下的所有退款单。queryType: 1=活动中, 3=退款成功（含退款中）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := apis.NewThirdAppApi().QueryAliRefundList(apis.AliRefundListQuery{
			OrderID:   args[0],
			QueryType: "3",
		})
		if err != nil {
			cmd.PrintErrf("查询 1688 退款单列表失败: %v\n", err)
			return
		}
		if result.ErrorMessage != "" {
			cmd.Printf("查询失败: %s (%s)\n", result.ErrorMessage, result.ErrorCode)
			return
		}
		if result.Result == nil || len(result.Result.OpOrderRefundModels) == 0 {
			cmd.Println("未查询到退款单记录")
			return
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "|退款单号|退款单ID|退款金额(元)|退款状态|产品名称|申请时间|完成时间|\n")
		fmt.Fprintf(&sb, "|----|----|----|----|----|----|----|\n")
		for _, r := range result.Result.OpOrderRefundModels {
			fmt.Fprintf(&sb, "|%s|%d|%.2f|%s|%s|%s|%s|\n",
				r.RefundID, r.ID, float64(r.RefundPayment)/100, r.Status, r.ProductName, r.GmtApply, r.GmtCompleted)
		}
		cmd.Print(sb.String())
	},
}

// ali-refund-operations 查询 1688 退款单操作记录
var aliRefundOperationsCmd = &cobra.Command{
	Use:   "ali-refund-operations [refundId]",
	Short: "查询 1688 退款单操作记录",
	Long:  "根据退款单逻辑主键查询退款单操作记录列表",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := apis.NewThirdAppApi().QueryAliRefundOperationList(apis.AliRefundOperationQuery{
			RefundID: args[0],
			PageNo:   "1",
			PageSize: "50",
		})
		if err != nil {
			cmd.PrintErrf("查询 1688 退款操作记录失败: %v\n", err)
			return
		}
		if result.ErrorMessage != "" {
			cmd.Printf("查询失败: %s (%s)\n", result.ErrorMessage, result.ErrorCode)
			return
		}
		if result.Result == nil || len(result.Result.OpOrderRefundOperationModels) == 0 {
			cmd.Println("未查询到操作记录")
			return
		}
		var sb strings.Builder
		fmt.Fprintln(&sb, "|时间|操作前状态|操作后状态|操作人|操作备注|")
		fmt.Fprintln(&sb, "|----|----|----|----|----|")
		for _, op := range result.Result.OpOrderRefundOperationModels {
			fmt.Fprintf(&sb, "|%s|%s|%s|%s|%s|\n",
				op.GmtCreate, op.BeforeOperateStatus, op.AfterOperateStatus, op.OperatorLoginID, op.OperateRemark)
		}
		cmd.Print(sb.String())
	},
}

// ali-product 查询 1688 产品信息
var aliProductCmd = &cobra.Command{
	Use:   "ali-product [idOrUrl]",
	Short: "查询 1688 产品信息",
	Long: `查询 1688 产品信息。支持两种模式：
  按产品ID:  hypersku-cli third-app-api ali-product --id 123456
  按产品URL: hypersku-cli third-app-api ali-product --url "https://detail.1688.com/offer/xxx.html"`,
	Run: func(cmd *cobra.Command, args []string) {
		productID, _ := cmd.Flags().GetInt64("id")
		productURL, _ := cmd.Flags().GetString("url")

		if productID == 0 && productURL == "" {
			cmd.PrintErr("请指定 --id 或 --url 参数")
			return
		}

		if productID > 0 {
			result, err := apis.NewThirdAppApi().GetAliProductInfoById(strconv.FormatInt(productID, 10))
			if err != nil {
				cmd.PrintErrf("查询 1688 产品信息失败: %v\n", err)
				return
			}
			printAliProduct(cmd, &result.Data, result.Message)
			return
		}

		result, err := apis.NewThirdAppApi().GetAliProductInfoByUrl(productURL, true)
		if err != nil {
			cmd.PrintErrf("查询 1688 产品信息失败: %v\n", err)
			return
		}
		printAliProduct(cmd, &result.Data, result.Message)
	},
}

func printAliProduct(cmd *cobra.Command, p *apis.AliGoodsBaseInfo, msg string) {
	if p == nil {
		cmd.Printf("未查询到产品信息: %s\n", msg)
		return
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "产品ID: %d\n", p.ID)
	fmt.Fprintf(&sb, "名称: %s\n", p.Name)
	fmt.Fprintf(&sb, "英文名: %s\n", p.EnName)
	fmt.Fprintf(&sb, "链接: %s\n", p.URL)
	fmt.Fprintf(&sb, "原价: %s %s\n", p.OriginalPrice, p.OriginalCurrencySymbol)
	fmt.Fprintf(&sb, "销售价: %.2f %s\n", p.SalePrice, p.CurrencySymbol)
	fmt.Fprintf(&sb, "重量: %.2f kg\n", p.Weight)
	fmt.Fprintf(&sb, "净重: %.2f kg\n", p.SuttleWeight)
	fmt.Fprintf(&sb, "发货地: %s\n", p.SendGoodsAddress)
	fmt.Fprintf(&sb, "库存: %d\n", p.TotalStore)
	fmt.Fprintf(&sb, "类目: %s / %s\n", p.CategoryName, p.CategoryEnName)
	fmt.Fprintf(&sb, "供应商: %s (%s)\n", p.SupplierName, p.SupplierID)
	if p.IsCombinedSku {
		fmt.Fprintf(&sb, "支持组合SKU: 是\n")
	}
	cmd.Print(sb.String())
}

// ali-supplier 查询 1688 供应商信息
var aliSupplierCmd = &cobra.Command{
	Use:   "ali-supplier",
	Short: "查询 1688 供应商信息",
	Long:  "根据 1688 登录ID或店铺地址查询供应商信息",
	Run: func(cmd *cobra.Command, args []string) {
		loginID, _ := cmd.Flags().GetString("login-id")
		domain, _ := cmd.Flags().GetString("domain")

		if loginID == "" && domain == "" {
			cmd.PrintErr("请指定 --login-id 或 --domain 参数")
			return
		}

		result, err := apis.NewThirdAppApi().GetAliSupplierInfo(loginID, domain)
		if err != nil {
			cmd.PrintErrf("查询 1688 供应商信息失败: %v\n", err)
			return
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "登录ID: %s\n", result.LoginID)
		fmt.Fprintf(&sb, "供应商名称: %s\n", result.SupplierName)
		fmt.Fprintf(&sb, "公司名称: %s\n", result.CompanyName)
		fmt.Fprintf(&sb, "类目: %s\n", result.CategoryName)
		fmt.Fprintf(&sb, "店铺地址: %s\n", result.ShopURL)
		if result.KuaJingBao {
			fmt.Fprintf(&sb, "跨境宝: 是\n")
		}
		cmd.Print(sb.String())
	},
}

// ali-mix-config 查询 1688 卖家混批设置
var aliMixConfigCmd = &cobra.Command{
	Use:   "ali-mix-config",
	Short: "查询 1688 卖家混批设置",
	Long:  "根据卖家 memberId 或 loginId 查询 1688 卖家混批设置",
	Run: func(cmd *cobra.Command, args []string) {
		memberID, _ := cmd.Flags().GetString("member-id")
		loginID, _ := cmd.Flags().GetString("login-id")

		result, err := apis.NewThirdAppApi().GetAliSellerMixConfig(memberID, loginID)
		if err != nil {
			cmd.PrintErrf("查询混批设置失败: %v\n", err)
			return
		}
		d := result.Data
		fmt.Printf("卖家: %s\n", d.MemberID)
		fmt.Printf("普通混批: %v\n", d.GeneralHunpi)
		fmt.Printf("混批金额: %d\n", d.MixAmount)
		fmt.Printf("混批数量: %d\n", d.MixNumber)
	},
}

// ali-sub-accounts 查询 1688 子账号列表
var aliSubAccountsCmd = &cobra.Command{
	Use:   "ali-sub-accounts",
	Short: "查询 1688 子账号列表",
	Long:  "查询当前 1688 账号下的所有子账号",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := apis.NewThirdAppApi().GetAliSubAccountList()
		if err != nil {
			cmd.PrintErrf("查询子账号列表失败: %v\n", err)
			return
		}
		fmt.Printf("主账号: %s (%s)\n", result.MainLoginID, result.MainMemberID)
		if len(result.SubAccountList) == 0 {
			fmt.Println("无子账号")
			return
		}
		fmt.Println("|子账号|MemberID|")
		fmt.Println("|----|----|")
		for _, acc := range result.SubAccountList {
			fmt.Printf("|%s|%s|\n", acc.LoginID, acc.MemberID)
		}
	},
}

// ============ 店铺查询子命令 ============

// shop-product 查询店铺产品
var shopProductCmd = &cobra.Command{
	Use:   "shop-product",
	Short: "查询店铺产品",
	Long:  "根据店铺ID查询产品信息",
	Run: func(cmd *cobra.Command, args []string) {
		storeID, _ := cmd.Flags().GetInt("store-id")
		if storeID == 0 {
			cmd.PrintErr("--store-id 是必需的")
			return
		}
		result, err := apis.NewThirdAppApi().GetShopProduct(apis.AppApiShopParam{
			StoreID:  storeID,
			PageNum:  1,
			PageSize: 20,
		})
		if err != nil {
			cmd.PrintErrf("查询店铺产品失败: %v\n", err)
			return
		}
		if result.Data != nil {
			cmd.Println(string(result.Data))
		} else {
			fmt.Printf("查询失败: %s\n", result.Message)
		}
	},
}

// shop-order-list 查询店铺订单列表
var shopOrderListCmd = &cobra.Command{
	Use:   "shop-order-list",
	Short: "查询店铺订单列表",
	Long:  "根据店铺ID查询订单列表，可按状态筛选",
	Run: func(cmd *cobra.Command, args []string) {
		storeID, _ := cmd.Flags().GetInt("store-id")
		status, _ := cmd.Flags().GetString("status")
		if storeID == 0 {
			cmd.PrintErr("--store-id 是必需的")
			return
		}
		result, err := apis.NewThirdAppApi().GetShopOrderList(apis.AppApiShopParam{
			StoreID:  storeID,
			Status:   status,
			PageNum:  1,
			PageSize: 20,
		})
		if err != nil {
			cmd.PrintErrf("查询店铺订单失败: %v\n", err)
			return
		}
		if result.Data != nil {
			cmd.Println(string(result.Data))
		} else {
			fmt.Printf("查询失败: %s\n", result.Message)
		}
	},
}

// shop-inventory 查询店铺库存
var shopInventoryCmd = &cobra.Command{
	Use:   "shop-inventory",
	Short: "查询店铺库存",
	Long:  "查询店铺库存位置和库存商品项的关联信息",
	Run: func(cmd *cobra.Command, args []string) {
		storeID, _ := cmd.Flags().GetInt("store-id")
		if storeID == 0 {
			cmd.PrintErr("--store-id 是必需的")
			return
		}
		result, err := apis.NewThirdAppApi().GetShopInventoryLevels(apis.AppApiShopParam{
			StoreID:  storeID,
			PageNum:  1,
			PageSize: 50,
		})
		if err != nil {
			cmd.PrintErrf("查询店铺库存失败: %v\n", err)
			return
		}
		if result.Data != nil {
			fmt.Println("|位置ID|库存项ID|可用数量|")
			fmt.Println("|----|----|----|")
			for _, l := range result.Data {
				fmt.Printf("|%d|%d|%d|\n", l.LocationID, l.InventoryItemID, l.Available)
			}
		} else {
			fmt.Printf("查询失败: %s\n", result.Message)
		}
	},
}

func init() {
	// 1688 子命令
	thirdAppApiCmd.AddCommand(aliLogisticsCmd)
	thirdAppApiCmd.AddCommand(aliLogisticsTraceCmd)
	thirdAppApiCmd.AddCommand(aliOrderDetailCmd)
	thirdAppApiCmd.AddCommand(aliRefundDetailCmd)
	thirdAppApiCmd.AddCommand(aliRefundListCmd)
	thirdAppApiCmd.AddCommand(aliRefundOperationsCmd)

	aliProductCmd.Flags().Int64("id", 0, "1688 产品ID")
	aliProductCmd.Flags().String("url", "", "1688 产品链接")
	thirdAppApiCmd.AddCommand(aliProductCmd)

	aliSupplierCmd.Flags().String("login-id", "", "1688 供应商登录ID")
	aliSupplierCmd.Flags().String("domain", "", "供应商店铺地址")
	thirdAppApiCmd.AddCommand(aliSupplierCmd)

	aliMixConfigCmd.Flags().String("member-id", "", "卖家 memberId")
	aliMixConfigCmd.Flags().String("login-id", "", "卖家 loginId")
	thirdAppApiCmd.AddCommand(aliMixConfigCmd)

	thirdAppApiCmd.AddCommand(aliSubAccountsCmd)

	// 店铺查询子命令
	shopProductCmd.Flags().Int("store-id", 0, "店铺ID（必填）")
	shopOrderListCmd.Flags().Int("store-id", 0, "店铺ID（必填）")
	shopOrderListCmd.Flags().String("status", "any", "订单状态: open/closed/cancelled/any")
	shopInventoryCmd.Flags().Int("store-id", 0, "店铺ID（必填）")

	thirdAppApiCmd.AddCommand(shopProductCmd)
	thirdAppApiCmd.AddCommand(shopOrderListCmd)
	thirdAppApiCmd.AddCommand(shopInventoryCmd)

	rootCmd.AddCommand(thirdAppApiCmd)
}
