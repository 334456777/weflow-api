package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/334456777/weflow-api/internal/client"
)

var (
	snsExportOutputDir  string
	snsExportFormat     string
	snsExportUsernames  string
	snsExportKeyword    string
	snsExportMedia      bool
	snsExportImages     bool
	snsExportLivePhotos bool
	snsExportVideos     bool
	snsExportStart      string
	snsExportEnd        string
	snsExportTimeout    time.Duration
)

var snsExportCmd = &cobra.Command{
	Use:   "export",
	Short: "POST 导出朋友圈（§7.5）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if snsExportOutputDir == "" {
			return fmt.Errorf("--output-dir 必填")
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		if snsExportTimeout > 0 {
			c.SetHTTPClient(&http.Client{Timeout: snsExportTimeout})
		}
		req := client.SnsExportRequest{
			OutputDir:        snsExportOutputDir,
			Format:           snsExportFormat,
			Usernames:        snsExportUsernames,
			Keyword:          snsExportKeyword,
			ExportMedia:      snsExportMedia,
			ExportImages:     snsExportImages,
			ExportLivePhotos: snsExportLivePhotos,
			ExportVideos:     snsExportVideos,
			Start:            snsExportStart,
			End:              snsExportEnd,
		}
		raw, err := c.SnsExport(cmd.Context(), req)
		if err != nil {
			return err
		}
		return newRenderer().Raw(raw)
	},
}

func init() {
	f := snsExportCmd.Flags()
	f.StringVar(&snsExportOutputDir, "output-dir", "", "导出目录（必填）")
	f.StringVar(&snsExportFormat, "format", "json", "导出格式：json|html|arkmejson")
	f.StringVar(&snsExportUsernames, "usernames", "", "发布者过滤，逗号分隔")
	f.StringVar(&snsExportKeyword, "keyword", "", "正文关键词")
	f.BoolVar(&snsExportMedia, "export-media", false, "导出媒体")
	f.BoolVar(&snsExportImages, "export-images", false, "导出图片")
	f.BoolVar(&snsExportLivePhotos, "export-live-photos", false, "导出实况图")
	f.BoolVar(&snsExportVideos, "export-videos", false, "导出视频")
	f.StringVar(&snsExportStart, "start", "", "起始时间")
	f.StringVar(&snsExportEnd, "end", "", "结束时间")
	f.DurationVar(&snsExportTimeout, "export-timeout", 10*time.Minute, "导出请求超时（覆盖全局 --timeout，默认 10m）")
	_ = snsExportCmd.MarkFlagRequired("output-dir")
	snsCmd.AddCommand(snsExportCmd)
}
