package cmd

import (
	"fmt"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

var logistics = &cobra.Command{
	Use:   "logistics",
	Short: "物流管理",
	Long:  "物流管理命令，提供物流查询等操作",
}

var getTracking = &cobra.Command{
	Use:   "tracking [trackingNumber]",
	Short: "查询物流轨迹",
	Long:  "根据物流单号查询物流轨迹",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		trackingNumber := args[0]
		trackList, err := apis.NewLogisticsApi().GetTracking(trackingNumber)
		if err != nil {
			cmd.PrintErr("查询物流轨迹发生错误", err)
			cmd.Help()
			return
		}

		if trackList == nil {
			cmd.Print("未查询到物流轨迹")
			return
		}

		trackContent := make([]string, len(*trackList)+2)
		trackContent[0] = "|时间|轨迹|"
		trackContent[1] = "|----|----|"
		for i, item := range *trackList {
			trackContent[i+2] = fmt.Sprintf("|%s|%s|", item.Time, item.Thing)
		}

		cmd.Print(strings.Join(trackContent, "\n"))
	},
}

func init() {
	logistics.AddCommand(getTracking)
	rootCmd.AddCommand(logistics)
}
