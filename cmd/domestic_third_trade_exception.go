package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

var (
	hyperskustatusMap = map[int]string{
		1:  "未发货",
		2:  "假发货",
		3:  "未到货",
		4:  "假签收",
		5:  "未签收",
		6:  "退件",
		7:  "丢件",
		8:  "未入库",
		9:  "丢包裹",
		10: "无货",
	}

	hyperskuSubStatusMap = map[int]string{
		1: "待处理",
		2: "处理中",
		3: "已处理",
		4: "已关闭",
		5: "已拒绝",
	}

	boolMap = map[bool]string{
		true:  "是",
		false: "否",
	}
)

var domesticThirdTradeExceptionCmd = &cobra.Command{
	Use:   "domestic-third-trade-exception",
	Short: "国内第三方交易异常订单管理",
	Long: `国内第三方交易异常订单管理命令集。HyperSKU 在第三方平台（1688、淘宝等）采购产生的采购单，其国内段物流发生异常（丢包裹、丢件、未签收等）时在此监控管理。

可用子命令：
  page-list      分页查询国内第三方交易异常订单（含物流明细）
  message-list   查询国内第三方交易异常订单的留言列表`,
}

// 分页查询国内第三方交易异常订单
var domesticThirdTradeExceptionPageListCmd = func() *cobra.Command {
	var query apis.DomesticThirdTradeExceptionQuery

	cmd := &cobra.Command{
		Use:   "page-list",
		Short: "分页查询国内第三方交易异常订单（含物流明细）",
		Long: `按 HyperSKU 异常状态分页查询国内第三方交易异常订单（第三方平台采购单，如 1688/淘宝），每条记录包含订单信息及对应的国内段物流明细（监控单号、物流单号、物流公司、最新状态等）。

必填参数：
  --hypersku-status       异常主状态
  --hypersku-sub-status   异常子状态列表

示例：
  hypersku-cli domestic-third-trade-exception page-list --hypersku-status 9 --hypersku-sub-status 1,2`,
		Run: func(cmd *cobra.Command, args []string) {
			if query.HyperskuStatus == 0 {
				cmd.PrintErr("参数 'hypersku-status' 是必需的，可选值：")
				for k, v := range hyperskustatusMap {
					cmd.PrintErrf("  %d - %s\n", k, v)
				}
				cmd.Help()
				return
			}

			res, err := apis.NewDomesticThirdTradeExceptionApi().PageList(query)
			if err != nil {
				cmd.PrintErrf("分页查询国内第三方交易异常订单失败: %v\n", err)
				return
			}

			if res == nil || res.Data == nil || len(res.Data.Rows) == 0 {
				total := int64(0)
				if res != nil && res.Data != nil {
					total = res.Data.Total
				}
				cmd.Printf("当前页码：%d，页大小：%d，总数：%d\n\n无数据", query.Page, query.Limit, total)
				return
			}

			total := res.Data.Total
			rows := res.Data.Rows

			warehouseStatusMap := map[int]string{
				1: "已签收",
				2: "已入库",
			}

			sourceTypeMap := map[int]string{
				1:   "1688",
				18:  "淘宝",
				19:  "天猫",
				17:  "京东",
				25:  "拼多多",
				40:  "小红书",
				99:  "Hyper From",
				999: "其他",
			}

			var sb strings.Builder

			columns := []string{
				"MonitorOrderId",
				"MonitorLogisticsId",
				"交易号",
				"仓库",
				"订单付款时间",
				"商品数量",
				"物流单号",
				"物流公司",
				"物流发货时间",
				"签收仓库",
				"仓库状态",
				"仓库签收时间",
				"创建时间",
				"更新时间",
				"状态",
				"采购来源",
				"采购订单类型",
				"是否有退款记录",
			}

			fmt.Fprintf(&sb, "当前页码：%d，页大小：%d，总数：%d\n", query.Page, query.Limit, total)
			fmt.Fprintf(&sb, "|%s|\n", strings.Join(columns, "|"))
			fmt.Fprintf(&sb, "%s|\n", strings.Repeat("|----", len(columns)))
			// 输出每个订单的物流信息
			for _, row := range rows {

				if len(row.LogisticsList) == 0 {
					continue
				}

				for _, lg := range row.LogisticsList {

					// 订单类型
					orderTypeStrs := make([]string, 0)
					for _, t := range strings.Split(lg.OrderTypes, ",") {
						if val, err := strconv.Atoi(t); err == nil {
							orderTypeStrs = append(orderTypeStrs, OrderType[val])
						}
					}

					vals := []string{
						fmt.Sprint(row.ID),
						fmt.Sprint(lg.ID),
						fmt.Sprint(lg.ThirdOrderID),
						row.WarehouseName,
						row.PaymentTime,
						fmt.Sprint(row.Quantity),
						lg.TrackingNumber,
						lg.LogisticsCompanyName,
						lg.DeliveredTime,
						lg.SignedWarehouseName,
						warehouseStatusMap[lg.WarehouseStatus],
						lg.WarehouseSignedTime,
						lg.CrtTime,
						lg.HyperskuStatusChangeTime,
						hyperskuSubStatusMap[lg.HyperskuSubStatus],
						sourceTypeMap[row.SourceType],
						strings.Join(orderTypeStrs, ","),
						boolMap[row.HasRefundRecords],
					}

					fmt.Fprintf(&sb, "|%s|\n", strings.Join(vals, "|"))
				}
			}

			cmd.Print(sb.String())
		},
	}

	cmd.Flags().IntVarP(&query.Page, "page", "p", 1, "页码（从 1 开始）")
	cmd.Flags().IntVarP(&query.Limit, "limit", "l", 10, "每页返回的条数")
	cmd.Flags().IntVarP(&query.HyperskuStatus, "hypersku-status", "s", 0, "HyperSKU 异常状态，可选值：1-未发货，2-假发货，3-未到货，4-假签收，5-未签收，6-退件，7-丢件，8-未入库，9-丢包裹，10-无货")
	cmd.Flags().IntSliceVarP(&query.HyperskuSubStatusList, "hypersku-sub-status", "c", []int{1, 2}, "HyperSKU 异常子状态列表，可选值：1-待处理，2-处理中，3-已处理，4-已关闭，5-已拒绝")
	cmd.Flags().StringVarP(&query.BuyerID, "buyer-id", "b", "", "买家 ID（可选，仅返回该买家的异常订单）")

	// cmd.MarkFlagRequired("page")
	// cmd.MarkFlagRequired("limit")
	cmd.MarkFlagRequired("hypersku-status")
	cmd.MarkFlagRequired("hypersku-sub-status")
	return cmd
}()

// 获取国内第三方交易异常订单留言列表
var domesticThirdTradeExceptionMessageListCmd = &cobra.Command{
	Use:   "message-list [monitorOrderId] [monitorLogisticsId]",
	Short: "查询国内第三方交易异常订单留言列表",
	Long:  `根据监控订单ID（MonitorOrderId）和监控物流ID（MonitorLogisticsId）查询国内第三方交易异常订单的留言记录，可用于跟进异常处理情况。`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		monitorOrderId := args[0]
		monitorLogisticsId := args[1]

		messages, err := apis.NewDomesticThirdTradeExceptionApi().MessageList(monitorOrderId, monitorLogisticsId)
		if err != nil {
			cmd.PrintErrf("获取订单 %s 物流 %s 留言列表失败: %v\n", monitorOrderId, monitorLogisticsId, err)
			return
		}

		if messages == nil || len(*messages) == 0 {
			cmd.Print("未查询到该订单的留言")
			return
		}

		contentList := make([]string, len(*messages)+2)
		contentList[0] = "|留言时间|留言人|留言|"
		contentList[1] = "|----|----|----|"
		for i, item := range *messages {
			contentList[i+2] = fmt.Sprintf("|%s|%s|%s|", item.CrtTime, item.CrtName, item.Remark)
		}

		cmd.Print(strings.Join(contentList, "\n"))
	},
}

func init() {
	domesticThirdTradeExceptionCmd.AddCommand(domesticThirdTradeExceptionPageListCmd, domesticThirdTradeExceptionMessageListCmd)
	rootCmd.AddCommand(domesticThirdTradeExceptionCmd)
}
